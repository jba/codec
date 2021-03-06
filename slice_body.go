// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by embed.sh. DO NOT EDIT.

package codec

const sliceBody = `
« $typeID := typeID .Type »
« $typeName := print $typeID "_codec" »
« $goName := goName .Type »
« $elTypeID := typeID .Type.Elem »

var «$typeID»_type = reflect.TypeOf((*«$goName»)(nil)).Elem()

type «$typeName» struct {
	codecapi.NonStruct
	«if .ElField»
		«$elTypeID»_codec *«$elTypeID»_codec
	«end -»
}

«if .ElField»
	func (c *«$typeName») TypesUsed() []reflect.Type {
		return []reflect.Type{«$elTypeID»_type}
	}

	func (c *«$typeName») SetCodecs(tcs []codecapi.TypeCodec) {
		c.«$elTypeID»_codec = tcs[0].(*«$elTypeID»_codec)
	}
«else»
	func (c *«$typeName») TypesUsed() []reflect.Type { return nil }
	func (c *«$typeName») SetCodecs([]codecapi.TypeCodec) {}
«end»

func (c *«$typeName») Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.(«$goName»)) }

func (c *«$typeName») encode(e *codecapi.Encoder, s «$goName») {
	if s == nil {
		e.EncodeNil()
		return
	}
	e.StartList(len(s))
	for _, x := range s {
		«encodeStmt .Type.Elem "x"»
	}
}

func (c *«$typeName») Decode(d *codecapi.Decoder) interface{} {
	var x «$goName»
	c.decode(d, &x)
	return x
}

func (c *«$typeName») decode(d *codecapi.Decoder, p *«$goName») {
	n := d.StartList()
	if n < 0 { return }
	s := make([]«goName .Type.Elem», n)
	for i := 0; i < n; i++ {
		«decodeStmt .Type.Elem "s[i]"»
	}
	*p = s
}

func init() {
  codecapi.Register(«$typeID»_type, func() codecapi.TypeCodec { return &«$typeName»{} })
}
`
