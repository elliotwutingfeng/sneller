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

package ion

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"

	"golang.org/x/exp/slices"

	"github.com/SnellerInc/sneller/date"
)

// Stop can be returned by the function passed
// to List.Each and Struct.Each to stop
// iterating and return early.
//
//lint:ignore ST1012 sentinel error
var Stop = errors.New("stop early")

// Datum represents any Ion datum.
//
// The Marshal and Unmarshal functions natively
// understand that Datum can be constructed and
// re-encoded from any ion value.
//
// A Datum should be a value returned by
//
//	Float, Int, Uint, Struct, List, Bool,
//	BigInt, Timestamp, Annotation, ..., or ReadDatum.
type Datum struct {
	st  []string
	buf []byte
}

func rawDatum(st *Symtab, b []byte) Datum {
	d := Datum{buf: b[:SizeOf(b)]}
	if st != nil {
		d.st = st.alias()
	}
	return d
}

// Empty is the zero value of a Datum.
var Empty = Datum{}

func (d Datum) Clone() Datum {
	return Datum{
		st:  d.st, // no need to clone
		buf: slices.Clone(d.buf),
	}
}

// Equal returns whether d and x are
// semantically equivalent.
func (d Datum) Equal(x Datum) bool {
	switch d.Type() {
	case NullType:
		return x.IsNull()
	case FloatType:
		if x.IsFloat() {
			d, _ := d.Float()
			x, _ := x.Float()
			return x == d || (math.IsNaN(d) && math.IsNaN(x))
		}
		if x.IsInt() {
			d, _ := d.Float()
			x, _ := x.Int()
			return float64(int64(d)) == float64(d) && int64(d) == int64(x)
		}
		if x.IsUint() {
			d, _ := d.Float()
			x, _ := x.Uint()
			return float64(uint64(d)) == float64(d) && uint64(d) == uint64(x)
		}
	case IntType:
		if x.IsInt() {
			d, _ := d.Int()
			x, _ := x.Int()
			return x == d
		}
		if x.IsUint() {
			d, _ := d.Int()
			x, _ := x.Uint()
			return d >= 0 && uint64(d) == x
		}
		if x.IsFloat() {
			d, _ := d.Int()
			x, _ := x.Float()
			return float64(int64(x)) == float64(x) && int64(x) == int64(d)
		}
	case UintType:
		if x.IsUint() {
			d, _ := d.Uint()
			x, _ := x.Uint()
			return d == x
		}
		if x.IsInt() {
			d, _ := d.Uint()
			x, _ := x.Int()
			return x >= 0 && uint64(x) == d
		}
		if x.IsFloat() {
			d, _ := d.Uint()
			x, _ := x.Float()
			return float64(uint64(x)) == float64(x) && uint64(x) == uint64(d)
		}
	case StructType:
		if x.IsStruct() {
			d, _ := d.Struct()
			x, _ := x.Struct()
			return d.Equal(x)
		}
		return false
	case ListType:
		if x.IsList() {
			d, _ := d.List()
			x, _ := x.List()
			return d.Equal(x)
		}
		return false
	case BoolType:
		if x.IsBool() {
			d, _ := d.Bool()
			x, _ := x.Bool()
			return d == x
		}
		return false
	case StringType:
		if x.IsString() {
			d, _ := d.String()
			x, _ := x.StringShared()
			return d == string(x)
		}
		if x.IsSymbol() {
			d, _ := d.String()
			x, _ := x.String()
			return d == x
		}
	case SymbolType:
		if x.IsString() {
			d, _ := d.String()
			x, _ := x.StringShared()
			return d == string(x)
		}
		if x.IsSymbol() {
			d, _ := d.String()
			x, _ := x.String()
			return d == x
		}
	case BlobType:
		if x.IsBlob() {
			d, _ := d.BlobShared()
			x, _ := x.BlobShared()
			return string(d) == string(x)
		}
	case TimestampType:
		if x.IsTimestamp() {
			d, _ := d.Timestamp()
			x, _ := x.Timestamp()
			return d.Equal(x)
		}
	}
	return false
}

// LessImprecise compares the raw Ion bytes.
//
// This method does not order correctly equal
// values having different binary representation.
// For instance a string can be saved literally,
// as a sequence of UTF-8 bytes, or be a symbol,
// that is a reference to the symbol table.
func (d Datum) LessImprecise(x Datum) bool {
	return bytes.Compare(d.buf, x.buf) < 0
}

