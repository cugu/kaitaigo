package spec

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"

	. "test_formats"
)

func TestExpr0(t *testing.T) {
	f, err := os.Open("../../src/str_encodings.bin")
	if err != nil {
		t.Fatal(err)
	}
	s := kaitai.NewStream(f)

	var r Expr0
	err = r.Read(s, &r, &r)
	if err != nil {
		t.Fatal(err)
	}

	mustBeF7, err := r.MustBeF7()
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, 0xf7, mustBeF7)

	mustBeAbc123, err := r.MustBeAbc123()
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, "abc123", mustBeAbc123)
}
