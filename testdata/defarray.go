// Code generated by the codec package. DO NOT EDIT.

package somepkg

import (
	"fmt"
	"github.com/jba/codec"
)

type codec_definedArray_codec struct{}

func (codec_definedArray_codec) Init() {}

func (c codec_definedArray_codec) Encode(e *codec.Encoder, x interface{}) {
	a := x.(codec.definedArray)
	c.encode(e, &a)
}

func (c codec_definedArray_codec) encode(e *codec.Encoder, s *codec.definedArray) {
	(slice_int_codec{}).encode(e, (*s)[:])
}

func (c codec_definedArray_codec) Decode(d *codec.Decoder) interface{} {
	var x codec.definedArray
	c.decode(d, &x)
	return x
}

func (c codec_definedArray_codec) decode(d *codec.Decoder, p *codec.definedArray) {
	n := d.StartList()
	if n < 0 {
		return
	}
	if n != 1 {
		d.Fail(fmt.Errorf("array size mismatch: got %d, want 1", n))
	}
	for i := 0; i < n; i++ {
		(*p)[i] = int(d.DecodeInt())
	}
}

func init() {
	codec.Register(codec.definedArray{}, codec_definedArray_codec{})
}

type slice_int_codec struct{}

func (slice_int_codec) Init() {}

func (c slice_int_codec) Encode(e *codec.Encoder, x interface{}) { c.encode(e, x.([]int)) }

func (c slice_int_codec) encode(e *codec.Encoder, s []int) {
	if s == nil {
		e.EncodeNil()
		return
	}
	e.StartList(len(s))
	for _, x := range s {
		e.EncodeInt(int64(x))
	}
}

func (c slice_int_codec) Decode(d *codec.Decoder) interface{} {
	var x []int
	c.decode(d, &x)
	return x
}

func (c slice_int_codec) decode(d *codec.Decoder, p *[]int) {
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
	codec.Register([]int(nil), slice_int_codec{})
}
