package tracer

import (
	"math/rand"

	"github.com/chewxy/math32"
	"github.com/robquant/tracer/pkg/geo"
)

// Material is an interface to describe light scattering on different materials
type Material interface {
	// Scatter takes the incident ray and a HitRecord and
	// returns true if the ray was not absorbed, the scattered fraction
	// of light in each color channel and the scattered ray
	Scatter(r *geo.Ray, h *HitRecord, rng *rand.Rand) (bool, geo.Vec3, geo.Ray)
}

// Lambertian holds albedo for a lambertian scattering surface
type Lambertian struct {
	albedo geo.Vec3
}

// NewLambertian creates new Lambertian from r,g,b albedo values
func NewLambertian(ar, ag, ab float32) *Lambertian {
	return &Lambertian{albedo: geo.NewVec3(ar, ag, ab)}
}

// Scatter implements Material Scatter interface for Lambertian
func (l *Lambertian) Scatter(r *geo.Ray, h *HitRecord, rng *rand.Rand) (bool, geo.Vec3, geo.Ray) {
	target := h.P().Add(h.Normal()).Add(RandomInUnitSphere(rng))
	return true, l.albedo, geo.NewRay(h.P(), target.Sub(h.P()))
}

// Metal hold albedo for a Metal surface
type Metal struct {
	albedo geo.Vec3
	fuzz   float32
}

func reflect(v, n geo.Vec3) geo.Vec3 {
	return v.Sub(n.Mul(2 * v.Dot(n)))
}

func refract(v, n geo.Vec3, refRatio float32) (bool, geo.Vec3) {
	uv := v.Normed()
	dt := uv.Dot(n)
	discriminant := 1.0 - refRatio*refRatio*(1-dt*dt)
	if discriminant > 0 {
		refracted := uv.Sub(n.Mul(dt)).Mul(refRatio)
		refracted = refracted.Sub(n.Mul(math32.Sqrt(discriminant)))
		return true, refracted
	}
	return false, geo.Vec3{}
}

func schlick(cosine, refIdx float32) float32 {
	r0 := (1 - refIdx) / (1 + refIdx)
	r0 = r0 * r0
	x := 1 - cosine
	x2 := x * x
	return r0 + (1-r0)*x2*x2*x
}

// NewMetal constructs a new Metal from r,g,b albedo values
func NewMetal(ar, ag, ab float32, fuzz float32) *Metal {
	if fuzz > 1 {
		fuzz = 1
	}
	return &Metal{albedo: geo.NewVec3(ar, ag, ab), fuzz: fuzz}
}

// Scatter implements the Material interface Scatter function for Metal
func (m *Metal) Scatter(r *geo.Ray, h *HitRecord, rng *rand.Rand) (bool, geo.Vec3, geo.Ray) {
	reflected := reflect(r.Dir().Normed(), h.Normal())
	scattered := geo.NewRay(h.P(), reflected.Add(RandomInUnitSphere(rng).Mul(m.fuzz)))
	return scattered.Dir().Dot(h.Normal()) > 0, m.albedo, scattered
}

type Dielectric struct {
	refIdx float32
}

func NewDielectric(refIdx float32) *Dielectric {
	return &Dielectric{refIdx: refIdx}
}

func (d *Dielectric) Scatter(r *geo.Ray, h *HitRecord, rng *rand.Rand) (bool, geo.Vec3, geo.Ray) {
	reflected := reflect(r.Dir(), h.Normal())
	attenuation := geo.NewVec3(1.0, 1.0, 1.0)
	var outwardNormal geo.Vec3
	var refRatio float32
	var cosine float32
	var reflectionProb float32 = 1.0
	s := r.Dir().Dot(h.Normal())
	if s > 0 {
		outwardNormal = h.Normal().Mul(-1)
		refRatio = d.refIdx
		cosine = d.refIdx * s / r.Dir().Len()
	} else {
		outwardNormal = h.Normal()
		refRatio = 1.0 / d.refIdx
		cosine = -s / r.Dir().Len()
	}
	var refracted bool
	var refractedDir geo.Vec3
	if refracted, refractedDir = refract(r.Dir(), outwardNormal, refRatio); refracted {
		reflectionProb = schlick(cosine, d.refIdx)
	}
	if rng.Float32() < reflectionProb {
		return true, attenuation, geo.NewRay(h.P(), reflected)
	}
	return true, attenuation, geo.NewRay(h.P(), refractedDir)
}
