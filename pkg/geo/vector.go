package geo

import "math"

var Origin = Vec3{0.0, 0.0, 0.0}
var UnitX = Vec3{1.0, 0.0, 0.0}
var UnitY = Vec3{0.0, 1.0, 0.0}
var UnitZ = Vec3{0.0, 0.0, .0}

// Vec3 is a three dimensional vector
type Vec3 struct {
	x, y, z float64
}

// NewVec constructs a new Vec3
func NewVec3(x, y, z float64) Vec3 {
	return Vec3{x, y, z}
}

// X returns the x/first component of v
func (v Vec3) X() float64 {
	return v.x
}

//Y returns the y/second component of v
func (v Vec3) Y() float64 {
	return v.y
}

// Z return the z/third component of v
func (v Vec3) Z() float64 {
	return v.z
}

// Neg returns a new Vec3 with a reversed direction
func (v Vec3) Neg() Vec3 {
	return Vec3{-v.x, -v.y, -v.z}
}

// Add return a new Vec3 which is the sum of o to v
func (v Vec3) Add(o Vec3) Vec3 {
	return Vec3{v.x + o.x, v.y + o.y, v.z + o.z}
}

//Sub returns a new Vec3 which is the difference of o and v
func (v Vec3) Sub(o Vec3) Vec3 {
	return Vec3{v.x - o.x, v.y - o.y, v.z - o.z}
}

//Scale scales a Vec3 in place
func (v *Vec3) Scale(t float64) {
	v.x *= t
	v.y *= t
	v.z *= t
}

//Mul returns a new Vec3 where each component is scaled
// by a factor t
func (v Vec3) Mul(t float64) Vec3 {
	return Vec3{v.x * t, v.y * t, v.z * t}
}

//Dot returns the dot product of v and o
func (v Vec3) Dot(o Vec3) float64 {
	return v.x*o.x + v.y*o.y + v.z*o.z
}

// Cross returns the cross product of v and o
func (v Vec3) Cross(o Vec3) Vec3 {
	return Vec3{v.y*o.z - v.z*o.y,
		v.z*o.x - v.x*o.z,
		v.x*o.y - v.y*o.x}
}

// LenSq returns the squared length of v
func (v Vec3) LenSq() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

// Len returns the length of v
func (v Vec3) Len() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

// Normed returns a new normalized (unit) Vec3
func (v Vec3) Normed() Vec3 {
	k := 1. / v.Len()
	return Vec3{v.x * k, v.y * k, v.z * k}
}

// Normalize normalizes v in-place
func (v *Vec3) Normalize() {
	k := 1 / v.Len()
	v.x *= k
	v.y *= k
	v.z *= k
}
