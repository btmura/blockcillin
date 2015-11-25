package main

import "math"

func linear(time, start, change, duration float32) float32 {
	t := time / duration
	return change*t + start
}

func easeInCubic(time, start, change, duration float32) float32 {
	t := time / duration
	return change*t*t*t + start
}

func pulse(time, start, amplitude, cycles float32) float32 {
	return start + amplitude*float32(math.Sin(float64(cycles*time)))
}

func toRadians(degrees float32) float32 {
	return degrees * float32(math.Pi) / 180
}
