// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codecapi

import (
	"fmt"
	"reflect"
)

type typeCodec interface {
	Init()
	Encode(*Encoder, interface{})
	Decode(*Decoder) interface{}
}

var (
	typeCodecsByName = map[string]typeCodec{}
	typeCodecsByType = map[reflect.Type]typeCodec{}
)

// typeName returns the full, qualified name for a type.
func typeName(t reflect.Type) string {
	if t.PkgPath() == "" {
		return t.String()
	}
	return t.PkgPath() + "." + t.Name()
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

type IntCodec struct{}

func (IntCodec) Init()                            {}
func (IntCodec) Encode(e *Encoder, x interface{}) { e.EncodeInt(int64(x.(int))) }
func (IntCodec) Decode(d *Decoder) interface{}    { return int(d.DecodeInt()) }

type int64Codec struct{}

func (int64Codec) Init()                            {}
func (int64Codec) Encode(e *Encoder, x interface{}) { e.EncodeInt(x.(int64)) }
func (int64Codec) Decode(d *Decoder) interface{}    { return d.DecodeInt() }

type int32Codec struct{}

func (int32Codec) Init()                            {}
func (int32Codec) Encode(e *Encoder, x interface{}) { e.EncodeInt(int64(x.(int32))) }
func (int32Codec) Decode(d *Decoder) interface{}    { return int32(d.DecodeInt()) }

type int16Codec struct{}

func (int16Codec) Init()                            {}
func (int16Codec) Encode(e *Encoder, x interface{}) { e.EncodeInt(int64(x.(int16))) }
func (int16Codec) Decode(d *Decoder) interface{}    { return int16(d.DecodeInt()) }

type int8Codec struct{}

func (int8Codec) Init()                            {}
func (int8Codec) Encode(e *Encoder, x interface{}) { e.writeByte(byte(x.(int8))) }
func (int8Codec) Decode(d *Decoder) interface{}    { return int8(d.readByte()) }

type stringCodec struct{}

func (stringCodec) Init()                            {}
func (stringCodec) Encode(e *Encoder, x interface{}) { e.EncodeString(x.(string)) }
func (stringCodec) Decode(d *Decoder) interface{}    { return d.DecodeString() }

type float64Codec struct{}

func (float64Codec) Init()                            {}
func (float64Codec) Encode(e *Encoder, x interface{}) { e.EncodeFloat(x.(float64)) }
func (float64Codec) Decode(d *Decoder) interface{}    { return d.DecodeFloat() }

type uint64Codec struct{}

func (uint64Codec) Init()                            {}
func (uint64Codec) Encode(e *Encoder, x interface{}) { e.EncodeUint(x.(uint64)) }
func (uint64Codec) Decode(d *Decoder) interface{}    { return d.DecodeUint() }

type boolCodec struct{}

func (boolCodec) Init()                            {}
func (boolCodec) Encode(e *Encoder, x interface{}) { e.EncodeBool(x.(bool)) }
func (boolCodec) Decode(d *Decoder) interface{}    { return d.DecodeBool() }

type bytesCodec struct{}

func (bytesCodec) Init()                            {}
func (bytesCodec) Encode(e *Encoder, x interface{}) { e.EncodeBytes(x.([]byte)) }
func (bytesCodec) Decode(d *Decoder) interface{}    { return d.DecodeBytes() }

func init() {
	Register(int64(0), int64Codec{})
	Register(int32(0), int32Codec{})
	Register(int16(0), int16Codec{})
	Register(int8(0), int8Codec{})
	Register("", stringCodec{})
	Register(int(0), IntCodec{})
	Register(float64(0), float64Codec{})
	Register(uint64(0), uint64Codec{})
	Register(false, boolCodec{})
	Register([]byte(nil), bytesCodec{})
}

var BuiltinTypes []reflect.Type

func init() {
	for t := range typeCodecsByType {
		BuiltinTypes = append(BuiltinTypes, t)
	}
}
