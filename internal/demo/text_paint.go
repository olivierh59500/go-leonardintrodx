package demo

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/go-leonardintrodx/internal/source"
)

// The original font effects place each lit bit in its own 3D/raster cell.
// DCK controls font order, speed and repetition; this painter supplies the
// source-specific cell geometry and keeps a two-pixel horizontal gap.
func newLargePainter(white *ebiten.Image) (scrolling.Painter, error) {
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
		batch.Begin(dst, white)
		for y := 0; y < source.LargeGlyphHeight; y++ {
			rowShiftX := 32 * math.Sin(sample.Time*0.68+float64(y)*0.18)
			rowShiftY := 23 * math.Sin(sample.Time*0.73+float64(y)*0.15)
			for x := 0; x < source.LargeGlyphWidth; x++ {
				if !font.Pixel(character, x, y) {
					continue
				}
				left, top := op.GeoM.Apply(float64(x), float64(y))
				right, bottom := op.GeoM.Apply(float64(x)+16.0/18.0, float64(y+1))
				batch.Rect(left+rowShiftX, top+rowShiftY, right-left, bottom-top-2,
					image.Rect(0, 0, 1, 1), color.White)
			}
		}
		batch.Flush()
	}, nil
}

func newSmallPainter(white *ebiten.Image) (scrolling.Painter, error) {
	font, err := readSmallFont()
	if err != nil {
		return nil, err
	}
	batch := render.NewBatch(4096)
	return func(dst *ebiten.Image, sample scrolling.Sample, op ebiten.DrawImageOptions) {
		if sample.Glyph.Rune > 255 {
			return
		}
		character := byte(sample.Glyph.Rune)
		batch.Begin(dst, white)
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				if !font.Pixel(character, x, y) {
					continue
				}
				left, top := op.GeoM.Apply(float64(x), float64(y))
				right, bottom := op.GeoM.Apply(float64(x+1), float64(y+1))
				// Face shading belongs to each lit bitmap bit, like the original
				// Direct3D cubelets. The surrounding warp adds motion and palette.
				left += 14 * math.Sin(sample.Time*0.79+left*0.011)
				top += 27 * math.Sin(sample.Time*0.71+left*0.009)
				width, height := right-left-2, bottom-top-3
				right, bottom = left+width, top+height
				if width <= 0 || height <= 0 || math.IsNaN(width) || math.IsNaN(height) {
					continue
				}
				front := [4]ebiten.Vertex{
					render.Vertex(left+5, top+5, 0, 0, color.RGBA{255, 255, 255, 255}),
					render.Vertex(right-2, top+5, 1, 0, color.RGBA{255, 250, 235, 255}),
					render.Vertex(right-2, bottom-2, 1, 1, color.RGBA{225, 240, 255, 255}),
					render.Vertex(left+5, bottom-2, 0, 1, color.RGBA{255, 230, 245, 255}),
				}
				topFace := [4]ebiten.Vertex{
					render.Vertex(left, top, 0, 0, color.RGBA{255, 238, 180, 255}),
					render.Vertex(right-5, top, 1, 0, color.RGBA{255, 238, 180, 255}),
					front[1], front[0],
				}
				side := [4]ebiten.Vertex{
					front[1], render.Vertex(right, top, 1, 0, color.RGBA{160, 210, 255, 255}),
					render.Vertex(right, bottom-7, 1, 1, color.RGBA{160, 210, 255, 255}), front[2],
				}
				batch.Quad(topFace)
				batch.Quad(side)
				batch.Quad(front)
			}
		}
		batch.Flush()
	}, nil
}
