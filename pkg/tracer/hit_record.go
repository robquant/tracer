package tracer

import "github.com/robquant/tracer/pkg/geo"

// HitRecord
type HitRecord struct {
	t      float64
	p      geo.Vec3
	normal geo.Vec3
}

func NewHitRecord(t float64, p, normal geo.Vec3) HitRecord {
	return HitRecord{t, p, normal}
}

func (h HitRecord) Normal() geo.Vec3 {
	return h.normal
}

func (h HitRecord) P() geo.Vec3 {
	return h.p
}

type Hitable interface {
	Hit(r *geo.Ray, tMin, tMax float64) (bool, HitRecord)
}
