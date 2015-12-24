package renderer

import "math"

func linear(time, start, change, duration float32) float32 {
	t := time / duration
	return change*t + start
}

func easeInCubic(time, start, change, duration float32) float32 {
	t := time / duration
	return change*t*t*t + start
}

func easeInExpo(time, start, change, duration float32) float32 {
	t := time / duration
	return change*float32(math.Pow(2, float64(10*(t-1)))) + start
}

func easeOutCubic(time, start, change, duration float32) float32 {
	t := time/duration - 1
	return change*(t*t*t+1) + start
}

func easeOutCubic2(t, start, change float32) float32 {
	t--
	return change*(t*t*t+1) + start
}

func pulse(t, start, amplitude, cycles float32) float32 {
	return start + amplitude*float32(math.Sin(float64(cycles*t)))
}

func toRadians(degrees float32) float32 {
	return degrees * float32(math.Pi) / 180
}
