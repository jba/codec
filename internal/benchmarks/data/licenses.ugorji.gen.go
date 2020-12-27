// +build go1.6

// Code generated by codecgen - DO NOT EDIT.

package data

import (
	"errors"
	pkg1_licensecheck "github.com/google/licensecheck"
	codec1978 "github.com/ugorji/go/codec"
	"runtime"
	"strconv"
)

const (
	// ----- content types ----
	codecSelferCcUTF81238 = 1
	codecSelferCcRAW1238  = 255
	// ----- value types used ----
	codecSelferValueTypeArray1238     = 10
	codecSelferValueTypeMap1238       = 9
	codecSelferValueTypeString1238    = 6
	codecSelferValueTypeInt1238       = 2
	codecSelferValueTypeUint1238      = 3
	codecSelferValueTypeFloat1238     = 4
	codecSelferValueTypeNil1238       = 1
	codecSelferBitsize1238            = uint8(32 << (^uint(0) >> 63))
	codecSelferDecContainerLenNil1238 = -2147483648
)

var (
	errCodecSelferOnlyMapOrArrayEncodeToStruct1238 = errors.New(`only encoded map or array can be decoded into a struct`)
)

type codecSelfer1238 struct{}

func codecSelfer1238False() bool { return false }
func codecSelfer1238True() bool  { return true }

func init() {
	if codec1978.GenVersion != 20 {
		_, file, _, _ := runtime.Caller(0)
		ver := strconv.FormatInt(int64(codec1978.GenVersion), 10)
		panic(errors.New("codecgen version mismatch: current: 20, need " + ver + ". Re-generate file: " + file))
	}
	if false { // reference the types, but skip this branch at build/run time
		var _ pkg1_licensecheck.Coverage
	}
}

func (x *LicenseData) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer1238
	z, r := codec1978.GenHelper().Encoder(e)
	_, _, _ = h, z, r
	if x == nil {
		r.EncodeNil()
	} else {
		yy2arr2 := z.EncBasicHandle().StructToArray
		_ = yy2arr2
		const yyr2 bool = false // struct tag has 'toArray'
		if yyr2 || yy2arr2 {
			z.EncWriteArrayStart(2)
			z.EncWriteArrayElem()
			if x.Files == nil {
				r.EncodeNil()
			} else {
				h.encSlicePtrtoLicenseFile(([]*LicenseFile)(x.Files), e)
			} // end block: if x.Files slice == nil
			z.EncWriteArrayElem()
			if x.Contents == nil {
				r.EncodeNil()
			} else {
				h.encSlicePtrtoLicenseContents(([]*LicenseContents)(x.Contents), e)
			} // end block: if x.Contents slice == nil
			z.EncWriteArrayEnd()
		} else {
			z.EncWriteMapStart(2)
			z.EncWriteMapElemKey()
			if z.IsJSONHandle() {
				z.WriteStr("\"Files\"")
			} else {
				r.EncodeString(`Files`)
			}
			z.EncWriteMapElemValue()
			if x.Files == nil {
				r.EncodeNil()
			} else {
				h.encSlicePtrtoLicenseFile(([]*LicenseFile)(x.Files), e)
			} // end block: if x.Files slice == nil
			z.EncWriteMapElemKey()
			if z.IsJSONHandle() {
				z.WriteStr("\"Contents\"")
			} else {
				r.EncodeString(`Contents`)
			}
			z.EncWriteMapElemValue()
			if x.Contents == nil {
				r.EncodeNil()
			} else {
				h.encSlicePtrtoLicenseContents(([]*LicenseContents)(x.Contents), e)
			} // end block: if x.Contents slice == nil
			z.EncWriteMapEnd()
		}
	}
}

