package main

import (
	"fmt"

	"github.com/agdt3/goray/cam"
	"github.com/agdt3/goray/obj"
	"github.com/agdt3/goray/read"
	"github.com/agdt3/goray/vec"
	//"github.com/agdt3/goray/track"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"os"
)

const INF_DIST float64 = 100000
const MESH_FILE_PATH string = "./meshes/test.mesh"

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
	Lights          []obj.Light
	RefractiveIndex float64
	Stats           CollisionStats
}

func NewWorld() *World {
	// World with sane defaults
	world := new(World)
	dir := vec.NewVec3(0, 0, -1)
	org := vec.NewVec3(0, 0, 0)

	world.Cam = cam.NewPerspectiveCamera(*org, *dir, 640, 480, 45, 45)
	world.Config = RayTraceConfig{true, true, false, 3}
	world.Img = image.NewRGBA(image.Rect(0, 0, world.Cam.Width, world.Cam.Height))
	world.RefractiveIndex = 1
	world.Stats = CollisionStats{0, 0}
	return world
}

func (w *World) MakeObjects() {
	// spheres
	center1 := vec.NewVec3(0, 0.5, -4)
	center2 := vec.NewVec3(3, 0, -7)
	sphere1 := obj.Sphere{"Sphere1", *center1, 1, color.RGBA{0, 0, 255, 1}, 1, 1.2}
	sphere2 := obj.Sphere{"Sphere2", *center2, 1, color.RGBA{0, 255, 0, 1}, 1, 1.2}

	// triangles
	v0 := vec.NewVec3(0, -1, -3)
	v1 := vec.NewVec3(1, 1, -3)
	v2 := vec.NewVec3(-1, 1, -3)
	triangle1 := obj.NewTriangle("Tri1", *v0, *v1, *v2, color.RGBA{255, 0, 0, 1}, 1, 1, false)

	// Slice of objects, 0 values, 3 capacity
	w.Objects = make([]obj.Object, 0, 3)
	w.Objects = append(w.Objects, obj.Object(sphere1))
	w.Objects = append(w.Objects, obj.Object(sphere2))
	w.Objects = append(w.Objects, obj.Object(triangle1))
}

func (w *World) MakeLights() {
	center := vec.NewVec3(0, 5, -2)
	light := obj.Light{"light1", *center, 1, 1, color.RGBA{255, 255, 255, 1}}

	w.Lights = make([]obj.Light, 0, 1)
	w.Lights = append(w.Lights, light)
}

func (w World) NewCameraRay(x, y int) *cam.Ray {
	px, py := w.Cam.ConvertPosToPixel(x, y)
	origin := vec.NewVec3(0, 0, 0)
	dir := vec.NewVec3(px, py, -1)
	dir.Normalize()
	ray := cam.NewRay("", "camera", origin, dir)
	return ray
}

func NewRefractionVector(l, n vec.Vec3, ref_index1, ref_index2 float64) vec.Vec3 {
	// l - Initial hit / light vector
	// n - n vector for surface
	// refIndex1 - refraction index of material of incoming ray
	// refIndex2 - refraction index of new material

	// TODO: Figure out how to make outgoing ray, deal with total internal
	// refraction
	// total internal reflection occurs when n2 < n1, so we can ignore this
	// for now
	inv_n := vec.Invert(n)
	r := ref_index1 / ref_index2
	c := vec.Dot(inv_n, l)
	v1 := vec.Multiply(l, r)
	modifier := (r * c) - math.Sqrt(1-((r*r)*(1-(c*c))))
	v2 := vec.Multiply(n, modifier)
	vr := vec.Add(v1, v2)
	return vr
}

func (w World) NewTransmittedRay(ray *cam.Ray, hit, n vec.Vec3, object obj.Object) (*cam.Ray, bool) {
	// TODO:
	// Deal with total internal refraction. Total internal reflection occurs
	// when n2 < n1, so we can ignore this for now
	external_ref_index := w.RefractiveIndex
	internal_ref_index := object.GetRefractiveIndex()

	irv := NewRefractionVector(ray.Direction, n, external_ref_index, internal_ref_index)
	irv.Normalize()
	internal_ray := cam.NewRay("", "refraction", &hit, &irv)
	is_hit2, hit2, n2, _, _ := object.Intersects(internal_ray)
	invn2 := vec.Invert(n2)
	if is_hit2 {
		erv := NewRefractionVector(irv, invn2, internal_ref_index, external_ref_index)
		erv.Normalize()
		return cam.NewRay("", "transmission", &hit2, &erv), true
	} else {
		fmt.Println("Did not hit object internally")
		return cam.NewRay("noid", "", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0, 0)), false
	}
}

func (w *World) NewShadowRay(incident *cam.Ray, n, hit vec.Vec3) *cam.Ray {
	// Creates specular shadow ray
	reflected_dir := vec.Reflect(incident.Direction, n)
	return cam.NewRay("", "shadow", &hit, &reflected_dir)
}

