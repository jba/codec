# Fast Encoding of Go Values

This project is an enhanced version of the package
[pkgsite/internal/godoc/codec](https://pkg.go.dev/golang.org/x/pkgsite/internal/godoc/codec).

The original motivation was fast decoding of parsed Go files, of type
go/ast.File. The pkg.go.dev site saves these when processing a module, and
decodes them on the serving path to render documentation. So decoding had to be
fast, and had to handle the cycles that these structures contain. It also had to
work with existing types that we did not control. We couldn't find any existing
encoders with these properties, so we wrote our own.

For usage, see the package documentation.

## Encoding Scheme

Go values are converted to byte sequences by mapping them to a low-level
wire protocol.

### Wire Protocol

The wire protocol is a virtual machine in which every encoded value begins with
a 1-byte code that describes what (if anything) follows. The encoding does not
preserve type information--for instance, the value `1` could be an int or a
bool-- but it does have enough information to skip values, since the decoder
must be able to do that if it encounters a struct field it doesn't know.

Most of the values of a value's initial byte can be devoted to small unsigned
integers. For example, the number 17 is represented by the single byte 17. Only
a few byte values have special meaning, as described below.

The `nil` code indicates that the value is nil. (We don't absolutely need this:
we could always represent the nil value for a type as something that couldn't
be mistaken for an encoded value of that type. For instance, we could use 0
for nil in the case of slices (which always begin with the nValues code), and
for pointers to numbers like *int, we could use something like "nBytes 0".
But it is simpler to have a reserved value for nil.)

The `nBytes` code indicates that an unsigned integer N is encoded next, followed
by N bytes of data. There are optimized codes `bytes0`, `bytes1`, etc. for
values of N from 0 to 4. These are used to represent strings and byte slices, as
well numbers bigger than can fit into the initial byte.

The `nValues` code is for sequences of values whose size is known beforehand,
like a Go slice or array.

The `ptr` and `refPtr` codes indicate a pointer to the encoded value. The latter
signals to the decoder that it should remember the pointer because it will be
referred to later in the stream.

The `ref` code is used to refer to an earlier encoded pointer. It is followed by
a uint denoting the relative offset to the position of the corresponding
`refPtr` code.

The `start` and `end` codes delimit a value whose length is unknown beforehand.
They are used for structs.

### Encoding Values

Small unsigned integers are encoded in a single byte, as described above. Those
that can't fit into the initial byte are encoded as byte sequences of length 1,
2, 4 or 8, holding big-endian values. For example, 255 is encoded as `bytes1
255`.

Signed integers are encoded as unsigned integers using zig-zag encoding:
positive numbers are encoded as twice their value, and negative numbers are
encoded as twice their negated value minus 1. This maps small negative values to
small unsigned numbers, since they tend to occur more frequently than large
values of either sign.

A boolean true is encoded as 1, false as 0.

Strings, byte slices and byte arrays are encoded as sequences of bytes. For
example, the string "hello" is represented as `nBytes 5 'h' 'e' 'l' 'l' 'o'`.

Floating-point values are encoded as unsigned integers, after reversing the
bits. Reversing makes small integer-valued floats take less space.

Complex values are encoded as two-element lists of floats, using the `nValues`
code.

Nil values are of course encoded with the `nil` code.

Arrays and slices of type other than `byte` are encoded with `nValues`. For
example, the slice `[]string{"hi", "bye"}` is encoded as
```
nValues 2 bytes2 'h' 'i' bytes3 'b' 'y' 'e'
```

Non-nil pointers are initially encoded with `ptr` followed by the encoding of
the value. For instance, the encoding of `p` in
```
i := 3
p := &i
```
is `ptr 3`. If pointer tracking is enabled and the pointer is encountered again,
then it is encoded with `ref` and the `ptr` code is backpatched to `refPtr`.

Interface values are encoded as a pair of a type number and the value. The
type numbers are assigned during encoding and stored at the beginning of the
output, so the decoder can set up the mapping before it begins.

To encode structs, the generator assigns a unique number to each field. An
encoded struct begins with the `start` code and ends with `end`. Each non-zero
field is encoded as its number followed by its value.

If a struct's fields are changed, the numbers can change. The encoder saves the
numbers assigned to to each field name in the encoded data, and the decoder maps
those numbers to the numbers assigned by the generated code. For example, if `F`
is the second field of a struct, it will be assigned 1. A struct value is
encoded, along with the association of "F" to 1. Then the struct is re-arranged
and `F` becomes the third field, where it is assigned 2. When the encoded data
is decoded, the decoder will map an encoded field value of 1 to 2.

The encoder recognizes types that implement encoding.BinaryMarshaler and
encoding.TextMarshaler, and uses those methods.

## Comparison with Other Encoders

This encoder uses code generation instead of reflection, so it is usually faster
than reflection-based encoders like encoding/gob and encoding/json. It is also
faster than github.com/ugorji/go/codec, even when that uses code generation. See
internal/benchmarks for comparison with gob and ugorji's code on a suite of
benchmarks.

Those benchmarks turn off this codec's ability to handle pointer sharing.
Turning on that feature slows it down noticeably. But the other encoders can't
handle sharing at all.

The gvisor project has their own encoder, which does handle sharing. See
https://pkg.go.dev/gvisor.dev/gvisor/pkg/state. I haven't benchmarked it.
From reading the code:

- It does have some features this encoder lacks, like some custom hooks.

- It doesn't seem to provide for skipping unknown struct fields.

- It seems that only types under programmer control can be encoded, because they
  have to implement certain methods.

- It appears that values must be converted to and from a set of types used by the wire
  protocol (in state/wire). For example, a `[]float64` must be first transformed
  into a `wire.Slice` that holds a `[]wire.Float64` before encoding, and the
  reverse must happen on decoding. It seems this is done cheaply with some
  unsafe code (pkg/state/encode_unsafe.go). But it's not clear why the extra
  level of abstraction is necessary.

- It handles more forms of sharing than this encoder. For example, consider
  a pointer to a struct field like
  ```
  type S struct { X int }
  s := &S{X: 1}
  p := &s.X
  x := []interface{}{s, p}
  ```
  Encoding `x` with gvisor.dev/gvisor/pkg/state will maintain the relationship
  between `p` and `s.X`. This encoder will not; it only recognizes sharing when
  the pointers are explicit, like
  ```
  s := &S{X: 1}
  p := s
  ```
