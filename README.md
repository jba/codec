# Fast Encoding of Go Values

This is a fork of pkgsite/internal/godoc/codec.

TODO:

- Add builtin codecs for all builtin types.

- Put benchmarks in separate module to avoid dependencies on GCS, GCP, etc.

# Features

- Add support for `foo:"name"`.

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

## Zero-copy DecodeBytes

For DecodeBytes to avoid a copy, we would have to be sure that the underlying
buffer isn't shared. Doing so would also pin the entire buffer, so it wouldn't
be a good idea unless a significant amount of the buffer were encoded as byte
slices.

## Allocations

Run
    go build -gcflags -S ./codecapi >& /tmp/asm
after changes and look for runtime.newobject to see what's escaping to the heap.

## Micro-Benchmarks

- Time the explicit x==nil arg to StartStruct against using reflection.

- Compare gob length-prefixed values with nValues and start/end.

- Compare the redundant (ptrCode; startCode) sequence with just startCode.

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

Generics doesn't really buy you much. The type codecs for slices, maps and pointers no
longer need to be generated, but those for structs still do.

This package no longer uses generics, but it used to. To see that version, check
out tag "uses-generics" and:

- Make sure ~/repos/go is at the dev.go2go branch.
- Use the `go2` alias, already defined in ~/.bash_aliases.



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
