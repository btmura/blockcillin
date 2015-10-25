package main

import "math"

// Matrix4 is a 4x4 matrix.
type Matrix4 [16]float32

func makePerspective(fovRadians, aspect, near, far float64) Matrix4 {
	f := math.Tan(math.Pi*0.5 - 0.5*fovRadians)
	rangeInv := 1.0 / (near - far)
	return Matrix4{
		float32(f / aspect), 0, 0, 0,
		0, float32(f), 0, 0,
		0, 0, float32((near + far) * rangeInv), -1,
		0, 0, float32(near * far * rangeInv * 2), 0,
	}
}

func makeTranslationMatrix(x, y, z float32) Matrix4 {
	return Matrix4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		x, y, z, 1,
	}
}

func makeXRotationMatrix(radians float64) Matrix4 {
	c := float32(math.Cos(radians))
	s := float32(math.Sin(radians))
	return Matrix4{
		1, 0, 0, 0,
		0, c, s, 0,
		0, -s, c, 0,
		0, 0, 0, 1,
	}
}

func makeYRotationMatrix(radians float64) Matrix4 {
	c := float32(math.Cos(radians))
	s := float32(math.Sin(radians))
	return Matrix4{
		c, 0, -s, 0,
		0, 1, 0, 0,
		s, 0, c, 0,
		0, 0, 0, 1,
	}
}

func makeZRotationMatrix(radians float64) Matrix4 {
	c := float32(math.Cos(radians))
	s := float32(math.Sin(radians))
	return Matrix4{
		c, s, 0, 0,
		-s, c, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

func makeScaleMatrix(sx, sy, sz float32) Matrix4 {
	return Matrix4{
		sx, 0, 0, 0,
		0, sy, 0, 0,
		0, 0, sz, 0,
		0, 0, 0, 1,
	}
}

func multipleMatrices(m, n Matrix4) Matrix4 {
	return Matrix4{
		m[0]*n[0] + m[1]*n[4] + m[2]*n[8] + m[3]*n[12],
		m[0]*n[1] + m[1]*n[5] + m[2]*n[9] + m[3]*n[13],
		m[0]*n[2] + m[1]*n[6] + m[2]*n[10] + m[3]*n[14],
		m[0]*n[3] + m[1]*n[7] + m[2]*n[11] + m[3]*n[15],

		m[4]*n[0] + m[5]*n[4] + m[6]*n[8] + m[7]*n[12],
		m[4]*n[1] + m[5]*n[5] + m[6]*n[9] + m[7]*n[13],
		m[4]*n[2] + m[5]*n[6] + m[6]*n[10] + m[7]*n[14],
		m[4]*n[3] + m[5]*n[7] + m[6]*n[11] + m[7]*n[15],

		m[8]*n[0] + m[9]*n[4] + m[10]*n[8] + m[11]*n[12],
		m[8]*n[1] + m[9]*n[5] + m[10]*n[9] + m[11]*n[13],
		m[8]*n[2] + m[9]*n[6] + m[10]*n[10] + m[11]*n[14],
		m[8]*n[3] + m[9]*n[7] + m[10]*n[11] + m[11]*n[15],

		m[12]*n[0] + m[13]*n[4] + m[14]*n[8] + m[15]*n[12],
		m[12]*n[1] + m[13]*n[5] + m[14]*n[9] + m[15]*n[13],
		m[12]*n[2] + m[13]*n[6] + m[14]*n[10] + m[15]*n[14],
		m[12]*n[3] + m[13]*n[7] + m[14]*n[11] + m[15]*n[15],
	}
}
