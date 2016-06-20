package vec

import (
	"fmt"
	"math"
)

// Vec3 is an R3 vector. Used to represent vectors and points.
type Vec3 struct {
	X                float64
	Y                float64
	Z                float64
	Magnitude        float64
	magnitudeInverse float64 // 1/Mag precomputed
}

// NewVec3 is a Vec3 constructor.
func NewVec3(x, y, z float64) *Vec3 {
	v := new(Vec3)
	v.X = x
	v.Y = y
	v.Z = z
	v.Magnitude, v.magnitudeInverse = v.CalculateMagnitude()

	return v
}

// Add adds one vector to calling vector
func (v *Vec3) Add(v2 *Vec3) {
	v.X += v2.X
	v.Y += v2.Y
	v.Z += v2.Z
}

// Subtract subtracts one vector from calling vector
func (v *Vec3) Subtract(v2 *Vec3) {
	v.X -= v2.X
	v.Y -= v2.Y
	v.Z -= v2.Z
}

// Multiply multiplies calling vector by a scalar
func (v *Vec3) Multiply(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
	v.Z *= scalar
}

// Divide divides calling vector by a scalar
func (v *Vec3) Divide(scalar float64) {
	v.X /= scalar
	v.Y /= scalar
	v.Z /= scalar
}

// Normalize normalizes the calling vector
func (v *Vec3) Normalize() {
	v.X = v.X * v.magnitudeInverse
	v.Y = v.Y * v.magnitudeInverse
	v.Z = v.Z * v.magnitudeInverse
	v.Magnitude = 1
	v.magnitudeInverse = 1 // Precomputed inverse to save a division
}

// CalculateMagnitude returns the vector's magnitude and inverse magnitude
func (v Vec3) CalculateMagnitude() (float64, float64) {
	// Recalculates magnitude
	// returns magnitude and inverse magnitude
	mag := math.Sqrt(
		(v.X * v.X) +
			(v.Y * v.Y) +
			(v.Z * v.Z))

	// prevent division by zero
	var invMag float64
	if mag != 0 {
		invMag = 1 / mag
	} else {
		invMag = 0
	}

	return mag, invMag
}

// String represents the calling vector as a string
func (v Vec3) String() string {
	return fmt.Sprintf(
		"Vec3: XYZ(%v, %v, %v) - Magnitude(%v)",
		v.X, v.Y, v.Z, v.Magnitude)
}

/* Utils */

// Add sums two vectors together and returns a third
func Add(v1, v2 Vec3) Vec3 {
	v3 := NewVec3(v1.X+v2.X, v1.Y+v2.Y, v1.Z+v2.Z)
	return *v3
}

// Subtract returns the difference between the first and second vector
// as a third vector
func Subtract(v1, v2 Vec3) Vec3 {
	v3 := NewVec3(v1.X-v2.X, v1.Y-v2.Y, v1.Z-v2.Z)
	return *v3
}

// Multiply multiplies a vector by a scalar and returns
// a new vector as a result
func Multiply(v1 Vec3, scalar float64) Vec3 {
	v2 := NewVec3(v1.X*scalar, v1.Y*scalar, v1.Z*scalar)
	return *v2
}

// Divide divides a vector by a scalar and returns
// a new vector as a result
func Divide(v1 Vec3, scalar float64) Vec3 {
	v2 := NewVec3(v1.X/scalar, v1.Y/scalar, v1.Z/scalar)
	return *v2
}

// Dot gets the dot product of two vectors and returns a scalar
func Dot(v1, v2 Vec3) float64 {
	return (v1.X * v2.X) + (v1.Y * v2.Y) + (v1.Z * v2.Z)
}

// Cross gets the cross product of two vectors and returns a third vector
func Cross(v1, v2 Vec3) Vec3 {
	xc := (v1.Y * v2.Z) - (v1.Z * v2.Y)
	yc := (v1.Z * v2.X) - (v1.X * v2.Z)
	zc := (v1.X * v2.Y) - (v1.Y * v2.X)
	v3 := NewVec3(xc, yc, zc)
	return *v3
}

// Invert changes the direction of a vector
func Invert(v1 Vec3) Vec3 {
	v2 := NewVec3(-1*v1.X, -1*v1.Y, -1*v1.Z)
	return *v2
}

// IsEqual compares two vectors on the bases of every property
func IsEqual(v1, v2 Vec3) bool {
	if v1.X != v2.X || v1.Y != v2.Y || v1.Z != v2.Z {
		return false
	}

	if v1.Magnitude != v2.Magnitude {
		return false
	}

	return true
}

// Reflect reflects an incedence vector about the n vector
func Reflect(i, n Vec3) Vec3 {
	// i - 2 * dot(n, i) * n
	return Subtract(i, Multiply(n, (Dot(n, i)*2.0)))
}
