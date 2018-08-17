package default_big_endian

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultBigEndian(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/enum_0.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r DefaultBigEndian
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 0x7000000, r.One())
}
