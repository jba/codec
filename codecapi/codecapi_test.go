// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codecapi

import (
	"math"
	"math/bits"
)

//func BenchmarkFloatEncoding(t *testing.T) {

// }

// func encodeFloatsDirect(w io.Writer, fs []float64) error {
// 	// Encode a 1-byte code for # of bytes to follow.

// 	for _, f := range fs {

// func floatBitsDirect(f float64) uint64     { return math.Float64bits(f) }
// func floatFromBitsDirect(u uint64) float64 { return math.Float64frombits(u) }

// // FloatBitsReverse returns a uint64 holding the bits of a floating-point number.
// // Floating-point numbers are transmitted as uint64s holding the bits
// // of the underlying representation. They are sent byte-reversed, with
// // the exponent end coming out first, so integer floating point numbers
// // (for example) transmit more compactly. This routine does the
// // swizzling.
// // From encoding/gob.
// func floatBitsReversed(f float64) uint64 {
// 	return bits.ReverseBytes64(math.Float64bits(f))
// }

// func floatFromBitsReversed(u uint64) float64 {
// 	return math.Float64frombits(bits.ReverseBytes64(u))
// }

// floatBits returns a uint64 holding the bits of a floating-point number.
// Floating-point numbers are transmitted as uint64s holding the bits
// of the underlying representation. They are sent byte-reversed, with
// the exponent end coming out first, so integer floating point numbers
// (for example) transmit more compactly. This routine does the
// swizzling.
func floatBits(f float64) uint64 {
	u := math.Float64bits(f)
	return bits.ReverseBytes64(u)
}
