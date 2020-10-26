// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"math/bits"
	"reflect"
	"runtime"
)

const uint64Size = 8

// An Encoder encodes Go values into a sequence of bytes.
// To use an Encoder:
// - Create one with NewEncoder.
// - Call the Encode method one or more times.
// - Retrieve the resulting bytes by calling Bytes.
type Encoder struct {
	w        io.Writer
	buf      []byte
	typeNums map[reflect.Type]int
	seen     map[uintptr]uint64 // for references; see StartStruct
	gobBuf   [1 + uint64Size]byte
}

type EncodeOptions struct {
	TrackPointers bool
}

// Header for an encoded stream.
// A 3-byte identifier and a version number.
var header = []byte("GJC1")

// NewEncoder returns an Encoder.
func NewEncoder(w io.Writer, opts *EncodeOptions) *Encoder {
	e := &Encoder{w: w}
	if opts != nil && opts.TrackPointers {
		e.seen = map[uintptr]uint64{}
	}
	return e
}

func totalAlloc() uint64 {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ms.TotalAlloc
}

// Encode encodes x.
func (e *Encoder) Encode(x interface{}) (err error) {
	if e.typeNums == nil {
		// First call to encode: write the header.
		if _, err := e.w.Write(header); err != nil {
			return err
		}
	}
	if e.buf != nil {
		e.buf = e.buf[0:0]
	} else {
		e.buf = make([]byte, 0, 1_000_000)
	}
	e.typeNums = map[reflect.Type]int{}
	defer handlePanic(&err)

	//start := totalAlloc()
	//fmt.Printf("before EncodeAny: total alloc = %d\n", start)
	e.EncodeAny(x)
	//fmt.Printf("after EncodeAny: total alloc = %d (%.1fM)\n", totalAlloc(), float64(totalAlloc()-start)/(1024*1024))
	data := e.buf // remember the data
	//fmt.Printf("len(data)=%d\n", len(data))
	e.buf = nil       // start with a fresh buffer
	e.encodeInitial() // encode metadata
	initial := e.buf  // remember that

	// Encode total size in 8 bytes.
	var szbuf [8]byte
	binary.LittleEndian.PutUint64(szbuf[:], uint64(len(initial)+len(data)))
	if _, err := e.w.Write(szbuf[:]); err != nil {
		return err
	}
	if _, err := e.w.Write(initial); err != nil {
		return err
	}
	_, err = e.w.Write(data)
	return err
}

// Bytes returns the encoded byte slice.
// func (e *Encoder) Bytes() []byte {
// 	data := e.buf                 // remember the data
// 	e.buf = nil                   // start with a fresh buffer
// 	e.encodeInitial()             // encode metadata
// 	return append(e.buf, data...) // concatenate metadata and data
// }

// A Decoder decodes a Go value encoded by an Encoder.
// To use a Decoder:
// - Pass NewDecoder the return value of Encoder.Bytes.
// - Call the Decode method once for each call to Encoder.Encode.
type Decoder struct {
	r          io.Reader
	buf        []byte
	i          int // offset
	typeCodecs []TypeCodec
	storeRefs  bool
	refs       []interface{} // list of struct pointers, in the order seen
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// NewDecoderBytes returns a Decoder for the given bytes.
func NewDecoderBytes(data []byte) *Decoder {
	return nil
}

// Decode decodes a value encoded with Encoder.Encode.
// It returns (nil, io.EOF) if there are no more values.
func (d *Decoder) Decode() (_ interface{}, err error) {
	if d.buf == nil {
		// First call to decode: read header.
		var buf [4]byte
		if _, err := io.ReadFull(d.r, buf[:]); err != nil {
			return nil, err
		}
		if !bytes.Equal(header, buf[:]) {
			return nil, fmt.Errorf("bad header: got %q, want %q", buf[:], header)
		}
	}
	var szbuf [8]byte
	if _, err := io.ReadFull(d.r, szbuf[:]); err != nil {
		return nil, err
	}
	sz := binary.LittleEndian.Uint64(szbuf[:])
	// TODO: We can reuse d.buf iff we didn't let slices from it escape previously.
	d.buf = make([]byte, sz)
	d.i = 0
	if _, err := io.ReadFull(d.r, d.buf); err != nil {
		return nil, err
	}
	defer handlePanic(&err)
	d.decodeInitial()
	return d.DecodeAny(), nil
}

// Expose Fail, but only as methods so they are more likely to be called during
// encoding and decoding.

// Fail causes e.Encode to stop and return with err.
func (*Encoder) Fail(err error) { fail(err) }

// Fail causes d.Decode to stop and return with err.
func (*Decoder) Fail(err error) { fail(err) }

func badcode(c byte) {
	failf("bad code: %d", c)
}

//////////////// Reading From and Writing To the Buffer

func (e *Encoder) writeByte(b byte) {
	e.buf = append(e.buf, b)
}

func (d *Decoder) readByte() byte {
	b := d.curByte()
	d.i++
	return b
}

// curByte returns the next byte to be read
// without actually consuming it.
func (d *Decoder) curByte() byte {
	return d.buf[d.i]
}

func (e *Encoder) writeBytes(b []byte) {
	e.buf = append(e.buf, b...)
}

// readBytes reads and returns the given number of bytes.
// It fails if there are not enough bytes in the input.
func (d *Decoder) readBytes(n int) []byte {
	d.i += n
	return d.buf[d.i-n : d.i]
}

func (e *Encoder) writeString(s string) {
	e.buf = append(e.buf, s...)
}

func (d *Decoder) readString(len int) string {
	return string(d.readBytes(len))
}

func (e *Encoder) writeUint32(u uint32) {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], u)
	e.writeBytes(buf[:])
}

