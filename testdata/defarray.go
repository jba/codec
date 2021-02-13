// Code generated by the codec package. DO NOT EDIT.

package codec

import (
	"reflect"

	"github.com/jba/codec/codecapi"
)

var slice_int_type = reflect.TypeOf((*[]int)(nil)).Elem()

type slice_int_codec struct {
}

func (c *slice_int_codec) Init(tcs map[reflect.Type]codecapi.TypeCodec, _ []int) {

}

func (c *slice_int_codec) Fields() []string { return nil }

func (c *slice_int_codec) TypesUsed() []reflect.Type {
	return nil
}

func (c *slice_int_codec) CodecsUsed([]codecapi.TypeCodec) {}

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
	codecapi.Register([]int(nil), func() codecapi.TypeCodec { return &slice_int_codec{} })
}

var definedArray_type = reflect.TypeOf((*definedArray)(nil)).Elem()

type definedArray_codec struct {
}

func (c *definedArray_codec) Init(tcs map[reflect.Type]codecapi.TypeCodec, _ []int) {

}

func (c *definedArray_codec) Fields() []string { return nil }

func (c *definedArray_codec) TypesUsed() []reflect.Type {
	return nil
}

func (c *definedArray_codec) CodecsUsed([]codecapi.TypeCodec) {}

func (c *definedArray_codec) Encode(e *codecapi.Encoder, x interface{}) {
	a := x.(definedArray)
	c.encode(e, &a)
}

func (c *definedArray_codec) encode(e *codecapi.Encoder, s *definedArray) {
	(&slice_int_codec{}).encode(e, (*s)[:])
}

func (c *definedArray_codec) Decode(d *codecapi.Decoder) interface{} {
	var x definedArray
	c.decode(d, &x)
	return x
}

func (c *definedArray_codec) decode(d *codecapi.Decoder, p *definedArray) {
	n := d.StartList()
	if n < 0 {
		return
	}
	if n != 1 {
		codecapi.Failf("array size mismatch: got %d, want 1", n)
	}
	for i := 0; i < n; i++ {
		(*p)[i] = int(d.DecodeInt())
	}
}

func init() {
	codecapi.Register(definedArray{}, func() codecapi.TypeCodec { return &definedArray_codec{} })
}
