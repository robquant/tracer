package tracer

import (
	"math"
	"sort"

	"github.com/robquant/tracer/pkg/geo"
)

type NodeKind uint8

const (
	NodeSphere NodeKind = iota
	NodeBvh
)

type HitableNode struct {
	kind    NodeKind
	sphere  *Sphere
	bvhNode *BvhNode
}

func (h *HitableNode) Hit(r *geo.Ray, tMin, tMax float32, rec *HitRecord) bool {
	switch h.kind {
	case NodeSphere:
		return h.sphere.Hit(r, tMin, tMax, rec)
	case NodeBvh:
		return h.bvhNode.Hit(r, tMin, tMax, rec)
	}
	return false
}

func (h *HitableNode) BoundingBox() (bool, geo.Aabb) {
	switch h.kind {
	case NodeSphere:
		return h.sphere.BoundingBox()
	case NodeBvh:
		return h.bvhNode.BoundingBox()
	}
	return false, geo.EmptyBox
}

type BvhNode struct {
	box         geo.Aabb
	left, right HitableNode
}

func NewBvhNode(box geo.Aabb) BvhNode {
	return BvhNode{box: box}
}

func hitableToNode(h Hitable) HitableNode {
	switch v := h.(type) {
	case *Sphere:
		return HitableNode{kind: NodeSphere, sphere: v}
	case *BvhNode:
		return HitableNode{kind: NodeBvh, bvhNode: v}
	}
	panic("unsupported Hitable type")
}

func NewBvhNodeFromList(l HitableList) BvhNode {
	_, mainBox := l.BoundingBox()
	n := len(l)
	leftArea := make([]float32, n)
	rightArea := make([]float32, n)

	// Pre-compute all bounding boxes once
	boxes := make([]geo.Aabb, n)
	for i, hitable := range l {
		_, boxes[i] = hitable.BoundingBox()
	}

	// Sort an index array by the chosen axis
	axis := mainBox.LongestAxis()
	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		ii, jj := indices[i], indices[j]
		switch axis {
		case 0:
			return boxes[ii].Min().X() < boxes[jj].Min().X()
		case 1:
			return boxes[ii].Min().Y() < boxes[jj].Min().Y()
		default:
			return boxes[ii].Min().Z() < boxes[jj].Min().Z()
		}
	})

	// Reorder both l and boxes to match sorted order
	sortedL := make(HitableList, n)
	sortedBoxes := make([]geo.Aabb, n)
	for i, idx := range indices {
		sortedL[i] = l[idx]
		sortedBoxes[i] = boxes[idx]
	}

	leftBox := geo.EmptyBox
	for i, box := range sortedBoxes {
		leftBox = geo.SurroundingBox(leftBox, box)
		leftArea[i] = leftBox.Area()
	}
	rightBox := geo.EmptyBox
	for i := n - 1; i >= 0; i-- {
		rightBox = geo.SurroundingBox(rightBox, sortedBoxes[i])
		rightArea[i] = rightBox.Area()
	}
	minSAH := float32(math.MaxFloat32)
	var minSAHIdx int
	for i := 0; i < n-1; i++ {
		SAH := float32(i)*leftArea[i] + float32(n-i-1)*rightArea[i+1]
		if SAH < minSAH {
			minSAH = SAH
			minSAHIdx = i
		}
	}
	var left, right HitableNode
	if minSAHIdx == 0 {
		left = hitableToNode(sortedL[0])
	} else {
		temp := NewBvhNodeFromList(sortedL[:minSAHIdx+1])
		left = HitableNode{kind: NodeBvh, bvhNode: &temp}
	}
	if minSAHIdx == n-2 {
		right = hitableToNode(sortedL[minSAHIdx+1])
	} else {
		temp := NewBvhNodeFromList(sortedL[minSAHIdx+1:])
		right = HitableNode{kind: NodeBvh, bvhNode: &temp}
	}
	return BvhNode{box: mainBox, left: left, right: right}
}

func (b *BvhNode) BoundingBox() (bool, geo.Aabb) {
	return true, b.box
}

func (b *BvhNode) Hit(r *geo.Ray, tMin, tMax float32, rec *HitRecord) bool {
	if b.box.Hit(r, tMin, tMax) {
		var leftRec, rightRec HitRecord
		if b.left.Hit(r, tMin, tMax, &leftRec) {
			if b.right.Hit(r, tMin, leftRec.t, &rightRec) {
				*rec = rightRec
				return true
			}
			*rec = leftRec
			return true
		}
		return b.right.Hit(r, tMin, tMax, rec)
	}
	return false
}
