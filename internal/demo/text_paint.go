package demo

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/font"
	"github.com/olivierh59500/democonstructionkit/geometry"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/go-leonardintrodx/internal/source"
)

var largeRowX = source.Oscillator{Amplitude: 32, Rate: .08, Spacing: .19}
var largeRowY = source.Oscillator{Amplitude: 32, Rate: .103, Spacing: .212}

func largeCells(bank *font.CellBank, white *ebiten.Image, wireframe func() bool) scrolling.CellPainterConfig {
	return scrolling.CellPainterConfig{
		Fonts: map[string]*font.CellBank{"original": bank}, White: white,
		Flat:      scrolling.FlatCellConfig{Size: geometry.Vec2{X: 16.0 / 18.0, Y: 1}, ScreenGap: geometry.Vec2{Y: 2}, Rectangles: true},
		Wireframe: func(scrolling.Sample) bool { return wireframe() },
		Rows: func(_ string, row int, seconds float64) geometry.Vec3 {
			return geometry.Vec3{X: largeRowX.At(seconds, row), Y: largeRowY.At(seconds, row)}
		},
	}
}

func strokeCell(batch *render.Batch, x, y, width, height float64) {
	if width <= 2 || height <= 2 {
		return
	}
	source := image.Rect(0, 0, 1, 1)
	batch.Rect(x, y, width, 1, source, color.White)
	batch.Rect(x, y+height-1, width, 1, source, color.White)
	batch.Rect(x, y+1, 1, height-2, source, color.White)
	batch.Rect(x+width-1, y+1, 1, height-2, source, color.White)
}
