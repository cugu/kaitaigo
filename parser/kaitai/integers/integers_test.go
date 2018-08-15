// Autogenerated from KST: please remove this line if doing any edits by hand!

package spec

import (
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegers(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/fixed_struct.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r Integers
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 255, r.Uint8)
	assert.EqualValues(t, 65535, r.Uint16)
	assert.EqualValues(t, uint32(4294967295), r.Uint32)
	assert.EqualValues(t, uint64(18446744073709551615), r.Uint64)
	assert.EqualValues(t, -1, r.Sint8)
	assert.EqualValues(t, -1, r.Sint16)
	assert.EqualValues(t, -1, r.Sint32)
	assert.EqualValues(t, -1, r.Sint64)
	assert.EqualValues(t, 66, r.Uint16le)
	assert.EqualValues(t, 66, r.Uint32le)
	assert.EqualValues(t, 66, r.Uint64le)
	assert.EqualValues(t, -66, r.Sint16le)
	assert.EqualValues(t, -66, r.Sint32le)
	assert.EqualValues(t, -66, r.Sint64le)
	assert.EqualValues(t, 66, r.Uint16be)
	assert.EqualValues(t, 66, r.Uint32be)
	assert.EqualValues(t, 66, r.Uint64be)
	assert.EqualValues(t, -66, r.Sint16be)
	assert.EqualValues(t, -66, r.Sint32be)
	assert.EqualValues(t, -66, r.Sint64be)
}
