// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codecapi

import (
	"fmt"
	"reflect"
	"strings"
)

// A typeCodec handles encoding and decoding of a particular type.
type typeCodec interface {
	Encode(*Encoder, interface{})
	Decode(*Decoder) interface{}
}

var (
	typeCodecsByName = map[string]typeCodec{}
	typeCodecsByType = map[reflect.Type]typeCodec{}
)

// typeName does its best to construct a unique name for a reflect.Type.
// We can't do anything about types defined inside functions,
// but we can make sure that types with the same name and package name
// are distinguished.
func typeName(t reflect.Type) string {
	if n := t.Name(); n != "" {
		pp := t.PkgPath()
		if pp == "" {
			return n
		}
		return pp + "." + n
	}
	switch t.Kind() {
	case reflect.Slice:
		return "[]" + typeName(t.Elem())
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", t.Len(), typeName(t.Elem()))
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", typeName(t.Key()), typeName(t.Elem()))
	case reflect.Ptr:
		return "*" + typeName(t.Elem())
	case reflect.Struct:
		fields := make([]string, t.NumField())
		for i := 0; i < len(fields); i++ {
			f := t.Field(i)
			tn := typeName(f.Type)
			if f.Anonymous {
				fields[i] = tn
			} else {
				fields[i] = f.Name + " " + tn
			}
		}
		return fmt.Sprintf("struct { %s }", strings.Join(fields, "; "))
	default:
		return t.String()
	}
}

// Register records the type of x for use by Encoders and Decoders.
// All types subject to encoding must be registered, even
// builtin types.
func Register(x interface{}, tc typeCodec) {
	t := reflect.TypeOf(x)
	tn := typeName(t)
	if _, ok := typeCodecsByName[tn]; ok {
		panic(fmt.Sprintf("codec.Register: duplicate type %s (typeName=%q)", t, tn))
	}
	typeCodecsByName[tn] = tc
	typeCodecsByType[t] = tc
}

type boolCodec struct{}

func (boolCodec) Encode(e *Encoder, x interface{}) { e.EncodeBool(x.(bool)) }
func (boolCodec) Decode(d *Decoder) interface{}    { return d.DecodeBool() }

type bytesCodec struct{}

func (bytesCodec) Encode(e *Encoder, x interface{}) { e.EncodeBytes(x.([]byte)) }
func (bytesCodec) Decode(d *Decoder) interface{}    { return d.DecodeBytes() }

type stringCodec struct{}

func (stringCodec) Encode(e *Encoder, x interface{}) { e.EncodeString(x.(string)) }
func (stringCodec) Decode(d *Decoder) interface{}    { return d.DecodeString() }

type intCodec struct{}

func (intCodec) Encode(e *Encoder, x interface{}) { e.EncodeInt(int64(x.(int))) }
func (intCodec) Decode(d *Decoder) interface{}    { return int(d.DecodeInt()) }

type int8Codec struct{}

func (int8Codec) Encode(e *Encoder, x interface{}) { e.writeByte(byte(x.(int8))) }
func (int8Codec) Decode(d *Decoder) interface{}    { return int8(d.readByte()) }

type int16Codec struct{}

func (int16Codec) Encode(e *Encoder, x interface{}) { e.EncodeInt(int64(x.(int16))) }
func (int16Codec) Decode(d *Decoder) interface{}    { return int16(d.DecodeInt()) }

type int32Codec struct{}

func (int32Codec) Encode(e *Encoder, x interface{}) { e.EncodeInt(int64(x.(int32))) }
func (int32Codec) Decode(d *Decoder) interface{}    { return int32(d.DecodeInt()) }

type int64Codec struct{}

func (int64Codec) Encode(e *Encoder, x interface{}) { e.EncodeInt(x.(int64)) }
func (int64Codec) Decode(d *Decoder) interface{}    { return d.DecodeInt() }

type float32Codec struct{}

func (float32Codec) Encode(e *Encoder, x interface{}) { e.EncodeFloat(float64(x.(float32))) }
func (float32Codec) Decode(d *Decoder) interface{}    { return float32(d.DecodeFloat()) }

type float64Codec struct{}

func (float64Codec) Encode(e *Encoder, x interface{}) { e.EncodeFloat(x.(float64)) }
func (float64Codec) Decode(d *Decoder) interface{}    { return d.DecodeFloat() }

type uintCodec struct{}

func (uintCodec) Encode(e *Encoder, x interface{}) { e.EncodeUint(uint64(x.(uint))) }
func (uintCodec) Decode(d *Decoder) interface{}    { return uint(d.DecodeUint()) }

type uintptrCodec struct{}

func (uintptrCodec) Encode(e *Encoder, x interface{}) { e.EncodeUint(uint64(x.(uintptr))) }
func (uintptrCodec) Decode(d *Decoder) interface{}    { return uintptr(d.DecodeUint()) }

type uint8Codec struct{}

func (uint8Codec) Encode(e *Encoder, x interface{}) { e.writeByte(byte(x.(uint8))) }
func (uint8Codec) Decode(d *Decoder) interface{}    { return d.readByte() }

type uint16Codec struct{}

func (uint16Codec) Encode(e *Encoder, x interface{}) { e.EncodeUint(uint64(x.(uint16))) }
func (uint16Codec) Decode(d *Decoder) interface{}    { return uint16(d.DecodeUint()) }

type uint32Codec struct{}

func (uint32Codec) Encode(e *Encoder, x interface{}) { e.EncodeUint(uint64(x.(uint32))) }
func (uint32Codec) Decode(d *Decoder) interface{}    { return uint32(d.DecodeUint()) }

type uint64Codec struct{}

func (uint64Codec) Encode(e *Encoder, x interface{}) { e.EncodeUint(x.(uint64)) }
func (uint64Codec) Decode(d *Decoder) interface{}    { return d.DecodeUint() }

type complex64Codec struct{}

func (complex64Codec) Encode(e *Encoder, x interface{}) { e.EncodeComplex(complex128(x.(complex64))) }
func (complex64Codec) Decode(d *Decoder) interface{}    { return complex64(d.DecodeComplex()) }

type complex128Codec struct{}

func (complex128Codec) Encode(e *Encoder, x interface{}) { e.EncodeComplex(x.(complex128)) }
func (complex128Codec) Decode(d *Decoder) interface{}    { return d.DecodeComplex() }

func init() {
	Register(false, boolCodec{})
	Register("", stringCodec{})
	Register([]byte(nil), bytesCodec{})
	Register(int(0), intCodec{})
	Register(int8(0), int8Codec{})
	Register(int16(0), int16Codec{})
	Register(int32(0), int32Codec{})
	Register(int64(0), int64Codec{})
	Register(float32(0), float32Codec{})
	Register(float64(0), float64Codec{})
	Register(uint(0), uintCodec{})
	Register(uintptr(0), uintptrCodec{})
	Register(uint8(0), uint8Codec{})
	Register(uint16(0), uint16Codec{})
	Register(uint32(0), uint32Codec{})
	Register(uint64(0), uint64Codec{})
	Register(complex64(0), complex64Codec{})
	Register(complex128(0), complex128Codec{})
}

var BuiltinTypes []reflect.Type

func init() {
	for t := range typeCodecsByType {
		BuiltinTypes = append(BuiltinTypes, t)
	}
}
