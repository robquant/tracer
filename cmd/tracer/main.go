package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"

	"github.com/chewxy/math32"
	"github.com/robquant/tracer/pkg/geo"
	"github.com/robquant/tracer/pkg/tracer"
)

func colorAt(r *geo.Ray, world tracer.Hitable, depth int) tracer.Color {
	if hit, rec := world.Hit(r, 0.001, math.MaxFloat32); hit {
		if depth < 50 {
			didscatter, attenuation, scattered := rec.Material().Scatter(r, &rec)
			if didscatter {
				return colorAt(&scattered, world, depth+1).MulVec(attenuation)
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
	choose := rand.Float32()
	if choose < 0.8 {
		ar := rand.Float32() * rand.Float32()
		ag := rand.Float32() * rand.Float32()
		ab := rand.Float32() * rand.Float32()
		return tracer.NewLambertian(ar, ag, ab)
	} else if choose < 0.95 {
		return tracer.NewMetal(0.5*(1+rand.Float32()), 0.5*(1+rand.Float32()), 0.5*(1+rand.Float32()), 0.5*rand.Float32())
	}
	return tracer.NewDielectric(1.5)
}

func randomScene() tracer.HitableList {
	scene := tracer.NewHitableList()
	scene = append(scene, tracer.NewSphere(geo.NewVec3(0, -1000, 0), 1000, tracer.NewLambertian(0.5, 0.5, 0.5)))
	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			center := geo.NewVec3(float32(a)+0.9*rand.Float32(), 0.2, float32(b)+rand.Float32())
			if center.Sub(geo.NewVec3(4, 0.2, 0)).Len() > 0.9 {
				scene = append(scene, tracer.NewSphere(center, 0.2, randomMaterial()))
			}
		}
	}
	scene = append(scene, tracer.NewSphere(geo.NewVec3(0, 1, 0), 1.0, tracer.NewDielectric(1.5)))
	scene = append(scene, tracer.NewSphere(geo.NewVec3(-4, 1, 0), 1.0, tracer.NewLambertian(0.4, 0.2, 0.1)))
	scene = append(scene, tracer.NewSphere(geo.NewVec3(4, 1, 0), 1.0, tracer.NewMetal(0.7, 0.6, 0.5, 0)))
	return scene
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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

	var nx, ny, ns, np int
	var outfname string
	flag.IntVar(&nx, "nx", 600, "X resolution")
	flag.IntVar(&ny, "ny", 400, "Y resolution")
	flag.IntVar(&ns, "ns", 10, "samples per pixel")
	flag.IntVar(&np, "np", runtime.NumCPU(), "number of parallel renderers")
	flag.StringVar(&outfname, "out", "image.png", "output file name")
	flag.Parse()

	img := image.NewRGBA(image.Rect(0, 0, nx, ny))
	lookAt := geo.NewVec3(0, 0, 0)
	lookFrom := geo.NewVec3(13, 2, 3)
	distToFocus := float32(10.0)
	aperture := float32(1 / 10.0)
	camera := tracer.NewCamera(lookFrom, lookAt, geo.UnitY, 20, float32(nx)/float32(ny), aperture, distToFocus)
	world := tracer.NewBvhNodeFromList(randomScene())

	wg := sync.WaitGroup{}
	blockQueue := make(chan image.Rectangle)
	for cpu := 0; cpu < np; cpu++ {
		wg.Add(1)
		go func(queue <-chan image.Rectangle) {
			for block := range queue {
				for y := block.Min.Y; y < block.Max.Y; y++ {
					for x := block.Min.X; x < block.Max.X; x++ {
						col := tracer.NewColor(0, 0, 0)
						for s := 0; s < ns; s++ {
							u := (float32(x) + rand.Float32()) / float32(nx)
							v := (float32(ny-y) + rand.Float32()) / float32(ny)
							ray := camera.GetRay(u, v)
							col = col.Add(colorAt(ray, &world, 0))
						}
						col.Scale(1. / float32(ns))
						col = tracer.NewColor(math32.Sqrt(col.R()), math32.Sqrt(col.G()), math32.Sqrt(col.B()))
						ir := uint8(math.Round(float64(255 * col.R())))
						ig := uint8(math.Round(float64(255 * col.G())))
						ib := uint8(math.Round(float64(255 * col.B())))
						img.SetRGBA(x, y, color.RGBA{ir, ig, ib, 255})
					}
				}
			}
			wg.Done()
		}(blockQueue)
	}
	for x := 0; x <= nx; x += 50 {
		for y := 0; y <= ny; y += 50 {
			r := image.Rect(x, y, min(x+50, nx), min(y+50, ny))
			blockQueue <- r
		}
	}
	close(blockQueue)
	wg.Wait()

	f, err := os.Create(outfname)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d\n", tracer.SphereCounter)
}