func (d Datum) Type() Type {
	if len(d.buf) == 0 {
		return InvalidType
	}
	return TypeOf(d.buf)
}

type resymbolizer struct {
	// idmap is a cache of old->new symbol mappings
	idmap  []Symbol
	srctab *Symtab
	dsttab *Symtab
}

func (r *resymbolizer) reset() {
	for i := range r.idmap {
		r.idmap[i] = 0
	}
}

func (r *resymbolizer) get(sym Symbol) Symbol {
	if int(sym) < len(r.idmap) && r.idmap[sym] != 0 {
		return r.idmap[sym]
	}
	if int(sym) >= len(r.idmap) {
		if cap(r.idmap) > int(sym) {
			r.idmap = r.idmap[:int(sym)+1]
		} else {
			newmap := make([]Symbol, int(sym)+1)
			copy(newmap, r.idmap)
			r.idmap = newmap
		}
	}
	r.idmap[sym] = r.dsttab.Intern(r.srctab.Get(sym))
	return r.idmap[sym]
}

func (d Datum) Encode(dst *Buffer, st *Symtab) {
	// fast path: no need to resymbolize
	if len(d.st) == 0 || st.contains(d.st) {
		dst.UnsafeAppend(d.buf)
		return
	}
	srcsyms := d.symtab()
	rs := &resymbolizer{
		srctab: &srcsyms,
		dsttab: st,
	}
	rs.resym(dst, d.buf)
}

// performance-sensitive resymbolization path
func (r *resymbolizer) resym(dst *Buffer, buf []byte) {
	switch TypeOf(buf) {
	case SymbolType:
		sym, _, _ := ReadSymbol(buf)
		dst.WriteSymbol(r.get(sym))
	case StructType:
		dst.BeginStruct(-1)
		body, _ := Contents(buf)
		var sym Symbol
		for len(body) > 0 {
			sym, body, _ = ReadLabel(body)
			dst.BeginField(r.get(sym))
			size := SizeOf(body)
			r.resym(dst, body[:size])
			body = body[size:]
		}
		dst.EndStruct()
	case ListType:
		dst.BeginList(-1)
		body, _ := Contents(buf)
		for len(body) > 0 {
			size := SizeOf(body)
			r.resym(dst, body[:size])
			body = body[size:]
		}
		dst.EndList()
	case AnnotationType:
		sym, body, _, _ := ReadAnnotation(buf)
		dst.BeginAnnotation(1)
		dst.BeginField(r.get(sym))
		r.resym(dst, body)
		dst.EndAnnotation()
	default:
		dst.UnsafeAppend(buf)
	}
}

func (d Datum) IsEmpty() bool      { return len(d.buf) == 0 }
func (d Datum) IsNull() bool       { return d.Type() == NullType }
func (d Datum) IsFloat() bool      { return d.Type() == FloatType }
func (d Datum) IsInt() bool        { return d.Type() == IntType }
func (d Datum) IsUint() bool       { return d.Type() == UintType }
func (d Datum) IsStruct() bool     { return d.Type() == StructType }
func (d Datum) IsList() bool       { return d.Type() == ListType }
func (d Datum) IsAnnotation() bool { return d.Type() == AnnotationType }
func (d Datum) IsBool() bool       { return d.Type() == BoolType }
func (d Datum) IsSymbol() bool     { return d.Type() == SymbolType }
func (d Datum) IsString() bool     { return d.Type() == StringType }
func (d Datum) IsBlob() bool       { return d.Type() == BlobType }
func (d Datum) IsTimestamp() bool  { return d.Type() == TimestampType }

func (d Datum) Null() error                        { return d.null("") }
func (d Datum) Float() (float64, error)            { return d.float("") }
func (d Datum) Int() (int64, error)                { return d.int("") }
func (d Datum) Uint() (uint64, error)              { return d.uint("") }
func (d Datum) Struct() (Struct, error)            { return d.struc("") }
func (d Datum) List() (List, error)                { return d.list("") }
func (d Datum) Annotation() (string, Datum, error) { return d.annotation("") }
func (d Datum) Bool() (bool, error)                { return d.bool("") }
func (d Datum) Symbol() (Symbol, error)            { return d.symbol("") }
func (d Datum) String() (string, error)            { return d.string("") }
func (d Datum) Blob() ([]byte, error)              { return d.blob("") }
func (d Datum) Timestamp() (date.Time, error)      { return d.timestamp("") }

