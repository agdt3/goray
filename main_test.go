package main

import (
	"fmt"
	"github.com/agdt3/goray/cam"
	"github.com/agdt3/goray/obj"
	"github.com/agdt3/goray/vec"
	"image/color"
	"testing"
)

func TestBlendColors(t *testing.T) {
	t.Parallel()

	c1 := color.RGBA{255, 255, 0, 1}
	c2 := color.RGBA{255, 0, 0, 1}
	c3 := BlendColors(c1, c2, 0.5)
	if c3.R != 255 || c3.G != 180 || c3.B != 0 || c3.A != 1 {
		t.Error("Color blending did not work correctly")
	}
}

func TestCameraRayIntersection1D(t *testing.T) {
	t.Parallel()

	center := vec.NewVec3(0, 0, -3)
	sphere := obj.Sphere{"sphere1", *center, 1, color.RGBA{0, 0, 255, 1}, 1, 1.0}
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
	t.Parallel()

	world := NewWorld()
	center := vec.NewVec3(0, 0, -3)
	sphere := obj.Sphere{"sphere1", *center, 1, color.RGBA{0, 0, 255, 1}, 1, 1.2}
	ray := cam.NewRay("A", "camera", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0, -1))
	_, hit, n, _, _ := sphere.Intersects(ray)

	internal_dir := NewRefractionVector(
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

	external_dir := NewRefractionVector(
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

	internal_dir := NewRefractionVector(
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

	external_dir := NewRefractionVector(
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

	/*
		fmt.Println(isHit2)
		fmt.Printf("hit %v\n", hit2)
		fmt.Printf("n %v\n", n2)
		fmt.Printf("t0 %v\n", t02)
		fmt.Printf("t1 %v\n", t12)
	*/
	//no-op
}
