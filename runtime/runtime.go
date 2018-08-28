package runtime

import (
	"io"
)

type KSYDecoder interface {
	Decode(io.ReadSeeker, interface{}, interface{}) error
}
