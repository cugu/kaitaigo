package runtime

import (
	"encoding/binary"
	"io"
)

type Decoder struct {
	io.ReadSeeker
	ByteOrder binary.ByteOrder
	Err       error
}

type KSYDecoder interface {
	// Decode(io.ReadSeeker) error
	DecodeAncestors(interface{}, interface{}) error
}
