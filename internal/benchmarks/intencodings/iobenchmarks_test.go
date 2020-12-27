package main

import (
	"math"
	"testing"
)

// func TestWeightedRandomUint(t *testing.T) {
// 	rand.Seed(time.Now().UnixNano())
// 	var size [10]int
// 	for i := 0; i < 1000; i++ {
// 		r := weightedRandomUint(0.5, 0.25, 0.20)
// 		buf := encodeGob(nil, r)
// 		n := len(buf)
// 		size[n]++
// 	}
// 	fmt.Println(size)
// }

func TestCodecs(t *testing.T) {
	var buf []byte
	want := []uint64{0, 1, 5, 127, 128, 129, 255, 2000, 4e6, 4e10, math.MaxUint32, math.MaxUint64}
	for _, c := range codecs {
		t.Run(c.name, func(t *testing.T) {
			for _, w := range want {
				buf = c.encoder(buf, w)
			}
			for _, w := range want {
				var g uint64
				g, buf = c.decoder(buf)
				if g != w {
					t.Errorf("got %d, want %d", g, w)
				}
			}
		})
	}
}
