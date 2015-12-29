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
	UseRefraction  bool
	MaxReflections uint
}

type World struct {
	Cam             *Camera
	Img             draw.Image // use the draw interface
	Config          RayTraceConfig
	Objects         []Object
	RefractiveIndex float64
}

func MakeWorld() *World {
	// World with sane defaults
	world := new(World)
	dir := vec.MakeVec3(0, 0, -1)
	org := vec.MakeVec3(0, 0, 0)

	world.Cam = MakePerspectiveCamera(*org, *dir, 640, 480, 45, 45)
	world.Config = RayTraceConfig{true, true, true, 1}
	world.Img = image.NewRGBA(image.Rect(0, 0, world.Cam.Width, world.Cam.Height))
	world.RefractiveIndex = 1
	world.addObjects()
	return world
}

func (w *World) addObjects() {
	center1 := vec.MakeVec3(0, 0, -4)
	center2 := vec.MakeVec3(2, 0, -5)
	sphere1 := Sphere{*center1, 1, color.RGBA{0, 0, 255, 1}, 1, 1.2}
	sphere2 := Sphere{*center2, 1, color.RGBA{0, 255, 0, 1}, 1, 1.2}

	w.Objects = make([]Object, 2)
	w.Objects[0] = Object(sphere2)
	w.Objects[1] = Object(sphere1)
}

func (w World) makeCameraRay(x, y int) Ray {
	px, py := w.Cam.ConvertPosToPixel(x, y)
	origin := vec.MakeVec3(0, 0, 0)
	dir := vec.MakeVec3(px, py, -1)
	dir.Normalize()
	ray := Ray{"camera", *origin, *dir}
	return ray
}

func (w World) makeRefractionRay(ray Ray, hit, n vec.Vec3, obj Object) Ray {
	InvN := vec.Invert(n)
	I := ray.Direction
	/*
		cosTheta := vec.Dot(ray.Direction, invN) / (ray.Magnitude * invN.Magnitude)
		theta1 := math.Acos(cosTheta)

		sinTheta2 := w.RefractiveIndex / obj.RefractiveIndex * math.Sin(theta1)
		theta2 := math.Asin(sinTheta2)
	*/
	// TODO: Figure out how to make outgoing ray, deal with total internal
	// refraction
	// total internal reflection occurs when n2 < n1, so we can ignore this
	// for now
	r := w.RefractiveIndex / obj.GetRefractiveIndex()
	c := vec.Dot(InvN, I)
	v1 := vec.Multiply(I, r)
	modifier := (r * c) - math.Sqrt(1-((r*r)*(1-(c*c))))
	v2 := vec.Multiply(n, modifier)
	vr := vec.Add(v1, v2)
	internalRay := Ray{"refraction", hit, vr}

	isHit, nhit, _, _, _ := obj.Intersects(internalRay)
	var nray Ray
	if isHit {
		// TODO: This is wrong, need a new direction
		nray = Ray{"refraction", nhit, internalRay.Direction}
	} else {
		fmt.Println("Error - Internal ray should hit")
		nray = Ray{"failed", *vec.MakeVec3(0, 0, 0), *vec.MakeVec3(0, 0, 0)}
	}
	return nray
}

func (w World) traceRay(ray Ray) color.RGBA {
	//isHit, hit, n, t0, t1 := w.sphere.Intersects(ray)
	var closest_dist float64
	closest_dist = 100000
	pixel_color := color.RGBA{0, 0, 0, 0}

	for _, obj := range w.Objects {
		isHit, hit, n, t0, _ := obj.Intersects(ray)
		if isHit {
			new_dist := math.Abs(t0)
			if new_dist < closest_dist {
				closest_dist = new_dist
				if w.Config.UseRefraction &&
					obj.GetRefractiveIndex() > w.RefractiveIndex {
					refRay := w.makeRefractionRay(ray, hit, n, obj)
					pixel_color = w.traceRay(refRay)
				} else {
					pixel_color = obj.GetColor()
				}
			}
		}
	}
	return pixel_color
}

func (w *World) Trace() {
	b := w.Img.Bounds()
	fmt.Println(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			ray := w.makeCameraRay(x, y)
			pixelColor := w.traceRay(ray)
			w.Img.Set(x, w.Cam.Height-y, pixelColor)
		}
	}

	// TODO: Export to function later
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
	GetRefractiveIndex() float64
}

type Sphere struct {
	Center          vec.Vec3
	Radius          float64
	Col             color.RGBA
	EasingDistance  float64
	RefractiveIndex float64
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

func (s Sphere) GetRefractiveIndex() float64 {
	return s.RefractiveIndex
}

func BlendColors(c1, c2 color.RGBA, t float64) color.RGBA {
	c3 := color.RGBA{0, 0, 0, 0}
	// TODO: Use bitwise manipulation here instead
	c3.R = uint8(math.Sqrt((1-t)*float64(c1.R*c1.R) + t*float64(c2.R*c2.R)))
	c3.G = uint8(math.Sqrt((1-t)*float64(c1.G*c1.G) + t*float64(c2.G*c2.G)))
	c3.B = uint8(math.Sqrt((1-t)*float64(c1.B*c1.B) + t*float64(c2.B*c2.B)))
	c3.A = uint8(((1 - t) * float64(c1.A)) + (t * float64(c2.A)))
	return c3
}

func main() {
	world := MakeWorld()
	world.Trace()
}
