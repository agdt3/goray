package track

import (
	//	"fmt"
	"github.com/agdt3/goray/cam"
	"github.com/agdt3/goray/obj"
	"github.com/agdt3/goray/vec"
	"image/color"
	//"reflect"
	"testing"
)

func TestNewRayTree(t *testing.T) {
	t.Parallel()

	ray1 := cam.NewRay("", "camera", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0, -1))
	ray2 := cam.NewRay("", "camera", vec.NewVec3(0, 0, 0), vec.NewVec3(0, -1, 0))
	tree := NewTree()
	tree.AddRoot(0, 0, 0.5, 0.5, ray1)
	tree.AddRoot(1, 1, 0.75, 0.75, ray2)

	if len(tree.Children) != 2 {
		t.Error("Root values were not added correctly")
	}
}

func TestMakeSubTree(t *testing.T) {
	t.Parallel()

	ray1 := cam.NewRay("", "camera", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0, -1))
	ray2 := cam.NewRay("", "camera", vec.NewVec3(0, 0, 0), vec.NewVec3(0, -1, 0))
	sphere1 := obj.Sphere{"sphere1", *vec.NewVec3(0, 0, -5), 1, color.RGBA{0, 0, 255, 1}, 1, 1}
	tree := NewTree()
	tree.AddRoot(0, 0, 0.5, 0.5, ray1)
	tree.AddRoot(1, 1, 0.75, 0.75, ray2)
	ray3 := cam.NewRay(cam.GenerateId(ray1), "reflection", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0, 0.5))
	tree.AddNode(ray3, ray1, sphere1, sphere1)

	if tree.NodeCount != 3 {
		t.Error("Nodes were not added correctly")
	}
}

func TestFindNodeById(t *testing.T) {
	t.Parallel()
	ray1 := cam.NewRay("", "camera", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0, -1))
	ray2 := cam.NewRay("", "camera", vec.NewVec3(0, 0, 0), vec.NewVec3(0, -1, 0))
	ray3 := cam.NewRay(cam.GenerateId(ray1), "reflection", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0, 0.5))
	ray4 := cam.NewRay(cam.GenerateId(ray3), "reflection", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0.5, 0.5))
	ray5 := cam.NewRay(cam.GenerateId(ray2), "reflection", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0.5, 0))
	sphere1 := obj.Sphere{"sphere1", *vec.NewVec3(0, 0, -5), 1, color.RGBA{0, 0, 255, 1}, 1, 1}
	sphere2 := obj.Sphere{"sphere2", *vec.NewVec3(0, 0, -8), 1, color.RGBA{0, 0, 255, 1}, 1, 1}

	tree := NewTree()

	tree.AddRoot(0, 0, 0.5, 0.5, ray1)
	tree.AddRoot(1, 1, 0.75, 0.75, ray2)
	tree.AddNode(ray3, ray1, sphere1, sphere1)
	tree.AddNode(ray4, ray3, sphere2, sphere2)
	tree.AddNode(ray5, ray2, sphere1, sphere1)

	if tree.NodeCount != 5 {
		t.Error("Nodes were not added correctly")
	}

	rn1, _ := tree.FindNodeByRayId(ray1.Id)
	rn2, _ := tree.FindNodeByRayId(ray2.Id)

	if rn1.RayId != ray1.Id || rn2.RayId != ray2.Id {
		t.Error("Root rays never added or found")
	}

	rn3, _ := tree.FindNodeByRayId(ray3.Id)
	rn4, _ := tree.FindNodeByRayId(ray4.Id)
	rn5, _ := tree.FindNodeByRayId(ray5.Id)

	if rn3.Parent.RayId != ray1.Id || rn3.Children[0].RayId != ray4.Id {
		t.Error("Ray tree not set up correctly")
	}

	if rn4.Parent.RayId != ray3.Id || len(rn4.Children) != 0 {
		t.Error("Ray tree not set up correctly")
	}

	if rn5.Parent.RayId != ray2.Id || len(rn5.Children) != 0 {
		t.Error("Ray tree not set up correctly")
	}
}
