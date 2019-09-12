<h1 align="center">kaitaigo</h1>

<p align="center">kaitaigo is a compiler and runtime to create Go parsers from <a href="http://kaitai.io/">Kaitai Struct</a> files.</p>

<p  align="center">
<a href="https://github.com/cugu/kaitaigo"><img src="https://img.shields.io/azure-devops/build/cugu/dfir/1" alt="build" /></a>
<a href="https://codecov.io/gh/cugu/kaitaigo"><img src="https://codecov.io/gh/cugu/kaitaigo/branch/master/graph/badge.svg" alt="coverage" /></a>
<a href="https://goreportcard.com/report/dfir.software/kaitaigo"><img src="https://goreportcard.com/badge/dfir.software/kaitaigo" alt="report" /></a>
</p>

## Installation

```sh
go get dfir.software/kaitaigo
```

## Usage

First we need a .ksy file. We take this simple example, but there are [many more](http://formats.kaitai.io/).

```yaml
# my_format.ksy
meta:
  id: my_format
seq:
  - id: data_size
    type: u1
  - id: data
    size: data_size
```

To create the Go code we use the kaitaigo command: `kaitaigo my_format.ksy`. This creates the ready to use `my_format.ksy.go`.

The parser can be used in other scripts like the following. Change package in `my_format.ksy.go` to main. Afterward you can run the script and use our new parser with `go run main.go my_format.ksy.go`.

```go
// main.go
package main

import (
	"bytes"
	"log"
)

func main() {
	f := bytes.NewReader([]byte("\x05Hello world!"))
	var r MyFormat
	err := r.Decode(f) // Decode takes any io.ReadSeeker
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(r.Data())) // Prints "Hello"
}
```

---

## kaitaigo features

This is not the official kaitai compiler. The official kaitai compiler contains go support as well.
More information in the [issue for Go language support](https://github.com/kaitai-io/kaitai_struct/issues/146).

### Supported kaitai features:

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

### Additional features:

#### whence

Can be used togher with `pos` the define the reference point of the position. Valid values are `seek_set`, `seek_end` and `seek_cur` (default).

### Limitations

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

## Licenses

The kaitaigo compiler is licensed as [GPLv3](licences/gpl-3.0.txt).
The runtime (/runtime) is licensed under [MIT license](licences/mit.txt).
Everything else is licensed as [GPLv3](licences/gpl-3.0.txt).
