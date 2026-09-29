// Package source reads the bitmap font encodings of the supplied production.
package source

import "fmt"

const (
	LargeGlyphWidth  = 32
	LargeGlyphHeight = 30
	LargeGlyphCount  = 48
	SmallGlyphWidth  = 8
	SmallGlyphHeight = 8
	SmallGlyphCount  = 100
)

// LargeFont retains the authored column-byte organization and lookup order.
type LargeFont struct {
	glyphs [LargeGlyphCount * 120]byte
	lookup [256]byte
}

func ReadLargeFont(glyphs, lookup []byte) (LargeFont, error) {
	var font LargeFont
	if len(glyphs) != len(font.glyphs) || len(lookup) != 512 {
		return font, fmt.Errorf("source: invalid large font or lookup size")
	}
	copy(font.glyphs[:], glyphs)
	for character := range font.lookup {
		font.lookup[character] = lookup[2*character+1]
	}
	return font, nil
}

// Index returns the original glyph index for an ASCII character.
func (f *LargeFont) Index(character byte) int { return int(f.lookup[character]) }

// GlyphPixel samples a 32 by 30 authored glyph by its original index.
func (f *LargeFont) GlyphPixel(index, x, y int) bool {
	if x < 0 || x >= LargeGlyphWidth || y < 0 || y >= LargeGlyphHeight {
		return false
	}
	if index < 0 || index >= LargeGlyphCount {
		return false
	}
	offset := index*120 + (x/16)*60 + (x%16)/8 + y*2
	return f.glyphs[offset]&(1<<uint(7-x%8)) != 0
}

func (f *LargeFont) Pixel(character byte, x, y int) bool {
	return f.GlyphPixel(f.Index(character), x, y)
}

// SmallFont stores 100 consecutive 8 by 8 ASCII glyphs starting at space.
type SmallFont struct {
	glyphs [SmallGlyphCount * SmallGlyphHeight]byte
}

func ReadSmallFont(glyphs []byte) (SmallFont, error) {
	var font SmallFont
	if len(glyphs) != len(font.glyphs) {
		return font, fmt.Errorf("source: invalid small font size")
	}
	copy(font.glyphs[:], glyphs)
	return font, nil
}

func (f *SmallFont) Pixel(character byte, x, y int) bool {
	index := int(character) - int(' ')
	if index < 0 || index >= SmallGlyphCount || x < 0 || x >= 8 || y < 0 || y >= 8 {
		return false
	}
	return f.glyphs[index*8+y]&(1<<uint(7-x)) != 0
}
