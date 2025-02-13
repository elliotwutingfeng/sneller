// Copyright (C) 2022 Sneller, Inc.
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package blockfmt

import (
	"encoding/base32"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/SnellerInc/sneller/aws/s3"
	"github.com/SnellerInc/sneller/fsutil"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/exp/maps"
)

// InputFS describes the FS implementation
// that is required for reading inputs.
type InputFS interface {
	fs.FS
	fsutil.ETagFS

	// Prefix should return a string
	// that is prepended to filesystem
	// paths to indicate the filesystem "origin."
	//
	// For example, an S3 bucket FS would have
	//   s3://bucket/
	// as its prefix.
	Prefix() string
}

// UploadFS describes the FS implementation
// that is required for writing outputs.
type UploadFS interface {
	InputFS

	// WriteFile should create the
	// file at path with the given contents.
	// If the file already exists, it should
	// be overwritten atomically.
	// WriteFile should return the ETag associated
	// with the written file along with the first encountered error.
	WriteFile(path string, buf []byte) (etag string, err error)

	// Create should create an Uploader
	// for the given path. The file should
	// not be visible at the provided path
	// until the Uploader has been closed
	// successfully.
	Create(path string) (Uploader, error)
}

// S3FS implements UploadFS and InputFS.
type S3FS struct {
	s3.BucketFS
}

// Prefix implements InputFS.Prefix
func (s *S3FS) Prefix() string {
	return "s3://" + s.Bucket + "/"
}

// ETag implements InputFS.ETag
func (s *S3FS) ETag(fullpath string, f fs.FileInfo) (string, error) {
	if rd, ok := f.(*s3.File); ok {
		return rd.ETag, nil
	}
	return "", fmt.Errorf("cannot produce ETag for %T", f)
}

// Create implements UploadFS.Create
func (s *S3FS) Create(path string) (Uploader, error) {
	up := &s3.Uploader{
		Key:    s.Key,
		Bucket: s.Bucket,
		Object: path,
	}
	err := up.Start()
	if err != nil {
		return nil, err
	}
	return up, nil
}

// WriteFile implements UploadFS.WriteFile
func (s *S3FS) WriteFile(path string, contents []byte) (string, error) {
	return s.Put(path, contents)
}

// NewDirFS creates a new DirFS in dir.
func NewDirFS(dir string) *DirFS {
	return &DirFS{
		FS:   os.DirFS(dir),
		Root: dir,
	}
}

// DirFS is an InputFS and UploadFS
// that is rooted in a particular directory.
type DirFS struct {
	fs.FS
	Root        string
	Log         func(f string, args ...interface{})
	MinPartSize int
}

func hashFile(r io.Reader) (string, error) {
	h, err := blake2b.New256(nil)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(h, r)
	if err != nil {
		return "", err
	}
	return "\"b2sum:" + base32.StdEncoding.EncodeToString(h.Sum(nil)) + `"`, nil
}

// Mmap maps the file given by [fullpath].
// The caller must ensure that the returned slice,
// if non-nil, is unmapped with [Unmap].
// NOTE: Mmap is not supported on all platforms.
// The caller should be prepared to handle the error and
// fall back to ordinary [Open] and [Read] calls.
func (d *DirFS) Mmap(fullpath string) ([]byte, error) {
	return mmap(filepath.Join(d.Root, filepath.FromSlash(fullpath)))
}

// Unmap unmaps a buffer returned by [Mmap].
func (d *DirFS) Unmap(buf []byte) error {
	return unmap(buf)
}

// Prefix implements InputFS.Prefix
func (d *DirFS) Prefix() string {
	return "file://"
}

// ETag implements InputFS.ETag
func (d *DirFS) ETag(fullpath string, info fs.FileInfo) (string, error) {
	fullpath = path.Clean(fullpath)
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("cannot get ETag of non-regular file %s", fullpath)
	}
	f, err := d.Open(fullpath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return hashFile(f)
}

// Remove removes the file at the specified path.
func (d *DirFS) Remove(fullpath string) error {
	fullpath = path.Clean(fullpath)
	if !fs.ValidPath(fullpath) {
		return fmt.Errorf("%s: %s", fullpath, fs.ErrInvalid)
	}
	return os.Remove(filepath.Join(d.Root, fullpath))
}

