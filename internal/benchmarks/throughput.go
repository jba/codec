// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
   This function computes the break-even throughput of two codecs, where second
   codec produces smaller output with possibly greater computation.

   Let
        p = throughput
        n = size
   Then
        I/O time = n / p

   Say variant 1 produces n1 bytes and takes t1 time with infinite throughput.
   Variant 2 produces n2 < n1 bytes and takes t2 > t1 time.

   Total time Ti = ti + ni / p

   When T1 = T2,
        t1 + n1/p = t2 + n2/p
        (n1 - n2)/p = t2 - t1
   Break-even throughput p = (n1 - n2) / (t2 - t1).

   Example: the more compact encoding saves 30K but takes 100ms longer.
   Break-even throughput is 30K/.1s = 300 K/s.

   If data can be read/written faster than the break-even throughput, then the
   more compact encoding isn't worth the extra time. For example, at 1 M/s,
   those 30K take ~30ms, much less than the 100ms extra compute time.

   Slower than the break-even throughput and it's worth it. For example, at 100
   K/s, the 30K take 300ms.
*/
package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jba/codec/internal/benchmarks/bench"
	"github.com/jba/codec/internal/benchmarks/data"
)

/*
Using codes for short byte sequences: bytes0Code .. bytes4Code:

jba/codec orig vs. jba/codec shortlen on hyperledger
encoding:
Space-optimized codec took less space (27400 vs. 27556) and was not slower (0ms vs. 0ms)
decoding:
Decoding was not slower (0ms vs. 0ms)

jba/codec orig vs. jba/codec shortlen on licenses
encoding:
Space-optimized codec took less space (275684053 vs. 277170253) and was not slower (1846ms vs. 2988ms)
decoding:
Decoding was not slower (775ms vs. 777ms)

jba/codec orig vs. jba/codec shortlen on licenses-small
encoding:
Space-optimized codec took less space (193983702 vs. 195469902) and was not slower (805ms vs. 819ms)
decoding:
Decoding break-even throughput: 107.0 M/s

jba/codec orig vs. jba/codec shortlen on stocks
encoding:
Space-optimized codec took less space (135783153 vs. 135783553) and was not slower (649ms vs. 657ms)
decoding:
Decoding break-even throughput: 0.1 M/s
Space savings not worth if for decoding.

jba/codec orig vs. jba/codec shortlen on scores
encoding:
Space-optimized codec took less space (2375708 vs. 2697131) and was not slower (18ms vs. 20ms)
decoding:
Decoding was not slower (20ms vs. 20ms)
*/

// If the break-even throughput is less than this, then the space optimization
// doesn't buy time even in low-throughput situations (reading from a cloud DB
// is about 60 M/s).
const throughputThreshold = 60

func runBreakEvenThroughput(dataNames []string) error {
	return errors.New("inactive")
	for _, bd := range datasToRun(dataNames) {
		if err := breakEvenThroughput(jbaCodec, gobCodec, bd); err != nil {
			log.Fatal(err)
		}
		fmt.Println()
	}
	return nil
}

func breakEvenThroughput(origCodec, spaceOptCodec Codec, bd data.BenchmarkData) error {
	fmt.Printf("%s vs. %s on %s\n", origCodec.name, spaceOptCodec.name, bd.Name)
	data, err := bd.Read()
	if err != nil {
		return err
	}
	origEnc, origTime, err := runEncoded(origCodec, data)
	if err != nil {
		return err
	}
	origLen := len(origEnc)

	optEnc, optTime, err := runEncoded(spaceOptCodec, data)
	if err != nil {
		return err
	}
	optLen := len(optEnc)

	fmt.Println("encoding:")
	savedBytes := origLen - optLen
	if savedBytes == 0 {
		fmt.Println("No space saved.")
		return nil
	}
	if savedBytes < 0 {
		fmt.Printf("Space-optimized codec took more space: %d vs. %d\n", optLen, origLen)
		return nil
	}
	deltaTime := optTime - origTime
	if deltaTime <= 0 {
		fmt.Printf("Space-optimized codec took less space (%d vs. %d) and was not slower (%dms vs. %dms)\n",
			optLen, origLen, optTime.Milliseconds(), origTime.Milliseconds())
	} else {
		fmt.Printf("saved bytes: %d, extra time: %dµs\n", savedBytes, deltaTime.Microseconds())
		bet := float64(savedBytes) / 1024 / 1024 / deltaTime.Seconds()
		fmt.Printf("Encoding break-even throughput: %.1f M/s\n", bet)
		if bet < throughputThreshold {
			fmt.Printf("Space savings not worth it for encoding.\n")
		}
	}

	fmt.Println("decoding:")
	origTime, err = runDecoded(origCodec, origEnc, bd.Newptr)
	if err != nil {
		return err
	}
	optTime, err = runDecoded(spaceOptCodec, optEnc, bd.Newptr)
	if err != nil {
		return err
	}
	deltaTime = optTime - origTime
	if deltaTime <= 0 {
		fmt.Printf("Decoding was not slower (%dms vs. %dms)\n", optTime.Milliseconds(), origTime.Milliseconds())
		return nil
	}
	bet := float64(savedBytes) / 1024 / 1024 / deltaTime.Seconds()
	fmt.Printf("Decoding break-even throughput: %.1f M/s\n", bet)
	if bet < throughputThreshold {
		fmt.Printf("Space savings not worth it for decoding.\n")
	}
	return nil
}

func runEncoded(c Codec, data interface{}) ([]byte, time.Duration, error) {
	var encoded []byte
	r, err := bench.Run1(newEncodeBenchmark(data, c, 0, &encoded))
	if err != nil {
		return nil, 0, fmt.Errorf("encoding with %s: %w", c.name, err)
	}
	return encoded, time.Duration(r.NsPerOp()), nil
}

func runDecoded(c Codec, enc []byte, newptr func() interface{}) (time.Duration, error) {
	r, err := bench.Run1(newDecodeBenchmark(enc, c, 0, newptr))
	if err != nil {
		return 0, fmt.Errorf("decoding with %s: %w", c.name, err)
	}
	return time.Duration(r.NsPerOp()), nil
}
