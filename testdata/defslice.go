// Code generated by the codec package. DO NOT EDIT.

package somepkg

import (
	"github.com/jba/codec"
)

type codec_definedSlice_codec struct{}

func (codec_definedSlice_codec) Init() {}

func (c codec_definedSlice_codec) Encode(e *codec.Encoder, x interface{}) {
	c.encode(e, x.(codec.definedSlice))
}

func (c codec_definedSlice_codec) encode(e *codec.Encoder, s codec.definedSlice) {
	if s == nil {
		e.EncodeNil()
		return
	}
	e.StartList(len(s))
	for _, x := range s {
		e.EncodeInt(int64(x))
	}
}

func (c codec_definedSlice_codec) Decode(d *codec.Decoder) interface{} {
	var x codec.definedSlice
	c.decode(d, &x)
	return x
}

func (c codec_definedSlice_codec) decode(d *codec.Decoder, p *codec.definedSlice) {
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
	codec.Register(codec.definedSlice(nil), codec_definedSlice_codec{})
}
