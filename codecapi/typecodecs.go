// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codecapi

import (
	"fmt"
	"reflect"
	"strings"
)

// A TypeCodec handles encoding and decoding of a particular type.
type TypeCodec interface {
	Fields() []string
	//Type() reflect.Type
	Init(typeCodecs map[reflect.Type]TypeCodec)
	Encode(*Encoder, interface{})
	Decode(*Decoder) interface{}
}

var (
	typeCodecBuildersByName = map[string]func() TypeCodec{}
	typeCodecBuildersByType = map[reflect.Type]func() TypeCodec{}

	nameToType = map[string]reflect.Type{}
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
func Register(x interface{}, tcb func() TypeCodec) {
	t := reflect.TypeOf(x)
	tn := typeName(t)
	if _, ok := typeCodecBuildersByName[tn]; ok {
		panic(fmt.Sprintf("codec.Register: duplicate type %s (typeName=%q)", t, tn))
	}
	typeCodecBuildersByName[tn] = tcb
	typeCodecBuildersByType[t] = tcb
	nameToType[tn] = t
}

type prim struct{}

func (prim) Fields() []string                { return nil }
func (prim) Init(map[reflect.Type]TypeCodec) {}

type boolCodec struct{ prim }

func (boolCodec) Encode(e *Encoder, x interface{}) { e.EncodeBool(x.(bool)) }
func (boolCodec) Decode(d *Decoder) interface{}    { return d.DecodeBool() }

type bytesCodec struct{ prim }

func (bytesCodec) Encode(e *Encoder, x interface{}) { e.EncodeBytes(x.([]byte)) }
func (bytesCodec) Decode(d *Decoder) interface{}    { return d.DecodeBytes() }

type stringCodec struct{ prim }

func (stringCodec) Encode(e *Encoder, x interface{}) { e.EncodeString(x.(string)) }
func (stringCodec) Decode(d *Decoder) interface{}    { return d.DecodeString() }

type intCodec struct{ prim }

func (intCodec) Encode(e *Encoder, x interface{}) { e.EncodeInt(int64(x.(int))) }
func (intCodec) Decode(d *Decoder) interface{}    { return int(d.DecodeInt()) }

type int8Codec struct{ prim }

func (int8Codec) Encode(e *Encoder, x interface{}) { e.writeByte(byte(x.(int8))) }
func (int8Codec) Decode(d *Decoder) interface{}    { return int8(d.readByte()) }

type int16Codec struct{ prim }

func (int16Codec) Encode(e *Encoder, x interface{}) { e.EncodeInt(int64(x.(int16))) }
func (int16Codec) Decode(d *Decoder) interface{}    { return int16(d.DecodeInt()) }

type int32Codec struct{ prim }

func (int32Codec) Encode(e *Encoder, x interface{}) { e.EncodeInt(int64(x.(int32))) }
func (int32Codec) Decode(d *Decoder) interface{}    { return int32(d.DecodeInt()) }

type int64Codec struct{ prim }

func (int64Codec) Encode(e *Encoder, x interface{}) { e.EncodeInt(x.(int64)) }
func (int64Codec) Decode(d *Decoder) interface{}    { return d.DecodeInt() }

type float32Codec struct{ prim }

func (float32Codec) Encode(e *Encoder, x interface{}) { e.EncodeFloat(float64(x.(float32))) }
func (float32Codec) Decode(d *Decoder) interface{}    { return float32(d.DecodeFloat()) }

type float64Codec struct{ prim }

func (float64Codec) Encode(e *Encoder, x interface{}) { e.EncodeFloat(x.(float64)) }
func (float64Codec) Decode(d *Decoder) interface{}    { return d.DecodeFloat() }

type uintCodec struct{ prim }

func (uintCodec) Encode(e *Encoder, x interface{}) { e.EncodeUint(uint64(x.(uint))) }
func (uintCodec) Decode(d *Decoder) interface{}    { return uint(d.DecodeUint()) }

type uintptrCodec struct{ prim }

func (uintptrCodec) Encode(e *Encoder, x interface{}) { e.EncodeUint(uint64(x.(uintptr))) }
func (uintptrCodec) Decode(d *Decoder) interface{}    { return uintptr(d.DecodeUint()) }

type uint8Codec struct{ prim }

func (uint8Codec) Encode(e *Encoder, x interface{}) { e.writeByte(byte(x.(uint8))) }
func (uint8Codec) Decode(d *Decoder) interface{}    { return d.readByte() }

type uint16Codec struct{ prim }

func (uint16Codec) Encode(e *Encoder, x interface{}) { e.EncodeUint(uint64(x.(uint16))) }
func (uint16Codec) Decode(d *Decoder) interface{}    { return uint16(d.DecodeUint()) }

type uint32Codec struct{ prim }

func (uint32Codec) Encode(e *Encoder, x interface{}) { e.EncodeUint(uint64(x.(uint32))) }
func (uint32Codec) Decode(d *Decoder) interface{}    { return uint32(d.DecodeUint()) }

type uint64Codec struct{ prim }

func (uint64Codec) Encode(e *Encoder, x interface{}) { e.EncodeUint(x.(uint64)) }
func (uint64Codec) Decode(d *Decoder) interface{}    { return d.DecodeUint() }

type complex64Codec struct{ prim }

func (complex64Codec) Encode(e *Encoder, x interface{}) { e.EncodeComplex(complex128(x.(complex64))) }
func (complex64Codec) Decode(d *Decoder) interface{}    { return complex64(d.DecodeComplex()) }

type complex128Codec struct{ prim }

func (complex128Codec) Encode(e *Encoder, x interface{}) { e.EncodeComplex(x.(complex128)) }
func (complex128Codec) Decode(d *Decoder) interface{}    { return d.DecodeComplex() }

func init() {
	Register(false, func() TypeCodec { return boolCodec{} })
	Register("", func() TypeCodec { return stringCodec{} })
	Register([]byte(nil), func() TypeCodec { return bytesCodec{} })
	Register(int(0), func() TypeCodec { return intCodec{} })
	Register(int8(0), func() TypeCodec { return int8Codec{} })
	Register(int16(0), func() TypeCodec { return int16Codec{} })
	Register(int32(0), func() TypeCodec { return int32Codec{} })
	Register(int64(0), func() TypeCodec { return int64Codec{} })
	Register(float32(0), func() TypeCodec { return float32Codec{} })
	Register(float64(0), func() TypeCodec { return float64Codec{} })
	Register(uint(0), func() TypeCodec { return uintCodec{} })
	Register(uintptr(0), func() TypeCodec { return uintptrCodec{} })
	Register(uint8(0), func() TypeCodec { return uint8Codec{} })
	Register(uint16(0), func() TypeCodec { return uint16Codec{} })
	Register(uint32(0), func() TypeCodec { return uint32Codec{} })
	Register(uint64(0), func() TypeCodec { return uint64Codec{} })
	Register(complex64(0), func() TypeCodec { return complex64Codec{} })
	Register(complex128(0), func() TypeCodec { return complex128Codec{} })
}

var BuiltinTypes []reflect.Type

func init() {
	for t := range typeCodecBuildersByType {
		BuiltinTypes = append(BuiltinTypes, t)
	}
}
