// Autogenerated from KST: please remove this line if doing any edits by hand!

package spec

import (
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeatEosStruct(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/repeat_eos_struct.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r RepeatEosStruct
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 2, len(r.Chunks))
	assert.EqualValues(t, 0, r.Chunks[0].Offset)
	assert.EqualValues(t, 66, r.Chunks[0].Len)
	assert.EqualValues(t, 66, r.Chunks[1].Offset)
	assert.EqualValues(t, 2069, r.Chunks[1].Len)
}
