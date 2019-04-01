package tracer

import (
	"github.com/robquant/tracer/pkg/geo"
)

var DefaultCamera = Camera{
	origin:          geo.Origin,
	lowerLeftCorner: geo.NewVec3(-2.0, -1.0, -1.0),
	horizontal:      geo.NewVec3(4.0, 0.0, 0.0),
	vertical:        geo.NewVec3(0.0, 2.0, 0.0)}

type Camera struct {
	origin          geo.Vec3
	lowerLeftCorner geo.Vec3
	horizontal      geo.Vec3
	vertical        geo.Vec3
}

func (c *Camera) GetRay(u, v float64) *geo.Ray {
	dir := c.lowerLeftCorner.Add(c.horizontal.Mul(u).Add(c.vertical.Mul(v)))
	r := geo.NewRay(geo.Origin, dir)
	return &r
}
