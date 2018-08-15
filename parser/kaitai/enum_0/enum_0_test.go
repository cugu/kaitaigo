package spec

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnum0(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/enum_0.bin")
	if err != nil {
		t.Fatal(err)
	}

	var h Enum0
	err = h.Read(s, &h, &h)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, Enum0_Animal__Cat, h.Pet1)
	assert.Equal(t, Enum0_Animal__Chicken, h.Pet2)
}
