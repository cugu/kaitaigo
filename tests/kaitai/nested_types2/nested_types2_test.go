// Autogenerated from KST: please remove this line if doing any edits by hand!

package nested_types2

import (
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNestedTypes2(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/fixed_struct.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r NestedTypes2
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 80, r.One().TypedAtRoot().ValueB())
	assert.EqualValues(t, 65, r.One().TypedHere1().ValueC())
	assert.EqualValues(t, 67, r.One().TypedHere1().TypedHere().ValueD())
	assert.EqualValues(t, 75, r.One().TypedHere1().TypedParent().ValueCc())
	assert.EqualValues(t, 45, r.One().TypedHere1().TypedRoot().ValueB())
	assert.EqualValues(t, 49, r.One().TypedHere2().ValueCc())
	assert.EqualValues(t, -1, r.Two().ValueB())
}