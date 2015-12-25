package vec

import (
	"math"
)

type Vec3 struct {
	X                float64
	Y                float64
	Z                float64
	Magnitude        float64
	magnitudeInverse float64 // 1/Mag precomputed
}

func MakeVec3(x, y, z float64) *Vec3 {
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

	return mag, 1 / mag
}

/* Utils */
func Add(v1, v2 Vec3) Vec3 {
	v3 := MakeVec3(v1.X+v2.X, v1.Y+v2.Y, v1.Z+v2.Z)
	return *v3
}

func Subtract(v1, v2 Vec3) Vec3 {
	v3 := MakeVec3(v1.X-v2.X, v1.Y-v2.Y, v1.Z-v2.Z)
	return *v3
}

func Dot(v1, v2 Vec3) float64 {
	return (v1.X * v2.X) + (v1.Y * v2.Y) + (v1.Z * v2.Z)
}

func Cross(v1, v2 Vec3) Vec3 {
	xc := (v1.Y * v2.Z) - (v1.Z * v2.Y)
	yc := (v1.Z * v2.X) - (v1.X * v2.Z)
	zc := (v1.X * v2.Y) - (v1.Y * v2.X)
	v3 := MakeVec3(xc, yc, zc)
	return *v3
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
