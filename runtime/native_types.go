package runtime

import (
	"encoding/binary"
	"io"
	"io/ioutil"
)

type Byte byte

func (v *Byte) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(&Decoder{reader, binary.LittleEndian, nil}, v, v)
}
func (v *Byte) DecodePos(dec *Decoder, offset int64, whence int, parent interface{}, root interface{}) (err error) {
	if dec.Err != nil {
		return dec.Err
	}
	_, dec.Err = dec.Seek(offset, whence)
	return v.DecodeAncestors(dec, parent, root)
}

func (v *Byte) DecodeAncestors(dec *Decoder, parent interface{}, root interface{}) (err error) {
	dec.Err = binary.Read(dec, dec.ByteOrder, v)
	return
}

type Bytes []byte

func (v *Bytes) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(&Decoder{reader, binary.LittleEndian, nil}, v, v)
}
func (v *Bytes) DecodePos(dec *Decoder, offset int64, whence int, parent interface{}, root interface{}) (err error) {
	if dec.Err != nil {
		return dec.Err
	}
	_, dec.Err = dec.Seek(offset, whence)
	return v.DecodeAncestors(dec, parent, root)
}

func (v *Bytes) DecodeAncestors(dec *Decoder, parent interface{}, root interface{}) (err error) {
	lv, err := ioutil.ReadAll(dec)
	if err != nil {
		dec.Err = err
		return err
	}
	*v = lv
	return
}

type ByteArray struct {
	Dec    *Decoder
	Start  int64
	Parent interface{}
	Root   interface{}

	Size int64
	Data []byte
}

func (v *ByteArray) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(&Decoder{reader, binary.LittleEndian, nil}, v, v)
}
func (v *ByteArray) DecodePos(dec *Decoder, offset int64, whence int, parent interface{}, root interface{}) (err error) {
	if dec.Err != nil {
		return dec.Err
	}
	_, dec.Err = dec.Seek(offset, whence)
	return v.DecodeAncestors(dec, parent, root)
}

func (v *ByteArray) DecodeAncestors(dec *Decoder, parent interface{}, root interface{}) (err error) {
	v.Data = make([]byte, v.Size)
	_, err = dec.Read(v.Data)
	if err != nil {
		dec.Err = err
	}
	return
}

type Uint8 uint8

func (v *Uint8) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(&Decoder{reader, binary.LittleEndian, nil}, v, v)
}
func (v *Uint8) DecodePos(dec *Decoder, offset int64, whence int, parent interface{}, root interface{}) (err error) {
	if dec.Err != nil {
		return dec.Err
	}
	_, dec.Err = dec.Seek(offset, whence)
	return v.DecodeAncestors(dec, parent, root)
}

func (v *Uint8) DecodeAncestors(dec *Decoder, parent interface{}, root interface{}) (err error) {
	dec.Err = binary.Read(dec, dec.ByteOrder, v)
	return
}

type Uint16 uint16

func (v *Uint16) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(&Decoder{reader, binary.LittleEndian, nil}, v, v)
}
func (v *Uint16) DecodePos(dec *Decoder, offset int64, whence int, parent interface{}, root interface{}) (err error) {
	if dec.Err != nil {
		return dec.Err
	}
	_, dec.Err = dec.Seek(offset, whence)
	return v.DecodeAncestors(dec, parent, root)
}

func (v *Uint16) DecodeAncestors(dec *Decoder, parent interface{}, root interface{}) (err error) {
	dec.Err = binary.Read(dec, dec.ByteOrder, v)
	return
}

type Uint32 uint32

func (v *Uint32) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(&Decoder{reader, binary.LittleEndian, nil}, v, v)
}
func (v *Uint32) DecodePos(dec *Decoder, offset int64, whence int, parent interface{}, root interface{}) (err error) {
	if dec.Err != nil {
		return dec.Err
	}
	_, dec.Err = dec.Seek(offset, whence)
	return v.DecodeAncestors(dec, parent, root)
}

func (v *Uint32) DecodeAncestors(dec *Decoder, parent interface{}, root interface{}) (err error) {
	dec.Err = binary.Read(dec, dec.ByteOrder, v)
	return
}

type Uint64 uint64

func (v *Uint64) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(&Decoder{reader, binary.LittleEndian, nil}, v, v)
}
func (v *Uint64) DecodePos(dec *Decoder, offset int64, whence int, parent interface{}, root interface{}) (err error) {
	if dec.Err != nil {
		return dec.Err
	}
	_, dec.Err = dec.Seek(offset, whence)
	return v.DecodeAncestors(dec, parent, root)
}

func (v *Uint64) DecodeAncestors(dec *Decoder, parent interface{}, root interface{}) (err error) {
	dec.Err = binary.Read(dec, dec.ByteOrder, v)
	return
}

type Int8 int8

func (v *Int8) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(&Decoder{reader, binary.LittleEndian, nil}, v, v)
}
func (v *Int8) DecodePos(dec *Decoder, offset int64, whence int, parent interface{}, root interface{}) (err error) {
	if dec.Err != nil {
		return dec.Err
	}
	_, dec.Err = dec.Seek(offset, whence)
	return v.DecodeAncestors(dec, parent, root)
}

func (v *Int8) DecodeAncestors(dec *Decoder, parent interface{}, root interface{}) (err error) {
	dec.Err = binary.Read(dec, dec.ByteOrder, v)
	return
}

type Int16 int16

func (v *Int16) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(&Decoder{reader, binary.LittleEndian, nil}, v, v)
}
func (v *Int16) DecodePos(dec *Decoder, offset int64, whence int, parent interface{}, root interface{}) (err error) {
	if dec.Err != nil {
		return dec.Err
	}
	_, dec.Err = dec.Seek(offset, whence)
	return v.DecodeAncestors(dec, parent, root)
}

func (v *Int16) DecodeAncestors(dec *Decoder, parent interface{}, root interface{}) (err error) {
	dec.Err = binary.Read(dec, dec.ByteOrder, v)
	return
}

type Int32 int32

func (v *Int32) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(&Decoder{reader, binary.LittleEndian, nil}, v, v)
}
func (v *Int32) DecodePos(dec *Decoder, offset int64, whence int, parent interface{}, root interface{}) (err error) {
	if dec.Err != nil {
		return dec.Err
	}
	_, dec.Err = dec.Seek(offset, whence)
	return v.DecodeAncestors(dec, parent, root)
}

func (v *Int32) DecodeAncestors(dec *Decoder, parent interface{}, root interface{}) (err error) {
	dec.Err = binary.Read(dec, dec.ByteOrder, v)
	return
}

type Int64 int64

func (v *Int64) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(&Decoder{reader, binary.LittleEndian, nil}, v, v)
}
func (v *Int64) DecodePos(dec *Decoder, offset int64, whence int, parent interface{}, root interface{}) (err error) {
	if dec.Err != nil {
		return dec.Err
	}
	_, dec.Err = dec.Seek(offset, whence)
	return v.DecodeAncestors(dec, parent, root)
}

func (v *Int64) DecodeAncestors(dec *Decoder, parent interface{}, root interface{}) (err error) {
	dec.Err = binary.Read(dec, dec.ByteOrder, v)
	return
}
