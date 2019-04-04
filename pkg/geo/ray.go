package geo

// Ray represents a ray with an origin orig
// and a direction dir
type Ray struct {
	orig, dir Vec3
	lenSq     float64
}

// NewRay constructs a new Ray from on origin and direction
func NewRay(orig, dir Vec3) Ray {
	lenSq := dir.LenSq()
	return Ray{orig: orig, dir: dir, lenSq: lenSq}
}

// Orig returns the rays direction
func (r *Ray) Orig() Vec3 {
	return r.orig
}

// Dir returns the rays direction
func (r *Ray) Dir() Vec3 {
	return r.dir
}

// At calculates a position on the ray
// as orig + t * dir
func (r *Ray) At(t float64) Vec3 {
	return r.orig.Add(r.dir.Mul(t))
}

func (r *Ray) LenSq() float64 {
	return r.lenSq
}
