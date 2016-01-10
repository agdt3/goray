package track

import (
	"fmt"
	"github.com/agdt3/goray/cam"
	//	"github.com/agdt3/goray/obj"
	"github.com/agdt3/goray/vec"
	"testing"
)

func TestNewRayTree(t *testing.T) {
	t.Parallel()

	ray1 := cam.NewRay("", "camera", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0, -1))
	ray2 := cam.NewRay("", "camera", vec.NewVec3(0, 0, 0), vec.NewVec3(0, -1, 0))
	tree := NewTree(0)
	tree.AddRoot(0, 0, 0.5, 0.5, ray1)
	tree.AddRoot(1, 1, 0.75, 0.75, ray2)

	fmt.Println(len(tree.Children))
	if len(tree.Children) != 2 {
		t.Error("Root values were not added correctly")
	}
}

func TestMakeSubTree(t *testing.T) {
	t.Parallel()

	ray1 := cam.NewRay("", "camera", vec.NewVec3(0, 0, 0), vec.NewVec3(0, 0, -1))
	ray2 := cam.NewRay("", "camera", vec.NewVec3(0, 0, 0), vec.NewVec3(0, -1, 0))
	tree := NewTree(0)
	tree.AddRoot(0, 0, 0.5, 0.5, ray1)
	tree.AddRoot(1, 1, 0.75, 0.75, ray2)

	if len(tree.Children) != 2 {
		t.Error("Root values were not added correctly")
	}
}
