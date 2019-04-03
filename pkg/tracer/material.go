package tracer

import "github.com/robquant/tracer/pkg/geo"

// Material is an interface to describe light scattering on different materials
type Material interface {
	// Scatter takes the incident ray and a HitRecord nad
	// returns true if the ray was not absorbed, the scattered fraction
	// of light in each color channel and the scattered ray
	Scatter(r *geo.Ray, h *HitRecord) (bool, geo.Vec3, geo.Ray)
}

// Lambertian holds albedo for a lambertian scattering surface
type Lambertian struct {
	albedo geo.Vec3
}

// NewLambertian creates new Lambertian from r,g,b albedo values
func NewLambertian(ar, ag, ab float64) *Lambertian {
	return &Lambertian{albedo: geo.NewVec3(ar, ag, ab)}
}

// Scatter implements Material Scatter interface for Lambertian
func (l *Lambertian) Scatter(r *geo.Ray, h *HitRecord) (bool, geo.Vec3, geo.Ray) {
	target := h.P().Add(h.Normal()).Add(RandomInUnitSphere())
	return true, l.albedo, geo.NewRay(h.P(), target.Sub(h.P()))
}

// Metal hold albedo for a Metal surface
type Metal struct {
	albedo geo.Vec3
	fuzz   float64
}

func reflect(v, n geo.Vec3) geo.Vec3 {
	return v.Sub(n.Mul(2 * v.Dot(n)))
}

// NewMetal constructs a new Metal from r,g,b albedo values
func NewMetal(ar, ag, ab float64, fuzz float64) *Metal {
	if fuzz > 1 {
		fuzz = 1
	}
	return &Metal{albedo: geo.NewVec3(ar, ag, ab), fuzz: fuzz}
}

// Scatter implements the Material interface Scatter function for Metal
func (m *Metal) Scatter(r *geo.Ray, h *HitRecord) (bool, geo.Vec3, geo.Ray) {
	reflected := reflect(r.Dir().Normed(), h.Normal())
	scattered := geo.NewRay(h.P(), reflected.Add(RandomInUnitSphere().Mul(m.fuzz)))
	return scattered.Dir().Dot(h.Normal()) > 0, m.albedo, scattered
}
