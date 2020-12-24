// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codec

import (
	"bytes"
	"flag"
	"fmt"
	"go/token"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	api "github.com/jba/codec/codecapi"
)

// After a change that affects generated code, run "go generate".
// Also run it if a test fails with 'unregistered type "*codec.node"'.

//go:generate ./generate_for_tests.sh

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
	DefSlice definedSlice
	DefArray definedArray
	DefMap   definedMap
	Pos      token.Pos

	Stocks []StockData
}
type StockData struct {
	Symbol    string
	Intervals []Interval
}

type Interval struct {
	Start, End                     time.Time
	Open, Close, Low, High, Volume float64
}

type node struct {
	Value int
	Next  *node
}

type structType struct {
	N node
}

func TestMain(m *testing.M) {
	flag.Parse()
	if *generateTestCodeFilename != "" {
		if err := GenerateFile(*generateTestCodeFilename, "codec", generatedTestTypes{}); err != nil {
			log.Fatal(err)
		}
		fmt.Println("generated file, now run tests again")
	} else {
		os.Exit(m.Run())
	}
}

func TestEncodeDecode(t *testing.T) {
	for _, opts := range []api.EncodeOptions{
		{TrackPointers: false},
		{TrackPointers: true},
		{TrackPointers: false, AltEncodedUints: true},
		{TrackPointers: false, GobEncodedUints: true},
	} {
		t.Run(fmt.Sprintf("%+v", opts), func(t *testing.T) {
			testEncodeDecode(t, opts)
		})
	}
}

func testEncodeDecode(t *testing.T, aopts api.EncodeOptions) {
	want := []interface{}{
		nil, 1, -5, 98.6, uint64(1 << 63), "Luke Luck likes lakes", true,
		time.Date(2020, time.January, 26, 0, 0, 0, 0, time.UTC),
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
		structType{N: node{1, nil}},
		definedSlice{1, 2, 3},
		definedArray{-7},
		definedMap{"true": true},
		[]StockData{{Symbol: "XVLD", Intervals: []Interval{
			{
				time.Date(2000, time.January, 1, 9, 0, 0, 0, time.UTC),
				time.Date(2000, time.January, 1, 17, 0, 0, 0, time.UTC),
				426.91, 430.06, 424.64, 459.29, 9.696951891448457,
			},
			{
				time.Date(2000, time.January, 1, 9, 0, 0, 0, time.UTC),
				time.Date(2000, time.January, 1, 17, 0, 0, 0, time.UTC),
				426.91, 430.06, 424.64, 459.29, 9.696951891448457,
			},
			{
				time.Date(2000, time.January, 1, 9, 0, 0, 0, time.UTC),
				time.Date(2000, time.January, 1, 17, 0, 0, 0, time.UTC),
				426.91, 430.06, 424.64, 459.29, 9.696951891448457,
			},
			{
				time.Date(2000, time.January, 1, 9, 0, 0, 0, time.UTC),
				time.Date(2000, time.January, 1, 17, 0, 0, 0, time.UTC),
				426.91, 430.06, 424.64, 459.29, 9.696951891448457,
			},
			{
				time.Date(2000, time.January, 1, 9, 0, 0, 0, time.UTC),
				time.Date(2000, time.January, 1, 17, 0, 0, 0, time.UTC),
				426.91, 430.06, 424.64, 459.29, 9.696951891448457,
			},
			{
				time.Date(2000, time.January, 1, 9, 0, 0, 0, time.UTC),
				time.Date(2000, time.January, 1, 17, 0, 0, 0, time.UTC),
				426.91, 430.06, 424.64, 459.29, 9.696951891448457,
			},
		}}},
	}
	var buf bytes.Buffer
	apiOpts = aopts
	e := NewEncoder(&buf, nil)
	for _, w := range want {
		if err := e.Encode(w); err != nil {
			t.Fatalf("%#v: %v", w, err)
		}
	}
	r := bytes.NewReader(buf.Bytes())
	d := NewDecoder(r)
	for _, w := range want {
		g, err := d.Decode()
		if err != nil {
			t.Fatalf("%#v: %v", w, err)
		}
		if !cmp.Equal(g, w) {
			t.Errorf("got %v, want %v", g, w)
		}
	}
}

func TestEncodeErrors(t *testing.T) {
	// The only encoding error is an unregistered type.
	e := NewEncoder(&bytes.Buffer{}, nil)
	type MyInt int
	checkMessage(t, e.Encode(MyInt(0)), "unregistered")
}

// TODO: use fuzzing to check for panics.

// func TestDecodeErrors(t *testing.T) {
// 	for _, test := range []struct {
// 		offset  int
// 		change  byte
// 		message string
// 	}{
// 		// d.buf[d.i:] should look like: nValues 2 0 nBytes 4 ...
// 		// Induce errors by changing some bytes.
// 		{0, startCode, "bad code"},   // mess with the initial code
// 		{1, 5, "bad list length"},    // mess with the list length
// 		{2, 1, "out of range"},       // mess with the type number
// 		{3, nValuesCode, "bad code"}, // mess with the uint code
// 		{4, 5, "bad length"},         // mess with the uint length
// 	} {
// 		d := NewDecoder(bytes.NewReader(mustEncode(t, uint64(3000))))
// 		d.decodeInitial()
// 		d.buf[d.i+test.offset] = test.change
// 		_, err := d.Decode()
// 		checkMessage(t, err, test.message)
// 	}
// }

func mustEncode(t *testing.T, x interface{}) []byte {
	t.Helper()
	var buf bytes.Buffer
	e := NewEncoder(&buf, &EncodeOptions{TrackPointers: true})
	if err := e.Encode(x); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
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

// TODO: figure out how to test skipping. We need two versions of the same struct.

// func TestSkip(t *testing.T) {
// 	var buf bytes.Buffer
// 	e := NewEncoder(&buf, nil)
// 	values := []interface{}{
// 		1,
// 		false,
// 		"yes",
// 		"no",
// 		65000,
// 	}
// 	for _, v := range values {
// 		if err := e.Encode(v); err != nil {
// 			t.Fatal(err)
// 		}
// 	}

// 	d := NewDecoder(bytes.NewReader(buf.Bytes()))
// 	// Skip odd indexes.
// 	for i, want := range values {
// 		if i%2 == 0 {
// 			got, err := d.Decode()
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			if !reflect.DeepEqual(got, want) {
// 				t.Errorf("got %v, want %v", got, want)
// 			}
// 		} else {
// 			d.skip()
// 		}
// 	}
// }

func TestSharing(t *testing.T) {
	n := &node{Value: 1, Next: &node{Value: 2}}
	n.Next.Next = n // create a cycle
	d := NewDecoder(bytes.NewReader(mustEncode(t, n)))
	g, err := d.Decode()
	if err != nil {
		t.Fatal(err)
	}
	got := g.(*node)
	if !cmp.Equal(got, n) {
		t.Error("did not preserve cycle")
	}
}
