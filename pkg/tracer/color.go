package tracer

import "github.com/robquant/tracer/pkg/geo"

var Red = Color{geo.UnitX}
var Green = Color{geo.UnitY}
var Blue = Color{geo.UnitZ}
var Black = Color{geo.Origin}

// Color is an RGB color triple using floats
type Color struct {
	geo.Vec3
}

// NewColor creates a new Color from RGB float values
func NewColor(r, g, b float32) Color {
	return Color{geo.NewVec3(r, g, b)}
}

// R returns the red component of c
func (c Color) R() float32 {
	return c.X()
}

// G returns the green component of c
func (c Color) G() float32 {
	return c.Y()
}

// B returns the blue component of c
func (c Color) B() float32 {
	return c.Z()
}

func (c Color) Add(o Color) Color {
	return Color{c.Vec3.Add(o.Vec3)}
}

func (c Color) Mul(t float32) Color {
	return Color{c.Vec3.Mul(t)}
}

func (c Color) MulVec(v geo.Vec3) Color {
	return Color{geo.NewVec3(c.R()*v.X(), c.G()*v.Y(), c.B()*v.Z())}
}
