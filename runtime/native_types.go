package runtime

import (
	"bufio"
	"encoding/binary"
	"io"
	"io/ioutil"
)

var RTDecoder *Decoder

type Byte byte

func (v *Byte) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
}

type Bytes []byte

func (v *Bytes) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	b, err := ioutil.ReadAll(RTDecoder)
	if err != nil {
		RTDecoder.Err = err
		return
	}
	*v = b
}

type BytesZ []byte

func (v *BytesZ) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	pos, err := RTDecoder.Seek(0, io.SeekCurrent)
	if err != nil {
		RTDecoder.Err = err
		return
	}
	b, err := bufio.NewReader(RTDecoder).ReadBytes(byte(0))
	if err != nil && err != io.EOF {
		RTDecoder.Err = err
		return
	}
	_, err = RTDecoder.Seek(pos+int64(len(b)), io.SeekStart)
	if err != nil {
		RTDecoder.Err = err
		return
	}
	if len(b) != 0 {
		*v = b[:len(b)-1]
	}
}

type ByteSlice []byte

func (v *ByteSlice) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	_, err := RTDecoder.Read(*v)
	if err != nil {
		RTDecoder.Err = err
	}
}

/*
type String ByteSlice //string

/*
func (v *String) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	b := make([]byte, len(*v))
	_, err := RTDecoder.Read(b)
	if err != nil {
		RTDecoder.Err = err
	}
	tmp := String(b)
	v = &tmp
}
*/

type Uint8 uint8

func (v *Uint8) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
}

type Uint16 uint16

func (v *Uint16) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
}

type Uint32 uint32

func (v *Uint32) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
}

type Uint64 uint64

func (v *Uint64) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
}

type Sint8 Int8
type Sint16 Int16
type Sint32 Int32
type Sint64 Int64

type Int8 int8

func (v *Int8) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
}

type Int16 int16

func (v *Int16) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
}

type Int32 int32

func (v *Int32) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
}

type Int64 int64

func (v *Int64) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
}

// Little endian

type Uint8Be uint8

func (v *Uint8Be) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, binary.BigEndian, v)
}

type Uint16Be uint16

func (v *Uint16Be) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, binary.BigEndian, v)
}

type Uint32Be uint32

func (v *Uint32Be) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, binary.BigEndian, v)
}

type Uint64Be uint64

func (v *Uint64Be) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, binary.BigEndian, v)
}

type Sint8Be Int8Be
type Sint16Be Int16Be
type Sint32Be Int32Be
type Sint64Be Int64Be

type Int8Be int8

func (v *Int8Be) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, binary.BigEndian, v)
}

type Int16Be int16

func (v *Int16Be) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, binary.BigEndian, v)
}

type Int32Be int32

func (v *Int32Be) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, binary.BigEndian, v)
}

type Int64Be int64

func (v *Int64Be) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, binary.BigEndian, v)
}

// Little endian

type Uint8Le uint8

func (v *Uint8Le) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, binary.LittleEndian, v)
}

type Uint16Le uint16

func (v *Uint16Le) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, binary.LittleEndian, v)
}

type Uint32Le uint32

func (v *Uint32Le) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, binary.LittleEndian, v)
}

type Uint64Le uint64

func (v *Uint64Le) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, binary.LittleEndian, v)
}

type Sint8Le Int8Le
type Sint16Le Int16Le
type Sint32Le Int32Le
type Sint64Le Int64Le

type Int8Le int8

func (v *Int8Le) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, binary.LittleEndian, v)
}

type Int16Le int16

func (v *Int16Le) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, binary.LittleEndian, v)
}

type Int32Le int32

func (v *Int32Le) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, binary.LittleEndian, v)
}

type Int64Le int64

func (v *Int64Le) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, binary.LittleEndian, v)
}
