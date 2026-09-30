package demo

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/font"
	"github.com/olivierh59500/democonstructionkit/geometry"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/go-leonardintrodx/internal/source"
	"image"
	"image/color"
	"math"
)

var smallCellColors = []color.NRGBA{
	{255, 255, 255, 255}, {255, 255, 255, 255}, {255, 255, 255, 255}, {255, 255, 255, 255},
	{255, 255, 0, 255}, {255, 128, 0, 255}, {0, 0, 255, 255}, {0, 255, 255, 255},
}

func newSmallScrolling(face scrolling.Face, bank *font.CellBank, text string, white *ebiten.Image, wireframe func() bool) (*scrolling.Scrolling, error) {
	var rows [source.SmallGlyphHeight]source.SmallScrollRow
	lastTime := math.NaN()
	cells := scrolling.CellPainterConfig{
		Fonts: map[string]*font.CellBank{"original": bank}, Shape: scrolling.CellCuboid,
		White: white, LineWidth: 1.2, CullBounds: image.Rect(0, 0, Width, Height),
		Wireframe: func(scrolling.Sample) bool { return wireframe() },
		Cuboid: scrolling.CuboidCellConfig{
			Size:   geometry.Vec3{X: source.SmallCellSide, Y: source.SmallCellSide, Z: source.SmallCellSide},
			Origin: geometry.Vec3{X: -source.SmallWorldOffset, Z: source.CubeViewZ},
			Colors: smallCellColors, FlipY: true, IgnoreMappedY: true,
			Camera: geometry.Camera{Center: geometry.Vec2{X: Width / 2, Y: Height / 2}, Focal: Height * math.Sqrt(3) / 2, Near: .1},
		},
		Rows: func(_ string, row int, seconds float64) geometry.Vec3 {
			if lastTime != seconds {
				rows = source.SmallScrollRows(seconds)
				lastTime = seconds
			}
			return geometry.Vec3{Y: rows[row].Y, Z: rows[row].Z}
		},
	}
	return scrolling.New(scrolling.Config{
		Text: text, Fonts: map[string]scrolling.Face{"original": face}, Font: "original",
		Speed: 85 * 7, X: 39 * source.SmallColumnPitch, Repeat: true,
		Map: func(sample scrolling.Sample, op *ebiten.DrawImageOptions) bool {
			if sample.Glyph.Image == nil {
				return false
			}
			left, _ := op.GeoM.Apply(0, 0)
			right, _ := op.GeoM.Apply(float64(sample.Glyph.Image.Bounds().Dx()), 0)
			return right > source.SmallWorldOffset-400 && left < source.SmallWorldOffset+400
		},
		Shape: "source-3d", Modes: map[string]scrolling.Mode{"source-3d": {Cells: &cells}},
	})
}
