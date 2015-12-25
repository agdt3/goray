package main

import (
	"fmt"
	"github.com/agdt3/goray/vec"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	//"math"
	"os"
)

type Ray struct {
	Type      string
	Origin    vec.Vec3
	Direction vec.Vec3
}

type CameraPerspective struct {
	Origin vec.Vec3
	Dir    vec.Vec3
}

type CameraTHEOTHERONE struct {
	Origin vec.Vec3
	Dir    vec.Vec3
}

type RayTraceConfig struct {
	UseLight       bool
	UseShadows     bool
	MaxReflections uint
	ImageWidth     int
	ImageHeight    int
}

type World struct {
	Cam CameraPerspective
	Img draw.Image // use the draw interface
	//    Obj []Objects
	Config RayTraceConfig
}

func NewWorld() *World {
	// World with sane defaults
	world := new(World)
	dir := vec.MakeVec3(0, 0, -1)
	org := vec.MakeVec3(0, 0, 0)

	world.Cam = CameraPerspective{*org, *dir}
	world.Config = RayTraceConfig{true, true, 1, 640, 480}
	world.Img = image.NewRGBA(image.Rect(0, 0, world.Config.ImageWidth, world.Config.ImageHeight))
	//    world.Obj = nil // TODO: vec.Make objects

	return world
}

func (w World) traceRay(origin vec.Vec3, dir vec.Vec3) color.RGBA {
	return color.RGBA{0, 0, 255, 0}
}

func (w World) generateRay(x, y, z float64) *vec.Vec3 {
	// generates a vector3 that is assumed
	// have an origin at (0,0,0)
	vec3 := vec.MakeVec3(x, y, z)
	return vec3
}

func (w *World) Trace() {
	b := w.Img.Bounds()
	fmt.Println(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			dir := vec.MakeVec3(0, 0, -1)
			origin := vec.MakeVec3(0, 0, 0)
			pixelColor := w.traceRay(*origin, *dir)
			w.Img.Set(x, y, pixelColor)
			//            fmt.Println(w.Img.At(x,y))
		}
	}

	// TODO: Get rid of this later
	f, err := os.Create("./test.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	jpeg.Encode(f, w.Img, &jpeg.Options{100})
}

type Object interface {
	Intersects(Ray) (bool, vec.Vec3, vec.Vec3, float64, float64)
}

type Sphere struct {
	Center vec.Vec3
	Radius float64
	Col    color.Color
}

func (s Sphere) Intersects(ray Ray) (bool, vec.Vec3, vec.Vec3, float64, float64) {
	sc := s.Center
	rd := ray.Direction
	rd.Normalize()

	srsq := s.Radius * s.Radius
	oc := vec.Subtract(sc, rd)
	l2oc := vec.Dot(oc, oc)
	t_ca := vec.Dot(oc, rd)

	//sphere located behind ray origin
	if t_ca < 0 {
		return false, *vec.MakeVec3(0, 0, 0), *vec.MakeVec3(0, 0, 0), 0, 0
	}

	d2 := l2oc - (t_ca * t_ca)

	// if the distance between the closest point to the sphere center on
	// the projected ray is greater than the radius, then the projected
	// ray is definitely outside the bounds of the sphere
	if d2 > srsq {
		return false, *vec.MakeVec3(0, 0, 0), *vec.MakeVec3(0, 0, 0), 0, 0
	}

	t2hc := srsq - d2

	if t2hc < 0 {
		return false, *vec.MakeVec3(0, 0, 0), *vec.MakeVec3(0, 0, 0), 0, 0
	}

	// if the origin is inside the sphere of light, it counts as a hit
	// the results aren't that useful
	if l2oc < srsq {
		t0 := 0
		t1 := 0
		hit := ray.Origin
		n := -1 * oc
		return true, hit, n, t0, t1
	}

	thc := math.Sqrt(t2hc)
	t0 := t_ca - thc
	t1 := t_ca + thc

	dist := easing_distance * t0
	//hit := Add(ray.Origin, (

	return false, *vec.MakeVec3(0, 0, 0), *vec.MakeVec3(0, 0, 0), 0, 0
}

func main() {
	world := NewWorld()
	world.Trace()
}
