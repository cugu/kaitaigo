package buffered_struct

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBufferedStruct(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/buffered_struct.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r BufferedStruct
	err = r.Decode(f)
	if err != nil {
		// t.Fatal(err)
	}

	assert.EqualValues(t, 0x10, r.Len1())
	assert.EqualValues(t, 0x42, r.Block1().Number1())
	assert.EqualValues(t, 0x43, r.Block1().Number2())
	assert.EqualValues(t, 0x8, r.Len2())
	assert.EqualValues(t, 0x44, r.Block2().Number1())
	assert.EqualValues(t, 0x45, r.Block2().Number2())
	assert.EqualValues(t, 0xee, r.Finisher())
}