// StringShared returns a []byte aliasing the
// contents of this Datum and should be copied
// as necessary to avoid issues that may arise
// from retaining aliased bytes.
//
// Unlike String, this method will not work with
// a symbol datum.
func (d Datum) StringShared() ([]byte, error) { return d.stringShared("") }

// BlobShared returns a []byte aliasing the
// contents of this Datum and should be copied
// as necessary to avoid issues that may arise
// from retaining aliased bytes.
func (d Datum) BlobShared() ([]byte, error) { return d.blobShared("") }

func (d Datum) null(field string) error {
	if !d.IsNull() {
		return d.bad(field, NullType)
	}
	return nil
}

func (d Datum) float(field string) (float64, error) {
	if !d.IsFloat() {
		return 0, d.bad(field, FloatType)
	}
	f, _, err := ReadFloat64(d.buf)
	if err != nil {
		panic(err)
	}
	return f, nil
}

func (d Datum) int(field string) (int64, error) {
	if !d.IsInt() {
		return 0, d.bad(field, IntType)
	}
	i, _, err := ReadInt(d.buf)
	if err != nil {
		panic(err)
	}
	return i, nil
}

func (d Datum) uint(field string) (uint64, error) {
	if !d.IsUint() {
		return 0, d.bad(field, UintType)
	}
	u, _, err := ReadUint(d.buf)
	if err != nil {
		panic(err)
	}
	return u, nil
}

func (d Datum) struc(field string) (Struct, error) {
	if !d.IsStruct() {
		return Struct{}, d.bad(field, StructType)
	}
	return Struct{st: d.st, buf: d.buf}, nil
}

// Field returns the value associated with the
// field with the given name if d is a struct.
// If d is not a struct or the field is not
// present, this returns Empty.
func (d Datum) Field(name string) Datum {
	if !d.IsStruct() {
		return Empty
	}
	s, _ := d.Struct()
	f, ok := s.FieldByName(name)
	if !ok {
		return Empty
	}
	return f.Datum
}

func (d Datum) list(field string) (List, error) {
	if !d.IsList() {
		return List{}, d.bad(field, ListType)
	}
	return List{st: d.st, buf: d.buf}, nil
}

func (d Datum) annotation(field string) (string, Datum, error) {
	if !d.IsAnnotation() {
		return "", Empty, d.bad(field, AnnotationType)
	}
	sym, body, _, err := ReadAnnotation(d.buf)
	if err != nil {
		panic(err)
	}
	st := d.symtab()
	s, ok := st.Lookup(sym)
	if !ok {
		panic("ion.Datum.Annotation: missing symbol")
	}
	return s, Datum{st: d.st, buf: body}, nil
}

func (d Datum) bool(field string) (bool, error) {
	if !d.IsBool() {
		return false, d.bad(field, BoolType)
	}
	b, _, err := ReadBool(d.buf)
	if err != nil {
		panic(err)
	}
	return b, nil
}

func (d Datum) symbol(field string) (Symbol, error) {
	if !d.IsSymbol() {
		return 0, d.bad(field, SymbolType)
	}
	sym, _, err := ReadSymbol(d.buf)
	if err != nil {
		panic(err)
	}
	return sym, nil
}

func (d Datum) string(field string) (string, error) {
	if d.IsSymbol() {
		sym, _ := d.Symbol()
		st := d.symtab()
		s, ok := st.Lookup(sym)
		if !ok {
			panic("ion.Datum.String: missing symbol")
		}
		return s, nil
	}
	s, err := d.stringShared(field)
	return string(s), err
}

func (d Datum) stringShared(field string) ([]byte, error) {
	if !d.IsString() {
		return nil, d.bad(field, StringType)
	}
	s, _, err := ReadStringShared(d.buf)
	if err != nil {
		panic(err)
	}
	return s, nil
}

func (d Datum) blob(field string) ([]byte, error) {
	b, err := d.blobShared(field)
	return slices.Clone(b), err
}

