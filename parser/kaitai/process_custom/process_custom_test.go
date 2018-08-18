package process_custom

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/cugu/kaitai.go/runtime"
)

func MyCustomFx(data runtime.ByteSlice, key runtime.Uint8, flag bool, someBytes runtime.Bytes) (out runtime.ByteSlice) {
	if !flag {
		key = -key
	}
	for i := 0; i < len(data); i++ {
		data[i] = (data[i] + byte(key))
	}
	return runtime.ByteSlice(data)
}

type CustomStruct struct{}

func (n CustomStruct) CustomFx(data runtime.ByteSlice, i int) (out runtime.ByteSlice) {
	return runtime.ByteSlice("_" + string(data) + "_")
}

type DeeplyStruct struct{}

func (n DeeplyStruct) Deeply() (c CustomStruct) {
	return
}

func Nested() (d DeeplyStruct) {
	return
}

func TestProcessCustom(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/process_rotate.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r ProcessCustom
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, []byte{0x10, 0xb3, 0x94, 0x94, 0xf4}, r.Buf1())
	assert.EqualValues(t, []byte{0x5f, 0xba, 0x7b, 0x93, 0x63, 0x23, 0x5f}, r.Buf2())
	assert.EqualValues(t, []byte{0x29, 0x33, 0xb1, 0x38, 0xb1}, r.Buf3())
}
