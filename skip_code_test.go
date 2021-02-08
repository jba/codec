// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was originally generated automatically, by adding a field of type
// Skip to generatedTestTypes in codec_test.go and then running `go generate`.
// Then the decode method was hand-edited to leave only the default case in
// the switch.

package codec

import (
	"reflect"

	"github.com/jba/codec/codecapi"
)

// Fields of Skip: U S1 S2 L

type ptr_Skip_codec struct{}

func (ptr_Skip_codec) Init(map[reflect.Type]codecapi.TypeCodec) {}

func (c ptr_Skip_codec) Fields() []string { return nil }

func (ptr_Skip_codec) TypesUsed() []reflect.Type { return nil }

func (c ptr_Skip_codec) Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.(*Skip)) }

func (c ptr_Skip_codec) encode(e *codecapi.Encoder, x *Skip) {
	if !e.StartPtr(x == nil, x) {
		return
	}
	(Skip_codec{}).encode(e, x)
}

func (c ptr_Skip_codec) Decode(d *codecapi.Decoder) interface{} {
	var x *Skip
	c.decode(d, &x)
	return x
}

func (c ptr_Skip_codec) decode(d *codecapi.Decoder, p **Skip) {
	proceed, ref := d.StartPtr()
	if !proceed {
		return
	}
	if ref != nil {
		*p = ref.(*Skip)
		return
	}
	var x Skip
	d.StoreRef(&x)
	(Skip_codec{}).decode(d, &x)
	*p = &x
}

type Skip_codec struct{}

func (Skip_codec) Init(map[reflect.Type]codecapi.TypeCodec) {}

func (Skip_codec) Fields() []string { return []string{"U", "S1", "S2", "L"} }

func (Skip_codec) TypesUsed() []reflect.Type { return nil }

func (c Skip_codec) Encode(e *codecapi.Encoder, x interface{}) {
	s := x.(Skip)
	c.encode(e, &s)
}

func (c Skip_codec) encode(e *codecapi.Encoder, x *Skip) {
	e.StartStruct()
	if x.U != 0 {
		e.EncodeUint(0)
		e.EncodeUint(x.U)
	}
	if x.S1 != "" {
		e.EncodeUint(1)
		e.EncodeString(x.S1)
	}
	if x.S2 != "" {
		e.EncodeUint(2)
		e.EncodeString(x.S2)
	}
	if x.L != nil {
		e.EncodeUint(3)
		(slice_ptr_Skip_codec{}).encode(e, x.L)
	}
	e.EndStruct()
}

func (c Skip_codec) Decode(d *codecapi.Decoder) interface{} {
	var x Skip
	c.decode(d, &x)
	return x
}

func (c Skip_codec) decode(d *codecapi.Decoder, x *Skip) {
	d.StartStruct()
	for {
		n := d.NextStructField()
		if n < 0 {
			break
		}
		switch n {
		default:
			d.UnknownField("Skip", n)
		}
	}
}

func init() {
	codecapi.Register(Skip{}, func() codecapi.TypeCodec {
		return Skip_codec{}
	})
	codecapi.Register(&Skip{}, func() codecapi.TypeCodec {
		return ptr_Skip_codec{}
	})
}

type slice_ptr_Skip_codec struct{}

func (slice_ptr_Skip_codec) Fields() []string { return nil }

func (slice_ptr_Skip_codec) TypesUsed() []reflect.Type { return nil }

func (slice_ptr_Skip_codec) Init(map[reflect.Type]codecapi.TypeCodec) {}

func (c slice_ptr_Skip_codec) Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.([]*Skip)) }

func (c slice_ptr_Skip_codec) encode(e *codecapi.Encoder, s []*Skip) {
	if s == nil {
		e.EncodeNil()
		return
	}
	e.StartList(len(s))
	for _, x := range s {
		(ptr_Skip_codec{}).encode(e, x)
	}
}

func (c slice_ptr_Skip_codec) Decode(d *codecapi.Decoder) interface{} {
	var x []*Skip
	c.decode(d, &x)
	return x
}

func (c slice_ptr_Skip_codec) decode(d *codecapi.Decoder, p *[]*Skip) {
	n := d.StartList()
	if n < 0 {
		return
	}
	s := make([]*Skip, n)
	for i := 0; i < n; i++ {
		(ptr_Skip_codec{}).decode(d, &s[i])
	}
	*p = s
}

func init() {
	codecapi.Register([]*Skip(nil), func() codecapi.TypeCodec {
		return slice_ptr_Skip_codec{}
	})
}
