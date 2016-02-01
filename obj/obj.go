package obj

import (
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