func (d Datum) blobShared(field string) ([]byte, error) {
	if !d.IsBlob() {
		return nil, d.bad(field, BlobType)
	}
	b, _ := Contents(d.buf)
	if b == nil {
		panic("ion.Datum.Blob: invalid ion")
	}
	return b, nil
}

func (d Datum) timestamp(field string) (date.Time, error) {
	if !d.IsTimestamp() {
		return date.Time{}, d.bad(field, TimestampType)
	}
	t, _, err := ReadTime(d.buf)
	if err != nil {
		panic(err)
	}
	return t, nil
}

func (d *Datum) bad(field string, want Type) error {
	return &TypeError{Wanted: want, Found: d.Type(), Field: field}
}

func (d *Datum) symtab() Symtab {
	return Symtab{interned: d.st}
}

func Float(f float64) Datum {
	var buf Buffer
	buf.WriteFloat64(f)
	return Datum{buf: buf.Bytes()}
}

// Null is the untyped null datum.
var Null = Datum{buf: []byte{0x0f}}

func Int(i int64) Datum {
	var buf Buffer
	buf.WriteInt(i)
	return Datum{buf: buf.Bytes()}
}

func Uint(u uint64) Datum {
	var buf Buffer
	buf.WriteUint(u)
	return Datum{buf: buf.Bytes()}
}

// Field is a structure field in a Struct or Annotation datum
type Field struct {
	Label string
	Datum
	Sym Symbol // symbol, if assigned
}

func ReadField(st *Symtab, body []byte) (Field, []byte, error) {
	sym, body, err := ReadLabel(body)
	if err != nil {
		return Field{}, nil, err
	}
	name, ok := st.Lookup(sym)
	if !ok {
		return Field{}, nil, fmt.Errorf("symbol %d not in symbol table", sym)
	}
	val, rest, err := ReadDatum(st, body)
	if err != nil {
		return Field{}, nil, err
	}
	return Field{Label: name, Datum: val, Sym: sym}, rest, nil
}

func (f *Field) Encode(dst *Buffer, st *Symtab) {
	dst.BeginField(st.Intern(f.Label))
	f.Datum.Encode(dst, st)
}

func (f *Field) Equal(f2 *Field) bool {
	return f.Label == f2.Label && f.Sym == f2.Sym && f.Datum.Equal(f2.Datum)
}

func (f Field) Clone() Field {
	return Field{
		Label: f.Label,
		Datum: f.Datum.Clone(),
		Sym:   f.Sym,
	}
}

func (f Field) Null() error                        { return f.null(f.Label) }
func (f Field) Float() (float64, error)            { return f.float(f.Label) }
func (f Field) Int() (int64, error)                { return f.int(f.Label) }
func (f Field) Uint() (uint64, error)              { return f.uint(f.Label) }
func (f Field) Struct() (Struct, error)            { return f.struc(f.Label) }
func (f Field) List() (List, error)                { return f.list(f.Label) }
func (f Field) Annotation() (string, Datum, error) { return f.annotation(f.Label) }
func (f Field) Bool() (bool, error)                { return f.bool(f.Label) }
func (f Field) Symbol() (Symbol, error)            { return f.symbol(f.Label) }
func (f Field) String() (string, error)            { return f.string(f.Label) }
func (f Field) Blob() ([]byte, error)              { return f.blob(f.Label) }
func (f Field) Timestamp() (date.Time, error)      { return f.timestamp(f.Label) }

// StringShared returns a []byte aliasing the
// contents of this Field and should be copied
// as necessary to avoid issues that may arise
// from retaining aliased bytes.
//
// Unlike String, this method will not work with
// a symbol datum.
func (f Field) StringShared() ([]byte, error) { return f.stringShared(f.Label) }

// BlobShared returns a []byte aliasing the
// contents of this Field and should be copied
// as necessary to avoid issues that may arise
// from retaining aliased bytes.
func (f Field) BlobShared() ([]byte, error) { return f.blobShared(f.Label) }

type composite struct {
	st  []string
	buf []byte
	_   struct{} // disallow conversion to Datum
}

var emptyStruct = []byte{0xd0}

// Struct is an ion structure datum
type Struct composite

func NewStruct(st *Symtab, f []Field) Struct {
	if len(f) == 0 {
		return Struct{}
	}
	var dst Buffer
	if st == nil {
		st = &Symtab{}
	}
	dst.WriteStruct(st, f)
	return Struct{st: st.alias(), buf: dst.Bytes()}
}

