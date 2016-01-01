package main

import (
	"fmt"
	"github.com/agdt3/goray/cam"
	"github.com/agdt3/goray/obj"
	"github.com/agdt3/goray/vec"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"os"
)

type RayTraceConfig struct {
	UseLight       bool
	UseShadows     bool
	UseRefraction  bool
	MaxReflections uint
}

type World struct {
	Cam             *cam.Camera
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

	world.Cam = cam.MakePerspectiveCamera(*org, *dir, 640, 480, 45, 45)
	world.Config = RayTraceConfig{true, true, true, 1}
	world.Img = image.NewRGBA(image.Rect(0, 0, world.Cam.Width, world.Cam.Height))
	world.RefractiveIndex = 1
	world.addObjects()
	return world
}

func (w *World) addObjects() {
	center1 := vec.MakeVec3(0, 0, -4)
	center2 := vec.MakeVec3(0.5, 0, -5)
	sphere1 := obj.Sphere{*center1, 1, color.RGBA{0, 0, 255, 1}, 1, 1.2}
	sphere2 := obj.Sphere{*center2, 1, color.RGBA{0, 255, 0, 1}, 1, 1.2}

	w.Objects = make([]Object, 2)
	w.Objects[0] = Object(sphere1)
	w.Objects[1] = Object(sphere2)
}

func (w World) makeCameraRay(x, y int) cam.Ray {
	px, py := w.Cam.ConvertPosToPixel(x, y)
	origin := vec.MakeVec3(0, 0, 0)
	dir := vec.MakeVec3(px, py, -1)
	dir.Normalize()
	ray := cam.Ray{"camera", *origin, *dir}
	return ray
}

func (w World) makeRefractionRay(ray cam.Ray, hit, n vec.Vec3, obj Object) (cam.Ray, bool) {
	fmt.Println("incoming ray")
	fmt.Println(ray)
	InvN := vec.Invert(n)
	I := ray.Direction
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
	internalRay := cam.Ray{"refraction", hit, vr}

	nisHit, nhit, nn, nt0, nt1 := obj.Intersects(internalRay)
	/*
		if nisHit {
			fmt.Println("internal hit result")
			fmt.Println(nhit)
			fmt.Println(nn)
			fmt.Println(nt0)
			fmt.Println(nt1)
		}
	*/
	return internalRay, false
	/*
		isHit, nhit, _, _, _ := obj.Intersects(internalRay)
		var nray cam.Ray
		nerr := false
		if isHit {
			fmt.Println("Something hit!")
			// TODO: This is wrong, need a new direction
			nray = cam.Ray{"refraction", nhit, ray.Direction}
		} else {
			//fmt.Println("Error - Internal ray should hit")
			nray = cam.Ray{"failed", *vec.MakeVec3(0, 0, 0), *vec.MakeVec3(0, 0, 0)}
			nerr = true
		}
		return nray, nerr
	*/
}

func (w World) traceRay(ray cam.Ray) color.RGBA {
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
					refRay, err := w.makeRefractionRay(ray, hit, n, obj)
					if !err {
						refColor := w.traceRay(refRay)
						if refColor.R != 0 && refColor.G != 0 && refColor.B != 0 {
							pixel_color = BlendColors(obj.GetColor(), refColor, 0.5)
						} else {
							pixel_color = obj.GetColor()
						}
						//fmt.Println("blended colors")
						//fmt.Println(pixel_color)
					}
					pixel_color = obj.GetColor()
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
	Intersects(cam.Ray) (bool, vec.Vec3, vec.Vec3, float64, float64)
	GetColor() color.RGBA
	GetRefractiveIndex() float64
}

func BlendColors(c1, c2 color.RGBA, t float64) color.RGBA {
	c3 := color.RGBA{0, 0, 0, 0}
	// TODO: Use bitwise manipulation here instead
	c1R16 := uint16(c1.R)
	c1G16 := uint16(c1.G)
	c1B16 := uint16(c1.B)

	c2R16 := uint16(c2.R)
	c2G16 := uint16(c2.G)
	c2B16 := uint16(c2.B)

	c1RSQ := float64(c1R16 * c1R16)
	c1GSQ := float64(c1G16 * c1G16)
	c1BSQ := float64(c1B16 * c1B16)

	c2RSQ := float64(c2R16 * c2R16)
	c2GSQ := float64(c2G16 * c2G16)
	c2BSQ := float64(c2B16 * c2B16)

	c3.R = uint8(math.Sqrt((1.0-t)*c1RSQ + t*c2RSQ))
	c3.G = uint8(math.Sqrt((1.0-t)*c1GSQ + t*c2GSQ))
	c3.B = uint8(math.Sqrt((1.0-t)*c1BSQ + t*c2BSQ))
	c3.A = uint8(((1 - t) * float64(c1.A)) + (t * float64(c2.A)))
	return c3
}

func main() {
	world := MakeWorld()
	world.Trace()
}
