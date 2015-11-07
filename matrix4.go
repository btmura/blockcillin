package main

import (
	"fmt"
	"math"
)

// Matrix4 is a 4x4 matrix.
type Matrix4 [16]float32

func NewPerspectiveMatrix(fovRadians, aspect, near, far float32) *Matrix4 {
	f := float32(math.Tan(math.Pi*0.5 - 0.5*float64(fovRadians)))
	rangeInv := 1.0 / (near - far)
	return &Matrix4{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (near + far) * rangeInv, -1,
		0, 0, near * far * rangeInv * 2, 0,
	}
}

func NewViewMatrix(cameraPosition, target, up *Vector3) *Matrix4 {
	cameraMatrix := NewLookAtMatrix(cameraPosition, target, up)
	return cameraMatrix.Inverse()
}

func NewLookAtMatrix(cameraPosition, target, up *Vector3) *Matrix4 {
	zAxis := cameraPosition.Sub(target).Normalize()
	xAxis := up.Cross(zAxis)
	yAxis := zAxis.Cross(xAxis)
	return &Matrix4{
		xAxis.X, xAxis.Y, xAxis.Z, 0,
		yAxis.X, yAxis.Y, yAxis.Z, 0,
		zAxis.X, zAxis.Y, zAxis.Z, 0,
		cameraPosition.X, cameraPosition.Y, cameraPosition.Z, 1,
	}
}

func NewTranslationMatrix(x, y, z float32) *Matrix4 {
	return &Matrix4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		x, y, z, 1,
	}
}

func NewXRotationMatrix(radians float32) *Matrix4 {
	c := float32(math.Cos(float64(radians)))
	s := float32(math.Sin(float64(radians)))
	return &Matrix4{
		1, 0, 0, 0,
		0, c, s, 0,
		0, -s, c, 0,
		0, 0, 0, 1,
	}
}

func NewYRotationMatrix(radians float32) *Matrix4 {
	c := float32(math.Cos(float64(radians)))
	s := float32(math.Sin(float64(radians)))
	return &Matrix4{
		c, 0, -s, 0,
		0, 1, 0, 0,
		s, 0, c, 0,
		0, 0, 0, 1,
	}
}

func NewZRotationMatrix(radians float32) *Matrix4 {
	c := float32(math.Cos(float64(radians)))
	s := float32(math.Sin(float64(radians)))
	return &Matrix4{
		c, s, 0, 0,
		-s, c, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

func NewScaleMatrix(sx, sy, sz float32) *Matrix4 {
	return &Matrix4{
		sx, 0, 0, 0,
		0, sy, 0, 0,
		0, 0, sz, 0,
		0, 0, 0, 1,
	}
}

func (m *Matrix4) Mult(n *Matrix4) *Matrix4 {
	return &Matrix4{
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

func (m *Matrix4) Inverse() *Matrix4 {
	m00 := m[0*4+0]
	m01 := m[0*4+1]
	m02 := m[0*4+2]
	m03 := m[0*4+3]
	m10 := m[1*4+0]
	m11 := m[1*4+1]
	m12 := m[1*4+2]
	m13 := m[1*4+3]
	m20 := m[2*4+0]
	m21 := m[2*4+1]
	m22 := m[2*4+2]
	m23 := m[2*4+3]
	m30 := m[3*4+0]
	m31 := m[3*4+1]
	m32 := m[3*4+2]
	m33 := m[3*4+3]
	tmp_0 := m22 * m33
	tmp_1 := m32 * m23
	tmp_2 := m12 * m33
	tmp_3 := m32 * m13
	tmp_4 := m12 * m23
	tmp_5 := m22 * m13
	tmp_6 := m02 * m33
	tmp_7 := m32 * m03
	tmp_8 := m02 * m23
	tmp_9 := m22 * m03
	tmp_10 := m02 * m13
	tmp_11 := m12 * m03
	tmp_12 := m20 * m31
	tmp_13 := m30 * m21
	tmp_14 := m10 * m31
	tmp_15 := m30 * m11
	tmp_16 := m10 * m21
	tmp_17 := m20 * m11
	tmp_18 := m00 * m31
	tmp_19 := m30 * m01
	tmp_20 := m00 * m21
	tmp_21 := m20 * m01
	tmp_22 := m00 * m11
	tmp_23 := m10 * m01

	t0 := (tmp_0*m11 + tmp_3*m21 + tmp_4*m31) - (tmp_1*m11 + tmp_2*m21 + tmp_5*m31)
	t1 := (tmp_1*m01 + tmp_6*m21 + tmp_9*m31) - (tmp_0*m01 + tmp_7*m21 + tmp_8*m31)
	t2 := (tmp_2*m01 + tmp_7*m11 + tmp_10*m31) - (tmp_3*m01 + tmp_6*m11 + tmp_11*m31)
	t3 := (tmp_5*m01 + tmp_8*m11 + tmp_11*m21) - (tmp_4*m01 + tmp_9*m11 + tmp_10*m21)

	d := 1.0 / (m00*t0 + m10*t1 + m20*t2 + m30*t3)

	return &Matrix4{
		d * t0,
		d * t1,
		d * t2,
		d * t3,
		d * ((tmp_1*m10 + tmp_2*m20 + tmp_5*m30) - (tmp_0*m10 + tmp_3*m20 + tmp_4*m30)),
		d * ((tmp_0*m00 + tmp_7*m20 + tmp_8*m30) - (tmp_1*m00 + tmp_6*m20 + tmp_9*m30)),
		d * ((tmp_3*m00 + tmp_6*m10 + tmp_11*m30) - (tmp_2*m00 + tmp_7*m10 + tmp_10*m30)),
		d * ((tmp_4*m00 + tmp_9*m10 + tmp_10*m20) - (tmp_5*m00 + tmp_8*m10 + tmp_11*m20)),
		d * ((tmp_12*m13 + tmp_15*m23 + tmp_16*m33) - (tmp_13*m13 + tmp_14*m23 + tmp_17*m33)),
		d * ((tmp_13*m03 + tmp_18*m23 + tmp_21*m33) - (tmp_12*m03 + tmp_19*m23 + tmp_20*m33)),
		d * ((tmp_14*m03 + tmp_19*m13 + tmp_22*m33) - (tmp_15*m03 + tmp_18*m13 + tmp_23*m33)),
		d * ((tmp_17*m03 + tmp_20*m13 + tmp_23*m23) - (tmp_16*m03 + tmp_21*m13 + tmp_22*m23)),
		d * ((tmp_14*m22 + tmp_17*m32 + tmp_13*m12) - (tmp_16*m32 + tmp_12*m12 + tmp_15*m22)),
		d * ((tmp_20*m32 + tmp_12*m02 + tmp_19*m22) - (tmp_18*m22 + tmp_21*m32 + tmp_13*m02)),
		d * ((tmp_18*m12 + tmp_23*m32 + tmp_15*m02) - (tmp_22*m32 + tmp_14*m02 + tmp_19*m12)),
		d * ((tmp_22*m22 + tmp_16*m02 + tmp_21*m12) - (tmp_20*m12 + tmp_23*m22 + tmp_17*m02)),
	}
}

func (m *Matrix4) String() string {
	return fmt.Sprintf("%8.2f %8.2f %8.2f %8.2f\n%8.2f %8.2f %8.2f %8.2f\n%8.2f %8.2f %8.2f %8.2f\n%8.2f %8.2f %8.2f %8.2f", m[0], m[1], m[2], m[3], m[4], m[5], m[6], m[7], m[8], m[9], m[10], m[11], m[12], m[13], m[14], m[15])
}
