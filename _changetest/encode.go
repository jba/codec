// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build version1,encode

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/go-cmp/cmp"
	"github.com/jba/codec"
)

func main() {
	v := Change{
		A: 1,
		B: 2,
		C: 3,
	}

	f, err := os.Create("change.enc")
	if err != nil {
		log.Fatal(err)
	}
	e := codec.NewEncoder(f, &codec.EncodeOptions{TrackPointers: true})
	if err := e.Encode(v); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	// Validate
	f, err = os.Open("change.enc")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	d := codec.NewDecoder(f, nil)
	var got Change
	if err := d.Decode(&got); err != nil {
		log.Fatal(err)
	}
	if !cmp.Equal(got, v) {
		log.Fatalf("got %+v, want %+v", got, v)
	}
	fmt.Println("encoded and verified change.enc")
}
