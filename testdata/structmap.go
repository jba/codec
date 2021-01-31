// Code generated by the codec package. DO NOT EDIT.

package codec

import (
	"github.com/jba/codec/codecapi"
	"reflect"
)

type map__1_int_structType_codec struct {
	array_1_int_codec *array_1_int_codec
	structType_codec  *structType_codec
}

func (c *map__1_int_structType_codec) Init(tcs map[reflect.Type]codecapi.TypeCodec) {
	c.array_1_int_codec = tcs[reflect.TypeOf((*[1]int)(nil)).Elem()].(*array_1_int_codec)
	c.structType_codec = tcs[reflect.TypeOf((*structType)(nil)).Elem()].(*structType_codec)
}

func (c *map__1_int_structType_codec) Fields() []string { return nil }

func (c *map__1_int_structType_codec) Encode(e *codecapi.Encoder, x interface{}) {
	c.encode(e, x.(map[[1]int]structType))
}

func (c *map__1_int_structType_codec) encode(e *codecapi.Encoder, m map[[1]int]structType) {
	if m == nil {
		e.EncodeNil()
		return
	}
	e.StartList(2 * len(m))
	for k, v := range m {
		(&array_1_int_codec{}).encode(e, &k)
		(&structType_codec{}).encode(e, &v)
	}
}

func (c *map__1_int_structType_codec) Decode(d *codecapi.Decoder) interface{} {
	var x map[[1]int]structType
	c.decode(d, &x)
	return x
}

func (c *map__1_int_structType_codec) decode(d *codecapi.Decoder, p *map[[1]int]structType) {
	n2 := d.StartList()
	if n2 < 0 {
		return
	}
	n := n2 / 2
	m := make(map[[1]int]structType, n)
	var k [1]int
	var v structType
	for i := 0; i < n; i++ {
		c.array_1_int_codec.decode(d, &k)
		c.structType_codec.decode(d, &v)
		m[k] = v
	}
	*p = m
}

func init() {
	codecapi.Register(map[[1]int]structType(nil), func() codecapi.TypeCodec { return &map__1_int_structType_codec{} })
}

type array_1_int_codec struct {
}

func (c *array_1_int_codec) Init(tcs map[reflect.Type]codecapi.TypeCodec) {

}

func (c *array_1_int_codec) Fields() []string { return nil }

func (c *array_1_int_codec) Encode(e *codecapi.Encoder, x interface{}) {
	a := x.([1]int)
	c.encode(e, &a)
}

func (c *array_1_int_codec) encode(e *codecapi.Encoder, s *[1]int) {
	(&slice_int_codec{}).encode(e, (*s)[:])
}

func (c *array_1_int_codec) Decode(d *codecapi.Decoder) interface{} {
	var x [1]int
	c.decode(d, &x)
	return x
}

func (c *array_1_int_codec) decode(d *codecapi.Decoder, p *[1]int) {
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
	codecapi.Register([1]int{}, func() codecapi.TypeCodec { return &array_1_int_codec{} })
}

// Fields of structType: N B unexported

type ptr_structType_codec struct {
	structType_codec *structType_codec
}

func (c *ptr_structType_codec) Init(tcs map[reflect.Type]codecapi.TypeCodec) {
	c.structType_codec = tcs[reflect.TypeOf((*structType)(nil)).Elem()].(*structType_codec)
}

func (c ptr_structType_codec) Fields() []string { return nil }

func (c ptr_structType_codec) Encode(e *codecapi.Encoder, x interface{}) {
	c.encode(e, x.(*structType))
}

func (c ptr_structType_codec) encode(e *codecapi.Encoder, x *structType) {
	if !e.StartPtr(x == nil, x) {
		return
	}
	(&structType_codec{}).encode(e, x)
}

func (c ptr_structType_codec) Decode(d *codecapi.Decoder) interface{} {
	var x *structType
	c.decode(d, &x)
	return x
}

func (c ptr_structType_codec) decode(d *codecapi.Decoder, p **structType) {
	proceed, ref := d.StartPtr()
	if !proceed {
		return
	}
	if ref != nil {
		*p = ref.(*structType)
		return
	}
	var x structType
	d.StoreRef(&x)
	c.structType_codec.decode(d, &x)
	*p = &x
}

type structType_codec struct {
	node_codec *node_codec
}

func (c *structType_codec) Init(tcs map[reflect.Type]codecapi.TypeCodec) {
	c.node_codec = tcs[reflect.TypeOf((*node)(nil)).Elem()].(*node_codec)
}

