# kaitai.go

kaitai.go is an [alternative](https://github.com/kaitai-io/kaitai_struct_compiler) compiler, that converts [Kaitai Struct](http://kaitai.io/) files into [Go](https://golang.org/) parsers.


## Supported kaitai features:

 - Type specification
    - meta
        - endianess*
    - doc
    - seq
    - instances
    - enums
 - Attribute specification
    - id
    - doc
    - contents
    - repeat, repeat-expr, repeat-until
    - if
    - size, size-eos
    - process
    - terminator
    - consume
    - include
    - pad
    - eos-error
 - Primitive data types
 - Processing specification
    - xor
    - rol
    - ror
    - zlib
 - Instance specification
    - pos
    - value

_*partially_

## Additional features:

### whence

Can be used togher with `pos` the define the reference point of the position. Valid values are `seek_set`, `seek_end` and `seek_cur` (default).

## Limitations

 - No _io (Most uses can be replaced with [whence](#whence))
 - Accessing nested types with `::` is not allowed
 - No fancy enums
 - No nested endianess
 - No encoding
 - No comparison of string, []byte or custom types
 - No min or max functions
 - fix type inference
 - -2 % 8 = -2
 - xor, ror, rol and zlib only work on bytes
 - float + int fails
