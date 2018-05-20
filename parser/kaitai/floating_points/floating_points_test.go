package spec

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"

	. "test_formats"
)

func TestFloatingPoints(t *testing.T) {
	f, err := os.Open("../../src/floating_points.bin")
	if err != nil {
		t.Fatal(err)
	}
	s := kaitai.NewStream(f)

	var h FloatingPoints
	err = h.Read(s, &h, &h)
	if err != nil {
		t.Fatal(err)
	}

	delta := 1e-6

	assert.Equal(t, float32(0.5), h.SingleValue, "They should be equal")
	assert.Equal(t, float32(0.5), h.SingleValueBe, "They should be equal")
	assert.Equal(t, 0.25, h.DoubleValue, "They should be equal")
	assert.Equal(t, 0.25, h.DoubleValueBe, "They should be equal")

	assert.InDelta(t, 1.2345, h.ApproximateValue, delta)

	singleValuePlusInt, err := h.SingleValuePlusInt()
	if err != nil {
		t.Fatal(err)
	}
	assert.InDelta(t, 1.5, singleValuePlusInt, delta)

	singleValuePlusFloat, err := h.SingleValuePlusFloat()
	if err != nil {
		t.Fatal(err)
	}
	assert.InDelta(t, 1.0, singleValuePlusFloat, delta)

	doubleValuePlusFloat, err := h.DoubleValuePlusFloat()
	if err != nil {
		t.Fatal(err)
	}
	assert.InDelta(t, 0.3, doubleValuePlusFloat, delta)
}
