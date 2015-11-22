package main

import "math"

func linear(time, start, change, duration float32) float32 {
	return change*time/duration + start
}

func toRadians(degrees float32) float32 {
	return degrees * float32(math.Pi) / 180
}
