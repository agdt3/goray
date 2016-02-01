package cam

import (
	"fmt"
	"github.com/agdt3/goray/vec"
	"github.com/satori/go.uuid"
	"math"
	"strings"
)

type Ray struct {
	Id        string
	Type      string
	Origin    vec.Vec3
	Direction vec.Vec3
}

func NewRay(id, typ string, orig, dir *vec.Vec3) *Ray {
	ray := new(Ray)
	if id == "" {
		ray.Id = GenerateId(nil)
	} else {
		ray.Id = id
	}
	ray.Type = typ
	ray.Origin = *orig
	ray.Direction = *dir
	return ray
}

func GenerateId(parent *Ray) string {
	// TODO: Replace with shorter string id manager
	// Otherwise id chains will be very long
	id := uuid.NewV4().String()
	if parent != nil {
		ids := []string{parent.Id, id}
		id = strings.Join(ids, "|")
	}
	return id
}

func (r *Ray) String() string {
	return fmt.Sprintf(
		"Ray: Type(%v) - Origin(%v, %v, %v) - Dir(%v, %v, %v)",
		r.Type, r.Origin.X, r.Origin.Y, r.Origin.Z,
		r.Direction.X, r.Direction.Y, r.Direction.Z)
}

func IsEqual(r1, r2 *Ray) bool {
	// Note: Does not consider Id as unique value
	if r1.Type != r2.Type {
		return false
	}

	if !vec.IsEqual(r1.Origin, r2.Origin) {
		return false
	}

	if !vec.IsEqual(r1.Direction, r2.Direction) {
		return false
	}

	return true
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

func NewPerspectiveCamera(org, dir vec.Vec3, w, h int, fovx, fovy float64) *Camera {
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