func (x *LicenseData) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer1238
	z, r := codec1978.GenHelper().Decoder(d)
	_, _, _ = h, z, r
	yyct2 := r.ContainerType()
	if yyct2 == codecSelferValueTypeNil1238 {
		*(x) = LicenseData{}
	} else if yyct2 == codecSelferValueTypeMap1238 {
		yyl2 := z.DecReadMapStart()
		if yyl2 == 0 {
		} else {
			x.codecDecodeSelfFromMap(yyl2, d)
		}
		z.DecReadMapEnd()
	} else if yyct2 == codecSelferValueTypeArray1238 {
		yyl2 := z.DecReadArrayStart()
		if yyl2 != 0 {
			x.codecDecodeSelfFromArray(yyl2, d)
		}
		z.DecReadArrayEnd()
	} else {
		panic(errCodecSelferOnlyMapOrArrayEncodeToStruct1238)
	}
}

func (x *LicenseData) codecDecodeSelfFromMap(l int, d *codec1978.Decoder) {
	var h codecSelfer1238
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
		case "Files":
			h.decSlicePtrtoLicenseFile((*[]*LicenseFile)(&x.Files), d)
		case "Contents":
			h.decSlicePtrtoLicenseContents((*[]*LicenseContents)(&x.Contents), d)
		default:
			z.DecStructFieldNotFound(-1, yys3)
		} // end switch yys3
	} // end for yyj3
}

func (x *LicenseData) codecDecodeSelfFromArray(l int, d *codec1978.Decoder) {
	var h codecSelfer1238
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
	h.decSlicePtrtoLicenseFile((*[]*LicenseFile)(&x.Files), d)
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
	h.decSlicePtrtoLicenseContents((*[]*LicenseContents)(&x.Contents), d)
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

func (x *LicenseData) IsCodecEmpty() bool {
	return !(len(x.Files) != 0 || len(x.Contents) != 0 || false)
}

func (x *LicenseFile) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer1238
	z, r := codec1978.GenHelper().Encoder(e)
	_, _, _ = h, z, r
	if x == nil {
		r.EncodeNil()
	} else {
		yy2arr2 := z.EncBasicHandle().StructToArray
		_ = yy2arr2
		const yyr2 bool = false // struct tag has 'toArray'
		if yyr2 || yy2arr2 {
			z.EncWriteArrayStart(4)
			z.EncWriteArrayElem()
			r.EncodeString(string(x.Module))
			z.EncWriteArrayElem()
			r.EncodeString(string(x.Version))
			z.EncWriteArrayElem()
			r.EncodeString(string(x.FilePath))
			z.EncWriteArrayElem()
			r.EncodeInt(int64(x.Contents))
			z.EncWriteArrayEnd()
		} else {
			z.EncWriteMapStart(4)
			z.EncWriteMapElemKey()
			if z.IsJSONHandle() {
				z.WriteStr("\"Module\"")
			} else {
				r.EncodeString(`Module`)
			}
			z.EncWriteMapElemValue()
			r.EncodeString(string(x.Module))
			z.EncWriteMapElemKey()
			if z.IsJSONHandle() {
				z.WriteStr("\"Version\"")
			} else {
				r.EncodeString(`Version`)
			}
			z.EncWriteMapElemValue()
			r.EncodeString(string(x.Version))
			z.EncWriteMapElemKey()
			if z.IsJSONHandle() {
				z.WriteStr("\"FilePath\"")
			} else {
				r.EncodeString(`FilePath`)
			}
			z.EncWriteMapElemValue()
			r.EncodeString(string(x.FilePath))
			z.EncWriteMapElemKey()
			if z.IsJSONHandle() {
				z.WriteStr("\"Contents\"")
			} else {
				r.EncodeString(`Contents`)
			}
			z.EncWriteMapElemValue()
			r.EncodeInt(int64(x.Contents))
			z.EncWriteMapEnd()
		}
	}
}