func (b *Buffer) WriteStruct(st *Symtab, f []Field) {
	if len(f) == 0 {
		b.UnsafeAppend(emptyStruct)
		return
	}
	b.BeginStruct(-1)
	for i := range f {
		f[i].Encode(b, st)
	}
	b.EndStruct()
}

func (s Struct) Datum() Datum {
	if len(s.buf) == 0 {
		return Datum{buf: emptyStruct}
	}
	return Datum{st: s.st, buf: s.buf}
}

func (s Struct) Encode(dst *Buffer, st *Symtab) {
	// fast path: we can avoid resym
	if s.Empty() || st.contains(s.st) {
		dst.UnsafeAppend(s.bytes())
		return
	}
	dst.BeginStruct(-1)
	s.Each(func(f Field) error {
		f.Encode(dst, st)
		return nil
	})
	dst.EndStruct()
}

func (s Struct) Equal(s2 Struct) bool {
	if s.Empty() {
		return s2.Empty()
	}
	if bytes.Equal(s.buf, s2.buf) && stoverlap(s.st, s2.st) {
		return true
	}
	// TODO: optimize this
	f1 := s.Fields(nil)
	f2 := s2.Fields(nil)
	if len(f1) != len(f2) {
		return false
	}
	for i := range f1 {
		f1[i].Sym = 0
		f2[i].Sym = 0
	}
	slices.SortFunc(f1, func(x, y Field) bool {
		return x.Label < y.Label
	})
	slices.SortFunc(f2, func(x, y Field) bool {
		return x.Label < y.Label
	})
	for i := range f1 {
		if f1[i].Label != f2[i].Label {
			return false
		}
		if !Equal(f1[i].Datum, f2[i].Datum) {
			return false
		}
	}
	return true
}

func (s Struct) Len() int {
	if s.Empty() {
		return 0
	}
	n := 0
	s.Each(func(Field) error {
		n++
		return nil
	})
	return n
}

func (s *Struct) Empty() bool {
	if len(s.buf) == 0 {
		return true
	}
	body, _ := Contents(s.buf)
	return len(body) == 0
}

func (s *Struct) bytes() []byte {
	if len(s.buf) == 0 {
		return emptyStruct
	}
	return s.buf
}

// Each calls fn for each field in the struct.
// If fn returns Stop, Each stops and returns nil.
// If fn returns any other non-nil error, Each
// stops and returns that error. If Each
// encounters a malformed field while unpacking
// the struct, Each returns a non-nil error.
func (s Struct) Each(fn func(Field) error) error {
	if s.Empty() {
		return nil
	}
	if TypeOf(s.buf) != StructType {
		return fmt.Errorf("expected a struct; found ion type %s", TypeOf(s.buf))
	}
	body, _ := Contents(s.buf)
	if body == nil {
		return errInvalidIon
	}
	st := s.symtab()
	for len(body) > 0 {
		f, rest, err := ReadField(&st, body)
		if err != nil {
			return err
		}
		err = fn(f)
		if err == Stop {
			break
		} else if err != nil {
			return err
		}
		body = rest
	}
	return nil
}

// Fields appends fields to the given slice and returns
// the appended slice.
func (s Struct) Fields(fields []Field) []Field {
	fields = slices.Grow(fields, s.Len())
	s.Each(func(f Field) error {
		fields = append(fields, f)
		return nil
	})
	return fields
}

func (s Struct) Field(x Symbol) (Field, bool) {
	var field Field
	var ok bool
	s.Each(func(f Field) error {
		if f.Sym == x {
			field, ok = f, true
			return Stop
		}
		return nil
	})
	return field, ok
}

func (s Struct) FieldByName(name string) (Field, bool) {
	var field Field
	var ok bool
	s.Each(func(f Field) error {
		if f.Label == name {
			field, ok = f, true
			return Stop
		}
		return nil
	})
	return field, ok
}

// mergeFields merges the given fields with the
// fields of this struct into a new struct,
// overwriting any previous fields with
// conflicting names.
//
// This should only be used for testing in this
// package.
func (s Struct) mergeFields(st *Symtab, fields []Field) Struct {
	into := make([]Field, 0, s.Len()+len(fields))
	add := func(f Field) {
		for i := range into {
			if into[i].Label == f.Label {
				into[i] = f
				return
			}
		}
		into = append(into, f)
	}
	s.Each(func(f Field) error {
		add(f)
		return nil
	})
	for i := range fields {
		add(fields[i])
	}
	return NewStruct(st, into)
}

