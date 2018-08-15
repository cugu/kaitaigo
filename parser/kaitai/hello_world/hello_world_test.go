package hello_world

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloWorld(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai//fixed_struct.bin")
	if err != nil {
		t.Fatal(err)
	}
	var r HelloWorld
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 80, *r.One())
}