// WriteFile implements UploadFS.WriteFile
func (d *DirFS) WriteFile(fullpath string, buf []byte) (string, error) {
	if d.Log != nil {
		d.Log("WriteFile %s", fullpath)
	}
	if !fs.ValidPath(fullpath) {
		return "", fs.ErrInvalid
	}
	fullpath = filepath.Clean(filepath.Join(d.Root, fullpath))
	dir, base := filepath.Split(fullpath)
	if dir == "" {
		dir = "."
	}
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp(dir, base)
	if err != nil {
		if d.Log != nil {
			d.Log("CreateTemp: %s", err)
		}
		return "", err
	}
	_, err = tmp.Write(buf)
	tmp.Close()
	if err != nil {
		os.Remove(tmp.Name())
		return "", err
	}
	err = os.Rename(tmp.Name(), fullpath)
	if err != nil {
		os.Remove(tmp.Name())
		return "", err
	}
	ret := blake2b.Sum256(buf)
	return "\"b2sum:" + base32.StdEncoding.EncodeToString(ret[:]) + `"`, nil
}

// fileUploader is a BufferUploader wrapper
// that simulates multi-part uploads to
// a single file (by buffering the whole
// output to memory and then performing
// a WriteFile)
type fileUploader struct {
	lock  sync.Mutex
	fp    string
	dir   *DirFS
	parts map[int64]string
	size  int64
	etag  string
}

func (f *fileUploader) MinPartSize() int {
	if f.dir.MinPartSize <= 0 {
		return 1024 * 1024
	}
	return f.dir.MinPartSize
}

func (f *fileUploader) Size() int64 { return f.size }

// Abort cleans up any state produced by calls to Upload.
// Abort will have no effect if Close has already been called
// and it returned successfully.
func (f *fileUploader) Abort() error {
	for _, name := range f.parts {
		os.Remove(name)
	}
	maps.Clear(f.parts)
	f.size = 0
	return nil
}

func (f *fileUploader) ETag() string {
	return f.etag
}

func (f *fileUploader) Close(final []byte) error {
	f.lock.Lock()
	defer f.lock.Unlock()
	fullpath := filepath.Join(f.dir.Root, f.fp)
	err := os.MkdirAll(filepath.Dir(fullpath), 0750)
	if err != nil {
		return err
	}
	file, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	abort := func() {
		file.Close()
		os.Remove(file.Name())
	}
	// emit parts in sorted order
	parts := maps.Keys(f.parts)
	slices.Sort(parts)
	h, err := blake2b.New256(nil)
	if err != nil {
		panic(err)
	}
	dst := io.MultiWriter(file, h)
	for _, part := range parts {
		name := f.parts[part]
		part, err := os.Open(name)
		if err != nil {
			abort()
			return fmt.Errorf("DirFS uploader: opening part: %w", err)
		}
		nn, err := io.Copy(dst, part)
		part.Close()
		os.Remove(name)
		if err != nil {
			abort()
			return fmt.Errorf("DirFS uploader: appending data: %w", err)
		}
		f.size += nn
	}
	f.size += int64(len(final))
	_, err = dst.Write(final)
	if err != nil {
		abort()
		return err
	}
	f.etag = "\"b2sum:" + base32.StdEncoding.EncodeToString(h.Sum(nil)) + `"`
	return file.Close()
}

func (f *fileUploader) Upload(part int64, contents []byte) error {
	if len(contents) < f.MinPartSize() {
		return fmt.Errorf("upload part %d: len(contents)=%d; MinPartSize = %d", part, len(contents), f.MinPartSize())
	}
	fullpath := filepath.Join(f.dir.Root, f.fp)
	dir, name := filepath.Split(fullpath)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}
	file, err := os.CreateTemp(dir, name+".part.*")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(contents)
	if err != nil {
		os.Remove(file.Name())
		return err
	}
	f.lock.Lock()
	f.parts[part] = file.Name()
	f.lock.Unlock()
	return nil
}

