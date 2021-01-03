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
	Node     *node
	Slice    []int
	Array    [1]int
	Map      map[string]bool
	Struct   structType
	Time     time.Time
	IP       net.IP
	DefSlice definedSlice
	DefArray definedArray
	DefMap   definedMap
	Pos      token.Pos
	T        foo.T
	//Skip     Skip // add this to re-generate code for TestSkip; see skip_code_test.go
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
	want := []interface{}{
		nil, "Luke Luck likes lakes", true,
		1, -5, 255, 65000, uint64(65000), uint64(1 << 63),
		0.0, 98.6, 100, 1.23e63, math.NaN(), math.Inf(1), math.Inf(-1),
		time.Date(2020, time.January, 26, 0, 0, 0, 0, time.UTC),
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
		map[string]bool{"a": true, "b": false},
		map[string]bool(nil),
		structType{B: 129, N: node{1, nil}, unexported: 23},
		definedSlice{1, 2, 3},
		definedArray{-7},
		definedMap{"true": true},
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
		g, err := d.Decode()
		if err != nil {
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

	var buf bytes.Buffer
	err := NewEncoder(&buf, &EncodeOptions{TrackPointers: true}).Encode(n)
	if err != nil {
		t.Fatal(err)
	}
	d := NewDecoder(bytes.NewReader(buf.Bytes()), nil)
	g, err := d.Decode()
	if err != nil {
		t.Fatal(err)
	}
	_ = g
	// got := g.(*node)
	// if !cmp.Equal(got, n) {
	// 	t.Error("did not preserve cycle")
	// }
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

type Skip struct {
	U  uint64
	S1 string
	S2 string
	L  []*Skip
}

func TestSkipUnknownFields(t *testing.T) {
	s := &Skip{
		U:  1,     // < endCode
		S1: "ab",  // bytes2
		S2: "abc", // bytes3
	}
	v := Skip{
		U:  255,                // bytes1
		S1: "abcd",             //bytes4
		S2: "elephantine",      // nBytes
		L:  []*Skip{s, s, nil}, // nValues, ptrCode, startCode, refCode, nilCode
	}

	var buf bytes.Buffer
	e := NewEncoder(&buf, &EncodeOptions{TrackPointers: true})
	if err := e.Encode(v); err != nil {
		t.Fatal(err)
	}

	d := NewDecoder(bytes.NewReader(buf.Bytes()), nil)

	got, err := d.Decode()
	if err != nil {
		t.Fatal(err)
	}
	var want Skip
	if !cmp.Equal(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestFailOnUnknownField(t *testing.T) {
	v := &Skip{U: 1}
	var buf bytes.Buffer
	e := NewEncoder(&buf, nil)
	if err := e.Encode(v); err != nil {
		t.Fatal(err)
	}
	d := NewDecoder(bytes.NewReader(buf.Bytes()), &DecodeOptions{FailOnUnknownField: true})
	_, err := d.Decode()
	if err == nil {
		t.Fatal("got nil, want error")
	}
}
