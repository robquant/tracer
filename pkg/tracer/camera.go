package tracer

import (
	"math"

	"github.com/robquant/tracer/pkg/geo"
)

type Camera struct {
	origin          geo.Vec3
	lowerLeftCorner geo.Vec3
	horizontal      geo.Vec3
	vertical        geo.Vec3
}

// NewCamera constructs a new Camera from the vertical
// field of view in degrees, and the aspect ratio
func NewCamera(lookFrom, lookAt, vUp geo.Vec3, vertFov, aspectRatio float64) *Camera {
	theta := math.Pi / 180 * vertFov
	halfHeight := math.Tan(theta / 2)
	halfWidth := aspectRatio * halfHeight
	w := lookFrom.Sub(lookAt).Normed()
	u := vUp.Cross(w).Normed()
	v := w.Cross(u)
	lowerLeftCorner := lookFrom.Sub(u.Mul(halfWidth)).Sub(v.Mul(halfHeight)).Sub(w)
	horizontal := u.Mul(2 * halfWidth)
	vertical := v.Mul(2 * halfHeight)

	return &Camera{lookFrom, lowerLeftCorner, horizontal, vertical}
}

func (c *Camera) GetRay(s, t float64) *geo.Ray {
	dir := c.lowerLeftCorner.Add(c.horizontal.Mul(s)).Add(c.vertical.Mul(t)).Sub(c.origin)
	r := geo.NewRay(c.origin, dir)
	return &r
}
