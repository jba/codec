/*
These benchmarks demonstrate that gob-style encoding of uints, where a number
consists of a count followed by that many bytes, is better overall than
encodings that represent uints as 1-, 2-, 4-, or 8-byte sequences.

These benchmarks use the following distribution of uints:
50%		1 byte
25%		2 bytes
12.5%	3 bytes
...

Gob uint encoding and decoding is slower by around 15% when the data is already
in memory. But when combined with various forms of I/O, its decreased size makes
it competitive.

Local files: about 7% slower for both encoding and decoding
local DB: faster encoding, same decoding.
cloud DB, GCS: faster for both.
*/
package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/bits"
	"math/rand"
	"testing"

	"cloud.google.com/go/storage"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/jba/codec/internal/benchmarks/bench"
	"github.com/jba/codec/internal/benchmarks/config"
	_ "github.com/lib/pq"
)

type codec struct {
	name    string
	encoder func([]byte, uint64) []byte
	decoder func([]byte) (uint64, []byte)
}

var codecs = []codec{
	{"gob2", encodeGob2, decodeGob2},
	{"gob", encodeGob, decodeGob},
	{"std48", encode48, decode48},
	{"std248", encode248, decode248},
	{"std248a", encode248a, decode248a},
	{"std1248a", encode1248a, decode1248a},
}

func main() {
	cfg, err := config.Read("../config.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("EXPONENTIALLY DECREASING PROBABILITY")
	runAllBenchmarks(cfg, randomUints(1e6))

	fmt.Println()
	fmt.Println("RANDOM uint16s")
	runAllBenchmarks(cfg, randomUint16s(1e6))
}

func runAllBenchmarks(cfg *config.Config, uints []uint64) {
	fmt.Println("Encoded sizes:")
	var s0 int
	for i, c := range codecs {
		s := len(encodeSlice(uints, c.encoder))
		if i == 0 {
			s0 = s
		}
		fmt.Printf("%s: %6d %.2fx\n", c.name, s, float64(s)/float64(s0))
	}

	benchmarkInMemory(uints)

	for _, kind := range []string{"local", "cloud"} {
		fmt.Println()
		db, err := cfg.Connect(kind)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		benchmarkDB(kind, uints, db)
	}

	fmt.Println()
	benchmarkFile(uints)

	fmt.Println()
	benchmarkGCS(cfg.GCSBucket, uints)
}

func runBenchmark(op string, uints []uint64, runner func(b *testing.B, uints []uint64, cd codec)) {
	fmt.Println(op)
	var bms []bench.Benchmark
	for _, codec := range codecs {
		codec := codec
		bm := bench.Benchmark{
			Name: codec.name,
			Func: func(b *testing.B) error { runner(b, uints, codec); return nil },
		}
		bms = append(bms, bm)
	}
	bench.Run(bms)
}

func encodeSlice(s []uint64, encoder func([]byte, uint64) []byte) []byte {
	var buf []byte
	for _, x := range s {
		buf = encoder(buf, x)
	}
	return buf
}

func decodeSlice(data []byte, decoder func([]byte) (uint64, []byte)) {
	for len(data) > 0 {
		_, data = decoder(data)
	}
}

////////////////////////////////////////////////////////////////

func benchmarkInMemory(uints []uint64) {
	fmt.Println("IN MEMORY")
	runBenchmark("encode", uints, encodeInMemoryBenchmark)
	runBenchmark("decode", uints, decodeInMemoryBenchmark)
}

func encodeInMemoryBenchmark(b *testing.B, uints []uint64, cd codec) {
	for i := 0; i < b.N; i++ {
		_ = encodeSlice(uints, cd.encoder)
	}
}

func decodeInMemoryBenchmark(b *testing.B, uints []uint64, cd codec) {
	buf := encodeSlice(uints, cd.encoder)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		decodeSlice(buf, cd.decoder)
	}
}

////////////////////////////////////////////////////////////////

func benchmarkDB(msg string, uints []uint64, db *sql.DB) {
	fmt.Println("DB BENCHMARKS", msg)

	makeTable(db)

	runBenchmark("insert", uints, func(b *testing.B, uints []uint64, cd codec) {
		encodeToDBBenchmark(b, db, uints, cd)
	})
	runBenchmark("select", uints, func(b *testing.B, uints []uint64, cd codec) {
		decodeFromDBBenchmark(b, db, cd)
	})
}

