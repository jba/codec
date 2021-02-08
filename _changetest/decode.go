// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build version2,decode

// This can't be a normal XXX_test.go file, because the test binary's package
// name (as in reflect.TypeOf, etc.) is not "main", but something constructed
// from the module path.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/go-cmp/cmp"
	"github.com/jba/codec"
)

func main() {
	f, err := os.Open("change.enc")
	if err != nil {
		log.Fatalf("FAIL: %v", err)
	}
	defer f.Close()

	d := codec.NewDecoder(f, nil)

	var got Change
	if err := d.Decode(&got); err != nil {
		log.Fatalf("FAIL: %v", err)
	}
	fmt.Println("decoded change.enc")

	// The Add field is in the struct but not the encoding, so it retains its zero value.
	want := Change{
		A:   1,
		B:   2,
		C:   3,
		Add: 0,
	}
	if !cmp.Equal(got, want) {
		log.Fatalf("FAIL: got %+v, want %+v", got, want)
	}
	fmt.Println("PASS")
}