func (d *Decoder) readUint32() uint32 {
	return binary.LittleEndian.Uint32(d.readBytes(4))
}

func (e *Encoder) writeUint64(u uint64) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], u)
	e.writeBytes(buf[:])
}

func (d *Decoder) readUint64() uint64 {
	return binary.LittleEndian.Uint64(d.readBytes(8))
}

//////////////// Encoding Scheme

// Byte codes that begin each encoded value.
// See the package doc for their descriptions.
const (
	nilCode = 255 - iota // a nil value
	// reserve a few values for future use
	reserved1
	reserved2
	reserved3
	reserved4
	reserved5
	reserved6
	ptrCode     // non-nil, non-ref pointer
	nBytesCode  // uint n follows, then n bytes
	nValuesCode // uint n follows, then n values
	refCode     // uint n follows, referring to a previous value
	startCode   // start of a value of indeterminate length
	endCode     // end of a value that began with start
	// Bytes less than endCode represent themselves.
)

// EncodeUint encodes a uint64.
func (e *Encoder) EncodeUint(u uint64) {
	switch {
	case u < endCode:
		// u fits into the initial byte.
		e.writeByte(byte(u))
	case u <= math.MaxUint32:
		// Encode as a sequence of 4 bytes, the little-endian representation of
		// a uint32.
		e.writeByte(nBytesCode)
		e.writeByte(4)
		e.writeUint32(uint32(u))
	default:
		// Encode as a sequence of 8 bytes, the little-endian representation of
		// a uint64.
		e.writeByte(nBytesCode)
		e.writeByte(8)
		e.writeUint64(u)
	}
}

// DecodeUint decodes a uint64.
func (d *Decoder) DecodeUint() uint64 {
	b := d.readByte()
	switch {
	case b < endCode:
		return uint64(b)
	case b == nBytesCode:
		switch n := d.readByte(); n {
		case 4:
			return uint64(d.readUint32())
		case 8:
			return d.readUint64()
		default:
			failf("DecodeUint: bad length %d", n)
		}
	default:
		badcode(b)
	}
	return 0
}

// for benchmarking
func (e *Encoder) encodeUintGob(x uint64) {
	// Code from encoding/gob/encode.go:encodeUint.
	if x <= 0x7F {
		e.writeByte(uint8(x))
		return
	}
	binary.BigEndian.PutUint64(e.gobBuf[1:], x)
	bc := bits.LeadingZeros64(x) >> 3     // 8 - bytelen(x)
	e.gobBuf[bc] = uint8(bc - uint64Size) // and then we subtract 8 to get -bytelen(x)
	e.writeBytes(e.gobBuf[bc : uint64Size+1])
}

func (d *Decoder) decodeUintGob() (x uint64) {
	b := d.readByte()
	if b <= 0x7f {
		return uint64(b)
	}
	n := -int(int8(b))
	if n > uint64Size {
		fail(errors.New("bad uint"))
	}
	buf := d.readBytes(n)
	// Don't need to check error; it's safe to loop regardless.
	// Could check that the high byte is zero but it's not worth it.
	for _, b := range buf {
		x = x<<8 | uint64(b)
	}
	return x
}

// EncodeInt encodes a signed integer.
func (e *Encoder) EncodeInt(i int64) {
	// Encode small negative as well as small positive integers efficiently.
	// Algorithm from gob; see "Encoding Details" at https://pkg.go.dev/encoding/gob.
	var u uint64
	if i < 0 {
		u = (^uint64(i) << 1) | 1 // complement i, bit 0 is 1
	} else {
		u = (uint64(i) << 1) // do not complement i, bit 0 is 0
	}
	e.EncodeUint(u)
}

// DecodeInt decodes a signed integer.
func (d *Decoder) DecodeInt() int64 {
	u := d.DecodeUint()
	if u&1 == 1 {
		return int64(^(u >> 1))
	}
	return int64(u >> 1)
}

