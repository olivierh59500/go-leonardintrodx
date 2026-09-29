package source

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/olivierh59500/go-leonardintrodx/assets"
)

func TestOriginalFontLayouts(t *testing.T) {
	largeData, err := assets.Files.ReadFile("original/font-large.bin")
	if err != nil {
		t.Fatal(err)
	}
	lookup, err := assets.Files.ReadFile("original/font-map.bin")
	if err != nil {
		t.Fatal(err)
	}
	large, err := ReadLargeFont(largeData, lookup)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 26; i++ {
		if got := large.Index(byte('A' + i)); got != i {
			t.Fatalf("letter %q points to glyph %d", byte('A'+i), got)
		}
	}
	if large.Index(' ') != 36 {
		t.Fatal("space is not the authored blank glyph")
	}
	var pixels [LargeGlyphWidth * LargeGlyphHeight]byte
	for y := 0; y < LargeGlyphHeight; y++ {
		for x := 0; x < LargeGlyphWidth; x++ {
			if large.Pixel('A', x, y) {
				pixels[y*LargeGlyphWidth+x] = 1
			}
			if large.Pixel(' ', x, y) {
				t.Fatal("space glyph is not blank")
			}
		}
	}
	if got := fmt.Sprintf("%x", sha256.Sum256(pixels[:])); got != "6cb4097abd36d6747364b9a6d2cd61e864653529efff6cf1c02667553b9ce1d0" {
		t.Fatalf("large glyph A changed: %s", got)
	}
	smallData, err := assets.Files.ReadFile("original/font-small.bin")
	if err != nil {
		t.Fatal(err)
	}
	small, err := ReadSmallFont(smallData)
	if err != nil {
		t.Fatal(err)
	}
	if !small.Pixel('A', 3, 1) || small.Pixel(' ', 3, 1) {
		t.Fatal("small font ASCII layout changed")
	}
}
