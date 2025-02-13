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

package plan

import (
	"fmt"
	"io"
	"slices"
	"sync"

	"github.com/SnellerInc/sneller/expr"
	"github.com/SnellerInc/sneller/ion"
	"github.com/SnellerInc/sneller/plan/pir"
)

type replacement struct {
	lock sync.Mutex

	rows []ion.Struct
}

func mustConst(d ion.Datum) expr.Constant {
	c, ok := expr.AsConstant(d)
	if !ok {
		panic("cannot convert value to constant")
	}
	return c
}

func first(s *ion.Struct) (ion.Field, bool) {
	var field ion.Field
	var ok bool
	s.Each(func(f ion.Field) error {
		field, ok = f, true
		return ion.Stop
	})
	return field, ok
}

func (r *replacement) toScalar() expr.Constant {
	if len(r.rows) == 0 {
		return expr.Null{}
	}
	s := &r.rows[0]
	f, ok := first(s)
	if !ok {
		return expr.Null{}
	}
	return mustConst(f.Datum)
}

func (r *replacement) toScalarList() ion.Bag {
	var ret ion.Bag
	for i := range r.rows {
		f, ok := first(&r.rows[i])
		if !ok {
			continue
		}
		ret.AddDatum(f.Datum)
	}
	return ret
}

func (r *replacement) toList() expr.Constant {
	lst := new(expr.List)
	for i := range r.rows {
		lst.Values = append(lst.Values, mustConst(r.rows[i].Datum()))
	}
	return lst
}

func (r *replacement) toStruct() expr.Constant {
	if len(r.rows) == 0 {
		return &expr.Struct{}
	}
	return mustConst(r.rows[0].Datum())
}

func (r *replacement) toHashLookup(kind, label string, x, elseval expr.Node) expr.Node {
	if len(r.rows) == 0 {
		return expr.Missing{}
	}
	var conv rowConverter
	switch kind {
	case "scalar":
		conv = &scalarConverter{label: label}
	case "struct":
		conv = &structConverter{label: label}
	case "list":
		conv = &listConverter{label: label}
	case "joinlist":
		conv = &joinListConverter{label: label}
	default:
		return expr.Null{}
	}
	for i := range r.rows {
		conv.add(&r.rows[i])
	}
	return conv.result(x, elseval)
}

type rowConverter interface {
	add(row *ion.Struct)
	result(key, elseval expr.Node) *expr.Lookup
}

type scalarConverter struct {
	label        string
	keys, values ion.Bag
}

func (c *scalarConverter) result(key, elseval expr.Node) *expr.Lookup {
	return &expr.Lookup{
		Expr:   key,
		Else:   elseval,
		Keys:   c.keys,
		Values: c.values,
	}
}

func (c *scalarConverter) add(row *ion.Struct) {
	if row.Len() != 2 {
		return
	}
	f := row.Fields(make([]ion.Field, 0, 2))
	if f[0].Label != c.label {
		f[0], f[1] = f[1], f[0]
		if f[0].Label != c.label {
			return
		}
	}
	c.keys.AddDatum(f[0].Datum)
	c.values.AddDatum(f[1].Datum)
}

type structConverter struct {
	label        string
	keys, values ion.Bag
}

func (c *structConverter) result(key, elseval expr.Node) *expr.Lookup {
	return &expr.Lookup{
		Expr:   key,
		Else:   elseval,
		Keys:   c.keys,
		Values: c.values,
	}
}

func (c *structConverter) add(row *ion.Struct) {
	if row.Len() == 0 {
		return
	}
	var key ion.Datum
	fields := make([]ion.Field, 0, row.Len()-1)
	row.Each(func(f ion.Field) error {
		if key.IsEmpty() && f.Label == c.label {
			key = f.Datum
			return nil
		}
		fields = append(fields, f)
		return nil
	})
	if key.IsEmpty() {
		return
	}
	c.keys.AddDatum(key)
	c.values.AddDatum(ion.NewStruct(nil, fields).Datum())
}

type listConverter struct {
	label string
	m     map[expr.Constant]*expr.List
}

func (c *listConverter) result(key, elseval expr.Node) *expr.Lookup {
	l := &expr.Lookup{Expr: key, Else: elseval}
	for k, v := range c.m {
		l.Keys.AddDatum(k.Datum())
		l.Values.AddDatum(v.Datum())
	}
	return l
}

