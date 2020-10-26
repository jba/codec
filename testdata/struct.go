// Code generated by the codec package. DO NOT EDIT.

package somepkg

import (
	"github.com/jba/codec"
	"go/ast"
	"go/token"
)

// Fields of ast_BasicLit: ValuePos Kind Value

type ptr_ast_BasicLit_codec struct{}

func (ptr_ast_BasicLit_codec) Init() {}

func (c ptr_ast_BasicLit_codec) Encode(e *codec.Encoder, x interface{}) {
	c.encode(e, x.(*ast.BasicLit))
}

func (c ptr_ast_BasicLit_codec) encode(e *codec.Encoder, x *ast.BasicLit) {
	if !e.StartPtr(x == nil, x) {
		return
	}
	(ast_BasicLit_codec{}).encode(e, x)
}

func (c ptr_ast_BasicLit_codec) Decode(d *codec.Decoder) interface{} {
	var x *ast.BasicLit
	c.decode(d, &x)
	return x
}

func (c ptr_ast_BasicLit_codec) decode(d *codec.Decoder, p **ast.BasicLit) {
	proceed, ref := d.StartPtr()
	if !proceed {
		return
	}
	if ref != nil {
		*p = ref.(*ast.BasicLit)
		return
	}
	var x ast.BasicLit
	d.StoreRef(&x)
	(ast_BasicLit_codec{}).decode(d, &x)
	*p = &x
}

type ast_BasicLit_codec struct{}

func (ast_BasicLit_codec) Init() {}

func (c ast_BasicLit_codec) Encode(e *codec.Encoder, x interface{}) {
	s := x.(ast.BasicLit)
	c.encode(e, &s)
}

func (c ast_BasicLit_codec) encode(e *codec.Encoder, x *ast.BasicLit) {
	e.StartStruct()
	if x.ValuePos != 0 {
		e.EncodeUint(0)
		e.EncodeInt(int64(x.ValuePos))
	}
	if x.Kind != 0 {
		e.EncodeUint(1)
		e.EncodeInt(int64(x.Kind))
	}
	if x.Value != "" {
		e.EncodeUint(2)
		e.EncodeString(x.Value)
	}
	e.EndStruct()
}

func (c ast_BasicLit_codec) Decode(d *codec.Decoder) interface{} {
	var x ast.BasicLit
	c.decode(d, &x)
	return x
}

func (c ast_BasicLit_codec) decode(d *codec.Decoder, x *ast.BasicLit) {
	d.StartStruct()
	for {
		n := d.NextStructField()
		if n < 0 {
			break
		}
		switch n {
		case 0:
			x.ValuePos = token.Pos(d.DecodeInt())
		case 1:
			x.Kind = token.Token(d.DecodeInt())
		case 2:
			x.Value = d.DecodeString()
		default:
			d.UnknownField("ast.BasicLit", n)
		}
	}
}

func init() {
	codec.Register(ast.BasicLit{}, ast_BasicLit_codec{})
	codec.Register(&ast.BasicLit{}, ptr_ast_BasicLit_codec{})
}
