// Autogenerated from KST: please remove this line if doing any edits by hand!

package spec

import (
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessXor4Value(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/process_xor_4.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r ProcessXor4Value
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, []uint8{236, 187, 163, 20}, r.Key)
	assert.EqualValues(t, []uint8{102, 111, 111, 32, 98, 97, 114}, r.Buf)
}
