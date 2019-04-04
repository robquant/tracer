package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"

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

func randomMaterial() tracer.Material {
	choose := rand.Float64()
	if choose < 0.8 {
		ar := rand.Float64() * rand.Float64()
		ag := rand.Float64() * rand.Float64()
		ab := rand.Float64() * rand.Float64()
		return tracer.NewLambertian(ar, ag, ab)
	} else if choose < 0.95 {
		return tracer.NewMetal(0.5*(1+rand.Float64()), 0.5*(1+rand.Float64()), 0.5*(1+rand.Float64()), 0.5*rand.Float64())
	}
	return tracer.NewDielectric(1.5)
}

func randomScene() *tracer.HitableList {
	scene := tracer.NewHitableList()
	scene = append(scene, tracer.NewSphere(geo.NewVec3(0, -1000, 0), 1000, tracer.NewLambertian(0.5, 0.5, 0.5)))
	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			center := geo.NewVec3(float64(a)+0.9*rand.Float64(), 0.2, float64(b)+rand.Float64())
			if center.Sub(geo.NewVec3(4, 0.2, 0)).Len() > 0.9 {
				scene = append(scene, tracer.NewSphere(center, 0.2, randomMaterial()))
			}
		}
	}
	scene = append(scene, tracer.NewSphere(geo.NewVec3(0, 1, 0), 1.0, tracer.NewDielectric(1.5)))
	scene = append(scene, tracer.NewSphere(geo.NewVec3(-4, 1, 0), 1.0, tracer.NewLambertian(0.4, 0.2, 0.1)))
	scene = append(scene, tracer.NewSphere(geo.NewVec3(4, 1, 0), 1.0, tracer.NewMetal(0.7, 0.6, 0.5, 0)))
	return &scene
}

func main() {
	if pr := os.Getenv("CPUPROFILE"); pr != "" {
		p, err := os.Create(pr)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer p.Close()
		if err := pprof.StartCPUProfile(p); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	nx := 600
	ny := 400
	ns := 10
	fmt.Printf("P3 %d %d\n255\n", nx, ny)
	lookAt := geo.NewVec3(0, 0, 0)
	lookFrom := geo.NewVec3(13, 2, 3)
	distToFocus := 10.0
	aperture := 1 / 10.0
	camera := tracer.NewCamera(lookFrom, lookAt, geo.UnitY, 20, float64(nx)/float64(ny), aperture, distToFocus)
	world := randomScene()
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
