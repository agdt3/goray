package track

import (
	//"fmt"
	"image/color"

	"github.com/agdt3/goray/cam"
	"github.com/agdt3/goray/obj"
	"github.com/agdt3/goray/vec"
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
	ray3 := cam.NewRay(cam.GenerateID(ray1), "reflection", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0, 0.5))
	tree.AddNode(ray3, ray1, sphere1, sphere1)

	if tree.NodeCount != 3 {
		t.Error("Nodes were not added correctly")
	}
}

func TestFindNodeById(t *testing.T) {
	t.Parallel()
	ray1 := cam.NewRay("", "camera", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0, -1))
	ray2 := cam.NewRay("", "camera", vec.NewVec3(0, 0, 0), vec.NewVec3(0, -1, 0))
	ray3 := cam.NewRay(cam.GenerateID(ray1), "reflection", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0, 0.5))
	ray4 := cam.NewRay(cam.GenerateID(ray3), "reflection", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0.5, 0.5))
	ray5 := cam.NewRay(cam.GenerateID(ray2), "reflection", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0.5, 0))
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

	rn1 := tree.FindNodeByRayId(ray1.ID)
	rn2 := tree.FindNodeByRayId(ray2.ID)

	if rn1 == nil || rn1.RayId != ray1.ID || rn2 == nil || rn2.RayId != ray2.ID {
		t.Error("Root rays never added or found")
	}

	rn3 := tree.FindNodeByRayId(ray3.ID)
	rn4 := tree.FindNodeByRayId(ray4.ID)
	rn5 := tree.FindNodeByRayId(ray5.ID)

	if rn3 == nil || rn3.Parent.RayId != ray1.ID || rn3.Children[0].RayId != ray4.ID {
		t.Error("Ray tree not set up correctly")
	}

	if rn4 == nil || rn4.Parent.RayId != ray3.ID || len(rn4.Children) != 0 {
		t.Error("Ray tree not set up correctly")
	}

	if rn5 == nil || rn5.Parent.RayId != ray2.ID || len(rn5.Children) != 0 {
		t.Error("Ray tree not set up correctly")
	}

	/*
		st1, _ := tree.GetSubTreeString(ray1, 0)
		st2, _ := tree.GetSubTreeString(ray2, 3)
		fmt.Println(st1)
		fmt.Println(st2)
	*/
}

func TestFindNodeByXY(t *testing.T) {
	t.Parallel()
	ray1 := cam.NewRay("", "camera", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0, -1))
	ray2 := cam.NewRay("", "camera", vec.NewVec3(0, 0, 0), vec.NewVec3(0, -1, 0))
	tree := NewTree()

	tree.AddRoot(0, 0, 0.5, 0.5, ray1)
	tree.AddRoot(1, 1, 0.75, 0.75, ray2)

	rn1 := tree.FindRootByPixel(0, 1)
	rn2 := tree.FindRootByPixel(1, 1)
	rn3 := tree.FindRootByPixel(0, 0)

	if rn1 != nil || rn2.RayId != ray2.ID || rn3.RayId != ray1.ID {
		t.Error("Could not correctly find rays")
	}

}
