package vec

import (
	"fmt"
	"math"
)

type Vec3 struct {
	X                float64
	Y                float64
	Z                float64
	Magnitude        float64
	magnitudeInverse float64 // 1/Mag precomputed
}

func NewVec3(x, y, z float64) *Vec3 {
	vec3 := new(Vec3)
	vec3.X = x
	vec3.Y = y
	vec3.Z = z
	vec3.Magnitude, vec3.magnitudeInverse = vec3.CalculateMagnitude()

	return vec3
}

func (v *Vec3) Add(v2 Vec3) {
	v.X += v2.X
	v.Y += v2.Y
	v.Z += v2.Z
}

func (v *Vec3) Subtract(v2 Vec3) {
	v.X -= v2.X
	v.Y -= v2.Y
	v.Z -= v2.Z
}

// Scalar multiplication/division only
func (v *Vec3) Multiply(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
	v.Z *= scalar
}

func (v *Vec3) Divide(scalar float64) {
	v.X /= scalar
	v.Y /= scalar
	v.Z /= scalar
}

func (v *Vec3) Normalize() {
	v.X = v.X * v.magnitudeInverse
	v.Y = v.Y * v.magnitudeInverse
	v.Z = v.Z * v.magnitudeInverse
	v.Magnitude = 1
	v.magnitudeInverse = 1 // Precomputed inverse to save a division
}

func (v Vec3) CalculateMagnitude() (float64, float64) {
	// Recalculates magnitude
	// returns magnitude and inverse magnitude
	mag := math.Sqrt(
		(v.X * v.X) +
			(v.Y * v.Y) +
			(v.Z * v.Z))

	// prevent division by zero
	var inv_mag float64
	if mag != 0 {
		inv_mag = 1 / mag
	} else {
		inv_mag = 0
	}

	return mag, inv_mag
}

func (v Vec3) String() string {
	return fmt.Sprintf(
		"Vec3: XYZ(%v, %v, %v) - Magnitude(%v)",
		v.X, v.Y, v.Z, v.Magnitude)
}

/* Utils */
func Add(v1, v2 Vec3) Vec3 {
	v3 := NewVec3(v1.X+v2.X, v1.Y+v2.Y, v1.Z+v2.Z)
	return *v3
}

func Subtract(v1, v2 Vec3) Vec3 {
	v3 := NewVec3(v1.X-v2.X, v1.Y-v2.Y, v1.Z-v2.Z)
	return *v3
}

func Multiply(v1 Vec3, scalar float64) Vec3 {
	v2 := NewVec3(v1.X*scalar, v1.Y*scalar, v1.Z*scalar)
	return *v2
}

func Divide(v1 Vec3, scalar float64) Vec3 {
	v2 := NewVec3(v1.X/scalar, v1.Y/scalar, v1.Z/scalar)
	return *v2
}

func Dot(v1, v2 Vec3) float64 {
	return (v1.X * v2.X) + (v1.Y * v2.Y) + (v1.Z * v2.Z)
}

func Cross(v1, v2 Vec3) Vec3 {
	xc := (v1.Y * v2.Z) - (v1.Z * v2.Y)
	yc := (v1.Z * v2.X) - (v1.X * v2.Z)
	zc := (v1.X * v2.Y) - (v1.Y * v2.X)
	v3 := NewVec3(xc, yc, zc)
	return *v3
}

func Invert(v1 Vec3) Vec3 {
	v2 := NewVec3(-1*v1.X, -1*v1.Y, -1*v1.Z)
	return *v2
}

func IsEqual(v1, v2 Vec3) bool {
	if v1.X != v2.X || v1.Y != v2.Y || v1.Z != v2.Z {
		return false
	}

	if v1.Magnitude != v2.Magnitude {
		return false
	}

	return true
}
