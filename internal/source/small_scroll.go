package source

import "github.com/olivierh59500/democonstructionkit/motion"

const (
	SmallCellSide    = 20
	SmallColumnPitch = 28
	SmallRowPitch    = 30
	SmallWorldOffset = 470
)

// SmallScrollRow is the world-space center of one row of the original
// eight-row cubelet scroller. The column transport supplies X separately.
type SmallScrollRow struct{ Y, Z float64 }

var smallGlobalY = [2]Oscillator{
	{Amplitude: 70, Rate: .0067}, {Amplitude: 70, Rate: .0239},
}
var smallGlobalZ = [2]Oscillator{
	{Amplitude: 30, Rate: .0097}, {Amplitude: 30, Rate: .0339},
}
var smallRowY = Oscillator{Amplitude: 30, Rate: .103, Spacing: .222}
var smallRowZ = Oscillator{Amplitude: 10, Rate: .1, Spacing: .17}

// NewSmallScrollProfile keeps the source's grouped global waves, row spacing
// and local waves as an editable recipe. Its two outputs are world Y and Z.
func NewSmallScrollProfile() (*motion.HarmonicRowProfile, error) {
	global := 0
	rowY := smallRowY
	rowY.Amplitude = -rowY.Amplitude
	return motion.NewHarmonicRowProfile(motion.HarmonicRowProfileConfig{
		Count: SmallGlyphHeight, ClockScale: [2]float64{85},
		Stages: []motion.HarmonicRowStage{
			{Motion: motion.HarmonicFormationConfig{
				X: []motion.IndexedHarmonic{smallGlobalY[0].Harmonic(), smallGlobalY[1].Harmonic()},
				Y: []motion.IndexedHarmonic{smallGlobalZ[0].Harmonic(), smallGlobalZ[1].Harmonic()},
			}, Index: &global},
			{Motion: motion.HarmonicFormationConfig{Spacing: motion.Point{X: -SmallRowPitch}}},
			{Motion: motion.HarmonicFormationConfig{
				X: []motion.IndexedHarmonic{rowY.Harmonic()}, Y: []motion.IndexedHarmonic{smallRowZ.Harmonic()},
			}},
		},
	})
}