// encodeLen encodes the length of a byte sequence.
func (e *Encoder) encodeLen(n int) {
	e.writeByte(nBytesCode)
	e.EncodeUint(uint64(n))
}

// decodeLen decodes the length of a byte sequence.
func (d *Decoder) decodeLen() int {
	if b := d.readByte(); b != nBytesCode {
		badcode(b)
	}
	return int(d.DecodeUint())
}

// EncodeBytes encodes a byte slice.
func (e *Encoder) EncodeBytes(b []byte) {
	e.encodeLen(len(b))
	e.writeBytes(b)
}

// DecodeBytes decodes a byte slice.
// It DOES copy.
// TODO: It does no copying.
func (d *Decoder) DecodeBytes() []byte {
	n := d.decodeLen()
	b := make([]byte, n)
	copy(b, d.readBytes(n))
	return b
}

// EncodeString encodes a string.
func (e *Encoder) EncodeString(s string) {
	e.encodeLen(len(s))
	e.writeString(s)
}

// DecodeString decodes a string.
func (d *Decoder) DecodeString() string {
	return d.readString(d.decodeLen())
}

// EncodeBool encodes a bool.
func (e *Encoder) EncodeBool(b bool) {
	if b {
		e.writeByte(1)
	} else {
		e.writeByte(0)
	}
}

// DecodeBool decodes a bool.
func (d *Decoder) DecodeBool() bool {
	b := d.readByte()
	switch b {
	case 0:
		return false
	case 1:
		return true
	default:
		failf("bad bool: %d", b)
		return false
	}
}

// EncodeFloat encodes a float64.
func (e *Encoder) EncodeFloat(f float64) {
	e.EncodeUint(math.Float64bits(f))
}

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

// DecodeFloat decodes a float64.
func (d *Decoder) DecodeFloat() float64 {
	return math.Float64frombits(d.DecodeUint())
}

func (e *Encoder) EncodeNil() {
	e.writeByte(nilCode)
}

// StartList should be called before encoding any sequence of variable-length
// values.
func (e *Encoder) StartList(len int) {
	e.writeByte(nValuesCode)
	e.EncodeUint(uint64(len))
}

// StartList should be called before decoding any sequence of variable-length
// values. It returns -1 if the encoded list was nil. Otherwise, it returns the
// length of the sequence.
func (d *Decoder) StartList() int {
	switch b := d.readByte(); b {
	case nilCode:
		return -1
	case nValuesCode:
		return int(d.DecodeUint())
	default:
		badcode(b)
		return 0
	}
}

//////////////// Pointer Support

// StartPtr should be called before encoding a pointer. The isNil
// argument says whether the pointer is nil. The p argument is the
// pointer. If StartPtr returns false, encoding should not proceed.
func (e *Encoder) StartPtr(isNil bool, p interface{}) bool {
	if isNil {
		e.EncodeNil()
		return false
	}
	if e.seen != nil {
		ptr := reflect.ValueOf(p).Pointer()
		if u, ok := e.seen[ptr]; ok {
			// If we have already seen this struct pointer,
			// encode a reference to it.
			e.writeByte(refCode)
			e.EncodeUint(u)
			return false // Caller should not encode the struct.
		}
		// Note that we have seen this pointer, and assign it
		// its position in the encoding.
		e.seen[ptr] = uint64(len(e.seen))
	}
	e.writeByte(ptrCode)
	return true
}

// StartPtr should be called before decoding a pointer.
// If the first return value is false, the destination should not be set.
// Otherwise, if the second return value is non-nil, assign it to the destination.
// Otherwise, proceed with decoding into the destination.
func (d *Decoder) StartPtr() (bool, interface{}) {
	b := d.readByte()
	switch b {
	case nilCode: // do not set the pointer
		return false, nil
	case refCode:
		u := d.DecodeUint()
		return true, d.refs[u]
	case ptrCode:
		return true, nil
	default:
		badcode(b)
		return false, nil // unreached, needed for compiler
	}
}

//////////////// Struct Support

// TODO: now for struct pointers we write
//    ptrCode
//    startCode
// which is a little redundant. What do we gain by optimizing the ptrCode away?

func (e *Encoder) StartStruct() {
	e.writeByte(startCode)
}

// StartStruct should be called before decoding a struct pointer. If it returns
// false, decoding should not proceed. If it returns true and the second return
// value is non-nil, it is a reference to a previous value and should be used
// instead of proceeding with decoding.
func (d *Decoder) StartStruct() {
	if b := d.readByte(); b != startCode {
		badcode(b)
	}
}