func (s *Struct) symtab() Symtab {
	return Symtab{interned: s.st}
}

var emptyList = []byte{0xb0}

// List is an ion list datum
type List composite

func NewList(st *Symtab, items []Datum) List {
	if len(items) == 0 {
		return List{}
	}
	var dst Buffer
	if st == nil {
		st = &Symtab{}
	}
	dst.WriteList(st, items)
	return List{
		st:  st.alias(),
		buf: dst.Bytes(),
	}
}

func (b *Buffer) WriteList(st *Symtab, items []Datum) {
	if len(items) == 0 {
		b.UnsafeAppend(emptyList)
		return
	}
	b.BeginList(-1)
	for i := range items {
		items[i].Encode(b, st)
	}
	b.EndList()
}

func (l List) Datum() Datum {
	if len(l.buf) == 0 {
		return Datum{buf: emptyList}
	}
	return Datum{st: l.st, buf: l.buf}
}

func (l List) Encode(dst *Buffer, st *Symtab) {
	// fast path: we can avoid resym
	if l.empty() || st.contains(l.st) {
		dst.UnsafeAppend(l.bytes())
		return
	}
	dst.BeginList(-1)
	l.Each(func(d Datum) error {
		d.Encode(dst, st)
		return nil
	})
	dst.EndList()
}

func (l List) Len() int {
	if l.empty() {
		return 0
	}
	n := 0
	l.Each(func(Datum) error {
		n++
		return nil
	})
	return n
}

func (l *List) empty() bool {
	if len(l.buf) == 0 {
		return true
	}
	body, _ := Contents(l.buf)
	return len(body) == 0
}

func (l *List) bytes() []byte {
	if l.empty() {
		return emptyList
	}
	return l.buf
}

// Each iterates over each datum in the
// list and calls fn on each datum in order.
// Each returns when it encounters an internal error
// (due to malformed ion) or when fn returns false.
func (l List) Each(fn func(Datum) error) error {
	if l.empty() {
		return nil
	}
	if TypeOf(l.buf) != ListType {
		return fmt.Errorf("expected a list; found ion type %s", TypeOf(l.buf))
	}
	body, _ := Contents(l.buf)
	if body == nil {
		return errInvalidIon
	}
	st := l.symtab()
	for len(body) > 0 {
		v, rest, err := ReadDatum(&st, body)
		if err != nil {
			return err
		}
		err = fn(v)
		if err == Stop {
			break
		} else if err != nil {
			return err
		}
		body = rest
	}
	return nil
}

func (l List) Items(items []Datum) []Datum {
	items = slices.Grow(items, l.Len())
	l.Each(func(d Datum) error {
		items = append(items, d)
		return nil
	})
	return items
}

func (l *List) symtab() Symtab {
	return Symtab{interned: l.st}
}

func (l List) Equal(l2 List) bool {
	if l.empty() {
		return l2.empty()
	}
	if bytes.Equal(l.buf, l2.buf) && stoverlap(l.st, l2.st) {
		return true
	}
	// TODO: optimize this
	i1 := l.Items(nil)
	i2 := l2.Items(nil)
	if len(i1) != len(i2) {
		return false
	}
	for i := range i1 {
		if !Equal(i1[i], i2[i]) {
			return false
		}
	}
	return true
}

var (
	False = Datum{buf: []byte{0x10}}
	True  = Datum{buf: []byte{0x11}}
)

func Bool(b bool) Datum {
	if b {
		return True
	}
	return False
}

func String(s string) Datum {
	var buf Buffer
	buf.WriteString(s)
	return Datum{buf: buf.Bytes()}
}

func Blob(b []byte) Datum {
	var buf Buffer
	buf.WriteBlob(b)
	return Datum{buf: buf.Bytes()}
}

// Interned returns a Datum that represents
// an interned string (a Symbol).
// Interned is always encoded as an ion symbol.
func Interned(st *Symtab, s string) Datum {
	if st == nil {
		st = new(Symtab)
	}
	var buf Buffer
	sym := st.Intern(s)
	buf.WriteSymbol(sym)
	return Datum{st: st.alias(), buf: buf.Bytes()}
}

