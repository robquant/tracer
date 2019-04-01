package tracer

import "github.com/robquant/tracer/pkg/geo"

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

type Hitable interface {
	Hit(r *geo.Ray, tMin, tMax float64) (bool, HitRecord)
}
