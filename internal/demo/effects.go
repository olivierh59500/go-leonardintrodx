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
	"github.com/olivierh59500/democonstructionkit/render"
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

func (g *Game) drawCubeWire(dst *ebiten.Image) {
	pose := g.cube.Transform
	rotation := geometry.RotateXYZ(pose.Rotation)
	var points [8]geometry.Vec3
	for i, point := range g.cube.Mesh.Points {
		points[i] = rotation.Apply(point.Scale(pose.Scale)).Add(pose.Position)
	}
	g.batch.Begin(dst, g.white)
	for _, face := range [6][4]int{
		{0, 3, 2, 1}, {4, 5, 6, 7}, {0, 4, 7, 3},
		{1, 2, 6, 5}, {0, 1, 5, 4}, {3, 7, 6, 2},
	} {
		a, b, c := points[face[0]], points[face[1]], points[face[2]]
		if b.Sub(a).Cross(c.Sub(a)).Dot(a) >= 0 {
			continue
		}
		var quad [4]geometry.Vec2
		for corner, index := range face {
			quad[corner], _, _ = g.cube.Camera.Project(points[index])
		}
		g.strokeQuad(quad, color.RGBA{255, 198, 255, 255})
	}
	g.batch.Flush()
}

func (g *Game) drawWireframe(dst *ebiten.Image, t float64) {
	if t >= BallsStart {
		g.bands.DrawOutline(dst, 1.2)
	}
	g.batch.Begin(dst, g.white)
	if t >= BallFieldStart {
		for _, sprite := range g.balls.Samples() {
			strokeCell(g.batch, sprite.X-16, sprite.Y-16, 32, 32)
		}
	}
	if t >= LogoStart {
		local := t - LogoStart
		paint := color.RGBA{235, 135, 245, 255}
		for row := 0; row < 10; row++ {
			top, bottom := float64(row)*10.3, float64(row+1)*10.3
			g.strokeQuad([4]geometry.Vec2{
				g.logo.Map(0, top, local), g.logo.Map(233, top, local),
				g.logo.Map(233, bottom, local), g.logo.Map(0, bottom, local),
			}, paint)
		}
	}
	g.batch.Flush()
	if t >= LargeTextStart {
		g.large.Draw(dst)
	}
	g.drawCubeWire(dst)
	if t >= SmallTextStart {
		g.small.Draw(dst)
	}
	g.drawMeters(dst)
	if g.showLoad {
		g.drawLoad(dst)
	}
}

func (g *Game) strokeQuad(quad [4]geometry.Vec2, paint color.Color) {
	g.batch.StrokePath(quad[:], render.PathStroke{Width: 1.2}, paint)
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
	g.batch.Begin(dst, g.white)
	for column := 0; column < 80; column++ {
		height := g.meter.Level(column)
		if height <= 0 {
			continue
		}
		x, top := float64(column*8), 479-height
		if g.wireframe {
			strokeCell(g.batch, x, top, 7, height)
			continue
		}
		yellow := color.RGBA{255, 255, 0, 255}
		red := color.RGBA{255, 0, 0, 255}
		g.batch.Quad([4]ebiten.Vertex{
			render.Vertex(x, top, 0, 0, yellow), render.Vertex(x+7, top, 1, 0, yellow),
			render.Vertex(x+7, 479, 1, 1, red), render.Vertex(x, 479, 0, 1, red),
		})
	}
	g.batch.Flush()
}

func (g *Game) drawLoad(dst *ebiten.Image) {
	load := math.Min(99, math.Max(0, ebiten.ActualTPS()/FPS*100))
	g.batch.Begin(dst, g.white)
	g.batch.Rect(400, 10, 160, 12, image.Rect(0, 0, 1, 1), color.RGBA{60, 20, 60, 255})
	g.batch.Rect(400, 10, 160*load/100, 12, image.Rect(0, 0, 1, 1), color.RGBA{180, 255, 80, 255})
	g.batch.Flush()
}
