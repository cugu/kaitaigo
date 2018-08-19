package floating_points

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloatingPoints(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/floating_points.bin")
	if err != nil {
		t.Fatal(err)
	}

	var h FloatingPoints
	err = h.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	delta := 1e-6

	assert.Equal(t, float32(0.5), h.SingleValue(), "They should be equal")
	assert.Equal(t, float32(0.5), h.SingleValueBe(), "They should be equal")
	assert.Equal(t, float64(0.25), h.DoubleValue(), "They should be equal")
	assert.Equal(t, float64(0.25), h.DoubleValueBe(), "They should be equal")

	assert.InDelta(t, float32(1.2345), h.ApproximateValue(), delta)

	singleValuePlusInt := h.SingleValuePlusInt()
	assert.InDelta(t, float32(1.5), singleValuePlusInt, delta)

	singleValuePlusFloat := h.SingleValuePlusFloat()
	assert.InDelta(t, float32(1.0), singleValuePlusFloat, delta)

	doubleValuePlusFloat := h.DoubleValuePlusFloat()
	assert.InDelta(t, float32(0.3), doubleValuePlusFloat, delta)

}
