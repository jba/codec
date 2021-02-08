// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/jba/codec"
)

func main() {
	err := codec.GenerateFile("skip.gen.go", "main", nil, []*Skip(nil))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("generated skip.gen.go")
}
