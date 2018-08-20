// Autogenerated from KST: please remove this line if doing any edits by hand!

package bcd_user_type_be

import (
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcdUserTypeBe(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/bcd_user_type_be.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r BcdUserTypeBe
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 12345678, r.Ltr().AsInt())
	assert.EqualValues(t, "12345678", r.Ltr().AsStr())
	assert.EqualValues(t, 87654321, r.Rtl().AsInt())
	assert.EqualValues(t, "87654321", r.Rtl().AsStr())
	assert.EqualValues(t, 123456, r.LeadingZeroLtr().AsInt())
	assert.EqualValues(t, "00123456", r.LeadingZeroLtr().AsStr())
}