// Create implements UploadFS.Create
func (d *DirFS) Create(fullpath string) (Uploader, error) {
	if d.Log != nil {
		d.Log("Create %s", fullpath)
	}
	fullpath = path.Clean(fullpath)
	if !fs.ValidPath(fullpath) {
		return nil, fs.ErrInvalid
	}
	return &fileUploader{
		parts: make(map[int64]string),
		fp:    fullpath,
		dir:   d,
	}, nil
}

// ETag gets the ETag for the provided Uploader.
// If the Uploader has an ETag() method, that method
// is used directly; otherwise ofs.ETag is used to
// determine the ETag.
func ETag(ofs UploadFS, up Uploader, fullpath string) (string, error) {
	type etagger interface {
		ETag() string
	}
	if et, ok := up.(etagger); ok {
		return et.ETag(), nil
	}
	info, err := fs.Stat(ofs, fullpath)
	if err != nil {
		return "", fmt.Errorf("blockfmt.ETag: %w", err)
	}
	return ofs.ETag(fullpath, info)
}

var (
	_ InputFS  = &DirFS{}
	_ UploadFS = &DirFS{}
	_ InputFS  = &S3FS{}
	_ UploadFS = &S3FS{}
)

func inferFormat(name string, fallback func(name string) RowFormat) RowFormat {
	for suff, cons := range SuffixToFormat {
		if strings.HasSuffix(name, suff) {
			f, _ := cons(nil)
			return f
		}
	}
	if fallback == nil {
		return nil
	}
	return fallback(name)
}

// Collector is a set of configuration
// parameters for collecting a list of objects.
type Collector struct {
	// Pattern is the glob pattern
	// that input objects should match.
	Pattern string
	// Start is a filename below which
	// all inputs are ignored. (Start can
	// be used to begin a Collect operation
	// where a previous one has left off
	// by using the last returned path as
	// the Start value for the next collection operation.)
	Start string
	// MaxItems, if non-zero, is the maximum number
	// of items to collect.
	MaxItems int
	// MaxSize, if non-zero, is the maximum size
	// of items to collect.
	MaxSize int64
	// Fallback is the function used to
	// determine the format of an input file.
	Fallback func(string) RowFormat
}

var errStop = errors.New("stop listing")

// Collect collects items from the provided
// InputFS and returns them as a list of Inputs,
// along with a boolean indicating whether or not
// the results are the complete list of files.
func (c *Collector) Collect(from InputFS) ([]Input, bool, error) {
	size := int64(0)
	prefix := from.Prefix()
	var have []Input
	walk := func(p string, f fs.File, err error) error {
		if err != nil {
			return err
		}
		info, err := f.Stat()
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				// race between readdir and stat
				return nil
			}
			return err
		}
		etag, err := from.ETag(p, info)
		if err != nil {
			return err
		}
		format := inferFormat(p, c.Fallback)
		have = append(have, Input{
			Path: prefix + p,
			ETag: etag,
			Size: info.Size(),
			R:    f,
			F:    format,
		})
		if c.MaxItems > 0 && len(have) >= c.MaxItems {
			return errStop
		}
		if c.MaxSize > 0 && size >= c.MaxSize {
			return errStop
		}
		return nil
	}
	err := fsutil.WalkGlob(from, c.Start, c.Pattern, walk)
	if err == errStop {
		return have, false, nil
	}
	if err != nil {
		return have, false, err
	}
	return have, true, nil
}

// CollectGlob turns a glob pattern
// into a list of Inputs, using fallback
// as the constructor for the RowFormat
// of each input object when the object
// suffix does not match any of the known
// format suffixes. If any of the files
// that match the glob pattern do
// not have known file suffixes and
// fallback does not return a non-nil RowFormat
// for those files, then CollectGlob will return
// an error indicating that the format for the file
// could not be determined.
func CollectGlob(ifs InputFS, fallback func(string) RowFormat, pattern string) ([]Input, error) {
	cl := Collector{
		Pattern:  pattern,
		Fallback: fallback,
	}
	ret, _, err := cl.Collect(ifs)
	return ret, err
}
