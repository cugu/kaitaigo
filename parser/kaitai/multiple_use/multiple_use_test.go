// Autogenerated from KST: please remove this line if doing any edits by hand!

package spec

import (
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultipleUse(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/position_abs.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r MultipleUse
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 32, r.T1.FirstUse.Value)
	tmp1, err := r.T2.SecondUse()
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, 32, tmp1.Value)
}