func (x *LicenseFile) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer1238
	z, r := codec1978.GenHelper().Decoder(d)
	_, _, _ = h, z, r
	yyct2 := r.ContainerType()
	if yyct2 == codecSelferValueTypeNil1238 {
		*(x) = LicenseFile{}
	} else if yyct2 == codecSelferValueTypeMap1238 {
		yyl2 := z.DecReadMapStart()
		if yyl2 == 0 {
		} else {
			x.codecDecodeSelfFromMap(yyl2, d)
		}
		z.DecReadMapEnd()
	} else if yyct2 == codecSelferValueTypeArray1238 {
		yyl2 := z.DecReadArrayStart()
		if yyl2 != 0 {
			x.codecDecodeSelfFromArray(yyl2, d)
		}
		z.DecReadArrayEnd()
	} else {
		panic(errCodecSelferOnlyMapOrArrayEncodeToStruct1238)
	}
}

func (x *LicenseFile) codecDecodeSelfFromMap(l int, d *codec1978.Decoder) {
	var h codecSelfer1238
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
		case "Module":
			x.Module = (string)(string(r.DecodeStringAsBytes()))
		case "Version":
			x.Version = (string)(string(r.DecodeStringAsBytes()))
		case "FilePath":
			x.FilePath = (string)(string(r.DecodeStringAsBytes()))
		case "Contents":
			x.Contents = (int)(z.C.IntV(r.DecodeInt64(), codecSelferBitsize1238))
		default:
			z.DecStructFieldNotFound(-1, yys3)
		} // end switch yys3
	} // end for yyj3
}

func (x *LicenseFile) codecDecodeSelfFromArray(l int, d *codec1978.Decoder) {
	var h codecSelfer1238
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
	x.Module = (string)(string(r.DecodeStringAsBytes()))
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
	x.Version = (string)(string(r.DecodeStringAsBytes()))
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
	x.FilePath = (string)(string(r.DecodeStringAsBytes()))
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
	x.Contents = (int)(z.C.IntV(r.DecodeInt64(), codecSelferBitsize1238))
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

func (x *LicenseFile) IsCodecEmpty() bool {
	return !(x.Module != "" || x.Version != "" || x.FilePath != "" || x.Contents != 0 || false)
}

func (x *LicenseContents) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer1238
	z, r := codec1978.GenHelper().Encoder(e)
	_, _, _ = h, z, r
	if x == nil {
		r.EncodeNil()
	} else {
		yy2arr2 := z.EncBasicHandle().StructToArray
		_ = yy2arr2
		const yyr2 bool = false // struct tag has 'toArray'
		if yyr2 || yy2arr2 {
			z.EncWriteArrayStart(6)
			z.EncWriteArrayElem()
			if x.Contents == nil {
				r.EncodeNil()
			} else {
				r.EncodeStringBytesRaw([]byte(x.Contents))
			} // end block: if x.Contents slice == nil
			z.EncWriteArrayElem()
			yy10 := &x.ContentsHash
			h.encArray32uint8((*[32]uint8)(yy10), e)
			z.EncWriteArrayElem()
			if x.OldTypes == nil {
				r.EncodeNil()
			} else {
				z.F.EncSliceStringV(x.OldTypes, e)
			} // end block: if x.OldTypes slice == nil
			z.EncWriteArrayElem()
			yy13 := &x.OldCoverage
			if yyxt14 := z.Extension(yy13); yyxt14 != nil {
				z.EncExtension(yy13, yyxt14)
			} else {
				z.EncFallback(yy13)
			}
			z.EncWriteArrayElem()
			if x.NewTypes == nil {
				r.EncodeNil()
			} else {
				z.F.EncSliceStringV(x.NewTypes, e)
			} // end block: if x.NewTypes slice == nil
			z.EncWriteArrayElem()
			yy16 := &x.NewCoverage
			if yyxt17 := z.Extension(yy16); yyxt17 != nil {
				z.EncExtension(yy16, yyxt17)
			} else {
				z.EncFallback(yy16)
			}
			z.EncWriteArrayEnd()
		} else {
			z.EncWriteMapStart(6)
			z.EncWriteMapElemKey()
			if z.IsJSONHandle() {
				z.WriteStr("\"Contents\"")
			} else {
				r.EncodeString(`Contents`)
			}
			z.EncWriteMapElemValue()
			if x.Contents == nil {
				r.EncodeNil()
			} else {
				r.EncodeStringBytesRaw([]byte(x.Contents))
			} // end block: if x.Contents slice == nil
			z.EncWriteMapElemKey()
			if z.IsJSONHandle() {
				z.WriteStr("\"ContentsHash\"")
			} else {
				r.EncodeString(`ContentsHash`)
			}
			z.EncWriteMapElemValue()
			yy19 := &x.ContentsHash
			h.encArray32uint8((*[32]uint8)(yy19), e)
			z.EncWriteMapElemKey()
			if z.IsJSONHandle() {
				z.WriteStr("\"OldTypes\"")
			} else {
				r.EncodeString(`OldTypes`)
			}
			z.EncWriteMapElemValue()
			if x.OldTypes == nil {
				r.EncodeNil()
			} else {
				z.F.EncSliceStringV(x.OldTypes, e)
			} // end block: if x.OldTypes slice == nil
			z.EncWriteMapElemKey()
			if z.IsJSONHandle() {
				z.WriteStr("\"OldCoverage\"")
			} else {
				r.EncodeString(`OldCoverage`)
			}
			z.EncWriteMapElemValue()
			yy22 := &x.OldCoverage
			if yyxt23 := z.Extension(yy22); yyxt23 != nil {
				z.EncExtension(yy22, yyxt23)
			} else {
				z.EncFallback(yy22)
			}
			z.EncWriteMapElemKey()
			if z.IsJSONHandle() {
				z.WriteStr("\"NewTypes\"")
			} else {
				r.EncodeString(`NewTypes`)
			}
			z.EncWriteMapElemValue()
			if x.NewTypes == nil {
				r.EncodeNil()
			} else {
				z.F.EncSliceStringV(x.NewTypes, e)
			} // end block: if x.NewTypes slice == nil
			z.EncWriteMapElemKey()
			if z.IsJSONHandle() {
				z.WriteStr("\"NewCoverage\"")
			} else {
				r.EncodeString(`NewCoverage`)
			}
			z.EncWriteMapElemValue()
			yy25 := &x.NewCoverage
			if yyxt26 := z.Extension(yy25); yyxt26 != nil {
				z.EncExtension(yy25, yyxt26)
			} else {
				z.EncFallback(yy25)
			}
			z.EncWriteMapEnd()
		}
	}
}

