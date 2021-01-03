// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package codecapi is used by the codec package and by code generated by
// codec.GenerateFile. It should NOT be used directly.
package codecapi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/bits"
	"reflect"
)

const uint64Size = 8

// Header for an encoded stream.
// A 3-byte identifier and a version number.
var header = []byte("GJC0")

type Encoder struct {
	opts     EncodeOptions
	w        io.Writer
	buf      []byte
	typeNums map[reflect.Type]int
	seen     map[uintptr]int // for references; see StartStruct
}

type EncodeOptions struct {
	TrackPointers bool
	Buffer        []byte
}

func NewEncoder(w io.Writer, opts EncodeOptions) *Encoder {
	e := &Encoder{w: w, opts: opts}
	if e.opts.TrackPointers {
		e.seen = make(map[uintptr]int, 1000)
	}
	return e
}

// Encode encodes x.
func (e *Encoder) Encode(x interface{}) (err error) {
	// Each call to Encode results in the following output:
	// - A size in bytes (uint64)
	// - Initial metadata
	// - The encoded value

	if e.typeNums == nil {
		// First call to encode: write the header.
		if _, err := e.w.Write(header); err != nil {
			return err
		}
	}
	if e.buf != nil {
		e.buf = e.buf[:0]
	} else if e.opts.Buffer != nil {
		e.buf = e.opts.Buffer[:0]
	} else {
		e.buf = make([]byte, 0, 64*1024)
	}
	e.typeNums = map[reflect.Type]int{}
	defer handlePanic(&err)

	e.EncodeAny(x)
	data := e.buf     // remember the data
	e.buf = nil       // start with a fresh buffer
	e.encodeInitial() // encode metadata
	initial := e.buf  // remember that
	e.buf = data      // restore e.buf for next call to Encode

	// Encode total size in a uint64.
	var buf [uint64Size]byte
	binary.BigEndian.PutUint64(buf[:], uint64(len(initial)+len(data)))
	if _, err := e.w.Write(buf[:]); err != nil {
		return err
	}
	if _, err := e.w.Write(initial); err != nil {
		return err
	}
	_, err = e.w.Write(data)
	return err
}

type Decoder struct {
	opts       DecodeOptions
	r          io.Reader
	buf        []byte
	i          int // offset into buf
	typeCodecs []typeCodec
	storeIndex int                 // for StartPtr to communicate with StoreRef
	refMap     map[int]interface{} // from buf offset to pointer
}

type DecodeOptions struct {
	FailOnUnknownField bool
}

