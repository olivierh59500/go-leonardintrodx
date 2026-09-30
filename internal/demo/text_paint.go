package demo

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/font"
	"github.com/olivierh59500/democonstructionkit/geometry"
	"github.com/olivierh59500/democonstructionkit/motion"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/go-leonardintrodx/internal/source"
)

var largeRowX = source.Oscillator{Amplitude: 32, Rate: .08, Spacing: .19}
var largeRowY = source.Oscillator{Amplitude: 32, Rate: .103, Spacing: .212}

func largeRowProfile() (*motion.HarmonicRowProfile, error) {
	return motion.NewHarmonicRowProfile(motion.HarmonicRowProfileConfig{
		Count: source.LargeGlyphHeight, ClockScale: [2]float64{85},
		Stages: []motion.HarmonicRowStage{{Motion: motion.HarmonicFormationConfig{
			X: []motion.IndexedHarmonic{largeRowX.Harmonic()}, Y: []motion.IndexedHarmonic{largeRowY.Harmonic()},
		}}},
	})
}

func largeCells(bank *font.CellBank, white *ebiten.Image, wireframe func() bool) (scrolling.CellPainterConfig, error) {
	rows, err := largeRowProfile()
	if err != nil {
		return scrolling.CellPainterConfig{}, err
	}
	return scrolling.CellPainterConfig{
		Fonts: map[string]*font.CellBank{"original": bank}, White: white,
		Flat:      scrolling.FlatCellConfig{Size: geometry.Vec2{X: 16.0 / 18.0, Y: 1}, ScreenGap: geometry.Vec2{Y: 2}, Rectangles: true},
		Wireframe: func(scrolling.Sample) bool { return wireframe() },
		Rows: func(_ string, row int, seconds float64) geometry.Vec3 {
			pose, _ := rows.Apply(row, seconds, motion.Point{})
			return geometry.Vec3{X: pose.X, Y: pose.Y}
		},
	}, nil
}
