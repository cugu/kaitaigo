// Autogenerated from KST: please remove this line if doing any edits by hand!

package zlib_with_header_78

import (
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZlibWithHeader78(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/zlib_with_header_78.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r ZlibWithHeader78
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, []uint8{97, 32, 113, 117, 105, 99, 107, 32, 98, 114, 111, 119, 110, 32, 102, 111, 120, 32, 106, 117, 109, 112, 115, 32, 111, 118, 101, 114}, r.Data())
}