func (c *listConverter) add(row *ion.Struct) {
	if row.Len() == 0 {
		return
	}
	var key expr.Constant
	fields := make([]expr.Field, 0, row.Len()-1)
	row.Each(func(f ion.Field) error {
		val := mustConst(f.Datum)
		if f.Label == c.label {
			key = val
			return nil
		}
		fields = append(fields, expr.Field{
			Label: f.Label,
			Value: val,
		})
		return nil
	})
	if key == nil {
		return
	}
	lst := c.m[key]
	if lst == nil {
		lst = &expr.List{}
		if c.m == nil {
			c.m = make(map[expr.Constant]*expr.List)
		}
		c.m[key] = lst
	}
	lst.Values = append(lst.Values, &expr.Struct{Fields: fields})
}

type joinListConverter struct {
	label string
	m     map[string][]ion.Datum
	st    ion.Symtab
	tmp   ion.Buffer
}

func (j *joinListConverter) stringify(d ion.Datum) []byte {
	j.tmp.Reset()
	d.Encode(&j.tmp, &j.st)
	return j.tmp.Bytes()
}

func (j *joinListConverter) add(row *ion.Struct) {
	var key, val ion.Datum
	row.Each(func(f ion.Field) error {
		if f.Label == j.label {
			key = f.Datum
		} else {
			val = f.Datum
		}
		return nil
	})
	if key.IsEmpty() || val.IsEmpty() {
		return
	}
	if j.m == nil {
		j.m = make(map[string][]ion.Datum)
	}
	str := j.stringify(key)
	j.m[string(str)] = append(j.m[string(str)], val)
}

func (j *joinListConverter) result(key, elseval expr.Node) *expr.Lookup {
	l := &expr.Lookup{Expr: key, Else: elseval}
	for k, v := range j.m {
		dat, _, err := ion.ReadDatum(&j.st, []byte(k))
		if err != nil {
			panic(err)
		}
		l.Keys.AddDatum(dat)
		lst := ion.NewList(nil, v).Datum()
		l.Values.AddDatum(lst)
	}
	return l
}

type subreplacement struct {
	parent *replacement
	curst  ion.Symtab
	tmp    []ion.Struct
}

func (s *subreplacement) Write(buf []byte) (int, error) {
	buf = slices.Clone(buf)
	orig := len(buf)
	s.tmp = s.tmp[:0]
	var err error
	var d ion.Datum
	for len(buf) > 0 {
		d, buf, err = ion.ReadDatum(&s.curst, buf)
		if err != nil {
			return orig - len(buf), err
		}
		if d.IsEmpty() || d.IsNull() {
			continue // symbol table or nop pad
		}
		st, _ := d.Struct()
		s.tmp = append(s.tmp, st)
	}
	s.parent.lock.Lock()
	defer s.parent.lock.Unlock()
	s.parent.rows = append(s.parent.rows, s.tmp...)
	s.tmp = s.tmp[:0]
	if len(s.parent.rows) > pir.LargeSize {
		return orig, fmt.Errorf("%d items in subreplacement exceeds limit", len(s.parent.rows))
	}
	return orig, nil
}

func (s *subreplacement) Close() error {
	return nil
}

func (r *replacement) Open() (io.WriteCloser, error) {
	return &subreplacement{
		parent: r,
	}, nil
}

func (r *replacement) Close() error {
	return nil
}

// replacer substitutes replacement tokens
// like IN_REPLACEMENT(expr, id)
// and SCALAR_REPLACMENT(id)
// with the appropriate constant from
// the replacement list
type replacer struct {
	inputs []replacement
	simpl  expr.Rewriter
}

// we perform simplification after substitution
// so that any constprop opportunities that appear
// after replacement get taken care of
func (r *replacer) simplify(e expr.Node) expr.Node {
	return r.simpl.Rewrite(e)
}

func (r *replacer) Walk(e expr.Node) expr.Rewriter {
	return r
}

func (r *replacer) Rewrite(e expr.Node) expr.Node {
	b, ok := e.(*expr.Builtin)
	if !ok {
		return r.simplify(e)
	}
	switch b.Func {
	default:
		return r.simplify(e)
	case expr.ListReplacement:
		id := int(b.Args[0].(expr.Integer))
		return r.inputs[id].toList()
	case expr.InReplacement:
		id := int(b.Args[1].(expr.Integer))
		return &expr.Member{
			Arg: b.Args[0],
			Set: r.inputs[id].toScalarList(),
		}
	case expr.HashReplacement:
		id := int(b.Args[0].(expr.Integer))
		kind := string(b.Args[1].(expr.String))
		label := string(b.Args[2].(expr.String))
		var elseval expr.Node
		if len(b.Args) == 5 {
			elseval = b.Args[4]
		}
		return r.inputs[id].toHashLookup(kind, label, b.Args[3], elseval)
	case expr.StructReplacement:
		id := int(b.Args[0].(expr.Integer))
		return r.inputs[id].toStruct()
	case expr.ScalarReplacement:
		id := int(b.Args[0].(expr.Integer))
		return r.inputs[id].toScalar()
	}
}
