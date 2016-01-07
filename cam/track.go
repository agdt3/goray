package cam

import (
	"fmt"
	"github.com/agdt3/goray/cam"
	"string"
)

type RayTreeNode struct {
	X               int
	Y               int
	PixelX          int
	PixelY          int
	RayId           string
	ObjectImpactId  string
	ObjectEmittedId string
	Type            string
	Parent          RayTreeNode
	Children        []RayTreeNode
}

type RayTree struct {
	Children []RayTreeNode
}

func MakeTree(num_roots uint) *RayTree {
	tree = make(RayTree)
	if num_roots > 0 {
		tree.Roots = make([]Object, num_roots)
	} else {
		tree.Roots = make([]Object, 0)
	}
	return tree
}

func (t *RayTree) AddRoot(x, y, px, py int, ray *Ray, obj *Object) {
	node := RayTreeNode{
		x,
		y,
		px,
		py,
		*ray.Id,
		*obj.Id,
		"",
		*ray.Type,
		RayTreeNode{},
		make([]RayTreeNode),
	}

	t.Children = append(t.Children, node)
}

func (t *RayTree) AddNode(ray *Ray, prodobj, hitobj *Object) {
	node := RayTreeNode{
		0,
		0,
		0,
		0,
		*ray.Id,
		*prodobj.Id,
		*hitobj.Id,
		*ray.Type,
		RayTreeNode{},
		make([]RayTreeNode),
		//[]RayTreeNode
	}

	parent := t.FindNode(ray.Id)
	t.Children = append(t.Children, node)
}

func (t RayTree) FindNodeByRayId(id string) (*RayTreeNode, bool) {
	/*
		if ChildNodeId: AAA|BBB|NNN
		then ParentNodeId: AAA|BBB
	*/
	ids := strings.Split(id, "|")
	node, found := t.Find(&t.Children, &ids)
	if found {
		return node, found
	} else {
		return node, found
	}
}

func (t RayTree) Find(children *[]RayTreeNode, ids *[]string) (*RayTreeNode, bool) {
	id := ids[0]
	// TODO: Replace with search term and comparator function, mayber
	if len(*children) > 0 {
		for _, v := range children {
			if v.RayId == id && len(ids) == 1 {
				if len(ids) == 1 {
					return &v, true
				} else {
					return t.Find(&v.Children, &ids[1:])
				}
				return &v, true
			} else {
				if len(v.Children) > 0 {
					// TODO: Should only pass a slice pointer
					return Find(&v.Children, id)
				}
			}
		}
	}

	return make(RayTreeNode), false
}

func (t RayTree) FindRootByPixel(x, y int) (*RayTreeNode, bool) {
	// TODO: Optimize this by structuring Children in a way that doesn't
	// force a linear O(theta) = n search
	for _, v := range t.Children {
		if v.X == x && v.Y == y {
			return &v, true
		}
	}

	return make(RayTreeNode), false
}
