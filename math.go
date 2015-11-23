package main

import "math"

func linear(time, start, change, duration float32) float32 {
	t := time / duration
	return change*t + start
}

func easeOutQuad(time, start, change, duration float32) float32 {
	t := time / duration
	return -change*t*(t-2) + start
}

func toRadians(degrees float32) float32 {
	return degrees * float32(math.Pi) / 180
}
