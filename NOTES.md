# Developer Notes

TODO:

- Put benchmarks in separate module to avoid dependencies on GCS, GCP, etc.

- Add EncodeBytes, DecodeBytes?


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

- Compare combined instructions ptrStartCode, refPtrStartCode with the pairs
  ([ref]PtrCode, startCode). Might make a small difference since there can be a
  lot of struct pointers.

## TypeCodec State

The TypeCodec classes are effectively singletons now, but the idea is that a
decoder carries around one per type, and it can have state. That state can hold,
e.g., the number of values of the type, for slice allocation (see below).

Register should store a function that creates a TypeDecoder.
TypeDecoder has type number, field names, allocator state.

There would be an Init function taking something that would allow the TypeCodec
to extract the TypeCodecs it depends on. So a []T codec would ask its arg for
the T codec, and store it in a field that had been generated:

```
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
```

Init would be called by the decoder after it had read the initial metadata
and created all the types' TypeCodecs.

## Slice Allocators

Slice allocators can speed decoding by allocating all values of a type at once.


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

Is this fast enough as a global function with a mutex and a fixed N? Compare
with decoder-specific state that pre-allocates the entire data's worth of Foos,
passed in the encoded header.


## Generics

Generics doesn't really buy you much. The type codecs for slices, maps and
pointers no longer need to be generated, but those for structs still do.
