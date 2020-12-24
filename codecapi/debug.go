// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build debug

package codecapi

import "fmt"

func (d *Decoder) dump() {
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
		fmt.Printf("%d bytes: %v\n", n, d.readBytes(n))
	case nValuesCode:
		// A uint n and n values follow.
		n := int(d.DecodeUint())
		for i := 0; i < n; i++ {
			d.dump()
		}
	case refCode:
		// A uint follows.
		fmt.Printf("ref: %d\n", d.DecodeUint())
	case startCode:
		fmt.Println("startCode")
		for d.curByte() != endCode {
			d.dump()
		}
		d.readByte() // consume the endCode byte
		fmt.Println("endCode")
	default:
		d.badcode(b)
	}
}