func (x *LicenseContents) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer1238
	z, r := codec1978.GenHelper().Decoder(d)
	_, _, _ = h, z, r
	yyct2 := r.ContainerType()
	if yyct2 == codecSelferValueTypeNil1238 {
		*(x) = LicenseContents{}
	} else if yyct2 == codecSelferValueTypeMap1238 {
		yyl2 := z.DecReadMapStart()
		if yyl2 == 0 {
		} else {
			x.codecDecodeSelfFromMap(yyl2, d)
		}
		z.DecReadMapEnd()
	} else if yyct2 == codecSelferValueTypeArray1238 {
		yyl2 := z.DecReadArrayStart()
		if yyl2 != 0 {
			x.codecDecodeSelfFromArray(yyl2, d)
		}
		z.DecReadArrayEnd()
	} else {
		panic(errCodecSelferOnlyMapOrArrayEncodeToStruct1238)
	}
}

func (x *LicenseContents) codecDecodeSelfFromMap(l int, d *codec1978.Decoder) {
	var h codecSelfer1238
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
		case "Contents":
			x.Contents = r.DecodeBytes(([]byte)(x.Contents), false)
		case "ContentsHash":
			h.decArray32uint8((*[32]uint8)(&x.ContentsHash), d)
		case "OldTypes":
			z.F.DecSliceStringX(&x.OldTypes, d)
		case "OldCoverage":
			if yyxt11 := z.Extension(x.OldCoverage); yyxt11 != nil {
				z.DecExtension(&x.OldCoverage, yyxt11)
			} else {
				z.DecFallback(&x.OldCoverage, false)
			}
		case "NewTypes":
			z.F.DecSliceStringX(&x.NewTypes, d)
		case "NewCoverage":
			if yyxt15 := z.Extension(x.NewCoverage); yyxt15 != nil {
				z.DecExtension(&x.NewCoverage, yyxt15)
			} else {
				z.DecFallback(&x.NewCoverage, false)
			}
		default:
			z.DecStructFieldNotFound(-1, yys3)
		} // end switch yys3
	} // end for yyj3
}

