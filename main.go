package main

import (
	"fmt"
	"github.com/agdt3/goray/vec"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"os"
)

type Ray struct {
	Type      string
	Origin    vec.Vec3
	Direction vec.Vec3
}

type Camera struct {
	Origin      vec.Vec3
	Dir         vec.Vec3
	Width       int
	Height      int
	FOVX        float64
	FOVY        float64
	AspectRatio float64
	Angle       float64
}

func MakePerspectiveCamera(org, dir vec.Vec3, w, h int, fovx, fovy float64) *Camera {
	cam := new(Camera)
	cam.Origin = org
	cam.Dir = dir
	cam.Width = w
	cam.Height = h
	cam.FOVX = fovx
	cam.FOVY = fovy
	cam.AspectRatio = float64(w) / float64(h)
	cam.Angle = math.Tan((fovx * 0.5) / 57.296) // convert degree to radians

	return cam
}

func (c Camera) ConvertPosToPixel(x, y int) (float64, float64) {
	px := (2.0*((float64(x)+0.5)/float64(c.Width)) - 1.0) * c.Angle * c.AspectRatio
	py := (1.0 - 2.0*((float64(y)+0.5)/float64(c.Height))) * c.Angle
	return px, py
}

type RayTraceConfig struct {
	UseLight       bool
	UseShadows     bool
	MaxReflections uint
}

type World struct {
	Cam     *Camera
	Img     draw.Image // use the draw interface
	Config  RayTraceConfig
	Objects []Object
}

func NewWorld() *World {
	// World with sane defaults
	world := new(World)
	dir := vec.MakeVec3(0, 0, -1)
	org := vec.MakeVec3(0, 0, 0)

	center := vec.MakeVec3(0, 0, -4)

	world.Cam = MakePerspectiveCamera(*org, *dir, 640, 480, 45, 45)
	world.Config = RayTraceConfig{true, true, 1}
	world.Img = image.NewRGBA(image.Rect(0, 0, world.Cam.Width, world.Cam.Height))

	sphere := Sphere{*center, 1, color.RGBA{0, 0, 255, 1}, 1}
	world.Objects = make([]Object, 1)
	world.Objects[0] = Object(sphere)

	return world
}

func (w World) traceRay(ray Ray) color.RGBA {
	//isHit, hit, n, t0, t1 := w.sphere.Intersects(ray)
	for _, obj := range w.Objects {
		isHit, _, _, _, _ := obj.Intersects(ray)
		if isHit {
			return obj.GetColor()
		}
	}
	return color.RGBA{0, 0, 0, 0}
}

func (w *World) Trace() {
	b := w.Img.Bounds()
	fmt.Println(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			px, py := w.Cam.ConvertPosToPixel(x, y)
			origin := vec.MakeVec3(0, 0, 0)
			dir := vec.MakeVec3(px, py, -1)
			dir.Normalize()
			ray := Ray{"camera", *origin, *dir}
			pixelColor := w.traceRay(ray)
			w.Img.Set(x, w.Cam.Height-y, pixelColor)
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
	GetColor() color.RGBA
}

type Sphere struct {
	Center         vec.Vec3
	Radius         float64
	Col            color.RGBA
	EasingDistance float64
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
		t0 := float64(0)
		t1 := float64(0)
		hit := ray.Origin
		oc.Multiply(-1)
		return true, hit, oc, t0, t1
	}

	thc := math.Sqrt(t2hc)
	t0 := t_ca - thc
	t1 := t_ca + thc

	dist := s.EasingDistance * t0
	rd.Multiply(dist)
	hit := vec.Add(ray.Origin, rd)

	n := vec.Subtract(hit, sc)
	n.Divide(s.Radius)

	/*
		fmt.Println(hit)
		fmt.Println(n)
		fmt.Println(t0)
		fmt.Println(t1)
	*/

	return true, hit, n, t0, t1
}

func (s Sphere) GetColor() color.RGBA {
	return s.Col
}

func main() {
	world := NewWorld()
	world.Trace()
}
