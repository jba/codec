// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package codec implements an encoder for Go values. It relies on code generation
rather than reflection, so it is significantly faster than reflection-based
encoders like gob. It can also preserve sharing among pointers (but not other
forms of sharing, like sub-slices).

Struct Tags



*/
package codec
