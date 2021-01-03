# Fast Encoding of Go Values

This is a fork of pkgsite/internal/godoc/codec.

TODO:

- Put benchmarks in separate module to avoid dependencies on GCS, GCP, etc.

- Support renaming with struct tags.

# Encoding Scheme

Every encoded value begins with a 1-byte code that describes what (if
anything) follows. There is enough information to skip over the value, since
the decoder must be able to do that if it encounters a struct field it
doesn't know.

Most of the values of that initial byte can be devoted to small unsigned
integers. For example, the number 17 is represented by the single byte 17.
Only a few byte values have special meaning.

The nil code indicates that the value is nil. We don't absolutely need this:
we could always represent the nil value for a type as something that couldn't
be mistaken for an encoded value of that type. For instance, we could use 0
for nil in the case of slices (which always begin with the nValues code), and
for pointers to numbers like *int, we could use something like "nBytes 0".
But it is simpler to have a reserved value for nil.

The nBytes code indicates that an unsigned integer N is encoded next,
followed by N bytes of data. There are optimized codes for values of N from 0 to
4. These are used to represent strings and byte slices, as well numbers bigger
than can fit into the initial byte. For example, the string "hello" is represented
as: nBytes 5 'h' 'e' 'l' 'l' 'o'.

Unsigned integers that can't fit into the initial byte are encoded as byte
sequences of length 1, 2, 4 or 8, holding big-endian values.

The nValues code is for sequences of values whose size is known beforehand,
like a Go slice or array. The slice []string{"hi", "bye"} is encoded as
  nValues 2 bytes2 'h' 'i' bytes3 'b' 'y' 'e'

The ptr and refPtr codes indicate a pointer to the encoded value. The latter
signals to the decoder that it should remember the pointer because it will be
referred to later in the stream.

The ref code is used to refer to an earlier encoded pointer. It is followed by a
uint denoting the relative offset to the position of corresponding refPtr code.

The start and end codes delimit a value whose length is unknown beforehand.
They are used for structs.


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

- Compare combined instructions ptrStartCode, refPtrStartCode with the pairs
  ([ref]PtrCode, startCode). Might make a small difference since there can be a
  lot of struct pointers.

# TypeCodec State

The TypeCodec classes are effectively singletons now, but the idea is that a
decoder carries around one per type, and it can have state. That state can hold,
e.g., the number of values of the type, for slice allocation (see below).

Register should store a function that creates a TypeDecoder.
TypeDecoder has type number, field names, allocator state.

There would be an Init function taking something that would allow the TypeCodec
to extract the TypeCodecs it depends on. So a []T codec would ask its arg for
the T codec, and store it in a field that had been generated:

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
Is this fast enough as a global function with a mutex
and a fixed N? Compare with decoder-specific state that
pre-allocates the entire data's worth of Foos, passed
in the encoded header.


# Generics

Generics doesn't really buy you much. The type codecs for slices, maps and
pointers no longer need to be generated, but those for structs still do.
