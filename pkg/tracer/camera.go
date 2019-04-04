package tracer

import (
	"math"
	"math/rand"

	"github.com/robquant/tracer/pkg/geo"
)

type Camera struct {
	origin          geo.Vec3
	lowerLeftCorner geo.Vec3
	horizontal      geo.Vec3
	vertical        geo.Vec3
	u, v, w         geo.Vec3
	lensRadius      float64
}

// NewCamera constructs a new Camera from the vertical
// field of view in degrees, and the aspect ratio
func NewCamera(lookFrom, lookAt, vUp geo.Vec3, vertFov, aspectRatio, aperture, focusDist float64) *Camera {
	lensRadius := aperture / 2
	theta := math.Pi / 180 * vertFov
	halfHeight := math.Tan(theta / 2)
	halfWidth := aspectRatio * halfHeight
	w := lookFrom.Sub(lookAt).Normed()
	u := vUp.Cross(w).Normed()
	v := w.Cross(u)
	lowerLeftCorner := lookFrom
	lowerLeftCorner = lowerLeftCorner.Sub(u.Mul(halfWidth * focusDist))
	lowerLeftCorner = lowerLeftCorner.Sub(v.Mul(halfHeight * focusDist))
	lowerLeftCorner = lowerLeftCorner.Sub(w.Mul(focusDist))
	horizontal := u.Mul(2 * halfWidth * focusDist)
	vertical := v.Mul(2 * halfHeight * focusDist)

	return &Camera{lookFrom, lowerLeftCorner, horizontal, vertical, u, v, w, lensRadius}
}

func randomInUnitDisk() geo.Vec3 {
	vec := geo.NewVec3(1.0, 1.0, 0.0)
	diag2D := geo.NewVec3(1, 1, 0)
	for vec.LenSq() >= 1.0 {
		vec = geo.NewVec3(rand.Float64(), rand.Float64(), 0).Mul(2.0).Sub(diag2D)
	}
	return vec
}

func (c *Camera) GetRay(s, t float64) *geo.Ray {
	rd := randomInUnitDisk().Mul(c.lensRadius)
	offset := c.u.Mul(rd.X()).Add(c.v.Mul(rd.Y()))
	dir := c.lowerLeftCorner.Add(c.horizontal.Mul(s)).Add(c.vertical.Mul(t)).Sub(c.origin).Sub(offset)
	r := geo.NewRay(c.origin.Add(offset), dir)
	return &r
}