func (c *structType_codec) Fields() []string {
	return []string{"N", "B", "unexported"}
}

func (c *structType_codec) Encode(e *codecapi.Encoder, x interface{}) {
	s := x.(structType)
	c.encode(e, &s)
}

func (c *structType_codec) encode(e *codecapi.Encoder, x *structType) {
	e.StartStruct()

	e.EncodeUint(0)
	(&node_codec{}).encode(e, &x.N)
	if x.B != 0 {
		e.EncodeUint(1)
		e.EncodeByte(x.B)
	}
	if x.unexported != 0 {
		e.EncodeUint(2)
		e.EncodeInt(int64(x.unexported))
	}
	e.EndStruct()
}

func (c *structType_codec) Decode(d *codecapi.Decoder) interface{} {
	var x structType
	c.decode(d, &x)
	return x
}

func (c *structType_codec) decode(d *codecapi.Decoder, x *structType) {
	d.StartStruct()
	for {
		n := d.NextStructField()
		if n < 0 {
			break
		}
		switch n {
		case 0:
			c.node_codec.decode(d, &x.N)
		case 1:
			x.B = d.DecodeByte()
		case 2:
			x.unexported = int(d.DecodeInt())
		default:
			d.UnknownField("structType", n)
		}
	}
}

func init() {
	codecapi.Register(structType{}, func() codecapi.TypeCodec { return &structType_codec{} })
	codecapi.Register(&structType{}, func() codecapi.TypeCodec { return &ptr_structType_codec{} })
}

type slice_int_codec struct {
}

func (c *slice_int_codec) Init(tcs map[reflect.Type]codecapi.TypeCodec) {

}

func (c *slice_int_codec) Fields() []string { return nil }

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

// Fields of node: Value Next

type ptr_node_codec struct {
	node_codec *node_codec
}

func (c *ptr_node_codec) Init(tcs map[reflect.Type]codecapi.TypeCodec) {
	c.node_codec = tcs[reflect.TypeOf((*node)(nil)).Elem()].(*node_codec)
}

func (c ptr_node_codec) Fields() []string { return nil }

func (c ptr_node_codec) Encode(e *codecapi.Encoder, x interface{}) { c.encode(e, x.(*node)) }

func (c ptr_node_codec) encode(e *codecapi.Encoder, x *node) {
	if !e.StartPtr(x == nil, x) {
		return
	}
	(&node_codec{}).encode(e, x)
}

func (c ptr_node_codec) Decode(d *codecapi.Decoder) interface{} {
	var x *node
	c.decode(d, &x)
	return x
}

func (c ptr_node_codec) decode(d *codecapi.Decoder, p **node) {
	proceed, ref := d.StartPtr()
	if !proceed {
		return
	}
	if ref != nil {
		*p = ref.(*node)
		return
	}
	var x node
	d.StoreRef(&x)
	c.node_codec.decode(d, &x)
	*p = &x
}

type node_codec struct {
	ptr_node_codec *ptr_node_codec
}

func (c *node_codec) Init(tcs map[reflect.Type]codecapi.TypeCodec) {
	c.ptr_node_codec = tcs[reflect.TypeOf((**node)(nil)).Elem()].(*ptr_node_codec)
}

func (c *node_codec) Fields() []string {
	return []string{"Value", "Next"}
}

func (c *node_codec) Encode(e *codecapi.Encoder, x interface{}) {
	s := x.(node)
	c.encode(e, &s)
}

func (c *node_codec) encode(e *codecapi.Encoder, x *node) {
	e.StartStruct()
	if x.Value != 0 {
		e.EncodeUint(0)
		e.EncodeInt(int64(x.Value))
	}
	if x.Next != nil {
		e.EncodeUint(1)
		(&ptr_node_codec{}).encode(e, x.Next)
	}
	e.EndStruct()
}

func (c *node_codec) Decode(d *codecapi.Decoder) interface{} {
	var x node
	c.decode(d, &x)
	return x
}

func (c *node_codec) decode(d *codecapi.Decoder, x *node) {
	d.StartStruct()
	for {
		n := d.NextStructField()
		if n < 0 {
			break
		}
		switch n {
		case 0:
			x.Value = int(d.DecodeInt())
		case 1:
			c.ptr_node_codec.decode(d, &x.Next)
		default:
			d.UnknownField("node", n)
		}
	}
}

func init() {
	codecapi.Register(node{}, func() codecapi.TypeCodec { return &node_codec{} })
	codecapi.Register(&node{}, func() codecapi.TypeCodec { return &ptr_node_codec{} })
}
