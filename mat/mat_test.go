package mat

import (
	"testing"

	"github.com/agdt3/goray/vec"
)

type SetupType struct {
	m22i *MatrixFloat64
	m33i *MatrixFloat64
	m22  *MatrixFloat64
	m23  *MatrixFloat64
	m32  *MatrixFloat64
	m33  *MatrixFloat64
}

func setupMatricies() *SetupType {
	setupObj := new(SetupType)

	// 2x2
	setupObj.m22i = NewMatrixFloat64(2, 2, []float64{1, 0, 0, 1})
	setupObj.m22 = NewMatrixFloat64(2, 2, []float64{1, 2, 3, 4})

	// 2x3
	setupObj.m23 = NewMatrixFloat64(2, 3, []float64{1, 2, 3, 4, 5, 6})

	// 2x3
	setupObj.m32 = NewMatrixFloat64(3, 2, []float64{1, 2, 3, 4, 5, 6})

	// 3x3
	setupObj.m33i = NewMatrixFloat64(
		3,
		3,
		[]float64{
			1, 0, 0,
			0, 1, 0,
			0, 0, 1})

	setupObj.m33 = NewMatrixFloat64(
		3,
		3,
		[]float64{
			1, 2, 3,
			4, 5, 6,
			7, 8, 9})

	return setupObj
}

func TestNewMatFloat64(t *testing.T) {
	t.Parallel()

	m := *(setupMatricies()).m23

	if m.R != 2 || m.C != 3 {
		t.Error("Row or Columen count is incorrect")
	}

	v1, _ := m.Get(0, 1)
	v2, _ := m.Get(1, 0)
	if v1 != 2 || v2 != 4 {
		t.Error("Values not set correctly")
	}

	if m.Determinant != 0 {
		t.Error("Default determinant not set")
	}
}

func TestIsEqual(t *testing.T) {
	t.Parallel()

	m1 := (setupMatricies()).m33
	m2 := (setupMatricies()).m33
	m3 := (setupMatricies()).m33
	m3.Set(0, 0, 3)

	if !IsEqual(m1, m2, 0.1) {
		t.Error("m1 and m2 should be equal")
	}

	if IsEqual(m1, m3, 0.1) {
		t.Error("m1 and m3 should not be equal")
	}
}

func TestNewMatFloat64_2x2(t *testing.T) {
	t.Parallel()

	m := *(setupMatricies()).m22

	if m.Determinant != -2 {
		t.Error("Determinant is incorrect")
	}
}

func TestNewMatFloat64_3x3(t *testing.T) {
	t.Parallel()

	m := *(setupMatricies()).m33
	m.Set(0, 0, 2)

	if m.Determinant != -3 {
		t.Error("Determinant is incorrect")
	}
}

func TestMatFloat64Get(t *testing.T) {
	t.Parallel()

	m := *(setupMatricies()).m33

	v1, _ := m.Get(0, 2)
	v2, _ := m.Get(1, 1)
	v3, _ := m.Get(2, 0)

	if v1 != 3 || v2 != 5 || v3 != 7 {
		t.Error("Did not get the correct values")
	}

	v4, err := m.Get(m.R, m.C)
	if err == nil && v4 != 0 {
		t.Error("Did not calculate out of bounds correctly")
	}
}

func TestMatFloat64Set(t *testing.T) {
	t.Parallel()

	m := *(setupMatricies()).m33

	m.Set(0, 2, 11.5)
	m.Set(1, 1, 12.3)
	m.Set(2, 0, 13.6)

	v1, _ := m.Get(0, 2)
	v2, _ := m.Get(1, 1)
	v3, _ := m.Get(2, 0)

	if v1 != 11.5 || v2 != 12.3 || v3 != 13.6 {
		t.Error("Did not set or get the correct values")
	}

	v4, err := m.Get(m.R, m.C)
	if err == nil && v4 != 0 {
		t.Error("Did not calculate out of bounds correctly")
	}
}

func TestMatFloat64Add(t *testing.T) {
	t.Parallel()

	m1 := (setupMatricies()).m33
	m2 := (setupMatricies()).m33i
	m1.Add(m2)

	v1, _ := m1.Get(0, 0)
	v2, _ := m1.Get(1, 1)
	v3, _ := m1.Get(2, 2)

	if v1 != 2 || v2 != 6 || v3 != 10 {
		t.Error("Did not set or get the correct values")
	}
}

func TestMatFloat64Subtract(t *testing.T) {
	t.Parallel()

	m1 := (setupMatricies()).m33
	m2 := (setupMatricies()).m33i
	m1.Subtract(m2)

	v1, _ := m1.Get(0, 0)
	v2, _ := m1.Get(1, 1)
	v3, _ := m1.Get(2, 2)

	if v1 != 0 || v2 != 4 || v3 != 8 {
		t.Error("Did not set or get the correct values")
	}
}

func TestMatFloat64MultiplyScalar(t *testing.T) {
	t.Parallel()

	m1 := *(setupMatricies()).m33
	m1.MultiplyScalar(3)

	v1, _ := m1.Get(0, 0)
	v2, _ := m1.Get(1, 1)
	v3, _ := m1.Get(2, 2)

	if v1 != 3 || v2 != 15 || v3 != 27 {
		t.Error("Matrix not multiplied correctly")
	}
}

func TestMatFloat64DivideScalar(t *testing.T) {
	t.Parallel()

	m1 := *(setupMatricies()).m33
	m1.DivideScalar(3)

	v1, _ := m1.Get(0, 2)
	v2, _ := m1.Get(1, 2)
	v3, _ := m1.Get(2, 2)

	if v1 != 1 || v2 != 2 || v3 != 3 {
		t.Error("Matrix not divided correctly")
	}
}

func TestMatFloat64MultiplyVec3(t *testing.T) {
	t.Parallel()

	m1 := (setupMatricies()).m33
	v1 := vec.NewVec3(1, 2, 3)
	v2, err := MultiplyVec3(m1, v1)

	if err != nil || v2.X != 14 || v2.Y != 32 || v2.Z != 50 {
		t.Error("Vector multiplication was incorrect")
	}
}

func TestMatFloat64MultiplyIdentity(t *testing.T) {
	t.Parallel()

	m1 := (setupMatricies()).m33
	m2 := (setupMatricies()).m33i
	m3, _ := Multiply(m1, m2)

	if !IsEqual(m1, m3, 0.1) {
		t.Error("m1 and m3 should be equal")
	}
}

func TestMatFloat64Multiply33By33(t *testing.T) {
	t.Parallel()

	m1 := (setupMatricies()).m33
	m2 := (setupMatricies()).m33
	m3, _ := Multiply(m1, m2)
	m4 := NewMatrixFloat64(3, 3, []float64{30, 36, 42, 66, 81, 96, 102, 126, 150})

	if !IsEqual(m3, m4, 0.1) {
		t.Error("m3 and m4 should be equal")
	}
}

func TestMatFloat64Multiply23By32(t *testing.T) {
	t.Parallel()

	m1 := (setupMatricies()).m23
	m2 := (setupMatricies()).m32
	m3, _ := Multiply(m1, m2)
	m4 := NewMatrixFloat64(2, 2, []float64{22, 28, 49, 64})

	if !IsEqual(m3, m4, 0.1) {
		t.Error("m3 and m4 should be equal")
	}
}
