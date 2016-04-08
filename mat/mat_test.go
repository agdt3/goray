package mat

import (
	"testing"
)

func TestNewMatFloat64(t *testing.T) {
	t.Parallel()

	m := NewMatrixFloat64(2, 3, []float64{1, 2, 3, 4, 5, 6})

	if m.R != 2 || m.C != 3 {
		t.Error("Row or Columen count is incorrect")
	}

	if m.Values[1] != 2 || m.Values[4] != 5 {
		t.Error("Values not set correctly")
	}

	if m.Determinant != 0 {
		t.Error("Default determinant not set")
	}
}

func TestNewMatFloat64_2x2(t *testing.T) {
	t.Parallel()

	m := NewMatrixFloat64(2, 2, []float64{1, 2, 3, 4})

	if m.Determinant != -2 {
		t.Error("Determinant is incorrect")
	}
}

func TestNewMatFloat64_3x3(t *testing.T) {
	t.Parallel()

	m := NewMatrixFloat64(
		3,
		3,
		[]float64{
			2, 2, 3,
			4, 5, 6,
			7, 8, 9})

	if m.Determinant != -3 {
		t.Error("Determinant is incorrect")
	}
}
