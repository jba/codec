This is a fork of pkgsite/internal/godoc/codec.

NEXT:

- Understand why gob allocates much less in the license benchmark.

- Maybe generate more stuff for the hyperledger benchmark. A lot of types and
  fields are unused?

- see if someone can slice a []byte field past the end to access the underlying
  buffer; and if so, use the 3-part slice operator.
- see if gob-encoded ints or floats helps in throughput-limited cases



TODO:
- put benchmarks in separate module to avoid dependencies on GCS, GCP, etc.

- create real-world-ish benchmarks:
  - go/ast.File for net/http (except no one else can do that due to cycles)

- ugorji/go/codec:
  - see codegen option at http://ugorji.net/blog/go-codecgen
  - I think that only works on types you own, because it adds methods to them?
- use gob encoding for uints

-  Test floats reversed vs. not, but only after we pick a uint encoding.
   Reversing only saves space for integer-valued floats, those whose fractional
   part is a power of two, perhaps others
   (https://play.golang.org/p/pYNbvRq1N2S). But your average float will not get
   shorter. It's not clear if it's worth it.

Possible benchmarks:
- https://github.com/robertkrimen/otto/blob/15f95af6e78dcd2030d8195a138bd88d4f403546/script.go
- https://github.com/gocolly/colly/blob/1cd684083cf9bf9a8e33b5dfd6414d8516ae63af/http_backend.go#L161

# Features

- Add support for `foo:"name"`.

- Handle MarshalText.




# TypeCodec state

The TypeCodec classes are effectively singletons now, but the idea is that
a decoder carries around one per type, and it can have state. That state
can hold, e.g., the number of values of the type, for slice allocation.

Register should store a function that creates a TypeDecoder.
TypeDecoder has type number, field names, allocator state.

The Init function now is just a stub, but it would take something that would
allow the TypeCodec to extract the TypeCodecs it depends on. So a []T codec
would ask its arg for the T codec, and store it in a field that had been
generated:

    type slice_Foo_codec struct {
        foo_codec TypeCodecFor[Foo]
    }

    func (s *slice_Foo_codec) Init(ds *codec.DecoderState) {
        s.foo_codec = GetTypeCodec[Foo](ds)
    }

    func (s *slice_Foo_codec) Encode(e *codec.Encoder, s []Foo) {
        codec.EncodeSlice[Foo](e, ???? What do to here?? DO we need codec state for encoding?
    }

    func (s *slice_Foo_codec) Decode(d *codec.Decoder, s *[]Foo) {
        codec.DecodeSlice(d, s, s.foo_codec)
    }


Init would be called by the decoder after it had read the initial metadata
and created all the types' TypeCodecs.

# Performance

Run
    go build -gcflags -S ./codecapi >& /tmp/asm
periodically and look for runtime.newobject to see what's escaping to the heap.

## Micro-Benchmarks

- Time the explicit x==nil arg to StartStruct against using reflection.



- Compare gob length-prefixed values with nValues and start/end.

- Compare the redundant (ptrCode; startCode) sequence with just startCode.

Add nBytes0Code ... nBytes4Code.

Add extCode: followed by uint, denotes an extension?

Faster []bool codec.

## Slice Allocators

Implement slice allocators:


```
var slice_Foo []*Foo

func alloc_Foo() *Foo {
  if len(slice_Foo) == 0 {
    slice_Foo = make([]*Foo, N)
  }
  x := slice_Foo[0]
  slice_Foo = slice_Foo[1]
  return x
}
```
Is this fast enough as a global function with a mutex
and a fixed N? Compare with decoder-specific state that
pre-allocates the entire data's worth of Foos, passed
in the encoded header.


# Generics

This package no longer uses generics, but it used to. To see that version, check
out tag "uses-generics" and:

- Make sure ~/repos/go is at the dev.go2go branch.
- Use the `go2` alias, already defined in ~/.bash_aliases.

# Notes

How does ugorji encoding allocate like one tenth of what we and gob do? See
the license benchmark.

I confirmed that if EncodeAny has a big enough buffer, it does no allocation. In
other words, all our alloc is from growing the buffer. With a magically big enough
buffer, we have 258950K/op; but ugorji, with no previous knowledge, gets
524295K/op, still half of us with an empty buffer and gob.

The answer is that ugorji has a fixed-size buffer that is flushed when full to
the io.Writer. See bufioEncWriter in ugorji/go/codec/writer.go. We can't do that
because we need to write initial metadata.
