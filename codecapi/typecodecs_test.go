// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codecapi

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	othercmp "github.com/jba/codec/internal/cmp"
)

func TestTypeStringUnique(t *testing.T) {
	check := func(in interface{}, want string, pkgPaths map[string]string) {
		t.Helper()
		typ := reflect.TypeOf(in)
		got := TypeString(typ, pkgPaths)
		if got != want {
			t.Errorf("%s: got %s, want %s", typ, got, want)
		}
	}

	pkgPaths := map[string]string{
		"github.com/jba/codec/codecapi":     "",
		"github.com/jba/codec/internal/cmp": "othercmp",
	}

	// Types that behave the same regardless of package path mapping.
	for _, test := range []struct {
		in   interface{}
		want string
	}{
		{0, "int"},
		{1.1, "float64"},
		{"", "string"},
		{[]int(nil), "[]int"},
		{[3]bool{}, "[3]bool"},
		{new(int), "*int"},
		{map[string]bool{}, "map[string]bool"},
		{map[[1]int]*struct{ C complex64 }{}, "map[[1]int]*struct { C complex64 }"},
		{new(interface{}), "*interface{}"},
	} {
		check(test.in, test.want, nil)
		check(test.in, test.want, pkgPaths)
	}

	// Types that behave differently.
	for _, test := range []struct {
		in                     interface{}
		wantUnique, wantMapped string
	}{
		{
			[]cmp.Option{},
			"[]github.com/google/go-cmp/cmp.Option",
			"[]cmp.Option",
		},
		{
			map[cmp.Option]*othercmp.Option{},
			"map[github.com/google/go-cmp/cmp.Option]*github.com/jba/codec/internal/cmp.Option",
			"map[cmp.Option]*othercmp.Option",
		},
		{
			new(cmp.Option),
			"*github.com/google/go-cmp/cmp.Option",
			"*cmp.Option",
		},
		{
			new([2]cmp.Option),
			"*[2]github.com/google/go-cmp/cmp.Option",
			"*[2]cmp.Option",
		},
		{
			struct{ X cmp.Option }{},
			"struct { X github.com/google/go-cmp/cmp.Option }",
			"struct { X cmp.Option }",
		},
		{
			Decoder{},
			"github.com/jba/codec/codecapi.Decoder",
			"Decoder",
		},
	} {
		check(test.in, test.wantUnique, nil)
		check(test.in, test.wantMapped, pkgPaths)
	}
}
