package obj

import (
	"image/color"
	"math"
	"testing"

	"github.com/agdt3/goray/cam"
	"github.com/agdt3/goray/vec"
)

// TODO: There is a better method than epsilon testing
func AlmostEqual(v1, v2, tolerance float64) bool {
	if math.Abs(v1)-math.Abs(v2) < tolerance {
		return true
	}
	return false
}

/*
func TestCameraRayIntersection1D(t *testing.T) {
	t.Skip()
	t.Parallel()

	center := vec.NewVec3(0, 0, -3)
	sphere := Sphere{"sphere1", *center, 1, color.RGBA{0, 0, 255, 1}, 1, 1.0}
	ray := cam.NewRay("A", "camera", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0, -1))
	isHit, hit, n, t0, t1 := sphere.Intersects(ray)

	if !isHit {
		t.Error("Ray did not hit")
	}

	if !vec.IsEqual(hit, *vec.NewVec3(0, 0, -2)) {
		t.Error("Hit vector is not correct")
	}

	if !vec.IsEqual(n, *vec.NewVec3(0, 0, 1)) {
		t.Error("N vector is not correct")
	}

	if t0 != 2 || t1 != 4 {
		t.Error("Distance is incorrect")
	}
}

func TestTransmissionRayIntersection1DHeadOn(t *testing.T) {
	t.Skip()
	t.Parallel()

	center := vec.NewVec3(0, 0, -3)
	sphere := Sphere{"sphere1", *center, 1, color.RGBA{0, 0, 255, 1}, 1, 1.2}
	ray := cam.NewRay("A", "camera", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0, -1))
	_, hit, n, _, _ := sphere.Intersects(ray)

	internal_dir := RefractionVector(
		ray.Direction,
		n,
		world.RefractiveIndex,
		sphere.GetRefractiveIndex())

	if !vec.IsEqual(internal_dir, *vec.NewVec3(0, 0, -1)) {
		t.Error("Internal refracted direction is incorrect")
	}

	ref_ray := cam.NewRay("B", "refraction", &hit, &internal_dir)
	isHit2, hit2, n2, _, _ := sphere.Intersects(ref_ray)

	if !isHit2 {
		t.Error("Internal ray should intersec with sphere")
	}

	if !vec.IsEqual(hit2, *vec.NewVec3(0, 0, -4)) {
		t.Error("Internal ray did not intersect at the right location")
	}

	external_dir := RefractionVector(
		ref_ray.Direction,
		vec.Invert(n2),
		sphere.GetRefractiveIndex(),
		world.RefractiveIndex)

	if !vec.IsEqual(external_dir, *vec.NewVec3(0, 0, -1)) {
		t.Error("External refracted direction is incorrect")
	}

	trans_ray := cam.NewRay("C", "transmission", &hit2, &external_dir)
	world_trans_ray, _ := world.NewTransmittedRay(ray, hit, n, sphere)

	if !cam.IsEqual(trans_ray, world_trans_ray) {
		t.Error("World function incorrectly constructed transmission ray")
	}
}

func TestTransmissionRayIntersection1DAtAngle(t *testing.T) {
	t.Skip()
	//t.Parallel()

	world := NewWorld()
	center := vec.NewVec3(0.5, 0, -3)
	sphere := obj.Sphere{"sphere1", *center, 1, color.RGBA{0, 0, 255, 1}, 1, 1.2}
	ray := cam.NewRay("A", "camera", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0, -1))
	_, hit, n, _, _ := sphere.Intersects(ray)

	internal_dir := RefractionVector(
		ray.Direction,
		n,
		world.RefractiveIndex,
		sphere.GetRefractiveIndex())

	if !vec.IsEqual(internal_dir, *vec.NewVec3(0, 0, -1)) {
		t.Error("Internal refracted direction is incorrect")
	}

	ref_ray := cam.NewRay("B", "refraction", &hit, &internal_dir)
	isHit2, hit2, n2, _, _ := sphere.Intersects(ref_ray)

	if !isHit2 {
		t.Error("Internal ray should intersec with sphere")
	}

	if !vec.IsEqual(hit2, *vec.NewVec3(0, 0, -4)) {
		t.Error("Internal ray did not intersect at the right location")
	}

	external_dir := RefractionVector(
		ref_ray.Direction,
		vec.Invert(n2),
		sphere.GetRefractiveIndex(),
		world.RefractiveIndex)

	if !vec.IsEqual(external_dir, *vec.NewVec3(0, 0, -1)) {
		t.Error("External refracted direction is incorrect")
	}

	trans_ray := cam.Ray{"C", "transmitted", hit2, external_dir}
	world_trans_ray, _ := world.NewTransmittedRay(ray, hit, n, sphere)
	fmt.Println(ray)
	fmt.Println(ref_ray)
	fmt.Println(trans_ray)
	fmt.Println(world_trans_ray)

	//no-op
}
*/

