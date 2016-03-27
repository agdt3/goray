package cam

import (
	"fmt"
	"math"
	"strings"

	"github.com/agdt3/goray/vec"
	"github.com/satori/go.uuid"
)

// Ray contains the basic paramters for a ray in a scene
type Ray struct {
	ID        string
	Type      string
	Origin    vec.Vec3
	Direction vec.Vec3
}

// NewRay is a constructor of Rays
func NewRay(id, typ string, orig, dir *vec.Vec3) *Ray {
	ray := new(Ray)
	if id == "" {
		ray.ID = GenerateID(nil)
	} else {
		ray.ID = id
	}
	ray.Type = typ
	ray.Origin = *orig
	ray.Direction = *dir
	return ray
}

// GenerateID creates chains of uuids, based on parent ID
func GenerateID(parent *Ray) string {
	// TODO: Replace with shorter string id manager
	// Otherwise id chains will be very long
	id := uuid.NewV4().String()
	if parent != nil {
		ids := []string{parent.ID, id}
		id = strings.Join(ids, "|")
	}
	return id
}

// String is the string representation of a Ray
func (r *Ray) String() string {
	return fmt.Sprintf(
		"Ray: Type(%v) - Origin(%v, %v, %v) - Dir(%v, %v, %v)",
		r.Type, r.Origin.X, r.Origin.Y, r.Origin.Z,
		r.Direction.X, r.Direction.Y, r.Direction.Z)
}

// IsEqual is a comparator of Rays. It ignores the ID
func IsEqual(r1, r2 *Ray) bool {
	// Note: Does not consider ID as unique value
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

// Camera is a representation of the camera in the scene.
// The starting position is canonically (0, 0, 1)
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

// NewPerspectiveCamera creates a perspective camera in the scene
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

// ConvertPosToPixel takes an (x, y) position on an image and
// converts this to a centered (px, py) value
func (c Camera) ConvertPosToPixel(x, y int) (float64, float64) {
	px := (2.0*((float64(x)+0.5)/float64(c.Width)) - 1.0) * c.Angle * c.AspectRatio
	py := (1.0 - 2.0*((float64(y)+0.5)/float64(c.Height))) * c.Angle
	return px, py
}
