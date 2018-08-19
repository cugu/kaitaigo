package runtime

import (
	"io"
	"log"
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
