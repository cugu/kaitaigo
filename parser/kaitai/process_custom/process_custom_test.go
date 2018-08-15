package spec

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	assert.EqualValues(t, []byte{0x10, 0xb3, 0x94, 0x94, 0xf4}, r.Buf1)
	assert.EqualValues(t, []byte{0x5f, 0xba, 0x7b, 0x93, 0x63, 0x23, 0x5f}, r.Buf2)
	assert.EqualValues(t, []byte{0x29, 0x33, 0xb1, 0x38, 0xb1}, r.Buf3)
}