func TestLightIntersection1D(t *testing.T) {
	t.Parallel()

	ray := cam.NewRay("noid", "shadow", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0, -1))
	center := vec.NewVec3(0, 0, -3)
	light := Light{"light1", *center, 1, 1, color.RGBA{255, 255, 255, 1}}

	hit, dist := light.Intersects(ray)
	if !hit || dist != 2.0 {
		t.Error("Shadow ray did not intersect with light")
	}
}

func TestLightIntersection2D(t *testing.T) {
	t.Parallel()

	dir := vec.NewVec3(0, 1, -1)
	dir.Normalize()
	ray := cam.NewRay("noid", "shadow", vec.NewVec3(0, 0, 0), dir)
	center := vec.NewVec3(0, 2, -3)
	light := Light{"light1", *center, 1, 1, color.RGBA{255, 255, 255, 1}}

	hit, dist := light.Intersects(ray)
	if !hit {
		t.Error("Shadow ray did not intersect with light")
	}

	if !AlmostEqual(dist, 2.828, 0.001) {
		t.Error("Shadow ray did not hit in the correct location")
	}
	/*
		if !hit || dist != 2.828427124746191 {
			t.Error("Shadow ray did not intersect with light")
		}
	*/
}

func TestLightIntersectionByReflectedRay(t *testing.T) {
	t.Parallel()

	dir := vec.NewVec3(0, 1, -1)
	dir.Normalize()
	ray := cam.NewRay("noid", "camera", vec.NewVec3(0, 0, 0), dir)
	center1 := vec.NewVec3(0, 5, 0)
	center2 := vec.NewVec3(0, 2, -3)
	light := Light{"light1", *center1, 1, 1, color.RGBA{255, 255, 255, 1}}
	sphere := Sphere{"sphere1", *center2, 1, color.RGBA{255, 255, 255, 1}, 1, 1}

	is_hit, hit, n, t0, _ := sphere.Intersects(ray)
	if !is_hit || !AlmostEqual(t0, 2.828, 0.001) {
		t.Error("Shadow ray did not intersect with light")
	}

	reflected_dir := vec.Reflect(hit, n)
	shadow_ray := cam.NewRay("noid", "shadow", &hit, &reflected_dir)
	is_hit2, dist := light.Intersects(shadow_ray)

	if !is_hit2 || !AlmostEqual(dist, 2.828, 0.001) {
		t.Error("Shadow ray did not intersect with light")
	}
}

func TestTriangleInit(t *testing.T) {
	// Points v1 and v2 are created CCW relative to v0
	// (in a left-handed system where camera points at -z)
	// In such a case, the N vector points towards the viewer
	v0 := vec.NewVec3(0, 0, -1)
	v1 := vec.NewVec3(1, 1, -1)
	v2 := vec.NewVec3(-1, 1, -1)
	tri := NewTriangle("tri1", *v0, *v1, *v2, color.RGBA{0, 0, 0, 1}, 1, 1, false)

	e0 := vec.Subtract(*v1, *v0)
	e1 := vec.Subtract(*v2, *v1)
	e2 := vec.Subtract(*v0, *v2)

	n := vec.NewVec3(0, 0, 1)

	if tri.E0 != e0 {
		t.Error("Edges not properly calculated")
	}

	if tri.E1 != e1 {
		t.Error("Edges not properly calculated")
	}

	if tri.E2 != e2 {
		t.Error("Edges not properly calculated")
	}

	if tri.N != *n {
		t.Error("N vector not properly calculated")
	}
}

func TestTriangleIntersects(t *testing.T) {
	v0 := vec.NewVec3(0, -1, -1)
	v1 := vec.NewVec3(1, 1, -1)
	v2 := vec.NewVec3(-1, 1, -1)
	tri := NewTriangle("tri1", *v0, *v1, *v2, color.RGBA{0, 0, 0, 1}, 1, 1, false)

	ray := cam.Ray{"", "camera", *vec.NewVec3(0, 0, 0), *vec.NewVec3(0, 0, -1)}

	is_hit, p, n, t0, _ := tri.Intersects(&ray)

	i_p := *vec.NewVec3(0, 0, -1)
	i_n := *vec.NewVec3(0, 0, 1)
	i_t0 := 1.0

	if is_hit != true {
		t.Error("Triangle was not hit")
	}

	if !vec.IsEqual(p, i_p) {
		t.Error("Hit location was not correct")
	}

	if !vec.IsEqual(n, i_n) {
		t.Error("N vector not correct")
	}

	if t0 != i_t0 {
		t.Error("Distance incorrect")
	}
}
