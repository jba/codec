// Code generated by the codec package. DO NOT EDIT.

package codec

import (
	"reflect"

	"github.com/jba/codec/codecapi"
)

//// [][]int

var slice_slice_int_type = reflect.TypeOf((*[][]int)(nil)).Elem()

type slice_slice_int_codec struct {
	codecapi.NonStruct

	slice_int_codec *slice_int_codec
}

func (c *slice_slice_int_codec) TypesUsed() []reflect.Type {
	return []reflect.Type{slice_int_type}
}

func (c *slice_slice_int_codec) SetCodecs(tcs []codecapi.TypeCodec) {
	c.slice_int_codec = tcs[0].(*slice_int_codec)
}

func (c *slice_slice_int_codec) Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.([][]int)) }

func (c *slice_slice_int_codec) encode(e *codecapi.Encoder, s [][]int) {
	if s == nil {
		e.EncodeNil()
		return
	}
	e.StartList(len(s))
	for _, x := range s {
		c.slice_int_codec.encode(e, x)
	}
}

func (c *slice_slice_int_codec) Decode(d *codecapi.Decoder) interface{} {
	var x [][]int
	c.decode(d, &x)
	return x
}

func (c *slice_slice_int_codec) decode(d *codecapi.Decoder, p *[][]int) {
	n := d.StartList()
	if n < 0 {
		return
	}
	s := make([][]int, n)
	for i := 0; i < n; i++ {
		c.slice_int_codec.decode(d, &s[i])
	}
	*p = s
}

func init() {
	codecapi.Register(slice_slice_int_type, func() codecapi.TypeCodec { return &slice_slice_int_codec{} })
}

//// []int

var slice_int_type = reflect.TypeOf((*[]int)(nil)).Elem()

type slice_int_codec struct {
	codecapi.NonStruct
}

func (c *slice_int_codec) TypesUsed() []reflect.Type      { return nil }
func (c *slice_int_codec) SetCodecs([]codecapi.TypeCodec) {}

func (c *slice_int_codec) Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.([]int)) }

func (c *slice_int_codec) encode(e *codecapi.Encoder, s []int) {
	if s == nil {
		e.EncodeNil()
		return
	}
	e.StartList(len(s))
	for _, x := range s {
		e.EncodeInt(int64(x))
	}
}

func (c *slice_int_codec) Decode(d *codecapi.Decoder) interface{} {
	var x []int
	c.decode(d, &x)
	return x
}

func (c *slice_int_codec) decode(d *codecapi.Decoder, p *[]int) {
	n := d.StartList()
	if n < 0 {
		return
	}
	s := make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = int(d.DecodeInt())
	}
	*p = s
}

func init() {
	codecapi.Register(slice_int_type, func() codecapi.TypeCodec { return &slice_int_codec{} })
}
