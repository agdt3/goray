package vec

import (
	"testing"
)

func TestNewVec3(t *testing.T) {
	t.Parallel()
	v3 := NewVec3(1, 4, 8)

	if v3.X != 1 || v3.Y != 4 || v3.Z != 8 {
		t.Error("Vector did not initialize correctly")
	}

	if v3.Magnitude != 9 {
		t.Error("Vector magnitude was not calculated correctly")
	}
}

func TestDot(t *testing.T) {
	t.Parallel()
	v1 := NewVec3(1, 2, 3)
	v2 := NewVec3(4, 5, 6)

	dp := Dot(*v1, *v2)
	if dp != 32 {
		t.Error("Vector dot product is incorrect")
	}
}

func TestCross(t *testing.T) {
	t.Parallel()
	i := NewVec3(1, 0, 0)
	j := NewVec3(0, 1, 0)
	k := NewVec3(0, 0, 1)

	v1 := Cross(*i, *j)
	if v1.X != 0 || v1.Y != 0 || v1.Z != 1 {
		t.Error("Vector cross product is incorrect")
	}

	v2 := Cross(*j, *i)
	if v2.X != 0 || v2.Y != 0 || v2.Z != -1 {
		t.Error("Vector cross product is incorrect")
	}

	v3 := Cross(*j, *k)
	if v3.X != 1 || v3.Y != 0 || v3.Z != 0 {
		t.Error("Vector cross product is incorrect")
	}

	v4 := Cross(*k, *j)
	if v4.X != -1 || v4.Y != 0 || v4.Z != 0 {
		t.Error("Vector cross product is incorrect")
	}

	v5 := Cross(*i, *k)
	if v5.X != 0 || v5.Y != -1 || v5.Z != 0 {
		t.Error("Vector cross product is incorrect")
	}

	v6 := Cross(*k, *i)
	if v6.X != 0 || v6.Y != 1 || v6.Z != 0 {
		t.Error("Vector cross product is incorrect")
	}

}
