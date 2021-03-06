package obj

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/agdt3/goray/cam"
	"github.com/agdt3/goray/vec"
)

// Object is the base interface to a number of other
// objects including Sphere and Triangle
type Object interface {
	GetID() string
	GetColor() color.RGBA
	GetRefractiveIndex() float64
	Intersects(*cam.Ray) (bool, vec.Vec3, vec.Vec3, float64, float64)
}

// FalseObject returns a failure-state object
func FalseObject() (bool, vec.Vec3, vec.Vec3, float64, float64) {
	return false, *vec.NewVec3(0, 0, 0), *vec.NewVec3(0, 0, 0), 0, 0
}

// Sphere object
type Sphere struct {
	ID              string
	Center          vec.Vec3
	Radius          float64
	Col             color.RGBA
	EasingDistance  float64
	RefractiveIndex float64
}

// Intersects checks for intersections with sphere using
// geometric method
func (s Sphere) Intersects(ray *cam.Ray) (bool, vec.Vec3, vec.Vec3, float64, float64) {
	sc := s.Center
	rd := ray.Direction
	rd.Normalize()

	srsq := s.Radius * s.Radius
	oc := vec.Subtract(sc, ray.Origin)
	l2oc := vec.Dot(oc, oc)
	tCa := vec.Dot(oc, rd)

	//sphere located behind ray origin
	if tCa < 0 {
		return false, *vec.NewVec3(0, 0, 0), *vec.NewVec3(0, 0, 0), 0, 0
	}

	d2 := l2oc - (tCa * tCa)

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
	t0 := tCa - thc
	t1 := tCa + thc

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

// GetColor is the object specific method to return the color
// as color.RGBA
func (s Sphere) GetColor() color.RGBA {
	return s.Col
}

// GetRefractiveIndex is the object specific method to return the
// refractive index
func (s Sphere) GetRefractiveIndex() float64 {
	return s.RefractiveIndex
}

// GetID is the object specific method to return the ID of the sphere
func (s Sphere) GetID() string {
	return s.ID
}

// Light is a basic spherical light
// TODO: Ironically, Light does not fit the Object interface
type Light struct {
	Type         string //point, directional, etc
	Center       vec.Vec3
	Radius       float64
	RadiusSquare float64
	Col          color.RGBA
}

// Intersects checks for intersections between cam.Ray and Sphere-like
// light, using the geometric method
func (l *Light) Intersects(ray *cam.Ray) (bool, float64) {
	rd := ray.Direction
	rd.Normalize()

	oc := vec.Subtract(l.Center, ray.Origin)
	l2oc := vec.Dot(oc, oc)
	tCa := vec.Dot(oc, rd)

	//sphere located behind ray origin
	if tCa < 0 {
		return false, 0
	}

	d2 := l2oc - (tCa * tCa)

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
	t0 := tCa - thc
	t1 := tCa + thc

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

// Triangle is fundamental container for the most basic renderable shape
type Triangle struct {
	ID              string
	V0              vec.Vec3
	V1              vec.Vec3
	V2              vec.Vec3
	E0              vec.Vec3
	E1              vec.Vec3
	E2              vec.Vec3
	v0v1            vec.Vec3
	v0v2            vec.Vec3
	N               vec.Vec3
	Col             color.RGBA
	EasingDistance  float64
	RefractiveIndex float64
	Culling         bool
}

// NewTriangle is a constructor for Triangles and precomputed values
func NewTriangle(id string, v0, v1, v2 vec.Vec3, col color.RGBA, easing, refractive float64, culling bool) *Triangle {
	t := new(Triangle)

	t.ID = id
	t.Col = col
	t.EasingDistance = easing
	t.RefractiveIndex = refractive
	t.Culling = culling

	// Verticies
	t.V0 = v0
	t.V1 = v1
	t.V2 = v2

	// Edges
	t.E0 = vec.Subtract(v1, v0)
	t.E1 = vec.Subtract(v2, v1)
	t.E2 = vec.Subtract(v0, v2)

	// Cross product vectors
	t.v0v1 = vec.Subtract(v1, v0)
	t.v0v2 = vec.Subtract(v2, v0)
	t.N = vec.Cross(t.v0v1, t.v0v2)

	// May not need to normalize this
	t.N.Normalize()
	return t
}

// GetID is the object specific method to return the ID of the sphere
func (t *Triangle) GetID() string {
	return t.ID
}

// GetColor is the object specific method to return the color
// as color.RGBA
func (t *Triangle) GetColor() color.RGBA {
	return t.Col
}

// GetRefractiveIndex is the object specific method to return the
// refractive index
func (t *Triangle) GetRefractiveIndex() float64 {
	return t.RefractiveIndex
}

// IntersectsImplicit checks for intersections between a ray the triangle
// using the implicit method
// TODO: Dead code
func (t *Triangle) IntersectsImplicit(ray *cam.Ray) (bool, vec.Vec3, vec.Vec3, float64, float64) {
	// This is the geometric solution
	// Plane intersection first
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

// IntersectsBarycentric checks for intersections between a ray the triangle
// using the barycentric method
// TODO: Dead code
func (t *Triangle) IntersectsBarycentric(ray *cam.Ray) (bool, vec.Vec3, vec.Vec3, float64, float64) {
	TOLERANCE := 0.001
	dir := ray.Direction
	dir.Normalize()

	if math.Abs((vec.Dot(t.N, dir))) < TOLERANCE {
		// Ray direction and N are perpendicular
		// Ray is parallel to plane and will not intersect
		return FalseObject()
	}

	D := vec.Dot(t.N, t.V0)
	t0 := (vec.Dot(t.N, ray.Origin) + D) / (vec.Dot(t.N, dir))

	if t0 < 0 {
		// Plane is behind ray
		return FalseObject()
	}

	dist := t.EasingDistance * t0
	dir.Multiply(dist)
	P := vec.Add(ray.Origin, dir)

	/* Barycentric */
	v0p := vec.Subtract(P, t.V0) // v0 to P
	v1p := vec.Subtract(P, t.V1) // v1 to P
	v2p := vec.Subtract(P, t.V2) // v2 to P

	denominator := vec.Dot(t.N, t.N)
	// Can compare parallelogram areas instead of triangle areas

	// Edge 0
	c := vec.Cross(t.E0, v0p)
	if vec.Dot(c, t.N) < 0 {
		return FalseObject()
	}

	// Edge 1
	c = vec.Cross(t.E1, v1p)
	u := vec.Dot(c, t.N)
	if u < 0 {
		return FalseObject()
	}

	// Edge 2
	c = vec.Cross(t.E2, v2p)
	v := vec.Dot(c, t.N)
	if v < 0 {
		return FalseObject()
	}

	u /= denominator
	v /= denominator

	return true, P, t.N, t0, t0
}

// Intersects checks for intersections between a ray the triangle
// using the Trombole-Muller method
func (t *Triangle) Intersects(ray *cam.Ray) (bool, vec.Vec3, vec.Vec3, float64, float64) {
	TOLERANCE := 0.001
	dir := ray.Direction
	dir.Normalize()

	/* Trombole-Muller */
	pvec := vec.Cross(dir, t.v0v2)
	denominator := vec.Dot(pvec, t.v0v1)
	invDenominator := 1 / denominator

	if t.Culling && (denominator < TOLERANCE) {
		return FalseObject()
	} else if math.Abs(denominator) < TOLERANCE {
		return FalseObject()
	}

	tvec := vec.Subtract(ray.Origin, t.V0)
	qvec := vec.Cross(tvec, t.v0v1)

	u := vec.Dot(pvec, tvec) * invDenominator
	if u < 0 || u > 1 {
		return FalseObject()
	}

	v := vec.Dot(qvec, dir) * invDenominator
	if v < 0 || u+v > 1 {
		return FalseObject()
	}

	// (t0, u, v) as opposed to (t, u, v)
	t0 := vec.Dot(qvec, t.v0v2) * invDenominator

	dist := t.EasingDistance * t0
	dir.Multiply(dist)
	P := vec.Add(ray.Origin, dir)

	// Note that may have dist < t0
	return true, P, t.N, t0, t0
}

// PolygonMesh is a container for mesh polygon data
// and is not a true object
type PolygonMesh struct {
	NumFaces         []int
	NumVerticies     []int
	VertexIndecies   []int
	Verticies        []float64
	VertexNormals    []float64
	TextureVertecies []float64
}

// MakePolygonMesh sets up mesh container
func MakePolygonMesh() *PolygonMesh {
	poly := new(PolygonMesh)
	poly.NumFaces = make([]int, 1, 1)

	return poly
}

// ConvertPolygonSerial converts data in PolygonMesh
// into an array of triangles using a serial strategy
func (p *PolygonMesh) ConvertPolygonSerial() []Triangle {
	start := time.Now()

	var vertexIndex int
	var numTriangles int

	totalTriangles := 0
	for _, v := range p.NumVerticies {
		totalTriangles += (v - 2)
	}

	triangles := make([]Triangle, totalTriangles, totalTriangles)

	indexIntoVertexIndex := 0
	triangleIndex := 0
	for i := 0; i < p.NumFaces[0]; i++ {
		numTriangles = p.NumVerticies[i] - 2
		for j := 0; j < numTriangles; j++ {
			vertexIndex = p.VertexIndecies[indexIntoVertexIndex] * 3
			v0 := vec.NewVec3(
				p.Verticies[vertexIndex],
				p.Verticies[vertexIndex+1],
				p.Verticies[vertexIndex+2])

			vertexIndex = p.VertexIndecies[indexIntoVertexIndex+1] * 3
			v1 := vec.NewVec3(
				p.Verticies[vertexIndex],
				p.Verticies[vertexIndex+1],
				p.Verticies[vertexIndex+2])

			vertexIndex = p.VertexIndecies[indexIntoVertexIndex+2] * 3
			v2 := vec.NewVec3(
				p.Verticies[vertexIndex],
				p.Verticies[vertexIndex+1],
				p.Verticies[vertexIndex+2])

			// TODO: Remove test color
			red := uint8(0)
			green := uint8(0)
			blue := uint8(0)
			if i == 0 || i == 3 {
				red = uint8(255 / (j + 1))
			}
			if i == 1 || i == 4 {
				green = uint8(255 / (j + 1))
			}
			if i == 2 || i == 5 {
				blue = uint8(255 / (j + 1))
			}

			triangle := NewTriangle("", *v0, *v1, *v2, color.RGBA{red, green, blue, 1}, 1, 1, false)
			triangles[triangleIndex] = *triangle

			// Increment our various indecies
			indexIntoVertexIndex += 3
			triangleIndex++
		}
	}
	elapsed := time.Since(start)
	fmt.Println(elapsed)
	return triangles
}

// ConvertPolygonParallel converts data in PolygonMesh
// into an array of triangles using a parallel strategy
func (p *PolygonMesh) ConvertPolygonParallel() []Triangle {
	start := time.Now()

	totalTriangles := 0
	for _, v := range p.NumVerticies {
		totalTriangles += (v - 2)
	}

	// Create slice
	triangles := make([]Triangle, 0, totalTriangles)

	// Create channel queue
	trianglesChan := make(chan *Triangle, totalTriangles)

	indexIntoVertexIndex := 0
	for i := 0; i < p.NumFaces[0]; i++ {
		numTriangles := p.NumVerticies[i] - 2
		for j := 0; j < numTriangles; j++ {

			go p.generateTriangle(i, j, indexIntoVertexIndex, trianglesChan)

			// Increment our various indecies
			indexIntoVertexIndex += 3
		}
	}

	for k := 0; k < totalTriangles; k++ {
		triangles = append(triangles, *<-trianglesChan)
	}
	//fmt.Println(triangles)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
	return triangles
}

func (p *PolygonMesh) generateTriangle(i, j, indexIntoVertexIndex int, trianglesChan chan *Triangle) {
	var vertexIndex int
	vertexIndex = p.VertexIndecies[indexIntoVertexIndex] * 3
	v0 := vec.NewVec3(
		p.Verticies[vertexIndex],
		p.Verticies[vertexIndex+1],
		p.Verticies[vertexIndex+2])

	vertexIndex = p.VertexIndecies[indexIntoVertexIndex+1] * 3
	v1 := vec.NewVec3(
		p.Verticies[vertexIndex],
		p.Verticies[vertexIndex+1],
		p.Verticies[vertexIndex+2])

	vertexIndex = p.VertexIndecies[indexIntoVertexIndex+2] * 3
	v2 := vec.NewVec3(
		p.Verticies[vertexIndex],
		p.Verticies[vertexIndex+1],
		p.Verticies[vertexIndex+2])

	// TODO: Remove test color
	red := uint8(0)
	green := uint8(0)
	blue := uint8(0)
	if i == 0 || i == 3 {
		red = uint8(255 / (j + 1))
	}
	if i == 1 || i == 4 {
		green = uint8(255 / (j + 1))
	}
	if i == 2 || i == 5 {
		blue = uint8(255 / (j + 1))
	}

	trianglesChan <- NewTriangle("", *v0, *v1, *v2, color.RGBA{red, green, blue, 1}, 1, 1, false)
}

// String stringifies triangles
func (t *Triangle) String() string {
	return fmt.Sprintf(
		"V0: %v, %v, %v - V1: %v, %v, %v - V2: %v, %v, %v\nN: %v, %v, %v\n",
		t.V0.X, t.V0.Y, t.V0.Z,
		t.V1.X, t.V1.Y, t.V1.Z,
		t.V2.X, t.V2.Y, t.V2.Z,
		t.N.X, t.N.Y, t.N.Z)
}
