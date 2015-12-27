package renderer

import "math"

func linear(t, start, change float32) float32 {
	return change*t + start
}

func easeInExpo(t, start, change float32) float32 {
	return change*float32(math.Pow(2, float64(10*(t-1)))) + start
}

func easeOutCubic(t, start, change float32) float32 {
	t--
	return change*(t*t*t+1) + start
}

func pulse(t, start, amplitude, cycles float32) float32 {
	return start + amplitude*float32(math.Sin(float64(cycles*t)))
}
