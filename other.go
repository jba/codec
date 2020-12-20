package codec

import (
	"runtime"
)

func totalAlloc() uint64 {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ms.TotalAlloc
}

// for benchmarking
// func (e *Encoder) encodeUintGob(x uint64) {
// 	// Code from encoding/gob/encode.go:encodeUint.
// 	if x <= 0x7F {
// 		e.writeByte(uint8(x))
// 		return
// 	}
// 	binary.BigEndian.PutUint64(e.gobBuf[1:], x)
// 	bc := bits.LeadingZeros64(x) >> 3     // 8 - bytelen(x)
// 	e.gobBuf[bc] = uint8(bc - uint64Size) // and then we subtract 8 to get -bytelen(x)
// 	e.writeBytes(e.gobBuf[bc : uint64Size+1])
// }

// func (d *Decoder) decodeUintGob() (x uint64) {
// 	b := d.readByte()
// 	if b <= 0x7f {
// 		return uint64(b)
// 	}
// 	n := -int(int8(b))
// 	if n > uint64Size {
// 		fail(errors.New("bad uint"))
// 	}
// 	buf := d.readBytes(n)
// 	// Don't need to check error; it's safe to loop regardless.
// 	// Could check that the high byte is zero but it's not worth it.
// 	for _, b := range buf {
// 		x = x<<8 | uint64(b)
// 	}
// 	return x
// }

// Bytes returns the encoded byte slice.
// func (e *Encoder) Bytes() []byte {
// 	data := e.buf                 // remember the data
// 	e.buf = nil                   // start with a fresh buffer
// 	e.encodeInitial()             // encode metadata
// 	return append(e.buf, data...) // concatenate metadata and data
// }

// NewDecoderBytes returns a Decoder for the given bytes.
func NewDecoderBytes(data []byte) *Decoder {
	return nil
}
