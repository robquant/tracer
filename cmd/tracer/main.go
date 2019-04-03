package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/robquant/tracer/pkg/geo"
	"github.com/robquant/tracer/pkg/tracer"
)

func color(r *geo.Ray, world tracer.Hitable, depth int) tracer.Color {
	if hit, rec := world.Hit(r, 0.001, math.MaxFloat64); hit {
		if depth < 50 {
			didscatter, attenuation, scattered := rec.Material().Scatter(r, &rec)
			if didscatter {
				return color(&scattered, world, depth+1).MulVec(attenuation)
			} else {
				return tracer.Black
			}
		}
	}
	unitDirection := r.Dir().Normed()
	t := 0.5 * (unitDirection.Y() + 1.0)
	c1 := geo.NewVec3(1.0, 1.0, 1.0).Mul(1.0 - t)
	c2 := geo.NewVec3(0.5, 0.7, 1.0).Mul(t)
	return tracer.Color{Vec3: c1.Add(c2)}
}

func main() {
	nx := 400
	ny := 200
	ns := 100
	fmt.Printf("P3 %d %d\n255\n", nx, ny)
	camera := tracer.DefaultCamera
	world := tracer.NewHitableList()
	world = append(world, tracer.NewSphere(geo.NewVec3(0, 0, -1), 0.5, tracer.NewLambertian(0.8, 0.3, 0.3)))
	world = append(world, tracer.NewSphere(geo.NewVec3(0, -100.5, -1), 100, tracer.NewLambertian(0.8, 0.8, 0)))
	world = append(world, tracer.NewSphere(geo.NewVec3(1, 0, -1), 0.5, tracer.NewMetal(0.8, 0.6, 0.2, 0.3)))
	world = append(world, tracer.NewSphere(geo.NewVec3(-1, 0, -1), 0.5, tracer.NewMetal(0.8, 0.8, 0.8, 1.0)))
	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			col := tracer.NewColor(0, 0, 0)
			for s := 0; s < ns; s++ {
				u := (float64(i) + rand.Float64()) / float64(nx)
				v := (float64(j) + rand.Float64()) / float64(ny)
				ray := camera.GetRay(u, v)
				col = col.Add(color(ray, world, 0))
			}
			col.Scale(1. / float64(ns))
			col = tracer.NewColor(math.Sqrt(col.R()), math.Sqrt(col.G()), math.Sqrt(col.B()))
			ir := int(math.Round(255 * col.R()))
			ig := int(math.Round(255 * col.G()))
			ib := int(math.Round(255 * col.B()))
			fmt.Printf("%d %d %d\n", ir, ig, ib)
		}
	}
}
