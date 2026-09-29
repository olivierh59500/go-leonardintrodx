package source

import "math"

// Oscillator is the DirectX program's spaced sine controller. Time is seconds;
// the program multiplies elapsed time by 85 before advancing each phase.
type Oscillator struct {
	Amplitude, Rate, Spacing float64
}

func (o Oscillator) At(seconds float64, item int) float64 {
	phase := seconds*85*o.Rate + float64(item)*o.Spacing
	return o.Amplitude * math.Sin(math.Mod(phase, 2*math.Pi))
}
