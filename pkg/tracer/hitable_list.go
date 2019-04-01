package tracer

import "github.com/robquant/tracer/pkg/geo"

type HitableList []Hitable

func NewHitableList() HitableList {
	return make([]Hitable, 0)
}

func (l HitableList) Hit(r *geo.Ray, tMin, tMax float64) (bool, HitRecord) {
	var closestHitRecord HitRecord
	hitAnything := false
	closestSoFar := tMax
	for _, h := range l {
		if hit, hitRecord := h.Hit(r, tMin, closestSoFar); hit {
			hitAnything = true
			closestSoFar = hitRecord.t
			closestHitRecord = hitRecord
		}
	}
	return hitAnything, closestHitRecord
}
