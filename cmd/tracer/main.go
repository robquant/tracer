package main

import (
	"fmt"
	"math"

	"github.com/robquant/tracer/pkg"
)

func hitSphere(center pkg.Vec3, radius float64, r *pkg.Ray) float64 {
	oc := r.Orig().Sub(center)
	a := r.Dir().Dot(r.Dir())
	b := 2.0 * r.Dir().Dot(oc)
	c := oc.Dot(oc) - radius*radius
	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return -1.0
	} else {
		return (-b - math.Sqrt(discriminant)) / (2.0 * a)
	}
}

func color(r *pkg.Ray) pkg.Color {
	negZ := pkg.NewVec3(0, 0, -1)
	t := hitSphere(negZ, 0.5, r)
	if t > 0 {
		N := r.At(t).Sub(negZ)
		(&N).Normalize()
		return pkg.NewColor((N.X()+1)*0.5, (N.Y()+1)*0.5, (N.Z()+1)*0.5)
	}
	unitDirection := r.Dir().Normed()
	t = 0.5 * (unitDirection.Y() + 1.0)
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
