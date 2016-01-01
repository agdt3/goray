package main

import (
	//"fmt"
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

func TestCameraRayIntersection(t *testing.T) {
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

func TestTransmissionRayIntersection(t *testing.T) {
	//no-op
}