func NewDecoder(r io.Reader, opts DecodeOptions) *Decoder {
	return &Decoder{r: r, opts: opts}
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
	var szbuf [uint64Size]byte
	if _, err := io.ReadFull(d.r, szbuf[:]); err != nil {
		return nil, err
	}
	sz := binary.BigEndian.Uint64(szbuf[:])
	// We can reuse d.buf because we didn't let slices from it escape previously.
	if cap(d.buf) >= int(sz) {
		d.buf = d.buf[:sz]
	} else {
		d.buf = make([]byte, sz)
	}
	d.i = 0
	if _, err := io.ReadFull(d.r, d.buf); err != nil {
		return nil, err
	}
	defer handlePanic(&err)
	d.decodeInitial()

	return d.DecodeAny(), nil
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
// It panics if there are not enough bytes in the input.
// It does not copy.
func (d *Decoder) readBytes(n int) []byte {
	d.i += n
	return d.buf[d.i-n : d.i]
}

func (e *Encoder) writeString(s string) {
	e.buf = append(e.buf, s...)
}

//////////////// Encoding Scheme

// Byte codes that begin each encoded value.
// See README.md for their descriptions.
const (
	nilCode    = 255 - iota // a nil value
	bytes0Code              // N bytes follow
	bytes1Code
	bytes2Code
	bytes3Code
	bytes4Code
	nBytesCode  // uint n follows, then n bytes
	nValuesCode // uint n follows, then n values
	ptrCode     // non-nil, non-ref pointer
	refPtrCode  // non-nil, non-ref pointer that has a later ref
	refCode     // uint n follows: relative offset of previous refPtrCode
	// reserve a few values for future use
	reserved1
	reserved2
	reserved3
	startCode // start of a value of indeterminate length
	endCode   // end of a value that began with start
	// Bytes less than endCode represent themselves.
)

// EncodeUint encodes a uint64.
func (e *Encoder) EncodeUint(u uint64) {
	var buf [uint64Size]byte
	switch {
	case u < endCode:
		// u fits into the initial byte.
		e.writeByte(byte(u))

	case u <= math.MaxUint8:
		e.writeByte(bytes1Code)
		e.writeByte(byte(u))

	case u <= math.MaxUint16:
		// Encode as a sequence of 2 bytes, the big-endian representation of
		// a uint16.
		e.writeByte(bytes2Code)
		binary.BigEndian.PutUint16(buf[:2], uint16(u))
		e.writeBytes(buf[:2])

	case u <= math.MaxUint32:
		// Encode as a sequence of 4 bytes, the big-endian representation of
		// a uint32.
		e.writeByte(bytes4Code)
		binary.BigEndian.PutUint32(buf[:4], uint32(u))
		e.writeBytes(buf[:4])

	default:
		// Encode as a sequence of 8 bytes, the big-endian representation of
		// a uint64.
		e.writeByte(nBytesCode)
		e.writeByte(uint64Size)
		binary.BigEndian.PutUint64(buf[:], u)
		e.writeBytes(buf[:])
	}
}

// DecodeUint decodes a uint64.
func (d *Decoder) DecodeUint() uint64 {
	b := d.readByte()
	switch {
	case b < endCode:
		return uint64(b)
	case b == bytes1Code:
		return uint64(d.readByte())
	case b == bytes2Code:
		return uint64(binary.BigEndian.Uint16(d.readBytes(2)))
	case b == bytes4Code:
		return uint64(binary.BigEndian.Uint32(d.readBytes(4)))
	case b == nBytesCode:
		if d.readByte() == uint64Size {
			return binary.BigEndian.Uint64(d.readBytes(uint64Size))
		}
		Failf("DecodeUint: bad length")
	default:
		d.badcode(b)
	}
	return 0
}

func (e *Encoder) EncodeByte(b byte) {
	e.writeByte(b)
}

func (d *Decoder) DecodeByte() byte {
	return d.readByte()
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
	if n <= 4 {
		e.writeByte(byte(bytes0Code - n))
	} else {
		e.writeByte(nBytesCode)
		e.EncodeUint(uint64(n))
	}
}

// decodeLen decodes the length of a byte sequence.
func (d *Decoder) decodeLen() int {
	b := d.readByte()
	return d.resolveLen(b)
}

func (d *Decoder) resolveLen(b byte) int {
	switch b {
	case nBytesCode:
		return int(d.DecodeUint())
	case bytes0Code:
		return 0
	case bytes1Code:
		return 1
	case bytes2Code:
		return 2
	case bytes3Code:
		return 3
	case bytes4Code:
		return 4
	default:
		d.badcode(b)
		panic("unreachable")
	}
}

// EncodeBytes encodes a byte slice.
func (e *Encoder) EncodeBytes(b []byte) {
	e.encodeLen(len(b))
	e.writeBytes(b)
}

// DecodeBytes decodes a byte slice.
// It makes a copy of a portion of the underlying buffer.
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
	return string(d.readBytes(d.decodeLen()))
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
		Failf("bad bool: %d", b)
		return false
	}
}

// EncodeFloat encodes a float64.
func (e *Encoder) EncodeFloat(f float64) {
	e.EncodeUint(bits.ReverseBytes64(math.Float64bits(f)))
}

// DecodeFloat decodes a float64.
func (d *Decoder) DecodeFloat() float64 {
	return math.Float64frombits(bits.ReverseBytes64(d.DecodeUint()))
}

// EncodeComplex encodes a complex128.
func (e *Encoder) EncodeComplex(c complex128) {
	e.StartList(2)
	e.EncodeFloat(real(c))
	e.EncodeFloat(imag(c))
}

// DecodeComplex decodes a complex128.
func (d *Decoder) DecodeComplex() complex128 {
	n := d.StartList()
	if n != 2 {
		Failf("DecodeComplex: bad list length %d", n)
	}
	return complex(d.DecodeFloat(), d.DecodeFloat())
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
		d.badcode(b)
		return 0
	}
}

