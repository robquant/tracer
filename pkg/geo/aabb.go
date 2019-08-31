package geo

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

type Aabb struct {
	min, max Vec3
}

var EmptyBox Aabb = Aabb{Origin, Origin}

func NewAabb(min, max Vec3) *Aabb {
	return &Aabb{min: min, max: max}
}

func (a *Aabb) Min() Vec3 {
	return a.min
}

func (a *Aabb) Max() Vec3 {
	return a.max
}

func (a *Aabb) Hit(r *Ray, tmin, tmax float32) bool {
	rDir := r.Dir()
	rOrig := r.Orig()
	var invD, t0, t1 float32

	invD = float32(1.0) / rDir.X()
	t0 = (a.min.X() - rOrig.X()) * invD
	t1 = (a.max.X() - rOrig.X()) * invD
	if invD < 0 {
		t0, t1 = t1, t0
	}
	if t0 > tmin {
		tmin = t0
	}
	if t1 < tmax {
		tmax = t1
	}
	if tmax <= tmin {
		return false
	}
	invD = float32(1.0) / rDir.Y()
	t0 = (a.min.Y() - rOrig.Y()) * invD
	t1 = (a.max.Y() - rOrig.Y()) * invD
	if invD < 0 {
		t0, t1 = t1, t0
	}
	if t0 > tmin {
		tmin = t0
	}
	if t1 < tmax {
		tmax = t1
	}
	if tmax <= tmin {
		return false
	}
	invD = float32(1.0) / rDir.Z()
	t0 = (a.min.Z() - rOrig.Z()) * invD
	t1 = (a.max.Z() - rOrig.Z()) * invD
	if invD < 0 {
		t0, t1 = t1, t0
	}
	if t0 > tmin {
		tmin = t0
	}
	if t1 < tmax {
		tmax = t1
	}
	if tmax <= tmin {
		return false
	}
	return true
}

func SurroundingBox(box0, box1 Aabb) Aabb {
	var small, big Vec3
	if box0 == EmptyBox {
		small = box1.min
		big = box1.max
	} else if box1 == EmptyBox {
		small = box0.min
		big = box0.max
	} else {
		small = NewVec3(min(box0.min.X(), box1.min.X()),
			min(box0.min.Y(), box1.min.Y()),
			min(box0.min.Z(), box1.min.Z()))
		big = NewVec3(max(box0.max.X(), box1.max.X()),
			max(box0.max.Y(), box1.max.Y()),
			max(box0.max.Z(), box1.max.Z()))
	}
	return Aabb{small, big}
}

func (a *Aabb) LongestAxis() int {
	lenX := a.max.X() - a.min.X()
	lenY := a.max.Y() - a.min.Y()
	lenZ := a.max.Z() - a.min.Z()
	if lenX >= lenY {
		if lenX >= lenZ {
			return 0
		}
		return 2
	} else if lenY >= lenZ {
		return 1
	}
	return 2
}

func (a *Aabb) Area() float32 {
	size := a.max.Sub(a.min)
	return 2 * (size.x*size.y + size.x*size.z + size.y*size.z)
}
