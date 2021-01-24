// +build go1.6

// Code generated by codecgen - DO NOT EDIT.

package data

import (
	"errors"
	codec1978 "github.com/ugorji/go/codec"
	"runtime"
	"strconv"
)

const (
	// ----- content types ----
	codecSelferCcUTF85981 = 1
	codecSelferCcRAW5981  = 255
	// ----- value types used ----
	codecSelferValueTypeArray5981     = 10
	codecSelferValueTypeMap5981       = 9
	codecSelferValueTypeString5981    = 6
	codecSelferValueTypeInt5981       = 2
	codecSelferValueTypeUint5981      = 3
	codecSelferValueTypeFloat5981     = 4
	codecSelferValueTypeNil5981       = 1
	codecSelferBitsize5981            = uint8(32 << (^uint(0) >> 63))
	codecSelferDecContainerLenNil5981 = -2147483648
)

var (
	errCodecSelferOnlyMapOrArrayEncodeToStruct5981 = errors.New(`only encoded map or array can be decoded into a struct`)
)

type codecSelfer5981 struct{}

func codecSelfer5981False() bool { return false }
func codecSelfer5981True() bool  { return true }

func init() {
	if codec1978.GenVersion != 20 {
		_, file, _, _ := runtime.Caller(0)
		ver := strconv.FormatInt(int64(codec1978.GenVersion), 10)
		panic(errors.New("codecgen version mismatch: current: 20, need " + ver + ". Re-generate file: " + file))
	}
}

func (x *Score) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer5981
	z, r := codec1978.GenHelper().Encoder(e)
	_, _, _ = h, z, r
	if x == nil {
		r.EncodeNil()
	} else {
		yy2arr2 := z.EncBasicHandle().StructToArray
		_ = yy2arr2
		const yyr2 bool = false // struct tag has 'toArray'
		if yyr2 || yy2arr2 {
			z.EncWriteArrayStart(3)
			z.EncWriteArrayElem()
			r.EncodeInt(int64(x.GameID))
			z.EncWriteArrayElem()
			r.EncodeInt(int64(x.PlayerID))
			z.EncWriteArrayElem()
			if x.Scores == nil {
				r.EncodeNil()
			} else {
				z.F.EncSliceIntV(x.Scores, e)
			} // end block: if x.Scores slice == nil
			z.EncWriteArrayEnd()
		} else {
			z.EncWriteMapStart(3)
			z.EncWriteMapElemKey()
			if z.IsJSONHandle() {
				z.WriteStr("\"GameID\"")
			} else {
				r.EncodeString(`GameID`)
			}
			z.EncWriteMapElemValue()
			r.EncodeInt(int64(x.GameID))
			z.EncWriteMapElemKey()
			if z.IsJSONHandle() {
				z.WriteStr("\"PlayerID\"")
			} else {
				r.EncodeString(`PlayerID`)
			}
			z.EncWriteMapElemValue()
			r.EncodeInt(int64(x.PlayerID))
			z.EncWriteMapElemKey()
			if z.IsJSONHandle() {
				z.WriteStr("\"Scores\"")
			} else {
				r.EncodeString(`Scores`)
			}
			z.EncWriteMapElemValue()
			if x.Scores == nil {
				r.EncodeNil()
			} else {
				z.F.EncSliceIntV(x.Scores, e)
			} // end block: if x.Scores slice == nil
			z.EncWriteMapEnd()
		}
	}
}

func (x *Score) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer5981
	z, r := codec1978.GenHelper().Decoder(d)
	_, _, _ = h, z, r
	yyct2 := r.ContainerType()
	if yyct2 == codecSelferValueTypeNil5981 {
		*(x) = Score{}
	} else if yyct2 == codecSelferValueTypeMap5981 {
		yyl2 := z.DecReadMapStart()
		if yyl2 == 0 {
		} else {
			x.codecDecodeSelfFromMap(yyl2, d)
		}
		z.DecReadMapEnd()
	} else if yyct2 == codecSelferValueTypeArray5981 {
		yyl2 := z.DecReadArrayStart()
		if yyl2 != 0 {
			x.codecDecodeSelfFromArray(yyl2, d)
		}
		z.DecReadArrayEnd()
	} else {
		panic(errCodecSelferOnlyMapOrArrayEncodeToStruct5981)
	}
}

func (x *Score) codecDecodeSelfFromMap(l int, d *codec1978.Decoder) {
	var h codecSelfer5981
	z, r := codec1978.GenHelper().Decoder(d)
	_, _, _ = h, z, r
	var yyhl3 bool = l >= 0
	for yyj3 := 0; ; yyj3++ {
		if yyhl3 {
			if yyj3 >= l {
				break
			}
		} else {
			if z.DecCheckBreak() {
				break
			}
		}
		z.DecReadMapElemKey()
		yys3 := z.StringView(r.DecodeStringAsBytes())
		z.DecReadMapElemValue()
		switch yys3 {
		case "GameID":
			x.GameID = (int)(z.C.IntV(r.DecodeInt64(), codecSelferBitsize5981))
		case "PlayerID":
			x.PlayerID = (int)(z.C.IntV(r.DecodeInt64(), codecSelferBitsize5981))
		case "Scores":
			z.F.DecSliceIntX(&x.Scores, d)
		default:
			z.DecStructFieldNotFound(-1, yys3)
		} // end switch yys3
	} // end for yyj3
}

func (x *Score) codecDecodeSelfFromArray(l int, d *codec1978.Decoder) {
	var h codecSelfer5981
	z, r := codec1978.GenHelper().Decoder(d)
	_, _, _ = h, z, r
	var yyj8 int
	var yyb8 bool
	var yyhl8 bool = l >= 0
	yyj8++
	if yyhl8 {
		yyb8 = yyj8 > l
	} else {
		yyb8 = z.DecCheckBreak()
	}
	if yyb8 {
		z.DecReadArrayEnd()
		return
	}
	z.DecReadArrayElem()
	x.GameID = (int)(z.C.IntV(r.DecodeInt64(), codecSelferBitsize5981))
	yyj8++
	if yyhl8 {
		yyb8 = yyj8 > l
	} else {
		yyb8 = z.DecCheckBreak()
	}
	if yyb8 {
		z.DecReadArrayEnd()
		return
	}
	z.DecReadArrayElem()
	x.PlayerID = (int)(z.C.IntV(r.DecodeInt64(), codecSelferBitsize5981))
	yyj8++
	if yyhl8 {
		yyb8 = yyj8 > l
	} else {
		yyb8 = z.DecCheckBreak()
	}
	if yyb8 {
		z.DecReadArrayEnd()
		return
	}
	z.DecReadArrayElem()
	z.F.DecSliceIntX(&x.Scores, d)
	for {
		yyj8++
		if yyhl8 {
			yyb8 = yyj8 > l
		} else {
			yyb8 = z.DecCheckBreak()
		}
		if yyb8 {
			break
		}
		z.DecReadArrayElem()
		z.DecStructFieldNotFound(yyj8-1, "")
	}
}

func (x *Score) IsCodecEmpty() bool {
	return !(x.GameID != 0 || x.PlayerID != 0 || len(x.Scores) != 0 || false)
}
