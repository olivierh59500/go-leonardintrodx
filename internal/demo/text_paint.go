package demo

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/go-leonardintrodx/internal/source"
)

var largeRowX = source.Oscillator{Amplitude: 32, Rate: .08, Spacing: .19}
var largeRowY = source.Oscillator{Amplitude: 32, Rate: .103, Spacing: .212}

// The original font effects place each lit bit in its own 3D/raster cell.
// DCK controls font order, speed and repetition; this painter supplies the
// source-specific cell geometry and keeps a two-pixel horizontal gap.
func newLargePainter(white *ebiten.Image, wireframe func() bool) (scrolling.Painter, error) {
	font, err := readLargeFont()
	if err != nil {
		return nil, err
	}
	batch := render.NewBatch(4096)
	return func(dst *ebiten.Image, sample scrolling.Sample, op ebiten.DrawImageOptions) {
		if sample.Glyph.Rune > 255 {
			return
		}
		character := byte(sample.Glyph.Rune)
		outlined := wireframe()
		batch.Begin(dst, white)
		for y := 0; y < source.LargeGlyphHeight; y++ {
			rowShiftX := largeRowX.At(sample.Time, y)
			rowShiftY := largeRowY.At(sample.Time, y)
			for x := 0; x < source.LargeGlyphWidth; x++ {
				if !font.Pixel(character, x, y) {
					continue
				}
				left, top := op.GeoM.Apply(float64(x), float64(y))
				right, bottom := op.GeoM.Apply(float64(x)+16.0/18.0, float64(y+1))
				if outlined {
					strokeCell(batch, left+rowShiftX, top+rowShiftY, right-left, bottom-top-2)
				} else {
					batch.Rect(left+rowShiftX, top+rowShiftY, right-left, bottom-top-2,
						image.Rect(0, 0, 1, 1), color.White)
				}
			}
		}
		batch.Flush()
	}, nil
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
