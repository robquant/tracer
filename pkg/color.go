package pkg

var Red = Color{UnitX}
var Green = Color{UnitY}
var Blue = Color{UnitZ}

// Color is an RGB color triple using floats
type Color struct {
	Vec3
}

// NewColor creates a new Color from RGB float values
func NewColor(r, g, b float64) Color {
	return Color{Vec3{r, g, b}}
}

// R returns the red component of c
func (c Color) R() float64 {
	return c.x
}

// G returns the green component of c
func (c Color) G() float64 {
	return c.y
}

// B returns the blue component of c
func (c Color) B() float64 {
	return c.z
}
