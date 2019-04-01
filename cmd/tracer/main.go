package main

import (
	"fmt"
	"math"

	"github.com/robquant/tracer/pkg"
)

func color(r *pkg.Ray) pkg.Color {
	unitDirection := r.Dir().Normed()
	t := 0.5 * (unitDirection.Y() + 1.0)
	c1 := pkg.NewVec3(1.0, 1.0, 1.0).Mul(1.0 - t)
	c2 := pkg.NewVec3(0.5, 0.7, 1.0).Mul(t)
	return pkg.Color{Vec3: c1.Add(c2)}
}

func main() {
	nx := 200
	ny := 100
	fmt.Printf("P3 %d %d\n255\n", nx, ny)
	lowerLeftCorner := pkg.NewVec3(-2.0, -1.0, -1.0)
	horizontal := pkg.NewVec3(4.0, 0.0, 0.0)
	vertical := pkg.NewVec3(0.0, 2.0, 0.0)
	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			u := float64(i) / float64(nx)
			v := float64(j) / float64(ny)
			dir := lowerLeftCorner.Add(horizontal.Mul(u).Add(vertical.Mul(v)))
			r := pkg.NewRay(pkg.Origin, dir)
			col := color(&r)
			ir := int(math.Round(255 * col.R()))
			ig := int(math.Round(255 * col.G()))
			ib := int(math.Round(255 * col.B()))
			fmt.Printf("%d %d %d\n", ir, ig, ib)
		}
	}
}