//////////////// Pointer Support

// StartPtr should be called before encoding a pointer. The isNil
// argument says whether the pointer is nil. The p argument is the
// pointer. If StartPtr returns false, encoding should not proceed.
func (e *Encoder) StartPtr(isNil bool, p interface{}) bool {
	if reflect.ValueOf(p).IsNil() {
		e.EncodeNil()
		return false
	}
	if e.seen != nil {
		ptr := reflect.ValueOf(p).Pointer()
		if u, ok := e.seen[ptr]; ok {
			// If we have already seen this struct pointer,
			// encode a reference to it.
			e.writeByte(refCode)
			// Encode the relative position, because the buffer
			// will have data prepended to it.
			e.EncodeUint(uint64(len(e.buf) - u))
			e.buf[u] = refPtrCode // Backpatch the ptrCode to a refPtrCode.
			return false          // Caller should not encode the struct.
		}
		// Note that we have seen this pointer, and remember the position of the ptrCode.
		e.seen[ptr] = len(e.buf)
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
		i := d.i
		u := d.DecodeUint()
		return true, d.refMap[i-int(u)]
	case ptrCode:
		d.storeIndex = -1
		return true, nil
	case refPtrCode:
		// d.i was incremented by d.readByte, so the actual position of the code is one before.
		d.storeIndex = d.i - 1
		return true, nil
	default:
		d.badcode(b)
		panic("unreachable")
	}
}

//////////////// Struct Support

func (e *Encoder) StartStruct() {
	e.writeByte(startCode)
}

// StartStruct should be called before decoding a struct.
func (d *Decoder) StartStruct() {
	if b := d.readByte(); b != startCode {
		d.badcode(b)
	}
}

// StoreRef should be called by a struct decoder immediately after it allocates
// a struct pointer.
func (d *Decoder) StoreRef(p interface{}) {
	if d.storeIndex > 0 {
		if d.refMap == nil {
			d.refMap = map[int]interface{}{}
		}
		d.refMap[d.storeIndex] = p
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
	if d.opts.FailOnUnknownField {
		Failf("unknown field number %d for type %s", num, typeName)
	} else {
		d.skip()
	}
}

// skip reads past a value in the input.
func (d *Decoder) skip() {
	b := d.readByte()
	if b < endCode {
		// Small integers represent themselves in a single byte.
		return
	}
	if b >= bytes4Code && b < bytes0Code {
		d.readBytes(int(bytes0Code - b))
		return
	}
	switch b {
	case nilCode, bytes0Code:
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
	case ptrCode, refPtrCode:
		// One value follows.
		d.skip()
	case startCode:
		// Skip until we see endCode.
		for d.curByte() != endCode {
			d.skip()
		}
		d.readByte() // consume the endCode byte
	default:
		d.badcode(b)
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
		Failf("unregistered type %q", t)
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
		Failf("DecodeAny: bad list length %d", n)
	}
	num := d.DecodeUint()
	if num >= uint64(len(d.typeCodecs)) {
		Failf("type number %d out of range", num)
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
}

// decodeInitial decodes metadata that appears at the start of the
// encoded byte slice.
func (d *Decoder) decodeInitial() {
	// Decode the list of type names. The number of a type is its position in
	// the list.
	n := d.StartList()
	d.typeCodecs = make([]typeCodec, n)
	for num := 0; num < n; num++ {
		name := d.DecodeString()
		tc := typeCodecsByName[name]
		if tc == nil {
			Failf("unregistered type: %s", name)
		}
		d.typeCodecs[num] = tc
	}
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

// Failf calls fmt.Errorf with the given arguments, then Fail.
// It never returns.
func Failf(format string, args ...interface{}) {
	Fail(fmt.Errorf(format, args...))
}

// Fail aborts the current encoding or decoding with the given error.
// It never returns.
func Fail(err error) {
	panic(codecError{err})
}

func (d *Decoder) badcode(c byte) {
	//Failf("bad code %d at %d", c, d.i-1)
	panic(fmt.Sprintf("bad code %d", c))
}

// codecError wraps errors from Fail so a recover
// can distinguish them.
type codecError struct {
	err error
}

func (c codecError) String() string { return c.err.Error() }
