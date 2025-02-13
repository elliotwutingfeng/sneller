// Copyright (C) 2023 Sneller, Inc.
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

package elastic_proxy

import (
	"testing"
)

func TestDataShapeToElasticMapping(t *testing.T) {
	// given
	datashape := testDatashape()

	// when
	m := DataShapeToElasticMapping(datashape)

	var p *Properties

	// then
	assertType := func(field, want string) {
		t.Helper()
		mv, ok := (*p)[field]
		if !ok {
			t.Errorf("field %q not found", field)
			return
		}

		got := mv.Type
		if got != want {
			t.Logf("got:  %s", got)
			t.Logf("want: %s", want)
			t.Errorf("wrong type of field %q", field)
		}
	}

	assertNoField := func(field string) {
		t.Helper()
		_, ok := m.Properties[field]
		if ok {
			t.Errorf("field %q should not be present", field)
		}
	}

	assertCount := func(want int) {
		t.Helper()
		got := len(*p)
		if got != want {
			t.Logf("got:  %d", got)
			t.Logf("want: %d", want)
			t.Errorf("wrong number of fields")
		}
	}

	// top-level
	p = &m.Properties
	assertCount(6)
	assertType("enabled", "boolean")
	assertType("test", "keyword")
	assertType("user", "object")
	assertType("tags", "list")
	assertType("bag", "keyword")
	assertType("avatar", "object")

	assertNoField("tags.$items")

	// `user` struct
	{
		tmp := m.Properties["user"].Properties
		p = &tmp
	}
	assertCount(5)
	assertType("display_name", "keyword")
	assertType("karma", "double")
	assertType("lastlogin", "date")
	assertType("username", "keyword")
	assertType("statistics", "object")

	// `user.statistics` struct
	{
		tmp := m.Properties["user"].Properties["statistics"].Properties
		p = &tmp
	}
	assertCount(2)
	assertType("posts", "long")
	assertType("postsperday", "double")

	// `avatar` struct
	{
		tmp := m.Properties["avatar"].Properties
		p = &tmp
	}
	assertCount(3)
	assertType("size", "long")
	assertType("alt-text", "keyword")
	assertType("url", "keyword")
}

func testDatashape() map[string]any {
	return map[string]any{
		"enabled": map[string]any{ // nulls + bool -> bool
			nullField: 10,
			boolField: 5,
		},
		"tags": map[string]any{
			listField: 1,
		},
		"tags.$items": map[string]any{ // content of list values should be skipped
			stringField: 20,
		},
		"user.display_name": map[string]any{ // untyped null -> default
			nullField: 10,
		},
		"user.karma": map[string]any{ // int & float -> float
			intField:   5,
			floatField: 1,
		},
		"user.statistics.posts": map[string]any{ // int
			intField: 6,
		},
		"user.statistics.postsperday": map[string]any{ // float
			floatField: 1,
		},
		"user.lastlogin": map[string]any{ // timestamp
			timestampField: 1,
		},
		"user.username": map[string]any{ // string
			stringField: 1,
		},
		"user.statistics": map[string]any{
			structField: 2,
		},
		"avatar": map[string]any{
			structField: 1,
		},
		"avatar.url": map[string]any{
			stringField: 1,
		},
		"avatar.alt-text": map[string]any{
			stringField: 61,
		},
		"avatar.size": map[string]any{
			intField: 4,
		},
		"test": map[string]any{ // unsupported Ion type -> default
			decimalField: 2,
		},
		"user": map[string]any{
			structField: 1,
		},
		"bag": map[string]any{ // mixed types -> default
			timestampField: 2,
			intField:       3,
			boolField:      4,
			sexpField:      true, // not an int
		},
		"wrong-type": 42,
	}
}
