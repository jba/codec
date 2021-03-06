// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codec

import (
	"bytes"
	"flag"
	"fmt"
	"go/token"
	"io"
	"io/ioutil"
	"net"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	othercmp "github.com/jba/codec/internal/cmp"
	foo "github.com/jba/codec/internal/testpkg"
)

var update = flag.Bool("update", false, "update goldens instead of checking against them")

type (
	definedSlice []int
	definedArray [1]int
	definedMap   map[string]bool
)

type genStruct struct {
	S string
	B bool

	I   int
	I8  int8
	I16 int16
	I32 int32
	I64 int64

	F32 float32
	F64 float64

	U8  uint8
	U16 uint16
	U32 uint32
	U64 uint64

	C64  complex64
	C128 complex128

	BS         []byte
	T          foo.T
	unexported int
	Omit       int `test:"-"`
}

// A small struct that doesn't get in the way when
// it's part of something larger.
type smallStruct struct{ X int }

// A type that is not a slice, array, map, pointer or struct, but does implement a Marshaler interface.
type marsh int

func (m marsh) MarshalText() ([]byte, error) { return []byte(strconv.Itoa(int(m))), nil }

func (m *marsh) UnmarshalText(d []byte) error {
	i, err := strconv.Atoi(string(d))
	if err != nil {
		return err
	}
	*m = marsh(i)
	return nil
}

func TestGenerate(t *testing.T) {
	testGenerate(t, "slice", [][]int(nil))
	testGenerate(t, "islice", []interface{}(nil))
	testGenerate(t, "map", map[string]bool(nil))
	testGenerate(t, "struct", genStruct{unexported: 0}) // suppress staticcheck warning
	testGenerate(t, "binmarsh", time.Time{})
	testGenerate(t, "textmarsh", net.IP{})
	testGenerate(t, "structslice", []smallStruct(nil))
	testGenerate(t, "structmap", map[[1]int]smallStruct{})
	testGenerate(t, "defslice", definedSlice{})
	testGenerate(t, "defarray", definedArray{})
	testGenerate(t, "defmap", definedMap{})
	testGenerate(t, "slicemarsh", []marsh{})
}

func testGenerate(t *testing.T, name string, x interface{}) {
	t.Run(name, func(t *testing.T) {
		var buf bytes.Buffer
		if err := generate(&buf, "github.com/jba/codec", "test", x); err != nil {
			t.Fatal(err)
		}
		got := buf.String()
		if *update {
			writeGolden(t, name, got)
		} else {
			want := readGolden(t, name)
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("%s: mismatch (-want, +got):\n%s", name, diff)
				filename, err := writeTempFile(fmt.Sprintf("generate-%s-*.go", name), []byte(got))
				if err != nil {
					t.Fatal(err)
				}
				t.Logf("wrote got to %s", filename)
			}
		}
	})
}

