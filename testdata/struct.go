// Code generated by the codec package. DO NOT EDIT.

package somepkg

import (
	"github.com/jba/codec"
	"github.com/jba/codec/codecapi"
)

// Fields of codec_genStruct: S B I I8 I16 I32 I64 F32 F64 U8 U16 U32 U64

type ptr_codec_genStruct_codec struct{}

func (ptr_codec_genStruct_codec) Init() {}

func (c ptr_codec_genStruct_codec) Encode(e *codecapi.Encoder, x interface{}) {
	c.encode(e, x.(*codec.genStruct))
}

func (c ptr_codec_genStruct_codec) encode(e *codecapi.Encoder, x *codec.genStruct) {
	if !e.StartPtr(x == nil, x) {
		return
	}
	(codec_genStruct_codec{}).encode(e, x)
}

func (c ptr_codec_genStruct_codec) Decode(d *codecapi.Decoder) interface{} {
	var x *codec.genStruct
	c.decode(d, &x)
	return x
}

func (c ptr_codec_genStruct_codec) decode(d *codecapi.Decoder, p **codec.genStruct) {
	proceed, ref := d.StartPtr()
	if !proceed {
		return
	}
	if ref != nil {
		*p = ref.(*codec.genStruct)
		return
	}
	var x codec.genStruct
	d.StoreRef(&x)
	(codec_genStruct_codec{}).decode(d, &x)
	*p = &x
}

type codec_genStruct_codec struct{}

func (codec_genStruct_codec) Init() {}

func (c codec_genStruct_codec) Encode(e *codecapi.Encoder, x interface{}) {
	s := x.(codec.genStruct)
	c.encode(e, &s)
}

func (c codec_genStruct_codec) encode(e *codecapi.Encoder, x *codec.genStruct) {
	e.StartStruct()
	if x.S != "" {
		e.EncodeUint(0)
		e.EncodeString(x.S)
	}
	if x.B != false {
		e.EncodeUint(1)
		e.EncodeBool(x.B)
	}
	if x.I != 0 {
		e.EncodeUint(2)
		e.EncodeInt(int64(x.I))
	}
	if x.I8 != 0 {
		e.EncodeUint(3)
		e.EncodeByte(uint8(x.I8))
	}
	if x.I16 != 0 {
		e.EncodeUint(4)
		e.EncodeInt(int64(x.I16))
	}
	if x.I32 != 0 {
		e.EncodeUint(5)
		e.EncodeInt(int64(x.I32))
	}
	if x.I64 != 0 {
		e.EncodeUint(6)
		e.EncodeInt(x.I64)
	}
	if x.F32 != 0 {
		e.EncodeUint(7)
		e.EncodeFloat(float64(x.F32))
	}
	if x.F64 != 0 {
		e.EncodeUint(8)
		e.EncodeFloat(x.F64)
	}
	if x.U8 != 0 {
		e.EncodeUint(9)
		e.EncodeByte(x.U8)
	}
	if x.U16 != 0 {
		e.EncodeUint(10)
		e.EncodeUint(uint64(x.U16))
	}
	if x.U32 != 0 {
		e.EncodeUint(11)
		e.EncodeUint(uint64(x.U32))
	}
	if x.U64 != 0 {
		e.EncodeUint(12)
		e.EncodeUint(x.U64)
	}
	e.EndStruct()
}

func (c codec_genStruct_codec) Decode(d *codecapi.Decoder) interface{} {
	var x codec.genStruct
	c.decode(d, &x)
	return x
}

func (c codec_genStruct_codec) decode(d *codecapi.Decoder, x *codec.genStruct) {
	d.StartStruct()
	for {
		n := d.NextStructField()
		if n < 0 {
			break
		}
		switch n {
		case 0:
			x.S = d.DecodeString()
		case 1:
			x.B = d.DecodeBool()
		case 2:
			x.I = int(d.DecodeInt())
		case 3:
			x.I8 = int8(d.DecodeByte())
		case 4:
			x.I16 = int16(d.DecodeInt())
		case 5:
			x.I32 = int32(d.DecodeInt())
		case 6:
			x.I64 = d.DecodeInt()
		case 7:
			x.F32 = float32(d.DecodeFloat())
		case 8:
			x.F64 = d.DecodeFloat()
		case 9:
			x.U8 = d.DecodeByte()
		case 10:
			x.U16 = uint16(d.DecodeUint())
		case 11:
			x.U32 = uint32(d.DecodeUint())
		case 12:
			x.U64 = d.DecodeUint()
		default:
			d.UnknownField("codec.genStruct", n)
		}
	}
}

func init() {
	codecapi.Register(codec.genStruct{}, codec_genStruct_codec{})
	codecapi.Register(&codec.genStruct{}, ptr_codec_genStruct_codec{})
}
