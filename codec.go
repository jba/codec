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

	// If non-nil, Encode will use this buffer instead of creating one. If the
	// encoding is large, providing a buffer of sufficient size can speed up
	// encoding by reducing allocation.
	Buffer []byte
}

// NewEncoder returns an Encoder that writes to w.
func NewEncoder(w io.Writer, opts *EncodeOptions) *Encoder {
	aopts := api.EncodeOptions{}
	if opts != nil {
		aopts.TrackPointers = opts.TrackPointers
		aopts.Buffer = opts.Buffer
	}
	return &Encoder{state: api.NewEncoder(w, aopts)}
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

// DecodeOptions holds options for Decoding.
type DecodeOptions struct {
	// FailOnUnknownField configures whether unknown struct fields are skipped
	// (the default) or cause decoding to fail immediately.
	FailOnUnknownField bool
}

// NewDecoder creates a Decoder that reads from r.
func NewDecoder(r io.Reader, opts *DecodeOptions) *Decoder {
	aopts := api.DecodeOptions{}
	if opts != nil {
		aopts.FailOnUnknownField = opts.FailOnUnknownField
	}
	return &Decoder{state: api.NewDecoder(r, aopts)}
}

// Decode decodes a value encoded with Encoder.Encode.
// It returns (nil, io.EOF) if there are no more values.
func (d *Decoder) Decode() (_ interface{}, err error) {
	return d.state.Decode()
}
