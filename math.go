package main

import "math"

func toRadians(degrees float32) float32 {
	return degrees * float32(math.Pi) / 180
}