func writeGolden(t *testing.T, name string, data string) {
	filename := filepath.Join("testdata", name+".go")
	if err := ioutil.WriteFile(filename, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
	t.Logf("wrote %s", filename)
}

func readGolden(t *testing.T, name string) string {
	data, err := ioutil.ReadFile(filepath.Join("testdata", name+".go"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestGenerateErrors(t *testing.T) {
	var buf bytes.Buffer
	err := generate(&buf, "github.com/jba/codec", "test", struct{ X int }{})
	const want = "unnamed struct"
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("got %v, want error containing %q", err, want)
	}
}

func TestStructFields(t *testing.T) {
	type ef struct {
		A int
		B bool
		I int `codec:"-"` // this field will be ignored
		C string
		D int `codec:"N"`
	}

	var (
		intType    = reflect.TypeOf(0)
		stringType = reflect.TypeOf("")
		boolType   = reflect.TypeOf(false)
	)

	g := &generator{pkgPath: "p", fieldTagKey: "codec"}
	got := g.structFields(reflect.TypeOf(ef{}))
	want := []field{
		{"A", intType, "0"},
		{"B", boolType, "false"},
		{"C", stringType, `""`},
		{"N", intType, "0"},
	}
	diff := cmp.Diff(want, got,
		cmp.Comparer(func(t1, t2 reflect.Type) bool { return t1 == t2 }))
	if diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

type parseTagStruct struct {
	NoTag    int
	Name     int `test:"tag"`
	Omit     int `test:"-"`
	Opts     int `test:",a,b,c"`
	NameOpts int `test:"name,d"`
}

func TestParseTag(t *testing.T) {
	typ := reflect.TypeOf(parseTagStruct{})
	for _, test := range []struct {
		field    string
		wantName string
		wantOmit bool
		wantOpts []string
	}{
		{field: "NoTag", wantName: ""},
		{field: "Name", wantName: "tag"},
		{field: "Omit", wantOmit: true},
		{field: "Opts", wantOpts: []string{"a", "b", "c"}},
		{field: "NameOpts", wantName: "name", wantOpts: []string{"d"}},
	} {

		f, ok := typ.FieldByName(test.field)
		if !ok {
			t.Fatalf("no field %q", test.field)
		}
		gotName, gotOmit := parseTag("test", f.Tag)
		if gotName != test.wantName {
			t.Errorf("name: got %q, want %q", gotName, test.wantName)
		}
		if gotOmit != test.wantOmit {
			t.Errorf("omit: got %t, want %t", gotOmit, test.wantOmit)
		}
	}

}

var (
	cmpType      = reflect.TypeOf(cmp.Indirect{})
	othercmpType = reflect.TypeOf(othercmp.Indirect{})
	fooType      = reflect.TypeOf(foo.T(nil))
)

func TestGoName(t *testing.T) {
	var r io.Reader
	g := &generator{pkgPath: "github.com/jba/codec"}
	g.buildImportMap([]reflect.Type{cmpType, othercmpType})
	for _, test := range []struct {
		v    interface{}
		want string
	}{
		{0, "int"},
		{uint(0), "uint"},
		{token.Pos(0), "token.Pos"},
		{Encoder{}, "Encoder"},
		{[][]Encoder{}, "[][]Encoder"},
		{bytes.Buffer{}, "bytes.Buffer"},
		{&r, "*io.Reader"},
		{[]int(nil), "[]int"},
		{[1]int{0}, "[1]int"},
		{map[*Decoder][]io.Writer{}, "map[*Decoder][]io.Writer"},
		{cmp.Indirect{}, "cmp.Indirect"},
		{othercmp.Indirect{}, "cmp1.Indirect"},
	} {
		got := g.goName(reflect.TypeOf(test.v))
		if got != test.want {
			t.Errorf("%T: got %q, want %q", test.v, got, test.want)
		}
	}
}

func TestPackageName(t *testing.T) {
	for _, test := range []struct {
		typ  reflect.Type
		want string
	}{
		{reflect.TypeOf(0), ""},               // builtin type
		{reflect.TypeOf(new(cmp.Option)), ""}, // unnamed type
		{cmpType, "cmp"},
		{othercmpType, "cmp"},
		{fooType, "foo"},
	} {
		got := packageName(test.typ)
		if got != test.want {
			t.Errorf("%s: got %q, want %q", test.typ, got, test.want)
		}
	}
}

func TestBuildImportMap(t *testing.T) {
	g := &generator{pkgPath: "github.com/jba/codec"}
	types := []reflect.Type{reflect.TypeOf(0), cmpType, othercmpType, fooType, reflect.TypeOf(Decoder{})}
	g.buildImportMap(types)
	want := map[string]string{
		"reflect":                               "",
		"github.com/jba/codec/codecapi":         "",
		"github.com/google/go-cmp/cmp":          "",
		"github.com/jba/codec/internal/cmp":     "cmp1",
		"github.com/jba/codec/internal/testpkg": "foo",
	}
	if diff := cmp.Diff(want, g.importMap); diff != "" {
		t.Errorf("diff (-want, +got):\n%s", diff)
	}
	want = map[string]string{
		"reflect":                               "reflect",
		"github.com/jba/codec/codecapi":         "codecapi",
		"github.com/jba/codec":                  "",
		"github.com/google/go-cmp/cmp":          "cmp",
		"github.com/jba/codec/internal/cmp":     "cmp1",
		"github.com/jba/codec/internal/testpkg": "foo",
	}
	if diff := cmp.Diff(want, g.pkgPathMap); diff != "" {
		t.Errorf("diff (-want, +got):\n%s", diff)
	}

}

func TestReferencedTypeList(t *testing.T) {
	g := &generator{
		pkgPath:     "github.com/jba/codec",
		fieldTagKey: "codec",
	}
	cmpType := reflect.TypeOf(new(cmp.Option)).Elem()
	got := g.referencedTypeList([]interface{}{0, []*cmp.Option{}, genStruct{}, token.Pos(0), net.IP{}})
	wantvals := []interface{}{
		new(cmp.Option), []*cmp.Option{}, cmpType, genStruct{},
		foo.T{}, token.Pos(0), net.IP{}}
	var want []reflect.Type
	for _, v := range wantvals {
		if t, ok := v.(reflect.Type); ok {
			want = append(want, t)
		} else {
			want = append(want, reflect.TypeOf(v))
		}
	}
	if !cmp.Equal(got, want, cmp.Comparer(func(t1, t2 reflect.Type) bool { return t1 == t2 })) {
		t.Errorf("\ngot  %v\nwant %v", got, want)
	}
}
