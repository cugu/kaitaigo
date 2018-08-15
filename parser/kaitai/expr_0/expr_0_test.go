package expr_0

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpr0(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/str_encodings.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r Expr0
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	mustBeF7 := *r.MustBeF7()
	assert.EqualValues(t, 0xf7, mustBeF7)

	mustBeAbc123 := *r.MustBeAbc123()
	assert.EqualValues(t, "abc123", mustBeAbc123)
}
