// Code generated by the codec package. DO NOT EDIT.

package codec

import (
	"github.com/jba/codec/codecapi"
	"net"
	"reflect"
)

type net_IP_codec struct{}

func (c *net_IP_codec) Fields() []string { return nil }

func (c *net_IP_codec) Init(map[reflect.Type]codecapi.TypeCodec) {}

func (c *net_IP_codec) Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.(net.IP)) }

func (c *net_IP_codec) encode(e *codecapi.Encoder, m net.IP) {
	data, err := m.MarshalText()
	if err != nil {
		codecapi.Fail(err)
	}
	e.EncodeBytes(data)
}

func (c *net_IP_codec) Decode(d *codecapi.Decoder) interface{} {
	var x net.IP
	c.decode(d, &x)
	return x
}

func (c *net_IP_codec) decode(d *codecapi.Decoder, p *net.IP) {
	data := d.DecodeBytes()
	if err := p.UnmarshalText(data); err != nil {
		codecapi.Fail(err)
	}
}

func init() { codecapi.Register(*new(net.IP), func() codecapi.TypeCodec { return &net_IP_codec{} }) }
