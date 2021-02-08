// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build version2

// In this version of the Change struct, fields A, B and C have been reordered
// and field Add has been added.
//
// Removing a field is tested in ../_skiptest. It can't be tested here, because
// for this to work the generated file must be kept around, and it can't compile
// if a field is removed.

package main

type Change struct {
	Add int
	C   int
	A   int
	B   int
}