func (x *LicenseContents) codecDecodeSelfFromArray(l int, d *codec1978.Decoder) {
	var h codecSelfer1238
	z, r := codec1978.GenHelper().Decoder(d)
	_, _, _ = h, z, r
	var yyj16 int
	var yyb16 bool
	var yyhl16 bool = l >= 0
	yyj16++
	if yyhl16 {
		yyb16 = yyj16 > l
	} else {
		yyb16 = z.DecCheckBreak()
	}
	if yyb16 {
		z.DecReadArrayEnd()
		return
	}
	z.DecReadArrayElem()
	x.Contents = r.DecodeBytes(([]byte)(x.Contents), false)
	yyj16++
	if yyhl16 {
		yyb16 = yyj16 > l
	} else {
		yyb16 = z.DecCheckBreak()
	}
	if yyb16 {
		z.DecReadArrayEnd()
		return
	}
	z.DecReadArrayElem()
	h.decArray32uint8((*[32]uint8)(&x.ContentsHash), d)
	yyj16++
	if yyhl16 {
		yyb16 = yyj16 > l
	} else {
		yyb16 = z.DecCheckBreak()
	}
	if yyb16 {
		z.DecReadArrayEnd()
		return
	}
	z.DecReadArrayElem()
	z.F.DecSliceStringX(&x.OldTypes, d)
	yyj16++
	if yyhl16 {
		yyb16 = yyj16 > l
	} else {
		yyb16 = z.DecCheckBreak()
	}
	if yyb16 {
		z.DecReadArrayEnd()
		return
	}
	z.DecReadArrayElem()
	if yyxt24 := z.Extension(x.OldCoverage); yyxt24 != nil {
		z.DecExtension(&x.OldCoverage, yyxt24)
	} else {
		z.DecFallback(&x.OldCoverage, false)
	}
	yyj16++
	if yyhl16 {
		yyb16 = yyj16 > l
	} else {
		yyb16 = z.DecCheckBreak()
	}
	if yyb16 {
		z.DecReadArrayEnd()
		return
	}
	z.DecReadArrayElem()
	z.F.DecSliceStringX(&x.NewTypes, d)
	yyj16++
	if yyhl16 {
		yyb16 = yyj16 > l
	} else {
		yyb16 = z.DecCheckBreak()
	}
	if yyb16 {
		z.DecReadArrayEnd()
		return
	}
	z.DecReadArrayElem()
	if yyxt28 := z.Extension(x.NewCoverage); yyxt28 != nil {
		z.DecExtension(&x.NewCoverage, yyxt28)
	} else {
		z.DecFallback(&x.NewCoverage, false)
	}
	for {
		yyj16++
		if yyhl16 {
			yyb16 = yyj16 > l
		} else {
			yyb16 = z.DecCheckBreak()
		}
		if yyb16 {
			break
		}
		z.DecReadArrayElem()
		z.DecStructFieldNotFound(yyj16-1, "")
	}
}

func (x *LicenseContents) IsCodecEmpty() bool {
	return !(len(x.Contents) != 0 || len(x.ContentsHash) != 0 || len(x.OldTypes) != 0 || false || x.OldCoverage.Percent != 0 || len(x.OldCoverage.Match) != 0 || len(x.NewTypes) != 0 || false || x.NewCoverage.Percent != 0 || len(x.NewCoverage.Match) != 0 || false)
}

