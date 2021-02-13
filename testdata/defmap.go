// Code generated by the codec package. DO NOT EDIT.

package codec

import (
	"reflect"

	"github.com/jba/codec/codecapi"
)

var definedMap_type = reflect.TypeOf((*definedMap)(nil)).Elem()

type definedMap_codec struct {
}

func (c *definedMap_codec) Init(tcs map[reflect.Type]codecapi.TypeCodec, _ []int) {
}

func (c *definedMap_codec) Fields() []string { return nil }

func (c *definedMap_codec) TypesUsed() []reflect.Type {
	var types []reflect.Type

	return types
}

func (c *definedMap_codec) Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.(definedMap)) }

func (c *definedMap_codec) encode(e *codecapi.Encoder, m definedMap) {
	if m == nil {
		e.EncodeNil()
		return
	}
	e.StartList(2 * len(m))
	for k, v := range m {
		e.EncodeString(k)
		e.EncodeBool(v)
	}
}

func (c *definedMap_codec) Decode(d *codecapi.Decoder) interface{} {
	var x definedMap
	c.decode(d, &x)
	return x
}

func (c *definedMap_codec) decode(d *codecapi.Decoder, p *definedMap) {
	n2 := d.StartList()
	if n2 < 0 {
		return
	}
	n := n2 / 2
	m := make(definedMap, n)
	var k string
	var v bool
	for i := 0; i < n; i++ {
		k = d.DecodeString()
		v = d.DecodeBool()
		m[k] = v
	}
	*p = m
}

func init() {
	codecapi.Register(definedMap(nil), func() codecapi.TypeCodec { return &definedMap_codec{} })
}
