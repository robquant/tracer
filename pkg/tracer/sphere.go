package tracer

import (
	"math/rand"

	"github.com/chewxy/math32"
	"github.com/robquant/tracer/pkg/geo"
)

// Sphere represented by center and radius
type Sphere struct {
	center   geo.Vec3
	radius   float32
	material Material
}

// NewSphere constructs a new Sphere
func NewSphere(center geo.Vec3, r float32, m Material) *Sphere {
	return &Sphere{center, r, m}
}

// Hit calculates if geo.Ray r hits the sphere betwenn tMin and tMax
func (s *Sphere) Hit(r *geo.Ray, tMin, tMax float32) (bool, HitRecord) {
	oc := r.Orig().Sub(s.center)
	a := r.LenSq()
	b := r.Dir().Dot(oc)
	c := oc.LenSq() - s.radius*s.radius
	discriminant := b*b - a*c
	if discriminant > 0 {
		sqrt := math32.Sqrt(discriminant)
		temp := (-b - sqrt) / a
		if temp < tMax && temp > tMin {
			p := r.At(temp)
			return true, NewHitRecord(temp, p, p.Sub(s.center).Mul(1.0/s.radius), s.material)
		}
		temp = (-b + sqrt) / a
		if temp < tMax && temp > tMin {
			p := r.At(temp)
			return true, NewHitRecord(temp, p, p.Sub(s.center).Mul(1.0/s.radius), s.material)
		}
	}
	return false, HitRecord{}
}

func RandomInUnitSphere() geo.Vec3 {
	vec := geo.NewVec3(1.0, 1.0, 1.0)
	for vec.LenSq() >= 1.0 {
		vec = geo.NewVec3(rand.Float32(), rand.Float32(), rand.Float32()).Mul(2.0).Sub(geo.Diag)
	}
	return vec
}
