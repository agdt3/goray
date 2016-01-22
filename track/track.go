package track

import (
	"errors"
	"fmt"
	"github.com/agdt3/goray/cam"
	"github.com/agdt3/goray/obj"
	"strings"
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

	//fmt.Printf("Add Root %p %v\n", &node, node)
	t.Children = append(t.Children, node)
	t.NodeCount += 1
}

func (t *RayTree) AddNode(ray, parentray *cam.Ray, parentobj, hitobj obj.Object) {
	parent, found := t.FindNodeByRayId(parentray.Id)
	if !found {
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

	//fmt.Printf("Add Node %p %v \n", &node, node)
	//fmt.Printf("To Parent %p %v \n", &parent, parent)
	node.Parent = parent
	parent.Children = append(parent.Children, *node)
	t.NodeCount += 1
}

func (t *RayTree) FindNodeByRayId(id string) (*RayTreeNode, bool) {
	/*
		if ChildNodeId: AAA|BBB|NNN
		then ParentNodeId: AAA|BBB
	*/
	return t.Find(t.Children, id)
}

func (t *RayTree) Find(children []RayTreeNode, id string) (*RayTreeNode, bool) {
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
					return &children[i], true
				}
				return &children[i], true
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

func (t *RayTree) PrintSubTree(ray *cam.Ray, verbosity uint) (string, error) {
	// 0 = Lowest
	// 3 = Highest TODO: Implement
	verbosity = 0
	node, found := t.FindNodeByRayId(ray.Id)
	if !found {
		return "", errors.New("Starting node could not be found")
	}

	accumulator := make([]string, 0)
	TraverseNodes(node, &accumulator)
	return strings.Join(accumulator, " "), nil
}

func TraverseNodes(node *RayTreeNode, accumulator *[]string) {
	// TODO: Does not keep track of verbosity level
	*accumulator = append(*accumulator, node.RayId)
	if len(node.Children) > 0 {
		for _, v := range node.Children {
			TraverseNodes(&v, accumulator)
		}
	}
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
