package demo

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/democonstructionkit/sprites"
	"github.com/olivierh59500/go-leonardintrodx/internal/source"
)

// The 16 ribbons use the source's four oscillators and pastel vertex palette.
var ribbonColors = [...]color.RGBA{
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

// The original program evaluates four sine banks for each of eighty 32-pixel
// textured sprites. DCK's shared painter receives the resulting positions.
func (g *Game) updateBallSamples(t float64) {
	for i := range g.ballSamples {
		x := 320 + ballWaves[0].At(t, i) + ballWaves[1].At(t, i)
		y := 240 + ballWaves[2].At(t, i) + ballWaves[3].At(t, i)
		g.ballSamples[i] = sprites.FieldSample{Index: i, X: math.Round(x), Y: math.Round(y), Scale: 1}
	}
}

func (g *Game) drawBands(dst *ebiten.Image, t float64) {
	g.batch.Begin(dst, g.white)
	phase := t - BallsStart
	for i, paint := range ribbonColors {
		left := math.Round(232 + ribbonWaves[0].At(phase, i) + ribbonWaves[1].At(phase, i))
		right := math.Round(232 + ribbonWaves[2].At(phase, i) + ribbonWaves[3].At(phase, i))
		left, right = max(0, min(463, left)), max(0, min(463, right))
		quad := [4]ebiten.Vertex{
			render.Vertex(0, left, 0, 0, paint), render.Vertex(640, right, 1, 0, paint),
			render.Vertex(640, right+8, 1, 1, paint), render.Vertex(0, left+8, 0, 1, paint),
		}
		g.batch.Quad(quad)
	}
	g.batch.Flush()
}

func (g *Game) drawCubeWire(dst *ebiten.Image, x, y float64) {
	vertices, _ := g.cube.Geometry(x, y)
	paint := color.RGBA{255, 198, 255, 255}
	for face := 0; face < 6; face++ {
		base := face * 20
		for corner := 0; corner < 4; corner++ {
			a := vertices[base+corner]
			b := vertices[base+(corner+1)%4]
			vector.StrokeLine(dst, a.DstX, a.DstY, b.DstX, b.DstY, 1.5, paint, false)
		}
	}
}

// The source combines two sine controllers per axis. Their time rates are
// the original phase increments multiplied by the 85-unit speed clock.
func cubePosition(t float64) (float64, float64) {
	x := 60*math.Sin(t*85*.0323) + 60*math.Sin(t*85*.0200)
	y := 50*math.Sin(t*85*.0321) + 50*math.Sin(t*85*.0257)
	z := 100 + 50*math.Sin(t*85*.0357) + 50*math.Sin(t*85*.0279)
	scale := 1.15 * 100 / max(40, z)
	return 320 + x*scale, 260 + y*scale
}

func (g *Game) drawMeters(dst *ebiten.Image) {
	g.batch.Begin(dst, g.white)
	for channel := 0; channel < 3; channel++ {
		amplitude := int(g.registers[8+channel] & 0x0f)
		period := int(g.registers[channel*2]) | int(g.registers[channel*2+1]&0xf)<<8
		if amplitude == 0 || period == 0 {
			continue
		}
		height := math.Min(120, float64(amplitude)*4+float64(2000/period))
		x := 17 + float64(channel)*8
		g.batch.Rect(x, 479-height, 5, height, image.Rect(0, 0, 1, 1), color.RGBA{255, 30 + uint8(60*channel), 20, 255})
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
