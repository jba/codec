// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codec

import (
	"bytes"
	"flag"
	"fmt"
	"go/token"
	"log"
	"math"
	"net"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	foo "github.com/jba/codec/internal/testpkg"
)

// After a change that affects generated code, run "go generate".
// Also run it if a test fails with 'unregistered type "*codec.node"'.

//go:generate rm -f types.gen_test.go
//go:generate go test -generate types.gen_test.go

var generateTestCodeFilename = flag.String("generate", "", "generate code for tests to filename")

// This struct exists just so we can pass it to GenerateFile and get all the
// types we need for these tests.
type generatedTestTypes struct {
	Node        *node
	Slice       []int
	Array       [1]int
	ByteSlice   []byte // should not generate a codec
	ByteArray   [2]byte
	Map         map[string]bool
	Struct      structType
	IP          net.IP
	StructSlice []structType
	StructArray [1]structType
	StructMap   map[[1]int]structType
	DefSlice    definedSlice
	DefArray    definedArray
	DefMap      definedMap
	Pos         token.Pos
	T           foo.T
	PtrSlice    *[]int
	PtrArray    *[1]int
	PtrMap      *map[int]int
	PtrTime     *time.Time
	SlicePtrInt []*int
}

// for testing sharing and cycles
type node struct {
	Value int
	Next  *node
}

type structType struct {
	N          node
	B          byte
	unexported int
	embed
}

type embed struct {
	E int
}

func TestMain(m *testing.M) {
	flag.Parse()
	if *generateTestCodeFilename != "" {
		if err := GenerateFile(*generateTestCodeFilename, "github.com/jba/codec", nil, generatedTestTypes{}); err != nil {
			log.Fatal(err)
		}
		fmt.Println("generated file, now run tests again")
	} else {
		os.Exit(m.Run())
	}
}

func TestEncodeDecode(t *testing.T) {
	for _, opts := range []EncodeOptions{
		{TrackPointers: false},
		{TrackPointers: true},
		{Buffer: make([]byte, 3)},
	} {
		t.Run(fmt.Sprintf("%+v", opts), func(t *testing.T) {
			testEncodeDecode(t, opts)
		})
	}
}

func testEncodeDecode(t *testing.T, opts EncodeOptions) {
	tm := time.Date(2020, time.March, 20, 0, 0, 0, 0, time.UTC)
	want := []interface{}{
		nil, "Luke Luck likes lakes", true,
		1, -5, 255, 65000, uint64(65000), uint64(1 << 63),
		0.0, 98.6, 100, 1.23e63, math.NaN(), math.Inf(1), math.Inf(-1),
		tm,
		&tm,
		net.IPv4(10, 128, 2, 18),
		&node{1, &node{2, &node{3, nil}}},
		(*node)(nil),
		[]int{},
		[]int{1},
		[]int{1, 2},
		[]int{1, 2, -3},
		[]int{1, 2, -3, -4},
		[]int{1, 2, -3, -4, 5, 6},
		[]int(nil),
		structType{B: 5, embed: embed{E: 6}},
		[]structType{{B: 3}},
		[1]structType{{B: 4}},
		map[string]bool{"a": true, "b": false},
		map[string]bool(nil),
		map[[1]int]structType{
			[1]int{7}: structType{B: 99, unexported: 8},
		},
		structType{B: 129, N: node{1, nil}, unexported: 23},
		definedSlice{1, 2, 3},
		definedArray{-7},
		definedMap{"true": true},
		[]byte{1, 2, 3},
		[2]byte{4, 5},
		func() *int { x := 6; return &x }(),
		&[]int{7, 8},
		&[1]int{9},
		&map[int]int{10: 11},
	}
	var buf bytes.Buffer
	e := NewEncoder(&buf, &opts)
	for _, w := range want {
		if err := e.Encode(w); err != nil {
			t.Fatalf("%#v: %v", w, err)
		}
	}
	r := bytes.NewReader(buf.Bytes())
	d := NewDecoder(r, nil)
	for _, w := range want {
		var g interface{}
		if err := d.Decode(&g); err != nil {
			t.Fatalf("%#v: %v", w, err)
		}
		if !cmp.Equal(g, w, cmpopts.EquateNaNs(), cmp.AllowUnexported(structType{})) {
			t.Errorf("got %v, want %v", g, w)
		}
	}
}

func TestSharing(t *testing.T) {
	n := &node{Value: 99, Next: &node{Value: 111}}
	n.Next.Next = n // create a cycle
	var g *node
	roundTripSharing(t, n, &g)
	if g.Value != 99 || g.Next.Value != 111 {
		t.Fatal("bad values")
	}
	if g.Next.Next != g {
		t.Error("did not preserve cycle")
	}

	two := 2
	seven := 7
	s := []*int{&two, &seven, &seven, &two, &two}
	var h []*int
	roundTripSharing(t, s, &h)
	if h[0] != h[3] || h[0] != h[4] || h[1] != h[2] {
		t.Error("did not preserve sharing")
	}
}

func roundTripSharing(t *testing.T, in, out interface{}) {
	t.Helper()
	var buf bytes.Buffer
	e := NewEncoder(&buf, &EncodeOptions{TrackPointers: true})
	if err := e.Encode(in); err != nil {
		t.Fatal(err)
	}
	d := NewDecoder(bytes.NewReader(buf.Bytes()), nil)
	if err := d.Decode(out); err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(in, reflect.ValueOf(out).Elem().Interface()) {
		t.Error("unequal")
	}
}

func TestEncodeErrors(t *testing.T) {
	// The only encoding error is an unregistered type.
	e := NewEncoder(&bytes.Buffer{}, nil)
	type MyInt int
	checkMessage(t, e.Encode(MyInt(0)), "unregistered")
}

func checkMessage(t *testing.T, err error, target string) {
	t.Helper()
	if err == nil {
		t.Error("want error, got nil")
	}
	if !strings.Contains(err.Error(), target) {
		t.Errorf("error %q does not contain %q", err, target)
	}
}