func (x codecSelfer1238) encSlicePtrtoLicenseFile(v []*LicenseFile, e *codec1978.Encoder) {
	var h codecSelfer1238
	z, r := codec1978.GenHelper().Encoder(e)
	_, _, _ = h, z, r
	if v == nil {
		r.EncodeNil()
		return
	}
	z.EncWriteArrayStart(len(v))
	for _, yyv1 := range v {
		z.EncWriteArrayElem()
		if yyv1 == nil {
			r.EncodeNil()
		} else {
			if yyxt2 := z.Extension(yyv1); yyxt2 != nil {
				z.EncExtension(yyv1, yyxt2)
			} else {
				yyv1.CodecEncodeSelf(e)
			}
		}
	}
	z.EncWriteArrayEnd()
}

func (x codecSelfer1238) decSlicePtrtoLicenseFile(v *[]*LicenseFile, d *codec1978.Decoder) {
	var h codecSelfer1238
	z, r := codec1978.GenHelper().Decoder(d)
	_, _, _ = h, z, r

	yyv1 := *v
	yyh1, yyl1 := z.DecSliceHelperStart()
	var yyc1 bool
	_ = yyc1
	if yyh1.IsNil {
		if yyv1 != nil {
			yyv1 = nil
			yyc1 = true
		}
	} else if yyl1 == 0 {
		if yyv1 == nil {
			yyv1 = []*LicenseFile{}
			yyc1 = true
		} else if len(yyv1) != 0 {
			yyv1 = yyv1[:0]
			yyc1 = true
		}
	} else {
		yyhl1 := yyl1 > 0
		var yyrl1 int
		_ = yyrl1
		if yyhl1 {
			if yyl1 > cap(yyv1) {
				yyrl1 = z.DecInferLen(yyl1, z.DecBasicHandle().MaxInitLen, 8)
				if yyrl1 <= cap(yyv1) {
					yyv1 = yyv1[:yyrl1]
				} else {
					yyv1 = make([]*LicenseFile, yyrl1)
				}
				yyc1 = true
			} else if yyl1 != len(yyv1) {
				yyv1 = yyv1[:yyl1]
				yyc1 = true
			}
		}
		var yyj1 int
		for yyj1 = 0; (yyhl1 && yyj1 < yyl1) || !(yyhl1 || z.DecCheckBreak()); yyj1++ { // bounds-check-elimination
			if yyj1 == 0 && yyv1 == nil {
				if yyhl1 {
					yyrl1 = z.DecInferLen(yyl1, z.DecBasicHandle().MaxInitLen, 8)
				} else {
					yyrl1 = 8
				}
				yyv1 = make([]*LicenseFile, yyrl1)
				yyc1 = true
			}
			yyh1.ElemContainerState(yyj1)
			var yydb1 bool
			if yyj1 >= len(yyv1) {
				yyv1 = append(yyv1, nil)
				yyc1 = true
			}
			if yydb1 {
				z.DecSwallow()
			} else {
				if r.TryNil() {
					yyv1[yyj1] = nil
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(LicenseFile)
					}
					if yyxt3 := z.Extension(yyv1[yyj1]); yyxt3 != nil {
						z.DecExtension(yyv1[yyj1], yyxt3)
					} else {
						yyv1[yyj1].CodecDecodeSelf(d)
					}
				}
			}
		}
		if yyj1 < len(yyv1) {
			yyv1 = yyv1[:yyj1]
			yyc1 = true
		} else if yyj1 == 0 && yyv1 == nil {
			yyv1 = make([]*LicenseFile, 0)
			yyc1 = true
		}
	}
	yyh1.End()
	if yyc1 {
		*v = yyv1
	}
}

func (x codecSelfer1238) encSlicePtrtoLicenseContents(v []*LicenseContents, e *codec1978.Encoder) {
	var h codecSelfer1238
	z, r := codec1978.GenHelper().Encoder(e)
	_, _, _ = h, z, r
	if v == nil {
		r.EncodeNil()
		return
	}
	z.EncWriteArrayStart(len(v))
	for _, yyv1 := range v {
		z.EncWriteArrayElem()
		if yyv1 == nil {
			r.EncodeNil()
		} else {
			if yyxt2 := z.Extension(yyv1); yyxt2 != nil {
				z.EncExtension(yyv1, yyxt2)
			} else {
				yyv1.CodecEncodeSelf(e)
			}
		}
	}
	z.EncWriteArrayEnd()
}

