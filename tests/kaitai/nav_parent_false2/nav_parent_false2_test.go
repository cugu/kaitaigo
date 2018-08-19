package nav_parent_false2

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNavParentFalse2(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/fixed_struct.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r NavParentFalse2
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 80, r.Parentless().Foo())
}
