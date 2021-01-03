// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build debug

package codecapi

import "fmt"

func (d *Decoder) dump() {
	for d.i < len(d.buf) {
		d.dump1(0)
	}
}

func (d *Decoder) dump1(level int) {
	fmt.Printf("%3d ", d.i)
	for i := 0; i < level; i++ {
		fmt.Print("    ")
	}
	b := d.readByte()
	if b < endCode {
		// Small integers represent themselves in a single byte.
		fmt.Printf("smallint %d\n", b)
		return
	}
	if b >= bytes0Code && b <= bytes4Code {
		n := int(b - bytes0Code)
		fmt.Printf("%d bytes: %v\n", n, d.readBytes(n))
		return
	}
	switch b {
	case nilCode:
		fmt.Println("nil")
	case nBytesCode:
		n := int(d.DecodeUint())
		fmt.Printf("nBytes %d", n)
		data := d.readBytes(n)
		if n < 10 {
			fmt.Printf(" %v", data)
		}
		fmt.Println()
	case nValuesCode:
		// A uint n and n values follow.
		n := int(d.DecodeUint())
		fmt.Printf("nValues %d\n", n)
		for i := 0; i < n; i++ {
			d.dump1(level + 1)
		}
	case ptrCode:
		fmt.Println("ptr")
		d.dump1(level)
	case refPtrCode:
		fmt.Println("refPtr")
		d.dump1(level)
	case refCode:
		// A uint follows.
		fmt.Printf("ref %d\n", d.DecodeUint())
	case startCode:
		fmt.Println("start")
		for d.curByte() != endCode {
			d.dump1(level + 1)
		}
		d.dump1(level) // endCode
	case endCode:
		fmt.Println("end")
	default:
		d.badcode(b)
	}
}
