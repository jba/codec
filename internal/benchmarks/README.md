# Benchmarks for jba/codec

This module compares some encoders using a collection of real-world and
simulated benchmarks.

The encoders are:

- github.com/jba/codec
- encoding/gob
- github.com/ugorji/go/codec

The benchmarks are:

- ast: The syntax trees for the net/http package, modified to remove cycles.
  Rapid decoding of syntax trees was the original motivation for the codec; see
  https://go.googlesource.com/pkgsite/+/master/internal/godoc/codec.

- scores: A synthetic, integer-heavy benchmark.

- stocks: Float-heavy simulated stock data.

- hyperledger: Some data from a blockchain package.

- licenses: A real collection of license file data from public Go modules.

- licenses-small: The same data as licenses, with truncated license file
  contents so the data isn't dominated by large byte slices.

Each benchmark is run with a set of simulated throughputs, to mimic real-world
situations like reading from a storage bucket or database.

## Running the Benchmarks

First, cd to the `data` directory and run `go generate`.

Then, in this directory, run `go run. bm`.

## Ugorji Code Generation

Generating code with the ugorji codec's `codecgen` program produced some
improvement:
>>>>>>> 59b9a78 (cleanup)

```
* = with codegen

---- stocks at max Mi/sec ----
encode
    jba/codec     1  1320897K/op  1.00s/op 1.00x
          gob     1  1046322K/op  1.25s/op 0.80x
  ugorji-cbor     1   523841K/op  1.22s/op 0.82x
* ugorji-cbor     2   261940K/op  0.77s/op 1.07x


decode
    jba/codec     3  303485K/op  0.37s/op 1.00x
          gob     2  246484K/op  0.73s/op 0.51x
  ugorji-cbor     1  406077K/op  9.12s/op 0.04x
* ugorji-cbor     1  476451K/op  8.09s/op 0.05x

---- stocks at 100 Mi/sec ----
encode
    jba/codec     1  1321311K/op  2.22s/op 1.00x
          gob     1  1046700K/op  2.40s/op 0.92x
  ugorji-cbor     1   524843K/op  1.55s/op 1.44x
* ugorji-cbor     1   526196K/op  1.50s/op 1.38x

decode
    jba/codec     1  303901K/op   1.67s/op 1.00x
          gob     1  246865K/op   1.90s/op 0.88x
  ugorji-cbor     1  406077K/op  11.15s/op 0.15x
* ugorji-cbor     1  476451K/op  10.27s/op 0.16x

---- scores at max Mi/sec ----
encode
    jba/codec    78  10367K/op  0.02s/op 1.00x
          gob    51   9463K/op  0.02s/op 0.66x
  ugorji-cbor    42    197K/op  0.03s/op 0.61x
* ugorji-cbor    44    188K/op  0.02s/op 0.64x

decode
    jba/codec    74  10950K/op  0.02s/op 1.00x
          gob    51  11554K/op  0.02s/op 0.69x
  ugorji-cbor     4  10625K/op  0.30s/op 0.05x
* ugorji-cbor     4  10624K/op  0.30s/op 0.05x

---- scores at 100 Mi/sec ----
encode
    jba/codec    36  10427K/op  0.03s/op 1.00x
          gob    30   9519K/op  0.04s/op 0.82x
  ugorji-cbor    43    193K/op  0.03s/op 1.29x
* ugorji-cbor    42    197K/op  0.03s/op 1.24x

decode
    jba/codec    33  10955K/op  0.03s/op 1.00x
          gob    30  11560K/op  0.04s/op 0.85x
  ugorji-cbor     3  10625K/op  0.41s/op 0.09x
* ugorji-cbor     3  10624K/op  0.37s/op 0.10x
```

## Numeric Encodings

See uint-encodings.txt for data supporting the choice of 1248 encoding.

See float-encodings.txt for data supporting the choice of reversing float bytes.

# Notes

## ugorji allocation
How does ugorji encoding allocate like one tenth of what we and gob do? See
the license benchmark.

I confirmed that if EncodeAny has a big enough buffer, it does no allocation. In
other words, all our alloc is from growing the buffer. With a magically big enough
buffer, we have 258950K/op; but ugorji, with no previous knowledge, gets
524295K/op, still half of us with an empty buffer and gob.

The answer is that ugorji has a fixed-size buffer that is flushed when full to
the io.Writer. See bufioEncWriter in ugorji/go/codec/writer.go. We can't do that
because we need to write initial metadata.

## gob allocation

Gob allocates less than jba/codec when encoding the licenses benchmark:
```
---- licenses at max Mi/sec ----
encode
  jba/codec 48     1  2516760K/op  2.67s/op 1.00x
           gob     1  2110611K/op  2.01s/op 1.33x
```

I'm not sure why, but since on most other benchmarks it allocates more and isn't
faster anyway, I don't think it's worth more investigation.
