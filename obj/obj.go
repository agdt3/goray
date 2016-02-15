package obj

import (
	//"fmt"
	"github.com/agdt3/goray/cam"
	"github.com/agdt3/goray/vec"
	"image/color"
	"math"
)

type Object interface {
	GetId() string
	GetColor() color.RGBA
	GetRefractiveIndex() float64
	Intersects(*cam.Ray) (bool, vec.Vec3, vec.Vec3, float64, float64)
}

type Sphere struct {
	Id              string
	Center          vec.Vec3
	Radius          float64
	Col             color.RGBA
	EasingDistance  float64
	RefractiveIndex float64
}

func (s Sphere) Intersects(ray *cam.Ray) (bool, vec.Vec3, vec.Vec3, float64, float64) {
	sc := s.Center
	rd := ray.Direction
	rd.Normalize()

	srsq := s.Radius * s.Radius
	oc := vec.Subtract(sc, ray.Origin)
	l2oc := vec.Dot(oc, oc)
	t_ca := vec.Dot(oc, rd)

	//sphere located behind ray origin
	if t_ca < 0 {
		return false, *vec.NewVec3(0, 0, 0), *vec.NewVec3(0, 0, 0), 0, 0
	}

	d2 := l2oc - (t_ca * t_ca)

	// if the distance between the closest point to the sphere center on
	// the projected ray is greater than the radius, then the projected
	// ray is definitely outside the bounds of the sphere
	if d2 > srsq {
		return false, *vec.NewVec3(0, 0, 0), *vec.NewVec3(0, 0, 0), 0, 0
	}

	t2hc := srsq - d2

	if t2hc < 0 {
		return false, *vec.NewVec3(0, 0, 0), *vec.NewVec3(0, 0, 0), 0, 0
	}

	// If the origin is inside the sphere of light, it counts as a hit
	// This is a useful result if the ray is refracted inside the sphere
	/*
		if l2oc < srsq {
			// TODO: These numbers make no sense
			t0 := float64(0)
			t1 := float64(0)
			hit := ray.Origin
			oc.Multiply(-1)
			return true, hit, oc, t0, t1
		}
	*/

	thc := math.Sqrt(t2hc)
	t0 := t_ca - thc
	t1 := t_ca + thc

	// Sphere is behind the point of origin
	if t0 < 0 && t1 < 0 {
		return false, *vec.NewVec3(0, 0, 0), *vec.NewVec3(0, 0, 0), 0, 0
	} else if t0 <= 0 && t1 > 0 {
		// Point of origin is inside the sphere or on/inside the surface
		t0 = t1
	}

	// Swap if reversed
	if t0 > t1 {
		tmp := t0
		t0 = t1
		t1 = tmp
	}

	dist := s.EasingDistance * t0
	rd.Multiply(dist)
	hit := vec.Add(ray.Origin, rd)

	n := vec.Subtract(hit, sc)
	n.Divide(s.Radius)

	return true, hit, n, t0, t1
}

func (s Sphere) GetColor() color.RGBA {
	return s.Col
}

func (s Sphere) GetRefractiveIndex() float64 {
	return s.RefractiveIndex
}

func (s Sphere) GetId() string {
	return s.Id
}

// TODO: Ironically, Light does not fit the Object interface
type Light struct {
	Type         string //point, directional, etc
	Center       vec.Vec3
	Radius       float64
	RadiusSquare float64
	Col          color.RGBA
}

func (l *Light) Intersects(ray *cam.Ray) (bool, float64) {
	rd := ray.Direction
	rd.Normalize()

	oc := vec.Subtract(l.Center, ray.Origin)
	l2oc := vec.Dot(oc, oc)
	t_ca := vec.Dot(oc, rd)

	//sphere located behind ray origin
	if t_ca < 0 {
		return false, 0
	}

	d2 := l2oc - (t_ca * t_ca)

	// if the distance between the closest point to the sphere center on
	// the projected ray is greater than the radius, then the projected
	// ray is definitely outside the bounds of the sphere
	if d2 > l.RadiusSquare {
		return false, 0
	}

	t2hc := l.RadiusSquare - d2

	if t2hc < 0 {
		return false, 0
	}

	thc := math.Sqrt(t2hc)
	t0 := t_ca - thc
	t1 := t_ca + thc

	// Sphere is behind the point of origin
	if t0 < 0 && t1 < 0 {
		return false, 0
	} else if t0 <= 0 && t1 > 0 {
		// Point of origin is inside the sphere or on/inside the surface
		t0 = t1
	}

	// Swap if reversed
	if t0 > t1 {
		tmp := t0
		t0 = t1
		t1 = tmp
	}

	return true, t0
}

