package source

import (
	"math"

	"github.com/olivierh59500/democonstructionkit/motion"
)

// Oscillator is the DirectX program's spaced sine controller. Time is seconds;
// the program multiplies elapsed time by 85 before advancing each phase.
type Oscillator struct {
	Amplitude, Rate, Spacing float64
}

func (o Oscillator) At(seconds float64, item int) float64 {
	phase := seconds*85*o.Rate + float64(item)*o.Spacing
	return o.Amplitude * math.Sin(math.Mod(phase, 2*math.Pi))
}

// Harmonic supplies the authored coefficients and native phase association to
// DCK. The caller scales seconds by 85 before sampling the shared controller.
func (o Oscillator) Harmonic() motion.IndexedHarmonic {
	return motion.IndexedHarmonic{Amplitude: o.Amplitude, Rate: o.Rate, IndexRate: o.Spacing,
		PhasePeriod: 2 * math.Pi, RoundProduct: true, FusedIndexPhase: true}
}
