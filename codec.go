// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codec

import (
	"io"

	api "github.com/jba/codec/codecapi"
)

// An Encoder encodes Go values into a sequence of bytes.
type Encoder struct {
	state *api.Encoder
}

// EncodeOptions holds options for encoding.
type EncodeOptions struct {
	// If TrackPointers is true, the encoder will keep track of pointers so it
	// can preserve the pointer topology of the encoded value. Cyclical and
	// shared values will decode to the same representation. If TrackPointers is
	// false, then shared pointers will decode to distinct values, and cycles
	// will result in stack overflow.
	//
	// Setting this to true will significantly slow down encoding.
	TrackPointers bool
}

// NewEncoder returns an Encoder that writes to w.
func NewEncoder(w io.Writer, opts *EncodeOptions) *Encoder {
	return &Encoder{state: api.NewEncoder(w, api.EncodeOptions{TrackPointers: opts != nil && opts.TrackPointers})}
}

// Encode encodes x.
func (e *Encoder) Encode(x interface{}) (err error) {
	return e.state.Encode(x)
}

// A Decoder decodes a Go value encoded by an Encoder.
// To use a Decoder:
// - Pass NewDecoder the return value of Encoder.Bytes.
// - Call the Decode method once for each call to Encoder.Encode.
type Decoder struct {
	state *api.Decoder
}

// NewDecoder creates a Decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{state: api.NewDecoder(r)}
}

// Decode decodes a value encoded with Encoder.Encode.
// It returns (nil, io.EOF) if there are no more values.
func (d *Decoder) Decode() (_ interface{}, err error) {
	return d.state.Decode()
}
