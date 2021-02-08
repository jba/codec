// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build gen

package main

import (
	"fmt"
	"log"

	"github.com/jba/codec"
)

func main() {
	err := codec.GenerateFile("change.gen.go", "main", nil, Change{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("generated change.gen.go")
}
