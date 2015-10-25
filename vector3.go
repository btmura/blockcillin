package main

import "math"

// Vector3 is a vector with x, y, and z.
type Vector3 struct {
	x float32
	y float32
	z float32
}

func (v *Vector3) Sub(w *Vector3) *Vector3 {
	return &Vector3{v.x - w.x, v.y - w.y, v.z - w.z}
}

func (v *Vector3) Length() float32 {
	return float32(math.Sqrt(float64(v.x*v.x + v.y*v.y + v.z*v.z)))
}

func (v *Vector3) Cross(w *Vector3) *Vector3 {
	return &Vector3{
		v.y*w.z - v.z*w.y,
		v.z*w.x - v.x*w.z,
		v.x*w.y - v.y*w.x,
	}
}

func (v *Vector3) Normalize() *Vector3 {
	l := v.Length()
	if l > 0.00001 {
		return &Vector3{v.x / l, v.y / l, v.z / l}
	}
	return &Vector3{}
}
