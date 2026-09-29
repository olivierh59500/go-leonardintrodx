package demo

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/geometry"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/democonstructionkit/sprites"
)

func newBallField(ball *ebiten.Image, count int, cloud bool) (*sprites.ProjectedField, error) {
	points := make([]sprites.Point, count)
	for i := range points {
		angle := 2 * math.Pi * float64(i%16) / 16
		radius := 132.0
		if cloud {
			switch {
			case i < 38:
				radius, angle = 250, 2*math.Pi*float64(i)/38
			case i < 66:
				radius, angle = 135, 2*math.Pi*float64(i-38)/28+0.35
			default:
				radius, angle = 70, 2*math.Pi*float64(i-66)/14+0.7
			}
		}
		points[i] = sprites.Point{X: radius * math.Cos(angle), Y: radius * 0.88 * math.Sin(angle), Z: 215 + 25*math.Sin(angle*2+float64(i/16)*0.83)}
	}
	camera := geometry.Camera{Center: geometry.Vec2{X: 320, Y: 240}, Focal: 230, Near: 1}
	return sprites.NewProjectedField(sprites.ProjectedFieldConfig{
		Field:            sprites.FieldConfig{Points: points},
		View:             sprites.FieldView{Camera: camera, SortDepth: true},
		Style:            sprites.FieldStyle{Image: ball, ScaleByDepth: true, Appearance: sprites.FieldAppearance{Width: 31, Height: 31, AnchorX: .5, AnchorY: .5}},
		RendererCapacity: count,
	})
}

func (g *Game) updateBalls(field *sprites.ProjectedField, t float64, cloud bool) error {
	view := field.View()
	view.Angle = t * 0.42
	view.Offset = geometry.Vec3{X: 90 * math.Sin(t*0.31), Y: 60 * math.Sin(t*0.23+1.2)}
	if cloud {
		view.Angle = t * 0.25
		view.Offset.X += 60 * math.Sin(t*0.19+.5)
	}
	field.SetView(view)
	return field.Update(kit.Frame{Time: t})
}

// The 16 ribbons use the source's four oscillators and pastel vertex palette.
var ribbonColors = [...]color.RGBA{
	{255, 255, 255, 255}, {255, 255, 208, 255}, {255, 255, 176, 255}, {255, 255, 144, 255},
	{255, 224, 128, 255}, {255, 192, 128, 255}, {255, 160, 128, 255}, {255, 128, 128, 255},
	{255, 128, 128, 255}, {208, 128, 160, 255}, {176, 128, 192, 255}, {144, 128, 224, 255},
	{128, 144, 240, 255}, {128, 176, 240, 255}, {128, 208, 240, 255}, {144, 240, 240, 255},
}

func (g *Game) drawBands(dst *ebiten.Image, t float64) {
	g.batch.Begin(dst, g.white)
	phase := (t - BallsStart) * 85
	for i, paint := range ribbonColors {
		column := float64(i)
		left := math.Round(232 + 116*(math.Sin(phase*.019+column*.108333)+math.Sin(phase*.0337+column*.0413)))
		right := math.Round(232 + 116*(math.Sin(phase*.023+column*.1125)+math.Sin(phase*.0232+column*.0539)))
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
