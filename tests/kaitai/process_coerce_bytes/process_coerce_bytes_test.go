// Autogenerated from KST: please remove this line if doing any edits by hand!

package process_coerce_bytes

import (
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessCoerceBytes(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/process_coerce_bytes.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r ProcessCoerceBytes
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	records := r.Records()
	assert.EqualValues(t, 0, records[0].Flag())
	assert.EqualValues(t, []uint8{65, 65, 65, 65}, records[0].Buf())
	assert.EqualValues(t, 1, records[1].Flag())
	assert.EqualValues(t, []uint8{66, 66, 66, 66}, records[1].Buf())
}