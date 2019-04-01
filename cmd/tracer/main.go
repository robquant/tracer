package main

import (
	"fmt"
	"math"

	"github.com/robquant/tracer/pkg/geo"
	"github.com/robquant/tracer/pkg/tracer"
)

func color(r *geo.Ray, world tracer.Hitable) tracer.Color {
	if hit, rec := world.Hit(r, 0.0, math.MaxFloat64); hit {
		normal := rec.Normal()
		return tracer.Color{geo.NewVec3(normal.X()+1, normal.Y()+1, normal.Z()+1).Mul(0.5)}
	}
	unitDirection := r.Dir().Normed()
	t := 0.5 * (unitDirection.Y() + 1.0)
	c1 := geo.NewVec3(1.0, 1.0, 1.0).Mul(1.0 - t)
	c2 := geo.NewVec3(0.5, 0.7, 1.0).Mul(t)
	return tracer.Color{Vec3: c1.Add(c2)}
}

func main() {
	nx := 200
	ny := 100
	fmt.Printf("P3 %d %d\n255\n", nx, ny)
	lowerLeftCorner := geo.NewVec3(-2.0, -1.0, -1.0)
	horizontal := geo.NewVec3(4.0, 0.0, 0.0)
	vertical := geo.NewVec3(0.0, 2.0, 0.0)
	world := tracer.NewHitableList()
	world = append(world, tracer.NewSphere(geo.NewVec3(0, 0, -1), 0.5))
	world = append(world, tracer.NewSphere(geo.NewVec3(0, -100.5, -1), 100))
	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			u := float64(i) / float64(nx)
			v := float64(j) / float64(ny)
			dir := lowerLeftCorner.Add(horizontal.Mul(u).Add(vertical.Mul(v)))
			r := geo.NewRay(geo.Origin, dir)
			col := color(&r, world)
			ir := int(math.Round(255 * col.R()))
			ig := int(math.Round(255 * col.G()))
			ib := int(math.Round(255 * col.B()))
			fmt.Printf("%d %d %d\n", ir, ig, ib)
		}
	}
}
