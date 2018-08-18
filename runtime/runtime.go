package runtime

import (
	"bytes"
	"compress/zlib"
	"io"
	"io/ioutil"
	"log"
	"math/bits"
	"runtime/debug"
)

type Decoder struct {
	io.ReadSeeker
	err error
}

func New(reader io.ReadSeeker) *Decoder {
	return &Decoder{reader, nil}
}

func (d *Decoder) Err() (err error) {
	return d.err
}

func (d *Decoder) SetErr(err error) {
	if err != nil {
		log.Println(err)
		debug.PrintStack()
		d.err = err
	}
}

func (d *Decoder) UnsetErr() {
	d.err = nil
}

type KSYDecoder interface {
	DecodeAncestors(interface{}, interface{})
}

// ProcessXOR returns data xored with the key.
func ProcessXOR(data []byte, key []byte) []byte {
	out := make([]byte, len(data))
	for i := range data {
		out[i] = data[i] ^ key[i%len(key)]
	}
	return out
}

// ProcessRotateLeft returns the single bytes in data rotated left by
// amount bits.
func ProcessRotateLeft(data []byte, amount int) []byte {
	out := make([]byte, len(data))
	for i := range data {
		out[i] = bits.RotateLeft8(data[i], amount)
	}
	return out
}

// ProcessRotateRight returns the single bytes in data rotated right by
// amount bits.
func ProcessRotateRight(data []byte, amount int) []byte {
	return ProcessRotateLeft(data, -amount)
}

// ProcessZlib decompresses the given bytes as specified in RFC 1950.
func ProcessZlib(in []byte) (out []byte, err error) {
	b := bytes.NewReader([]uint8(in))

	// FIXME zlib.NewReader allocates a bunch of memory.  In the future
	// we could reuse it by using a sync.Pool if this is called in a tight loop.
	r, err := zlib.NewReader(b)
	if err != nil {
		return
	}

	return ioutil.ReadAll(r)
}