// StoreRef should be called by a struct decoder immediately after it allocates
// a struct pointer.
func (d *Decoder) StoreRef(p interface{}) {
	if d.storeRefs {
		d.refs = append(d.refs, p)
	}
}

// EndStruct should be called after encoding a struct.
func (e *Encoder) EndStruct() {
	e.writeByte(endCode)
}

// NextStructField should be called by a struct decoder in a loop.
// It returns the field number of the next encoded field, or -1
// if there are no more fields.
func (d *Decoder) NextStructField() int {
	if d.curByte() == endCode {
		d.readByte() // consume the end byte
		return -1
	}
	return int(d.DecodeUint())
}

// UnknownField should be called by a struct decoder
// when it sees a field number that it doesn't know.
func (d *Decoder) UnknownField(typeName string, num int) {
	d.skip()
}

// skip reads past a value in the input.
func (d *Decoder) skip() {
	b := d.readByte()
	if b < endCode {
		// Small integers represent themselves in a single byte.
		return
	}
	switch b {
	case nilCode:
		// Nothing follows.
	case nBytesCode:
		// A uint n and n bytes follow. It is efficient to call readBytes here
		// because it does no allocation.
		d.readBytes(int(d.DecodeUint()))
	case nValuesCode:
		// A uint n and n values follow.
		n := int(d.DecodeUint())
		for i := 0; i < n; i++ {
			d.skip()
		}
	case refCode:
		// A uint follows.
		d.DecodeUint()
	case startCode:
		// Skip until we see endCode.
		for d.curByte() != endCode {
			d.skip()
		}
		d.readByte() // consume the endCode byte
	default:
		badcode(b)
	}
}

//////////////// Encoding Arbitrary Values

// EncodeAny encodes a Go type. The type must have
// been registered with Register.
func (e *Encoder) EncodeAny(x interface{}) {
	// Encode a nil interface value with a zero.
	if x == nil {
		e.writeByte(0)
		return
	}
	// Find the TypeCodec for the type, which has the encoder.
	t := reflect.TypeOf(x)
	tc := typeCodecsByType[t]
	if tc == nil {
		failf("unregistered type %q", t)
	}
	// Assign a number to the type if we haven't already.
	num, ok := e.typeNums[t]
	if !ok {
		num = len(e.typeNums)
		e.typeNums[t] = num
	}
	// Encode a 2-element list of the type number and the encoded value.
	e.StartList(2)
	e.EncodeUint(uint64(num))
	tc.Encode(e, x)
}

// DecodeAny decodes a value encoded by EncodeAny.
func (d *Decoder) DecodeAny() interface{} {
	// If we're looking at a zero, this is a nil interface.
	if d.curByte() == 0 {
		d.readByte() // consume the byte
		return nil
	}
	// Otherwise, we should have a two-item list: type number and value.
	n := d.StartList()
	if n != 2 {
		failf("DecodeAny: bad list length %d", n)
	}
	num := d.DecodeUint()
	if num >= uint64(len(d.typeCodecs)) {
		failf("type number %d out of range", num)
	}
	tc := d.typeCodecs[num]
	return tc.Decode(d)
}

// encodeInitial encodes metadata that appears at the start of the
// encoded byte slice.
func (e *Encoder) encodeInitial() {
	// Encode the list of type names we saw, in the order we
	// assigned numbers to them.
	names := make([]string, len(e.typeNums))
	for t, num := range e.typeNums {
		names[num] = typeName(t)
	}
	e.StartList(len(names))
	for _, n := range names {
		e.EncodeString(n)
	}
	// Encode whether we tracked pointers.
	e.EncodeBool(e.seen != nil)
}

// decodeInitial decodes metadata that appears at the start of the
// encoded byte slice.
func (d *Decoder) decodeInitial() {
	// Decode the list of type names. The number of a type is its position in
	// the list.
	n := d.StartList()
	d.typeCodecs = make([]TypeCodec, n)
	for num := 0; num < n; num++ {
		name := d.DecodeString()
		tc := typeCodecsByName[name]
		if tc == nil {
			failf("unregistered type: %s", name)
		}
		d.typeCodecs[num] = tc
	}
	d.storeRefs = d.DecodeBool()
}

//////////////// Errors

func handlePanic(errp *error) {
	r := recover()
	if r == nil {
		// No panic; do nothing.
		return
	}
	// If the panic is not from this package, re-panic.
	cerr, ok := r.(codecError)
	if !ok {
		panic(r)
	}
	// Otherwise, set errp.
	*errp = cerr.err
}

func fail(err error) {
	panic(codecError{err})
}

func failf(format string, args ...interface{}) {
	fail(fmt.Errorf(format, args...))
}

// codecError wraps errors from this package so handlePanic
// can distinguish them.
type codecError struct {
	err error
}

func (c codecError) String() string { return c.err.Error() }
