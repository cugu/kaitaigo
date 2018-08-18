package runtime

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
	"io/ioutil"
	"math/bits"
)

type Decoder struct {
	io.ReadSeeker
	ByteOrder binary.ByteOrder
	Err       error
}

type KSYDecoder interface {
	DecodeAncestors(interface{}, interface{})
}

// ProcessXOR returns data xored with the key.
func ProcessXOR(data Bytes, key Bytes) Bytes {
	out := make(Bytes, len(data))
	for i := range data {
		out[i] = data[i] ^ key[i%len(key)]
	}
	return out
}

// ProcessRotateLeft returns the single bytes in data rotated left by
// amount bits.
func ProcessRotateLeft(data ByteSlice, amount int) ByteSlice {
	out := make(ByteSlice, len(data))
	for i := range data {
		out[i] = bits.RotateLeft8(data[i], amount)
	}
	return out
}

// ProcessRotateRight returns the single bytes in data rotated right by
// amount bits.
func ProcessRotateRight(data ByteSlice, amount int) ByteSlice {
	return ProcessRotateLeft(data, -amount)
}

// ProcessZlib decompresses the given bytes as specified in RFC 1950.
func ProcessZlib(in Bytes) (out Bytes) {
	b := bytes.NewReader([]uint8(in))

	// FIXME zlib.NewReader allocates a bunch of memory.  In the future
	// we could reuse it by using a sync.Pool if this is called in a tight loop.
	r, err := zlib.NewReader(b)
	if err != nil {
		RTDecoder.Err = err
		return
	}

	out, err = ioutil.ReadAll(r)
	if err != nil {
		RTDecoder.Err = err
		return
	}
	return Bytes(out)
}
