// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codec

const marshalBody = `
«/*» Template body for a type that implements encoding.BinaryMarshaler or encoding.TextMarshaler. «*/»

« $typeID := typeID .Type »
« $typeName := print $typeID "_codec" »
« $goName := goName .Type »

var «$typeID»_type = reflect.TypeOf((*«$goName»)(nil)).Elem()

type «$typeName» struct{
	codecapi.NonStruct
}

func (c *«$typeName») TypesUsed() []reflect.Type { return nil }

func (c *«$typeName») SetCodecs([]codecapi.TypeCodec) {}

func (c *«$typeName») Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.(«$goName»)) }

func (c *«$typeName») encode(e *codecapi.Encoder, m «$goName») {
	data, err := m.Marshal«.Kind»()
	if err != nil {
		codecapi.Fail(err)
	}
	e.EncodeBytes(data)
}

func (c *«$typeName») Decode(d *codecapi.Decoder) interface{} {
	var x «$goName»
	c.decode(d, &x)
	return x
}

func (c *«$typeName») decode(d *codecapi.Decoder, p *«$goName») {
	data := d.DecodeBytes()
	if err := p.Unmarshal«.Kind»(data); err != nil {
		codecapi.Fail(err)
	}
}

func init() {
	codecapi.Register(*new(«$goName»), func() codecapi.TypeCodec { return &«$typeName»{} })
}

`