type Triangle struct {
	Id              string
	V0              vec.Vec3
	V1              vec.Vec3
	V2              vec.Vec3
	E0              vec.Vec3
	E1              vec.Vec3
	E2              vec.Vec3
	N               vec.Vec3
	Col             color.RGBA
	EasingDistance  float64
	RefractiveIndex float64
}

func NewTriangle(id string, v0, v1, v2 vec.Vec3, col color.RGBA, easing, refractive float64) *Triangle {
	t := new(Triangle)

	t.Id = id
	t.Col = col
	t.EasingDistance = easing
	t.RefractiveIndex = refractive

	// Verticies
	t.V0 = v0
	t.V1 = v1
	t.V2 = v2

	// Edges
	t.E0 = vec.Subtract(v1, v0)
	t.E1 = vec.Subtract(v2, v1)
	t.E2 = vec.Subtract(v0, v2)

	// Cross product vectors
	v1v0 := vec.Subtract(v1, v0)
	v2v0 := vec.Subtract(v2, v0)
	t.N = vec.Cross(v1v0, v2v0)
	t.N.Normalize()
	return t
}

func (t *Triangle) GetId() string {
	return t.Id
}

func (t *Triangle) GetColor() color.RGBA {
	return t.Col
}

func (t *Triangle) GetRefractiveIndex() float64 {
	return t.RefractiveIndex
}

func (t *Triangle) Intersects(ray *cam.Ray) (bool, vec.Vec3, vec.Vec3, float64, float64) {
	// This is the geometric solution
	TOLERANCE := 0.001
	dir := ray.Direction
	dir.Normalize()

	if math.Abs((vec.Dot(t.N, dir))) < TOLERANCE {
		// Ray direction and N are perpendicular
		// Ray is parallel to plane and will not intersect
		return false, *vec.NewVec3(0, 0, 0), *vec.NewVec3(0, 0, 0), 0, 0
	}

	D := vec.Dot(t.N, t.V0)
	t0 := (vec.Dot(t.N, ray.Origin) + D) / (vec.Dot(t.N, dir))

	if t0 < 0 {
		// Plane is behind ray
		return false, *vec.NewVec3(0, 0, 0), *vec.NewVec3(0, 0, 0), 0, 0
	}

	// Inside-outside test to see if point is inside triangle, not just plane
	// If the vector formed by the dot product of N and the cross product of
	// edge (P,v0) and edge (v1,v0) is positive, P is to the left of
	// edge (v1,v0) - the result of the cross is pointing in the same
	// direction as the triangle's N
	// If all the possible edges formed from all the possible verticies of the
	// triangle and P are to the left of the triangle's edges, then all the
	// edges formed by P are inside the triangle and so is P

	dist := t.EasingDistance * t0
	dir.Multiply(dist)
	P := vec.Add(ray.Origin, dir)

	pe0 := vec.Subtract(P, t.V0)
	pe1 := vec.Subtract(P, t.V1)
	pe2 := vec.Subtract(P, t.V2)
	if (vec.Dot(vec.Cross(t.E0, pe0), t.N) <= 0) ||
		(vec.Dot(vec.Cross(t.E1, pe1), t.N) <= 0) ||
		(vec.Dot(vec.Cross(t.E2, pe2), t.N) <= 0) {

		return false, *vec.NewVec3(0, 0, 0), *vec.NewVec3(0, 0, 0), 0, 0
	}

	return true, P, t.N, t0, t0
}
