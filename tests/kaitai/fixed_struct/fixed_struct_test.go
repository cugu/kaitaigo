// Autogenerated from KST: please remove this line if doing any edits by hand!

package fixed_struct

import (
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixedStruct(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/fixed_struct.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r FixedStruct
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 255, r.Hdr().Uint8())
	assert.EqualValues(t, 65535, r.Hdr().Uint16())
	assert.EqualValues(t, uint32(4294967295), r.Hdr().Uint32())
	assert.EqualValues(t, uint64(18446744073709551615), r.Hdr().Uint64())
	assert.EqualValues(t, -1, r.Hdr().Sint8())
	assert.EqualValues(t, -1, r.Hdr().Sint16())
	assert.EqualValues(t, -1, r.Hdr().Sint32())
	assert.EqualValues(t, -1, r.Hdr().Sint64())
	assert.EqualValues(t, 66, r.Hdr().Uint16Le())
	assert.EqualValues(t, 66, r.Hdr().Uint32Le())
	assert.EqualValues(t, 66, r.Hdr().Uint64Le())
	assert.EqualValues(t, -66, r.Hdr().Sint16Le())
	assert.EqualValues(t, -66, r.Hdr().Sint32Le())
	assert.EqualValues(t, -66, r.Hdr().Sint64Le())
	assert.EqualValues(t, 66, r.Hdr().Uint16Be())
	assert.EqualValues(t, 66, r.Hdr().Uint32Be())
	assert.EqualValues(t, 66, r.Hdr().Uint64Be())
	assert.EqualValues(t, -66, r.Hdr().Sint16Be())
	assert.EqualValues(t, -66, r.Hdr().Sint32Be())
	assert.EqualValues(t, -66, r.Hdr().Sint64Be())
}