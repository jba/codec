// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build ignore

package main

import (
	"log"
	"os"

	"github.com/jba/codec/internal/benchmarks/data"
)

func main() {
	things := os.Args[1:]
	if len(things) == 0 {
		log.Fatal("need things to generate")
	}
	if err := data.Generate(things); err != nil {
		log.Fatal(err)
	}
}
