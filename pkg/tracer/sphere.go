package tracer

import (
	"math"

	"github.com/robquant/tracer/pkg/geo"
)

// Sphere represented by center and radius
type Sphere struct {
	center geo.Vec3
	radius float64
}

// NewSphere constructs a new Sphere
func NewSphere(center geo.Vec3, r float64) Sphere {
	return Sphere{center, r}
}

// Hit calculates if geo.Ray r hits the sphere betwenn tMin and tMax
func (s Sphere) Hit(r *geo.Ray, tMin, tMax float64) (bool, HitRecord) {
	oc := r.Orig().Sub(s.center)
	a := r.Dir().Dot(r.Dir())
	b := r.Dir().Dot(oc)
	c := oc.Dot(oc) - s.radius*s.radius
	discriminant := b*b - a*c
	if discriminant > 0 {
		sqrt := math.Sqrt(b*b - a*c)
		temp := (-b - sqrt) / a
		if temp < tMax && temp > tMin {
			p := r.At(temp)
			return true, NewHitRecord(temp, p, p.Sub(s.center).Mul(1.0/s.radius))
		}
		temp = (-b + sqrt) / a
		if temp < tMax && temp > tMin {
			p := r.At(temp)
			return true, NewHitRecord(temp, p, p.Sub(s.center).Mul(1.0/s.radius))
		}
	}
	return false, HitRecord{}
}
