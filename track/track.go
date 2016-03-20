package track

import (
	"errors"
	"fmt"
	"strings"

	"github.com/agdt3/goray/cam"
	"github.com/agdt3/goray/obj"
)

type RayTreeNode struct {
	X         int
	Y         int
	PixelX    float64
	PixelY    float64
	RayId     string
	ObjEmitId string
	ObjHitId  string
	RayType   string
	Parent    *RayTreeNode
	Children  []RayTreeNode
}

type RayTree struct {
	NodeCount int
	Children  []RayTreeNode
}

// Comparator type
// type comparator func(param string, node RayTreeNode) bool

func NewTree() *RayTree {
	tree := new(RayTree)
	tree.NodeCount = 0
	return tree
}

func (t *RayTree) String() string {
	return fmt.Sprintf(
		"RayTree: Number of roots(%v) - Number of nodes(%v)",
		len(t.Children),
		t.NodeCount)
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
	t.NodeCount += 1
}

func (t *RayTree) AddNode(ray, parentray *cam.Ray, parentobj, hitobj obj.Object) {
	parent := t.FindNodeByRayId(parentray.Id)
	if parent == nil {
		// TODO: Turn this into an error? or use nil on parent pointer
		fmt.Println("Could not find parent node. Cannot add current node")
		return
	}

	node := new(RayTreeNode)
	node.X = parent.X
	node.Y = parent.Y
	node.PixelX = parent.PixelX
	node.PixelY = parent.PixelY
	node.RayId = ray.Id
	node.RayType = ray.Type

	if parentobj != nil {
		node.ObjEmitId = parentobj.GetId()
	} else {
		node.ObjEmitId = ""
	}

	if hitobj != nil {
		node.ObjHitId = hitobj.GetId()
	} else {
		node.ObjHitId = ""
	}

	node.Parent = parent
	parent.Children = append(parent.Children, *node)
	t.NodeCount += 1
}

func (t *RayTree) FindNodeByRayId(id string) *RayTreeNode {
	/*
		if ChildNodeId: AAA|BBB|NNN
		then ParentNodeId: AAA|BBB
	*/
	return t.Find(t.Children, id)
}

func (t *RayTree) Find(children []RayTreeNode, id string) *RayTreeNode {
	// TODO: Replace with search term and comparator function, maybe
	if len(children) > 0 {
		for i, v := range children {
			// id may be longer than RayId
			if strings.HasPrefix(id, v.RayId) {
				if (len(id) > len(v.RayId)) && (len(v.Children) > 0) {
					// TODO: Potential bug:
					// len(v.Chil) == 0 but id is searchable

					// id is longer, keep searching
					return t.Find(children[i].Children, id)
				} else {
					// id lengths match, this is the correct ray node
					return &children[i]
				}
				return &children[i]
			}
		}
	}

	return nil
}

func (t *RayTree) FindRootByPixel(x, y int) *RayTreeNode {
	// TODO: Optimize this by structuring Children in a way that doesn't
	// force a linear O(theta) = n search

	for i, _ := range t.Children {
		if t.Children[i].X == x && t.Children[i].Y == y {
			return &t.Children[i]
		}
	}

	return nil
}

func (t *RayTree) GetSubTreeString(ray *cam.Ray, verbosity int) (string, error) {
	// 0 = Lowest
	// 3 = Highest TODO: Implement
	node := t.FindNodeByRayId(ray.Id)
	if node == nil {
		return "", errors.New("Starting node could not be found")
	}

	accumulator := make([]string, 0)
	t.accumulateNodes(node, &accumulator, verbosity)
	return strings.Join(accumulator, " --> "), nil
}

func (t *RayTree) accumulateNodes(node *RayTreeNode, accumulator *[]string, verbosity int) {
	// TODO: An alternative approach is to collect all data and
	// use verbosity to filter on display
	switch verbosity {
	case 0:
		shortId := t.shortenId(node)
		*accumulator = append(*accumulator, shortId)
		break
	case 1:
		*accumulator = append(*accumulator, node.RayId)
		break
	case 2:
		break
	case 3:
		*accumulator = append(*accumulator, node.String())
		break
	default:
		*accumulator = append(*accumulator, node.RayId)
	}

	if len(node.Children) > 0 {
		for _, v := range node.Children {
			t.accumulateNodes(&v, accumulator, verbosity)
		}
	}
}

func (t *RayTree) shortenId(node *RayTreeNode) string {
	parts := strings.Split(node.RayId, "|")
	for i, v := range parts {
		parts[i] = v[:8]
	}
	return strings.Join(parts, "|")
}

func (rtn *RayTreeNode) String() string {
	isRoot := true
	if rtn.Parent != nil {
		isRoot = false
	}

	lineone := fmt.Sprintf("RayTreeNode - ID(%v)\n", rtn.RayId)
	linetwo := fmt.Sprintf("Root(%v) - Type(%v) - XY | PXPY(%v, %v, | %v, %v)\n",
		isRoot, rtn.RayType, rtn.X, rtn.Y, rtn.PixelX, rtn.PixelY)
	linethree := fmt.Sprintf("Emitted By(%v) - Hit(%v)\n", rtn.ObjEmitId, rtn.ObjHitId)
	linefour := fmt.Sprintf("Children(%v)\n", len(rtn.Children))
	lines := strings.Join([]string{lineone, linetwo, linethree, linefour}, "")
	return lines
}
