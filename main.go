package main

import (
	"fmt"
	"github.com/agdt3/goray/cam"
	"github.com/agdt3/goray/obj"
	//"github.com/agdt3/goray/track"
	"github.com/agdt3/goray/vec"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"os"
)

type CollisionStats struct {
	Successes uint
	Failures  uint
}

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
	Objects         []obj.Object
	RefractiveIndex float64
	Stats           CollisionStats
}

func NewWorld() *World {
	// World with sane defaults
	world := new(World)
	dir := vec.NewVec3(0, 0, -1)
	org := vec.NewVec3(0, 0, 0)

	world.Cam = cam.MakePerspectiveCamera(*org, *dir, 640, 480, 45, 45)
	world.Config = RayTraceConfig{true, true, true, 5}
	world.Img = image.NewRGBA(image.Rect(0, 0, world.Cam.Width, world.Cam.Height))
	world.RefractiveIndex = 1
	world.Stats = CollisionStats{0, 0}
	return world
}

func (w *World) MakeObjects() {
	center1 := vec.NewVec3(0, 0, -4)
	center2 := vec.NewVec3(0, 0, -10)
	sphere1 := obj.Sphere{"Sphere1", *center1, 1, color.RGBA{0, 0, 255, 1}, 1, 1.2}
	sphere2 := obj.Sphere{"Sphere2", *center2, 1, color.RGBA{0, 255, 0, 1}, 1, 1.2}

	// Known objects
	w.Objects = make([]obj.Object, 2)
	w.Objects[0] = obj.Object(sphere1)
	w.Objects[1] = obj.Object(sphere2)
	// Dynamic number
	//w.Objects = make([]Object, 0)
	//w.Objects = append(w.Objects, Object(sphere1))
	//w.Objects = append(w.Objects, Object(sphere2))
}

func (w World) makeCameraRay(x, y int) *cam.Ray {
	px, py := w.Cam.ConvertPosToPixel(x, y)
	origin := vec.NewVec3(0, 0, 0)
	dir := vec.NewVec3(px, py, -1)
	dir.Normalize()
	//ray := cam.Ray{"A", "camera", *origin, *dir}
	ray := cam.NewRay("", "camera", origin, dir)
	return ray
}

func RefractionVector(l, n vec.Vec3, refIndex1, refIndex2 float64) vec.Vec3 {
	// l - Initial hit / light vector
	// n - n vector for surface
	// refIndex1 - refraction index of material of incoming ray
	// refIndex2 - refraction index of new material

	// TODO: Figure out how to make outgoing ray, deal with total internal
	// refraction
	// total internal reflection occurs when n2 < n1, so we can ignore this
	// for now
	invN := vec.Invert(n)
	r := refIndex1 / refIndex2
	c := vec.Dot(invN, l)
	v1 := vec.Multiply(l, r)
	modifier := (r * c) - math.Sqrt(1-((r*r)*(1-(c*c))))
	v2 := vec.Multiply(n, modifier)
	vr := vec.Add(v1, v2)
	return vr
}

func (w World) MakeTransmittedRay(ray *cam.Ray, hit, n vec.Vec3, object obj.Object) (*cam.Ray, bool) {
	// TODO:
	// Deal with total internal refraction. Total internal reflection occurs
	// when n2 < n1, so we can ignore this for now
	externalRefIndex := w.RefractiveIndex
	internalRefIndex := object.GetRefractiveIndex()

	irv := RefractionVector(ray.Direction, n, externalRefIndex, internalRefIndex)
	irv.Normalize()
	//internalRay := cam.Ray{"refraction", hit, irv}
	internalRay := cam.NewRay("", "refraction", &hit, &irv)
	isHit2, hit2, n2, _, _ := object.Intersects(internalRay)
	invn2 := vec.Invert(n2)
	if isHit2 {
		erv := RefractionVector(irv, invn2, internalRefIndex, externalRefIndex)
		erv.Normalize()
		return cam.NewRay("", "transmission", &hit2, &erv), true
	} else {
		fmt.Println("Did not hit object internally")
		return cam.NewRay("noid", "", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0, 0)), false
	}
}

func (w *World) traceRay(ray *cam.Ray, reflection uint) (color.RGBA, bool) {
	pixel_color := color.RGBA{0, 0, 0, 0}
	if reflection > w.Config.MaxReflections {
		return pixel_color, false
	} else {
		reflection += 1
	}

	var closest_dist float64
	closest_dist = 100000
	did_hit := false

	for _, obj := range w.Objects {
		isHit, hit, n, t0, _ := obj.Intersects(ray)
		if isHit {
			did_hit = true
			new_dist := math.Abs(t0)
			if new_dist < closest_dist {
				closest_dist = new_dist
				pixel_color = obj.GetColor()
				if w.Config.UseRefraction {
					transRay, success := w.MakeTransmittedRay(ray, hit, n, obj)
					if success {
						refColor, hit := w.traceRay(transRay, reflection)
						if hit {
							pixel_color = BlendColors(obj.GetColor(), refColor, 0.5)
							w.Stats.Successes += 1
						} else {
							if transRay.Origin.Z > -9.0 {
								w.Stats.Failures += 1
								//fmt.Println(transRay)
							}
						}
					}
				}
			}
		}
	}
	return pixel_color, did_hit
}

func (w *World) Trace() {
	b := w.Img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			ray := w.makeCameraRay(x, y)
			pixelColor, _ := w.traceRay(ray, 0)
			w.Img.Set(x, w.Cam.Height-y, pixelColor)
		}
	}

	// TODO: Export to separate function
	f, err := os.Create("./test.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	jpeg.Encode(f, w.Img, &jpeg.Options{100})
}

func (w World) ShowStats() {
	fmt.Printf("Successes %v\n", w.Stats.Successes)
	fmt.Printf("Failures %v\n", w.Stats.Failures)
	total := float64(w.Stats.Successes) + float64(w.Stats.Failures)
	ratio := float64(w.Stats.Failures) / total
	fmt.Printf("Ratio %v\n", ratio)
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
	world := NewWorld()
	world.MakeObjects()
	world.Trace()
	world.ShowStats()
}
