package spec

import (
	"os"
	"testing"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
	"github.com/stretchr/testify/assert"

	. "test_formats"
)

func TestProcessCustom(t *testing.T) {
	f, err := os.Open("../../src/process_rotate.bin")
	if err != nil {
		t.Fatal(err)
	}
	s := kaitai.NewStream(f)

	var r ProcessCustom
	err = r.Read(s, &r, &r)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, []byte{0x10, 0xb3, 0x94, 0x94, 0xf4}, r.Buf1)
	assert.EqualValues(t, []byte{0x5f, 0xba, 0x7b, 0x93, 0x63, 0x23, 0x5f}, r.Buf2)
	assert.EqualValues(t, []byte{0x29, 0x33, 0xb1, 0x38, 0xb1}, r.Buf3)
}
