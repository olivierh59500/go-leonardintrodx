package demo

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/effects"
	"github.com/olivierh59500/democonstructionkit/geometry"
	"github.com/olivierh59500/democonstructionkit/motion"
	"github.com/olivierh59500/democonstructionkit/sprites"
	"github.com/olivierh59500/go-leonardintrodx/internal/source"
)

// The 16 ribbons use the source's four oscillators and pastel vertex palette.
var ribbonColors = [...]color.NRGBA{
	{255, 255, 255, 255}, {255, 255, 208, 255}, {255, 255, 176, 255}, {255, 255, 144, 255},
	{255, 224, 128, 255}, {255, 192, 128, 255}, {255, 160, 128, 255}, {255, 128, 128, 255},
	{255, 128, 128, 255}, {208, 128, 160, 255}, {176, 128, 192, 255}, {144, 128, 224, 255},
	{128, 144, 240, 255}, {128, 176, 240, 255}, {128, 208, 240, 255}, {144, 240, 240, 255},
}

var ribbonWaves = [4]source.Oscillator{
	{Amplitude: 116, Rate: .0190, Spacing: .108333},
	{Amplitude: 116, Rate: .0337, Spacing: .0413},
	{Amplitude: 116, Rate: .0230, Spacing: .1125},
	{Amplitude: 116, Rate: .0232, Spacing: .0539},
}

var ballWaves = [4]source.Oscillator{
	{Amplitude: 151, Rate: .0242, Spacing: .1254},
	{Amplitude: 151, Rate: .0193, Spacing: .2530},
	{Amplitude: 111, Rate: .0242, Spacing: .1145},
	{Amplitude: 111, Rate: .0203, Spacing: .2230},
}

// nativeFormation preserves the executable's oscillator constants and signed
// phase reduction. The DCK controllers own sampling and cached poses.
func nativeFormation(waves [4]source.Oscillator, x, y float64) motion.HarmonicFormationConfig {
	term := func(o source.Oscillator) motion.IndexedHarmonic {
		return motion.IndexedHarmonic{Amplitude: o.Amplitude, Rate: o.Rate, IndexRate: o.Spacing, PhasePeriod: 2 * math.Pi}
	}
	return motion.HarmonicFormationConfig{Origin: motion.Point{X: x, Y: y},
		X: []motion.IndexedHarmonic{term(waves[0]), term(waves[1])},
		Y: []motion.IndexedHarmonic{term(waves[2]), term(waves[3])}}
}

func newNativeBands(white *ebiten.Image) (*composite.HarmonicBands, error) {
	return composite.NewHarmonicBands(composite.HarmonicBandsConfig{
		LeftX: 0, RightX: 640, Thickness: 8, Colors: ribbonColors[:], White: white,
		Motion: nativeFormation(ribbonWaves, 232, 232), ClockScale: [2]float64{85, 0}, PixelSnap: true,
		Bounds: &motion.FormationBounds{Min: motion.Point{}, Max: motion.Point{X: 463, Y: 463}},
	})
}

var outlinedBalls = sprites.FieldOutline{Width: 1}

func (g *Game) drawWireframe(dst *ebiten.Image, t float64) {
	if t >= BallsStart {
		g.bands.DrawOutline(dst, 1.2)
	}
	if t >= BallFieldStart {
		style := g.balls.Style
		style.Outline = &outlinedBalls
		g.balls.DrawStyle(dst, style)
	}
	if t >= LogoStart {
		g.logo.DrawOutline(dst)
	}
	if t >= LargeTextStart {
		g.large.Draw(dst)
	}
	g.cube.DrawOutline(dst)
	if t >= SmallTextStart {
		g.small.Draw(dst)
	}
	g.drawMeters(dst)
	if g.showLoad {
		g.drawLoad(dst)
	}
}

// The source rotates a 60-unit cube around X and Y and translates it through
// the six stored oscillators. The Direct3D view moves the camera 200 units
// back; its projection matrix maps to a 415.69-pixel focal length here.
func cubeTransform(t float64) effects.Transform {
	x, y, z := source.CubePosition(t)
	angle := t * math.Pi / 2
	return effects.Transform{
		Position: geometry.Vec3{X: x, Y: -y, Z: source.CubeViewZ + z},
		Rotation: geometry.Vec3{X: -angle, Y: angle * 1.1}, Scale: 1,
	}
}

func (g *Game) drawMeters(dst *ebiten.Image) {
	if g.wireframe {
		g.meterBars.DrawOutline(dst)
	} else {
		g.meterBars.Draw(dst)
	}
}

func (g *Game) drawLoad(dst *ebiten.Image) {
	load := math.Min(99, math.Max(0, ebiten.ActualTPS()/FPS*100))
	g.batch.Begin(dst, g.white)
	g.batch.Rect(400, 10, 160, 12, image.Rect(0, 0, 1, 1), color.RGBA{60, 20, 60, 255})
	g.batch.Rect(400, 10, 160*load/100, 12, image.Rect(0, 0, 1, 1), color.RGBA{180, 255, 80, 255})
	g.batch.Flush()
}
