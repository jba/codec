// Code generated by the codec package. DO NOT EDIT.

package codec

import (
	"reflect"
	"time"

	"github.com/jba/codec/codecapi"
)

var time_Time_type = reflect.TypeOf((*time.Time)(nil)).Elem()

type time_Time_codec struct{}

func (c *time_Time_codec) Fields() []string { return nil }

func (c *time_Time_codec) Init(map[reflect.Type]codecapi.TypeCodec) {}

func (c *time_Time_codec) TypesUsed() []reflect.Type { return nil }

func (c *time_Time_codec) Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.(time.Time)) }

func (c *time_Time_codec) encode(e *codecapi.Encoder, m time.Time) {
	data, err := m.MarshalBinary()
	if err != nil {
		codecapi.Fail(err)
	}
	e.EncodeBytes(data)
}

func (c *time_Time_codec) Decode(d *codecapi.Decoder) interface{} {
	var x time.Time
	c.decode(d, &x)
	return x
}

func (c *time_Time_codec) decode(d *codecapi.Decoder, p *time.Time) {
	data := d.DecodeBytes()
	if err := p.UnmarshalBinary(data); err != nil {
		codecapi.Fail(err)
	}
}

func init() {
	codecapi.Register(*new(time.Time), func() codecapi.TypeCodec { return &time_Time_codec{} })
}