func (x codecSelfer1238) decSlicePtrtoLicenseContents(v *[]*LicenseContents, d *codec1978.Decoder) {
	var h codecSelfer1238
	z, r := codec1978.GenHelper().Decoder(d)
	_, _, _ = h, z, r

	yyv1 := *v
	yyh1, yyl1 := z.DecSliceHelperStart()
	var yyc1 bool
	_ = yyc1
	if yyh1.IsNil {
		if yyv1 != nil {
			yyv1 = nil
			yyc1 = true
		}
	} else if yyl1 == 0 {
		if yyv1 == nil {
			yyv1 = []*LicenseContents{}
			yyc1 = true
		} else if len(yyv1) != 0 {
			yyv1 = yyv1[:0]
			yyc1 = true
		}
	} else {
		yyhl1 := yyl1 > 0
		var yyrl1 int
		_ = yyrl1
		if yyhl1 {
			if yyl1 > cap(yyv1) {
				yyrl1 = z.DecInferLen(yyl1, z.DecBasicHandle().MaxInitLen, 8)
				if yyrl1 <= cap(yyv1) {
					yyv1 = yyv1[:yyrl1]
				} else {
					yyv1 = make([]*LicenseContents, yyrl1)
				}
				yyc1 = true
			} else if yyl1 != len(yyv1) {
				yyv1 = yyv1[:yyl1]
				yyc1 = true
			}
		}
		var yyj1 int
		for yyj1 = 0; (yyhl1 && yyj1 < yyl1) || !(yyhl1 || z.DecCheckBreak()); yyj1++ { // bounds-check-elimination
			if yyj1 == 0 && yyv1 == nil {
				if yyhl1 {
					yyrl1 = z.DecInferLen(yyl1, z.DecBasicHandle().MaxInitLen, 8)
				} else {
					yyrl1 = 8
				}
				yyv1 = make([]*LicenseContents, yyrl1)
				yyc1 = true
			}
			yyh1.ElemContainerState(yyj1)
			var yydb1 bool
			if yyj1 >= len(yyv1) {
				yyv1 = append(yyv1, nil)
				yyc1 = true
			}
			if yydb1 {
				z.DecSwallow()
			} else {
				if r.TryNil() {
					yyv1[yyj1] = nil
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(LicenseContents)
					}
					if yyxt3 := z.Extension(yyv1[yyj1]); yyxt3 != nil {
						z.DecExtension(yyv1[yyj1], yyxt3)
					} else {
						yyv1[yyj1].CodecDecodeSelf(d)
					}
				}
			}
		}
		if yyj1 < len(yyv1) {
			yyv1 = yyv1[:yyj1]
			yyc1 = true
		} else if yyj1 == 0 && yyv1 == nil {
			yyv1 = make([]*LicenseContents, 0)
			yyc1 = true
		}
	}
	yyh1.End()
	if yyc1 {
		*v = yyv1
	}
}

func (x codecSelfer1238) encArray32uint8(v *[32]uint8, e *codec1978.Encoder) {
	var h codecSelfer1238
	z, r := codec1978.GenHelper().Encoder(e)
	_, _, _ = h, z, r
	if v == nil {
		r.EncodeNil()
		return
	}
	r.EncodeStringBytesRaw(((*[32]byte)(v))[:])
}

func (x codecSelfer1238) decArray32uint8(v *[32]uint8, d *codec1978.Decoder) {
	var h codecSelfer1238
	z, r := codec1978.GenHelper().Decoder(d)
	_, _, _ = h, z, r
	r.DecodeBytes(((*[32]byte)(v))[:], true)
}
