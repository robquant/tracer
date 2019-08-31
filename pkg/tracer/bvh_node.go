package tracer

import (
	"fmt"
	"math"
	"sort"

	"github.com/robquant/tracer/pkg/geo"
)

type BvhNode struct {
	box         geo.Aabb
	left, right Hitable
}

func NewBvhNode(box geo.Aabb) BvhNode {
	return BvhNode{box: box}
}

func NewBvhNodeFromList(l HitableList) BvhNode {
	fmt.Printf("%d\n", len(l))
	_, mainBox := l.BoundingBox()
	leftArea := make([]float32, len(l))
	rightArea := make([]float32, len(l))
	axis := mainBox.LongestAxis()
	if axis == 0 {
		sort.Slice(l, func(i, j int) bool {
			_, box0 := l[i].BoundingBox()
			_, box1 := l[j].BoundingBox()
			return box0.Min().X() < box1.Min().X()
		})
	} else if axis == 1 {
		sort.Slice(l, func(i, j int) bool {
			_, box0 := l[i].BoundingBox()
			_, box1 := l[j].BoundingBox()
			return box0.Min().Y() < box1.Min().Y()
		})
	} else {
		sort.Slice(l, func(i, j int) bool {
			_, box0 := l[i].BoundingBox()
			_, box1 := l[j].BoundingBox()
			return box0.Min().Z() < box1.Min().Z()
		})
	}
	boxes := make([]geo.Aabb, len(l))
	for i, hitable := range l {
		_, boxes[i] = hitable.BoundingBox()
	}
	leftBox := geo.EmptyBox
	for i, box := range boxes {
		leftBox = geo.SurroundingBox(leftBox, box)
		leftArea[i] = leftBox.Area()
	}
	rightBox := geo.EmptyBox
	for i := len(boxes) - 1; i >= 0; i-- {
		rightBox = geo.SurroundingBox(rightBox, boxes[i])
		rightArea[i] = rightBox.Area()
	}
	minSAH := float32(math.MaxFloat32)
	var minSAHIdx int
	for i := 0; i < len(l)-1; i++ {
		SAH := float32(i)*leftArea[i] + float32(len(l)-i-1)*rightArea[i+1]
		if SAH < minSAH {
			minSAH = SAH
			minSAHIdx = i
		}
	}
	var left, right Hitable
	if minSAHIdx == 0 {
		left = l
	} else {
		temp := NewBvhNodeFromList(l[:minSAHIdx+1])
		left = &temp
	}
	if minSAHIdx == len(l)-2 {
		right = l[minSAHIdx+1:]
	} else {
		temp := NewBvhNodeFromList(l[minSAHIdx+1:])
		right = &temp
	}
	return BvhNode{box: mainBox, left: left, right: right}
}

func (b *BvhNode) BoundingBox() (bool, geo.Aabb) {
	return true, b.box
}

func (b *BvhNode) Hit(r *geo.Ray, tMin, tMax float32) (bool, HitRecord) {

	if b.box.Hit(r, tMin, tMax) {
		if hitLeft, leftRecord := b.left.Hit(r, tMin, tMax); hitLeft {
			if hitRight, rightRecord := b.right.Hit(r, tMin, leftRecord.t); hitRight {
				return true, rightRecord
			}
			return true, leftRecord
		}
		return b.right.Hit(r, tMin, tMax)
	}
	return false, HitRecord{}
}
