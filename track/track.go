package track

import (
	//"fmt"
	"github.com/agdt3/goray/cam"
	"github.com/agdt3/goray/obj"
	"strings"
)

type RayTreeNode struct {
	X               int
	Y               int
	PixelX          float64
	PixelY          float64
	RayId           string
	ObjectImpactId  string
	ObjectEmittedId string
	Type            string
	Parent          *RayTreeNode
	Children        []RayTreeNode
}

type RayTree struct {
	Children []RayTreeNode
}

func NewTree(num_roots uint) *RayTree {
	tree := new(RayTree)
	if num_roots > 0 {
		tree.Children = make([]RayTreeNode, num_roots)
	} else {
		tree.Children = make([]RayTreeNode, 0)
	}
	return tree
}

func (t *RayTree) AddRoot(x, y int, px, py float64, ray *cam.Ray) {
	children := make([]RayTreeNode, 0)
	node := RayTreeNode{
		x,
		y,
		px,
		py,
		ray.Id,
		"",
		"",
		ray.Type,
		nil,
		children,
	}

	t.Children = append(t.Children, node)
}

func (t *RayTree) AddNode(ray, parentray *cam.Ray, parentobj, hitobj obj.Object) {
	children := make([]RayTreeNode, 0)
	node := RayTreeNode{
		0,
		0,
		0,
		0,
		ray.Id,
		parentobj.GetId(),
		hitobj.GetId(),
		ray.Type,
		nil,
		children,
	}

	parent, _ := t.FindNodeByRayId(parentray.Id)
	parent.Children = append(parent.Children, node)
}

func (t *RayTree) FindNodeByRayId(id string) (*RayTreeNode, bool) {
	/*
		if ChildNodeId: AAA|BBB|NNN
		then ParentNodeId: AAA|BBB
	*/
	//ids := strings.Split(id, "|")
	return t.Find(t.Children, id)
}

func (t *RayTree) Find(children []RayTreeNode, id string) (*RayTreeNode, bool) {
	// TODO: Replace with search term and comparator function, maybe
	if len(children) > 0 {
		for _, v := range children {
			// id may be longer than RayId
			if strings.HasPrefix(id, v.RayId) {
				if (len(id) > len(v.RayId)) && (len(v.Children) > 0) {
					// TODO: Potential bug:
					// len(v.Chil) == 0 but id is searchable

					// id is longer, keep searching
					return t.Find(v.Children, id)
				} else {
					// id lengths match, this is the correct ray node
					return &v, true
				}
				return &v, true
			}
		}
	}

	return nil, false
}

func (t RayTree) FindRootByPixel(x, y int) (*RayTreeNode, bool) {
	// TODO: Optimize this by structuring Children in a way that doesn't
	// force a linear O(theta) = n search
	for _, v := range t.Children {
		if v.X == x && v.Y == y {
			return &v, true
		}
	}

	return nil, false
}
