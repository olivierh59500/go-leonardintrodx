package demo

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/font"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/go-leonardintrodx/assets"
	"github.com/olivierh59500/go-leonardintrodx/internal/source"
)

type artwork struct {
	ball, logo, largeAtlas, smallAtlas *ebiten.Image
	largeFace, smallFace               scrolling.Face
	largeText, smallText               string
	largeCells, smallCells             *font.CellBank
}

func loadArtwork() (*artwork, error) {
	a := &artwork{}
	var err error
	if a.ball, err = loadTransparent("original/ball.png"); err != nil {
		return nil, err
	}
	if a.logo, err = loadTransparent("original/logo.png"); err != nil {
		a.close()
		return nil, err
	}
	large, err := readLargeFont()
	if err != nil {
		a.close()
		return nil, err
	}
	a.largeAtlas = ebiten.NewImage(8*source.LargeGlyphWidth, 6*source.LargeGlyphHeight)
	pixels := image.NewNRGBA(a.largeAtlas.Bounds())
	for glyph := 0; glyph < source.LargeGlyphCount; glyph++ {
		for y := 0; y < source.LargeGlyphHeight; y++ {
			for x := 0; x < source.LargeGlyphWidth; x++ {
				if large.GlyphPixel(glyph, x, y) {
					pixels.SetNRGBA((glyph%8)*source.LargeGlyphWidth+x, (glyph/8)*source.LargeGlyphHeight+y, color.NRGBA{255, 255, 255, 255})
				}
			}
		}
	}
	a.largeAtlas.WritePixels(pixels.Pix)
	characters := map[rune]font.Glyph{}
	for letter := 32; letter < 127; letter++ {
		index := large.Index(byte(letter))
		if index >= source.LargeGlyphCount {
			continue
		}
		x, y := (index%8)*source.LargeGlyphWidth, (index/8)*source.LargeGlyphHeight
		rect := image.Rect(x, y, x+source.LargeGlyphWidth, y+source.LargeGlyphHeight)
		if letter == ' ' {
			rect = image.Rectangle{}
		}
		characters[rune(letter)] = font.Glyph{Rect: rect, Advance: source.LargeGlyphWidth}
	}
	face, err := font.New(font.Config{Bounds: a.largeAtlas.Bounds(), Glyphs: characters, LineHeight: source.LargeGlyphHeight, SpaceAdvance: source.LargeGlyphWidth, Fallback: '?'})
	if err != nil {
		a.close()
		return nil, err
	}
	a.largeFace = scrolling.Face{Atlas: a.largeAtlas, Metrics: face, ScaleX: 18, ScaleY: 16}
	var largeCharacters []rune
	for r := 0; r < 256; r++ {
		largeCharacters = append(largeCharacters, rune(r))
	}
	a.largeCells, err = font.NewCellBank(font.CellBankConfig{
		Width: source.LargeGlyphWidth, Height: source.LargeGlyphHeight, Characters: string(largeCharacters),
		Key:   func(r rune) int { return large.Index(byte(r)) },
		Pixel: func(r rune, x, y int) bool { return large.Pixel(byte(r), x, y) },
	})
	if err != nil {
		a.close()
		return nil, err
	}
	small, err := readSmallFont()
	if err != nil {
		a.close()
		return nil, err
	}
	a.smallAtlas = ebiten.NewImage(16*8, 7*8)
	pixels = image.NewNRGBA(a.smallAtlas.Bounds())
	var order []rune
	for character := 32; character < 32+source.SmallGlyphCount; character++ {
		order = append(order, rune(character))
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				if small.Pixel(byte(character), x, y) {
					index := character - 32
					pixels.SetNRGBA((index%16)*8+x, (index/16)*8+y, color.NRGBA{255, 255, 255, 255})
				}
			}
		}
	}
	a.smallAtlas.WritePixels(pixels.Pix)
	smallMetrics, err := font.NewGrid(font.Grid{Bounds: a.smallAtlas.Bounds(), Cell: image.Pt(8, 8), Columns: 16, Order: string(order), Advance: 8, LineHeight: 8, SpaceAdvance: 8})
	if err != nil {
		a.close()
		return nil, err
	}
	a.smallFace = scrolling.Face{Atlas: a.smallAtlas, Metrics: smallMetrics, ScaleX: source.SmallColumnPitch, ScaleY: source.SmallRowPitch}
	a.smallCells, err = font.NewCellBank(font.CellBankConfig{
		Width: source.SmallGlyphWidth, Height: source.SmallGlyphHeight, Characters: string(order),
		Pixel: func(r rune, x, y int) bool { return small.Pixel(byte(r), x, y) },
	})
	if err != nil {
		a.close()
		return nil, err
	}
	if data, readErr := assets.Files.ReadFile("original/large-message.txt"); readErr == nil {
		a.largeText = string(data)
	} else {
		a.close()
		return nil, readErr
	}
	if data, readErr := assets.Files.ReadFile("original/small-message.txt"); readErr == nil {
		a.smallText = string(data)
	} else {
		a.close()
		return nil, readErr
	}
	return a, nil
}

func readLargeFont() (source.LargeFont, error) {
	glyphs, err := assets.Files.ReadFile("original/font-large.bin")
	if err != nil {
		return source.LargeFont{}, err
	}
	lookup, err := assets.Files.ReadFile("original/font-map.bin")
	if err != nil {
		return source.LargeFont{}, err
	}
	return source.ReadLargeFont(glyphs, lookup)
}
func readSmallFont() (source.SmallFont, error) {
	data, err := assets.Files.ReadFile("original/font-small.bin")
	if err != nil {
		return source.SmallFont{}, err
	}
	return source.ReadSmallFont(data)
}

func loadTransparent(name string) (*ebiten.Image, error) {
	data, err := assets.Files.ReadFile(name)
	if err != nil {
		return nil, err
	}
	decoded, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	paletted, ok := decoded.(*image.Paletted)
	if !ok {
		return nil, fmt.Errorf("%s: expected source palette", name)
	}
	imageData := image.NewNRGBA(paletted.Bounds())
	for y := 0; y < paletted.Bounds().Dy(); y++ {
		for x := 0; x < paletted.Bounds().Dx(); x++ {
			index := paletted.ColorIndexAt(x, y)
			if index == 0 {
				continue
			}
			colour := color.NRGBAModel.Convert(paletted.Palette[index]).(color.NRGBA)
			colour.A = 255
			imageData.SetNRGBA(x, y, colour)
		}
	}
	return ebiten.NewImageFromImage(imageData), nil
}

func (a *artwork) close() {
	for _, image := range []*ebiten.Image{a.ball, a.logo, a.largeAtlas, a.smallAtlas} {
		if image != nil {
			image.Deallocate()
		}
	}
}