// Annotation objects represent
// ion annotation datums.
func Annotation(st *Symtab, label string, val Datum) Datum {
	var dst Buffer
	if st == nil {
		st = &Symtab{}
	}
	dst.BeginAnnotation(1)
	dst.BeginField(st.Intern(label))
	if val.IsEmpty() {
		dst.WriteNull()
	} else {
		val.Encode(&dst, st)
	}
	dst.EndAnnotation()
	return Datum{
		st:  st.alias(),
		buf: dst.Bytes(),
	}
}

func Timestamp(t date.Time) Datum {
	var buf Buffer
	buf.WriteTime(t)
	return Datum{buf: buf.Bytes()}
}

func decodeNullDatum(_ *Symtab, b []byte) (Datum, []byte, error) {
	s := SizeOf(b)
	if s <= 0 || s > len(b) {
		return Empty, b, errInvalidIon
	}
	// note: we're skipping the whole datum here
	// so that a multi-byte nop pad is skipped appropriately
	return Null, b[s:], nil
}

func decodeBoolDatum(_ *Symtab, b []byte) (Datum, []byte, error) {
	_, rest, err := ReadBool(b)
	if err != nil {
		return Empty, rest, err
	}
	return rawDatum(nil, b), rest, nil
}

func decodeUintDatum(_ *Symtab, b []byte) (Datum, []byte, error) {
	if SizeOf(b) > 9 {
		return Empty, b, fmt.Errorf("int size %d out of range", SizeOf(b))
	}
	_, rest, err := ReadUint(b)
	if err != nil {
		return Empty, rest, err
	}
	return rawDatum(nil, b), rest, nil
}

func decodeIntDatum(_ *Symtab, b []byte) (Datum, []byte, error) {
	if SizeOf(b) > 9 {
		return Empty, b, fmt.Errorf("int size %d out of range", SizeOf(b))
	}
	_, rest, err := ReadInt(b)
	if err != nil {
		return Empty, rest, err
	}
	return rawDatum(nil, b), rest, nil
}

func decodeFloatDatum(_ *Symtab, b []byte) (Datum, []byte, error) {
	_, rest, err := ReadFloat64(b)
	if err != nil {
		return Empty, rest, err
	}
	return rawDatum(nil, b), rest, nil
}

func decodeDecimalDatum(_ *Symtab, b []byte) (Datum, []byte, error) {
	return Empty, nil, fmt.Errorf("ion: decimal decoding unimplemented")
}

func decodeTimestampDatum(_ *Symtab, b []byte) (Datum, []byte, error) {
	_, rest, err := ReadTime(b)
	if err != nil {
		return Empty, rest, err
	}
	return rawDatum(nil, b), rest, nil
}

func decodeSymbolDatum(st *Symtab, b []byte) (Datum, []byte, error) {
	sym, rest, err := ReadSymbol(b)
	if err != nil {
		return Empty, rest, err
	}
	if _, ok := st.Lookup(sym); !ok {
		return Empty, rest, fmt.Errorf("symbol %d not in symbol table", sym)
	}
	return rawDatum(st, b), rest, nil
}

func decodeBytesDatum(_ *Symtab, b []byte) (Datum, []byte, error) {
	buf, rest := Contents(b)
	if buf == nil {
		return Empty, b, errInvalidIon
	}
	return rawDatum(nil, b), rest, nil
}

func decodeListDatum(st *Symtab, b []byte) (Datum, []byte, error) {
	size := SizeOf(b)
	if size <= 0 || size > len(b) {
		return Empty, nil, fmt.Errorf("size %d exceeds buffer size %d", size, len(b))
	}
	body, rest := Contents(b)
	if body == nil {
		return Empty, nil, errInvalidIon
	}
	for len(body) > 0 {
		var err error
		body, err = validateDatum(st, body)
		if err != nil {
			return Empty, nil, err
		}
	}
	return rawDatum(st, b), rest, nil
}