func (w *World) intersectLightsOld(ray *cam.Ray) (color.RGBA, bool) {
	hit_light := false
	hit_dist := INF_DIST
	hit_color := color.RGBA{0, 0, 0, 0}
	for _, v := range w.Lights {
		if hit, dist := v.Intersects(ray); hit && dist < hit_dist {
			hit_light = true
			hit_dist = dist
			hit_color = v.Col
		}
	}

	if hit_light {
		for _, v := range w.Objects {
			if is_hit, _, _, t0, _ := v.Intersects(ray); is_hit && math.Abs(t0) < hit_dist {
				fmt.Println("hit an object")
				return hit_color, false
			}
		}
	}
	return hit_color, hit_light
}

func (w *World) intersectLights(ray *cam.Ray, dist float64) (*obj.Light, float64) {
	closest_dist := dist
	var closest_light *obj.Light
	closest_light = nil
	for i, v := range w.Lights {
		is_hit, new_dist := v.Intersects(ray)
		if is_hit && (new_dist < closest_dist) {
			closest_dist = new_dist
			closest_light = &w.Lights[i]
		}
	}
	return closest_light, closest_dist
}

func (w *World) intersectObjects(ray *cam.Ray, dist float64) (*obj.Object, vec.Vec3, vec.Vec3, float64) {
	closest_dist := dist
	closest_hit_location := *vec.NewVec3(0, 0, 0)
	closest_n_vector := *vec.NewVec3(0, 0, 0)
	var closest_obj *obj.Object
	closest_obj = nil
	for i, obj := range w.Objects {
		is_hit, hit, n, t0, _ := obj.Intersects(ray)
		if new_dist := math.Abs(t0); is_hit && (new_dist < closest_dist) {
			closest_dist = new_dist
			closest_obj = &w.Objects[i]
			closest_hit_location = hit
			closest_n_vector = n
		}
	}
	return closest_obj, closest_hit_location, closest_n_vector, closest_dist
}

// TODO: Figure out whether this should return obj pointer or color val
func (w *World) TraceRay(ray *cam.Ray, reflection uint) (color.RGBA, bool) {
	current_color := color.RGBA{0, 0, 0, 0}
	if reflection > w.Config.MaxReflections {
		return current_color, false
	} else {
		reflection += 1
	}

	closest_dist := INF_DIST

	hit_obj, hit, n, dist := w.intersectObjects(ray, closest_dist)
	if hit_obj != nil {
		closest_dist = dist
	}

	// Smack into some lights
	var light *obj.Light
	if w.Config.UseLight {
		light, _ = w.intersectLights(ray, closest_dist)
	}

	// If light is the closest thing we hit, return light
	// Lights terminate all rays
	var trans_color color.RGBA
	var reflected_color color.RGBA
	if light != nil && ray.Type != "camera" {
		return light.Col, true
	} else if hit_obj != nil {
		current_color = (*hit_obj).GetColor()

		// transmitted ray
		trans_hit := false
		if w.Config.UseRefraction {
			trans_ray, _ := w.NewTransmittedRay(ray, hit, n, *hit_obj)
			trans_color, trans_hit = w.TraceRay(trans_ray, reflection)
		}

		// shadow ray
		reflect_hit := false
		if w.Config.UseShadows {
			reflected_ray := w.NewShadowRay(ray, n, hit)
			reflected_color, reflect_hit = w.TraceRay(reflected_ray, reflection)
		}

		if trans_hit {
			current_color = BlendColors(current_color, trans_color, 0.5)
		}

		if reflect_hit {
			current_color = BlendColors(current_color, reflected_color, 0.5)
		}

		return current_color, true
	}

	return current_color, false
}

func (w *World) traceRay(ray *cam.Ray, reflection uint) (color.RGBA, bool) {
	pixel_color := color.RGBA{0, 0, 0, 0}
	if reflection > w.Config.MaxReflections {
		return pixel_color, false
	} else {
		reflection += 1
	}

	closest_dist := INF_DIST
	did_hit := false

	for _, obj := range w.Objects {
		isHit, hit, n, t0, _ := obj.Intersects(ray)
		if isHit {
			did_hit = true
			new_dist := math.Abs(t0)
			if new_dist < closest_dist {
				closest_dist = new_dist
				pixel_color = obj.GetColor()

				if w.Config.UseLight {
					// Gather up direct lights
					// Create shadow ray
					shadow_ray := w.NewShadowRay(ray, n, hit)
					light_color, did_hit_light := w.intersectLightsOld(shadow_ray)

					if did_hit_light {
						pixel_color = BlendColors(pixel_color, light_color, 0.5)
					}
					// Gather up indirect lights
				}

				if w.Config.UseRefraction {
					trans_ray, success := w.NewTransmittedRay(ray, hit, n, obj)
					if success {
						ref_color, hit := w.traceRay(trans_ray, reflection)
						if hit {
							pixel_color = BlendColors(obj.GetColor(), ref_color, 0.5)
							w.Stats.Successes += 1
						} else {
							if trans_ray.Origin.Z > -9.0 {
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
			ray := w.NewCameraRay(x, y)
			pixel_color, _ := w.TraceRay(ray, 0)
			w.Img.Set(x, y, pixel_color)
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

	read.ReadMeshFile(MESH_FILE_PATH)
	world := NewWorld()
	world.MakeObjects()
	world.MakeLights()
	world.Trace()
	//world.ShowStats()
}
