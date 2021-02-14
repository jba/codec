// Code generated by the codec package. DO NOT EDIT.

package codec

import (
	"reflect"

	"github.com/jba/codec/codecapi"
)

//// codec.definedSlice

var definedSlice_type = reflect.TypeOf((*definedSlice)(nil)).Elem()

type definedSlice_codec struct {
	codecapi.NonStruct
}

func (c *definedSlice_codec) TypesUsed() []reflect.Type      { return nil }
func (c *definedSlice_codec) SetCodecs([]codecapi.TypeCodec) {}

func (c *definedSlice_codec) Encode(e *codecapi.Encoder, x interface{}) {
	c.encode(e, x.(definedSlice))
}

func (c *definedSlice_codec) encode(e *codecapi.Encoder, s definedSlice) {
	if s == nil {
		e.EncodeNil()
		return
	}
	e.StartList(len(s))
	for _, x := range s {
		e.EncodeInt(int64(x))
	}
}

func (c *definedSlice_codec) Decode(d *codecapi.Decoder) interface{} {
	var x definedSlice
	c.decode(d, &x)
	return x
}

func (c *definedSlice_codec) decode(d *codecapi.Decoder, p *definedSlice) {
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
	codecapi.Register(definedSlice(nil), func() codecapi.TypeCodec { return &definedSlice_codec{} })
}
