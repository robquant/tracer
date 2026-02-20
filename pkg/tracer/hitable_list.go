package tracer

import "github.com/robquant/tracer/pkg/geo"

type HitableList []Hitable

func NewHitableList() HitableList {
	return make([]Hitable, 0)
}

func (l HitableList) Hit(r *geo.Ray, tMin, tMax float32, rec *HitRecord) bool {
	var tempRec HitRecord
	hitAnything := false
	closestSoFar := tMax
	for _, hitable := range l {
		if hitable.Hit(r, tMin, closestSoFar, &tempRec) {
			hitAnything = true
			closestSoFar = tempRec.t
			*rec = tempRec
		}
	}
	return hitAnything
}

func (l HitableList) BoundingBox() (bool, geo.Aabb) {
	if len(l) < 1 {
		return false, geo.EmptyBox
	}
	first, box := l[0].BoundingBox()
	if !first {
		return false, geo.EmptyBox
	}
	for i := 1; i < len(l); i++ {
		hasBox, tempBox := l[i].BoundingBox()
		if hasBox {
			box = geo.SurroundingBox(box, tempBox)
		} else {
			return false, geo.EmptyBox
		}
	}
	return true, box
}