func decodeStructDatum(st *Symtab, b []byte) (Datum, []byte, error) {
	size := SizeOf(b)
	if size <= 0 || size > len(b) {
		return Empty, nil, fmt.Errorf("size %d exceeds buffer size %d", size, len(b))
	}
	fields, rest := Contents(b)
	if fields == nil {
		return Empty, nil, errInvalidIon
	}
	for len(fields) > 0 {
		var sym Symbol
		var err error
		sym, fields, err = ReadLabel(fields)
		if err != nil {
			return Empty, nil, err
		}
		if len(fields) == 0 {
			return Empty, nil, io.ErrUnexpectedEOF
		}
		_, ok := st.Lookup(sym)
		if !ok {
			return Empty, nil, fmt.Errorf("symbol %d not in symbol table", sym)
		}
		fields, err = validateDatum(st, fields)
		if err != nil {
			return Empty, nil, err
		}
	}
	return rawDatum(st, b), rest, nil
}

func decodeReserved(_ *Symtab, b []byte) (Datum, []byte, error) {
	return Empty, b, fmt.Errorf("decoding error: tag %x is reserved", b[0])
}

func decodeAnnotationDatum(st *Symtab, b []byte) (Datum, []byte, error) {
	sym, body, rest, err := ReadAnnotation(b)
	if err != nil {
		return Empty, rest, err
	}
	if _, ok := st.Lookup(sym); !ok {
		return Empty, rest, fmt.Errorf("symbol %d not in symbol table", sym)
	}
	_, err = validateDatum(st, body)
	if err != nil {
		return Empty, rest, err
	}
	return Datum{
		st:  st.alias(),
		buf: b[:SizeOf(b)],
	}, rest, nil
}

var _datumTable = [...](func(*Symtab, []byte) (Datum, []byte, error)){
	NullType:       decodeNullDatum,
	BoolType:       decodeBoolDatum,
	UintType:       decodeUintDatum,
	IntType:        decodeIntDatum,
	FloatType:      decodeFloatDatum,
	DecimalType:    decodeDecimalDatum,
	TimestampType:  decodeTimestampDatum,
	SymbolType:     decodeSymbolDatum,
	StringType:     decodeBytesDatum,
	ClobType:       decodeBytesDatum, // fixme: treat clob differently than blob?
	BlobType:       decodeBytesDatum,
	ListType:       decodeListDatum,
	SexpType:       decodeListDatum, // fixme: treat sexp differently than list?
	StructType:     decodeStructDatum,
	AnnotationType: decodeAnnotationDatum,
	ReservedType:   decodeReserved,
}

var datumTable [16](func(*Symtab, []byte) (Datum, []byte, error))

func init() {
	copy(datumTable[:], _datumTable[:])
}

// ReadDatum reads the next datum from buf
// and returns it. ReadDatum does not return
// symbol tables directly; instead it unmarshals
// them into st and continues reading. It may
// return a nil datum if buf points to a symbol
// table followed by zero bytes of actual ion data.
//
// Any Symbol datums in buf are translated into
// Interned datums rather than Symbol datums,
// as this makes the returned Datum safe to
// re-encode with a new symbol table.
//
// The returned datum will share memory with buf and so
// the caller must guarantee that the contents of buf
// will not be modified until it is no longer needed.
func ReadDatum(st *Symtab, buf []byte) (Datum, []byte, error) {
	var err error
	if IsBVM(buf) || TypeOf(buf) == AnnotationType {
		buf, err = st.Unmarshal(buf)
		if err != nil {
			return Empty, nil, err
		}
		if len(buf) == 0 {
			return Empty, buf, nil
		}
	}
	return datumTable[TypeOf(buf)](st, buf)
}

// validateDatum validates that the next datum in buf
// does not exceed the bounds of buf without actually
// interpretting it. This also handles symbol tables
// the same way that ReadDatum does.
func validateDatum(st *Symtab, buf []byte) (next []byte, err error) {
	if IsBVM(buf) || TypeOf(buf) == AnnotationType {
		buf, err = st.Unmarshal(buf)
		if err != nil {
			return nil, err
		}
		if len(buf) == 0 {
			return nil, nil
		}
	}
	size := SizeOf(buf)
	if size <= 0 || size > len(buf) {
		return nil, fmt.Errorf("size %d exceeds buffer size %d", size, len(buf))
	}
	return buf[size:], nil
}

// Equal returns whether a and b are
// semantically equivalent.
func Equal(a, b Datum) bool {
	return a.Equal(b)
}

func stoverlap(st1, st2 []string) bool {
	return stcontains(st1, st2) || stcontains(st2, st1)
}
