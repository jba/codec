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
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/jba/codec"
)

func main() {
	f, err := os.Open("skip.enc")
	if err != nil {
		log.Fatalf("FAIL: %v", err)
	}
	d := codec.NewDecoder(f, nil)

	var got Skip
	if err := d.Decode(&got); err != nil {
		log.Fatalf("FAIL: %v", err)
	}
	fmt.Println("decoded skip.enc")
	var want Skip
	if !cmp.Equal(got, want) {
		log.Fatalf("FAIL: got %+v, want %+v", got, want)
	}
	f.Close()

	f, err = os.Open("skip.enc")
	if err != nil {
		log.Fatal("FAIL: %v", err)
	}
	defer f.Close()
	d = codec.NewDecoder(f, &codec.DecodeOptions{DisallowUnknownFields: true})
	if err := d.Decode(&got); err == nil {
		log.Fatal("FAIL: got nil, want error")
	} else if !strings.Contains(err.Error(), "unknown field") {
		log.Fatalf("FAIL: Decode error is unexpected: %q", err)
	}
	fmt.Println("PASS")
}
