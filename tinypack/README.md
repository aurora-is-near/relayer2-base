# TinyPack

TinyPack is a self-made serialization protocol+library designed specially for keeping data as concise as possible.
In contrast to protobuf/cbor/flatbuffers/capnp/etc:
- There's no field tags in resulting data (schema is static, though supports nullable fields).
- No reflection is used.
- Fixed-length arrays are supported.

In order to keep data as short as possible, two tricks are applied:
- Integers are encoded as [varints](https://developers.google.com/protocol-buffers/docs/encoding#varints) - the smaller the value, the less memory it will take.
- Final message goes through [capnp "packing"](https://capnproto.org/encoding.html) encoding - it compresses a byte-array efficiently and very quickly, utilizing the fact that statistically a lot of bytes are actually zeros.

### Supported types

| Type class | Golang type | Requirements | Binary representation |
| ----------- | ----------- | ----------- | ----------- |
| Boolean | `bool` | | `[0;1]`-valued byte |
| Integer | `int64` | | [varint](https://developers.google.com/protocol-buffers/docs/encoding#varints) |
| Unsigned integer | `uint64` | | [varint](https://developers.google.com/protocol-buffers/docs/encoding#varints) |
| Float | `float64` | | Straight binary representation |
| Pointer | `tinypack.Pointer[T]` | `pointer.Ptr` must never be nil | Binary representation of `T` |
| Nullable | `tinypack.Nullable[T]` | `nullable.Ptr` can be nil | `[0;1]`-valued byte, followed by binary representation of `T` |
| Fixed list | `tinypack.List[tinypack.LengthDescriptor, T]` | `list.Content` length must be same length that provided descriptor gives | Ordered concatenation of `T` items |
| Variadic list | `tinypack.VarList[T]` | | Varint-encoded length, followed by ordered concatenation of `T` items |
| Fixed data | `tinypack.Data[tinypack.LengthDescriptor]` | `data.Content` length must be same length that provided descriptor gives | Content-bytes |
| Variadic data | `tinypack.VarData` | | Varint-encoded length, followed by content-bytes |
| Composite | Any type that implements `tinypack.Composite` interface | | Ordered concatenation of provided fields |
| Custom | Any type that implements `tinypack.TinyPackable` interface | | User-implemented |