func encodeToDBBenchmark(b *testing.B, db *sql.DB, uints []uint64, cd codec) {
	for i := 0; i < b.N; i++ {
		buf := encodeSlice(uints, cd.encoder)
		_, err := db.Exec(`
			INSERT INTO iobench (name, data) VALUES($1, $2)
			ON CONFLICT (name) DO UPDATE SET data=excluded.data
		`, cd.name, buf)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func decodeFromDBBenchmark(b *testing.B, db *sql.DB, cd codec) {
	for i := 0; i < b.N; i++ {
		var data []byte
		if err := db.QueryRow(`SELECT data FROM iobench WHERE name = $1`, cd.name).Scan(&data); err != nil {
			log.Fatal(err)
		}
		decodeSlice(data, cd.decoder)
	}
}

func makeTable(db *sql.DB) {
	exec(db, `DROP TABLE IF EXISTS iobench`)
	exec(db, `CREATE TABLE iobench (name TEXT PRIMARY KEY, data BYTEA)`)
}

func exec(db *sql.DB, query string, args ...interface{}) {
	_, err := db.Exec(query, args...)
	if err != nil {
		log.Fatalf("%s: %v", query, err)
	}
}

////////////////////////////////////////////////////////////////

const filePrefix = "/tmp/iobench_"

func benchmarkFile(uints []uint64) {
	fmt.Println("FILE BENCHMARKS")
	runBenchmark("write", uints, encodeToFileBenchmark)
	runBenchmark("read", uints, decodeFromFileBenchmark)
}

func encodeToFileBenchmark(b *testing.B, uints []uint64, cd codec) {
	for i := 0; i < b.N; i++ {
		data := encodeSlice(uints, cd.encoder)
		if err := ioutil.WriteFile(filePrefix+cd.name, data, 0644); err != nil {
			log.Fatal(err)
		}
	}
}

func decodeFromFileBenchmark(b *testing.B, _ []uint64, cd codec) {
	for i := 0; i < b.N; i++ {
		data, err := ioutil.ReadFile(filePrefix + cd.name)
		if err != nil {
			log.Fatal(err)
		}
		decodeSlice(data, cd.decoder)
	}
}

////////////////////////////////////////////////////////////////

const objectPrefix = "iobench_"

func benchmarkGCS(bucket string, uints []uint64) {
	fmt.Println("GCS BENCHMARKS")
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	runBenchmark("gcs write", uints, func(b *testing.B, uints []uint64, cd codec) {
		encodeToBucketBenchmark(b, client, bucket, uints, cd)
	})

	runBenchmark("gcs read", uints, func(b *testing.B, _ []uint64, cd codec) {
		decodeFromBucketBenchmark(b, client, bucket, cd)
	})
}

func encodeToBucketBenchmark(b *testing.B, client *storage.Client, bucket string, uints []uint64, cd codec) {
	for i := 0; i < b.N; i++ {
		data := encodeSlice(uints, cd.encoder)
		w := client.Bucket(bucket).Object(objectPrefix + cd.name).NewWriter(context.Background())
		_, err := w.Write(data)
		err2 := w.Close()
		if err != nil || err2 != nil {
			log.Fatal(err, err2)
		}
	}
}

func decodeFromBucketBenchmark(b *testing.B, client *storage.Client, bucket string, cd codec) {
	for i := 0; i < b.N; i++ {
		r, err := client.Bucket(bucket).Object(objectPrefix + cd.name).NewReader(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()
		var buf bytes.Buffer
		_, err = io.Copy(&buf, r)
		if err != nil {
			log.Fatal(err)
		}
		decodeSlice(buf.Bytes(), cd.decoder)
	}
}

////////////////////////////////////////////////////////////////

// randomUints produces a slice of n randomly generated uint64s.
// We generate them randomly to try to reduce the effect of Postgres compression.
// We use the following distribution:
// - Half are 1-byte numbers, which take one or two gob bytes.
// - Half of the remaining take 3 gob bytes.
// -  "                         4  "
// and so on.
func randomUints(n int) []uint64 {
	var us []uint64
	c := n / 2
	for i := 0; i < c; i++ {
		u := uint64(rand.Intn(256))
		if bytelen(u) != 1 {
			panic("fail")
		}
		us = append(us, u)
	}
	for b := 2; b <= 8; b++ {
		c /= 2
		// Generate c b-byte uints.
		fmt.Printf("generating %d %d-byte numbers.\n", c, b)
		mask := ^uint64(0) >> ((8 - b) * 8)
		for i := 0; i < c; i++ {
			r := rand.Uint64() & mask
			for j := 0; j < 10; j++ {
				if bytelen(r) == b {
					break
				}
				r = rand.Uint64() & mask
			}
			if bytelen(r) != b {
				panic("fail")
			}
			us = append(us, r)
		}
	}
	rand.Shuffle(len(us), func(i, j int) { us[i], us[j] = us[j], us[i] })
	return us
}

func randomUint16s(n int) []uint64 {
	var us []uint64
	for i := 0; i < n; i++ {
		us = append(us, uint64(uint16(rand.Uint32())))
	}
	return us
}

func randomUint64s(n int) []uint64 {
	var us []uint64
	for i := 0; i < n; i++ {
		us = append(us, rand.Uint64())
	}
	return us
}

func bytelen(x uint64) int {
	if x == 0 {
		return 1
	}
	return 8 - (bits.LeadingZeros64(x) >> 3)
}

func weightedRandomUint(prob8, prob16, prob32 float64) uint64 {
	p := rand.Float64()
	x := rand.Uint64()
	if p < prob8 {
		return uint64(uint8(x))
	}
	p -= prob8
	if p < prob16 {
		if x < 256 {
			x += 256
		}
		return uint64(uint16(x))
	}
	p -= prob16
	if p < prob32 {
		if x < math.MaxUint16 {
			x += math.MaxUint16
		}
		return uint64(uint32(x))
	}
	return x
}

////////////////////////////////////////////////////////////////
// Encodings.

const (
	nBytesCode = 245
	endCode    = 243
	uint64Size = 8
)

func encodeUints(xs []uint64, encode func([]byte, uint64) []byte) []byte {
	var buf []byte
	for _, x := range xs {
		buf = encode(buf, x)
	}
	return buf
}

var gobBuf [9]byte

func encodeGob(buf []byte, x uint64) []byte {
	// Code from encoding/gob/encode.go:encodeUint.
	if x <= 0x7F {
		return append(buf, uint8(x))
	}
	binary.BigEndian.PutUint64(gobBuf[1:], x)
	bc := bits.LeadingZeros64(x) >> 3   // 8 - bytelen(x)
	gobBuf[bc] = uint8(bc - uint64Size) // and then we subtract 8 to get -bytelen(x)
	return append(buf, gobBuf[bc:uint64Size+1]...)
}

func decodeGob(buf []byte) (uint64, []byte) {
	b := buf[0]
	if b <= 0x7f {
		return uint64(b), buf[1:]
	}
	n := -int(int8(b))
	if n > uint64Size {
		panic("bad uint")
	}
	// Don't need to check error; it's safe to loop regardless.
	// Could check that the high byte is zero but it's not worth it.
	var x uint64
	for _, b := range buf[1 : n+1] {
		x = x<<8 | uint64(b)
	}
	return x, buf[n+1:]
}

func encodeGob2(buf []byte, x uint64) []byte {
	if x < endCode {
		return append(buf, uint8(x))
	}
	binary.BigEndian.PutUint64(gobBuf[1:], x)
	bc := bits.LeadingZeros64(x) >> 3 // 8 - bytelen(x)
	gobBuf[bc] = uint8(uint64Size - bc + nBytesCode)
	return append(buf, gobBuf[bc:]...)
}

func decodeGob2(buf []byte) (uint64, []byte) {
	b := buf[0]
	var u uint64
	if b < endCode {
		return uint64(b), buf[1:]
	}
	n := b - nBytesCode
	if n > uint64Size {
		panic("bad uint")
	}
	for _, b := range buf[1 : n+1] {
		u = u<<8 | uint64(b)
	}
	return u, buf[n+1:]
}

var stdBuf [8]byte

func encode48(buf []byte, u uint64) []byte {
	switch {
	case u < endCode:
		return append(buf, uint8(u))
	case u <= math.MaxUint32:
		// Encode as a sequence of 4 bytes, the little-endian representation of
		// a uint32.
		buf = append(buf, nBytesCode, 4)
		binary.LittleEndian.PutUint32(stdBuf[:4], uint32(u))
		return append(buf, stdBuf[:4]...)
	default:
		// Encode as a sequence of 8 bytes, the little-endian representation of
		// a uint64.
		buf = append(buf, nBytesCode, 8)
		binary.LittleEndian.PutUint64(stdBuf[:8], u)
		return append(buf, stdBuf[:8]...)
	}
}

func decode48(buf []byte) (uint64, []byte) {
	b := buf[0]
	if b < endCode {
		return uint64(b), buf[1:]
	}
	if b != nBytesCode {
		panic("did not see nBytesCode")
	}
	switch s := buf[1]; s {
	case 4:
		return uint64(binary.LittleEndian.Uint32(buf[2 : 2+s])), buf[2+s:]
	case 8:
		return binary.LittleEndian.Uint64(buf[2 : 2+s]), buf[2+s:]
	default:
		panic("bad size")
	}
}

func encode248(buf []byte, u uint64) []byte {
	switch {
	case u < endCode:
		return append(buf, uint8(u))
	case u <= math.MaxUint16:
		buf = append(buf, nBytesCode, 2)
		binary.LittleEndian.PutUint16(stdBuf[:2], uint16(u))
		return append(buf, stdBuf[:2]...)
	case u <= math.MaxUint32:
		// Encode as a sequence of 4 bytes, the little-endian representation of
		// a uint32.
		buf = append(buf, nBytesCode, 4)
		binary.LittleEndian.PutUint32(stdBuf[:4], uint32(u))
		return append(buf, stdBuf[:4]...)
	default:
		// Encode as a sequence of 8 bytes, the little-endian representation of
		// a uint64.
		buf = append(buf, nBytesCode, 8)
		binary.LittleEndian.PutUint64(stdBuf[:8], u)
		return append(buf, stdBuf[:8]...)
	}
}

func decode248(buf []byte) (uint64, []byte) {
	b := buf[0]
	if b < endCode {
		return uint64(b), buf[1:]
	}
	if b != nBytesCode {
		panic("did not see nBytesCode")
	}
	switch s := buf[1]; s {
	case 2:
		return uint64(binary.LittleEndian.Uint16(buf[2 : 2+s])), buf[2+s:]
	case 4:
		return uint64(binary.LittleEndian.Uint32(buf[2 : 2+s])), buf[2+s:]
	case 8:
		return binary.LittleEndian.Uint64(buf[2 : 2+s]), buf[2+s:]
	default:
		panic("bad size")
	}
}

func encode248a(buf []byte, u uint64) []byte {
	switch {
	case u < endCode:
		return append(buf, uint8(u))
	case u <= math.MaxUint16:
		buf = append(buf, nBytesCode+2)
		binary.LittleEndian.PutUint16(stdBuf[:2], uint16(u))
		return append(buf, stdBuf[:2]...)
	case u <= math.MaxUint32:
		// Encode as a sequence of 4 bytes, the little-endian representation of
		// a uint32.
		buf = append(buf, nBytesCode+4)
		binary.LittleEndian.PutUint32(stdBuf[:4], uint32(u))
		return append(buf, stdBuf[:4]...)
	default:
		// Encode as a sequence of 8 bytes, the little-endian representation of
		// a uint64.
		buf = append(buf, nBytesCode+8)
		binary.LittleEndian.PutUint64(stdBuf[:8], u)
		return append(buf, stdBuf[:8]...)
	}
}

func decode248a(buf []byte) (uint64, []byte) {
	b := buf[0]
	if b < endCode {
		return uint64(b), buf[1:]
	}
	if b < nBytesCode || b > nBytesCode+8 {
		panic("bad byte")
	}
	s := b - nBytesCode
	switch s {
	case 2:
		return uint64(binary.LittleEndian.Uint16(buf[1 : 1+s])), buf[1+s:]
	case 4:
		return uint64(binary.LittleEndian.Uint32(buf[1 : 1+s])), buf[1+s:]
	case 8:
		return binary.LittleEndian.Uint64(buf[1 : 1+s]), buf[1+s:]
	default:
		panic("bad size")
	}
}

func encode1248a(buf []byte, u uint64) []byte {
	switch {
	case u < endCode:
		return append(buf, uint8(u))
	case u <= math.MaxUint8:
		return append(buf, nBytesCode+1, uint8(u))
	case u <= math.MaxUint16:
		buf = append(buf, nBytesCode+2)
		binary.LittleEndian.PutUint16(stdBuf[:2], uint16(u))
		return append(buf, stdBuf[:2]...)
	case u <= math.MaxUint32:
		// Encode as a sequence of 4 bytes, the little-endian representation of
		// a uint32.
		buf = append(buf, nBytesCode+4)
		binary.LittleEndian.PutUint32(stdBuf[:4], uint32(u))
		return append(buf, stdBuf[:4]...)
	default:
		// Encode as a sequence of 8 bytes, the little-endian representation of
		// a uint64.
		buf = append(buf, nBytesCode+8)
		binary.LittleEndian.PutUint64(stdBuf[:8], u)
		return append(buf, stdBuf[:8]...)
	}
}

func decode1248a(buf []byte) (uint64, []byte) {
	b := buf[0]
	if b < endCode {
		return uint64(b), buf[1:]
	}
	if b < nBytesCode || b > nBytesCode+8 {
		fmt.Println(b)
		panic("bad byte")
	}
	s := b - nBytesCode
	switch s {
	case 1:
		return uint64(buf[1]), buf[1+s:]
	case 2:
		return uint64(binary.LittleEndian.Uint16(buf[1 : 1+s])), buf[1+s:]
	case 4:
		return uint64(binary.LittleEndian.Uint32(buf[1 : 1+s])), buf[1+s:]
	case 8:
		return binary.LittleEndian.Uint64(buf[1 : 1+s]), buf[1+s:]
	default:
		panic("bad size")
	}
}
