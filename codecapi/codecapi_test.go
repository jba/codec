// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codecapi

import (
	"bytes"
	"io"
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestEncodeDecode(t *testing.T) {
	want := []interface{}{
		nil, "Luke Luck likes lakes", true,
		[]byte{},
		[]byte{1},
		[]byte{1, 2},
		[]byte{1, 2, 3},
		[]byte{1, 2, 3, 4},
		[]byte{1, 2, 3, 4, 5},
		1, -5, 255, 65000, 130_000,
		int8(-11), int16(-32000), int32(-7676767), int64(-392032393),
		uint(17), uint8(11), uint8(255), uint16(32000), uint32(7676767), uint64(392032393), uint64(1 << 63),
		uintptr(123456),
		float32(98.1234), float64(98.1234), 1.23e63,
		complex(float32(1), float32(2)), complex(3, 4),
		math.NaN(), math.Inf(1), math.Inf(-1),
	}
	var buf bytes.Buffer
	e := NewEncoder(&buf, EncodeOptions{})
	for _, w := range want {
		if err := e.Encode(w); err != nil {
			t.Fatalf("%#v: %v", w, err)
		}
	}
	d := NewDecoder(bytes.NewReader(buf.Bytes()), DecodeOptions{})
	for _, w := range want {
		var g interface{}
		if err := d.Decode(&g); err != nil {
			t.Fatalf("%#v: %v", w, err)
		}
		if !cmp.Equal(g, w, cmpopts.EquateNaNs()) {
			t.Errorf("got %v (%[1]T), want %v (%[2]T)", g, w)
		}
	}
}

func TestDecodeEOF(t *testing.T) {
	var dopts DecodeOptions
	var got interface{}
	err := NewDecoder(bytes.NewReader(nil), dopts).Decode(&got)
	if got != nil || err != io.EOF {
		t.Errorf("got (%v, %v), want (nil, io.EOF)", got, err)
	}
	err = NewDecoder(bytes.NewReader(header), dopts).Decode(&got)
	if got != nil || err != io.EOF {
		t.Errorf("got (%v, %v), want (nil, io.EOF)", got, err)
	}

	var buf bytes.Buffer
	e := NewEncoder(&buf, EncodeOptions{})
	for i := 0; i < 3; i++ {
		if err := e.Encode(i); err != nil {
			t.Fatal(err)
		}
	}
	var igot int
	d := NewDecoder(bytes.NewReader(buf.Bytes()), DecodeOptions{})
	for i := 0; i < 3; i++ {
		err := d.Decode(&igot)
		if err != nil {
			t.Fatal(err)
		}
		if igot != i {
			t.Fatalf("got %d, want %d", igot, i)
		}
	}
	igot = -1
	err = d.Decode(&igot)
	if igot != -1 || err != io.EOF {
		t.Errorf("got (%v, %v), want (-1, io.EOF)", igot, err)
	}
}
