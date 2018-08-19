package enum_fancy

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnumFancy(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/enum_0.bin")
	if err != nil {
		t.Fatal(err)
	}

	var h EnumFancy
	err = h.Read(s, &h, &h)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, EnumFancy_Animal__Cat, h.Pet1())
	assert.Equal(t, EnumFancy_Animal__Chicken, h.Pet2())
}
