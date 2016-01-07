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

	center := vec.MakeVec3(0, 0, -3)
	sphere := obj.Sphere{*center, 1, color.RGBA{0, 0, 255, 1}, 1, 1.0}
	ray := cam.Ray{"camera", *vec.MakeVec3(0, 0, 0), *vec.MakeVec3(0, 0, -1)}
	isHit, hit, n, t0, t1 := sphere.Intersects(ray)

	if !isHit {
		t.Error("Ray did not hit")
	}

	if !vec.IsEqual(hit, *vec.MakeVec3(0, 0, -2)) {
		t.Error("Hit vector is not correct")
	}

	if !vec.IsEqual(n, *vec.MakeVec3(0, 0, 1)) {
		t.Error("N vector is not correct")
	}

	if t0 != 2 || t1 != 4 {
		t.Error("Distance is incorrect")
	}

	// TODO: Remove this
	/*
		fmt.Println(ray)
		fmt.Println(isHit)
		fmt.Printf("hit %v\n", hit)
		fmt.Printf("n %v\n", n)
		fmt.Printf("t0 %v\n", t0)
		fmt.Printf("t1 %v\n", t1)
	*/
}

func TestTransmissionRayIntersection1DHeadOn(t *testing.T) {
	t.Parallel()

	world := MakeWorld()
	center := vec.MakeVec3(0, 0, -3)
	sphere := obj.Sphere{*center, 1, color.RGBA{0, 0, 255, 1}, 1, 1.2}
	ray := cam.Ray{"camera", *vec.MakeVec3(0, 0, 0), *vec.MakeVec3(0, 0, -1)}
	_, hit, n, _, _ := sphere.Intersects(ray)

	internal_dir := RefractionVector(
		ray.Direction,
		n,
		world.RefractiveIndex,
		sphere.GetRefractiveIndex())

	if !vec.IsEqual(internal_dir, *vec.MakeVec3(0, 0, -1)) {
		t.Error("Internal refracted direction is incorrect")
	}

	ref_ray := cam.Ray{"refraction", hit, internal_dir}
	isHit2, hit2, n2, _, _ := sphere.Intersects(ref_ray)

	if !isHit2 {
		t.Error("Internal ray should intersec with sphere")
	}

	if !vec.IsEqual(hit2, *vec.MakeVec3(0, 0, -4)) {
		t.Error("Internal ray did not intersect at the right location")
	}

	external_dir := RefractionVector(
		ref_ray.Direction,
		vec.Invert(n2),
		sphere.GetRefractiveIndex(),
		world.RefractiveIndex)

	if !vec.IsEqual(external_dir, *vec.MakeVec3(0, 0, -1)) {
		t.Error("External refracted direction is incorrect")
	}

	trans_ray := cam.Ray{"transmission", hit2, external_dir}
	world_trans_ray, _ := world.MakeTransmittedRay(ray, hit, n, sphere)

	if !cam.IsEqual(trans_ray, world_trans_ray) {
		t.Error("World function incorrectly constructed transmission ray")
	}
}

func TestTransmissionRayIntersection1DAtAngle(t *testing.T) {
	t.Parallel()

	world := MakeWorld()
	center := vec.MakeVec3(0.5, 0, -3)
	sphere := obj.Sphere{*center, 1, color.RGBA{0, 0, 255, 1}, 1, 1.2}
	ray := cam.Ray{"camera", *vec.MakeVec3(0, 0, 0), *vec.MakeVec3(0, 0, -1)}
	_, hit, n, _, _ := sphere.Intersects(ray)

	internal_dir := RefractionVector(
		ray.Direction,
		n,
		world.RefractiveIndex,
		sphere.GetRefractiveIndex())

	if !vec.IsEqual(internal_dir, *vec.MakeVec3(0, 0, -1)) {
		t.Error("Internal refracted direction is incorrect")
	}

	ref_ray := cam.Ray{"refraction", hit, internal_dir}
	isHit2, hit2, n2, _, _ := sphere.Intersects(ref_ray)

	if !isHit2 {
		t.Error("Internal ray should intersec with sphere")
	}

	if !vec.IsEqual(hit2, *vec.MakeVec3(0, 0, -4)) {
		t.Error("Internal ray did not intersect at the right location")
	}

	external_dir := RefractionVector(
		ref_ray.Direction,
		vec.Invert(n2),
		sphere.GetRefractiveIndex(),
		world.RefractiveIndex)

	if !vec.IsEqual(external_dir, *vec.MakeVec3(0, 0, -1)) {
		t.Error("External refracted direction is incorrect")
	}

	trans_ray := cam.Ray{"transmitted", hit2, external_dir}
	world_trans_ray, _ := world.MakeTransmittedRay(ray, hit, n, sphere)
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
