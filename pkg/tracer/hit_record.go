package tracer

import "github.com/robquant/tracer/pkg/geo"

// HitRecord
type HitRecord struct {
	t        float32
	p        geo.Vec3
	normal   geo.Vec3
	material Material
}

func NewHitRecord(t float32, p, normal geo.Vec3, material Material) HitRecord {
	return HitRecord{t, p, normal, material}
}

func (h HitRecord) Normal() geo.Vec3 {
	return h.normal
}

func (h HitRecord) P() geo.Vec3 {
	return h.p
}

func (h HitRecord) Material() Material {
	return h.material
}

type Hitable interface {
	Hit(r *geo.Ray, tMin, tMax float32) (bool, HitRecord)
	BoundingBox() (bool, geo.Aabb)
}
