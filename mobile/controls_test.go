package mobile

import (
	"image"
	"testing"

	"github.com/olivierh59500/go-leonardintrodx/internal/demo"
)

func TestPixelLayoutKeepsSceneAndTouchControlsSeparate(t *testing.T) {
	width := logicalWidth(2424, 1080)
	if width != 1078 {
		t.Fatalf("Pixel logical width = %d, want 1078", width)
	}
	scene := image.Rect((width-demo.Width)/2, 0, (width+demo.Width)/2, demo.Height)
	for index, rect := range controlRects(width) {
		if rect.Empty() || rect.Overlaps(scene) {
			t.Fatalf("control %d overlaps the original scene: %v", index, rect)
		}
		center := rect.Min.Add(rect.Size().Div(2))
		if got := controlAt(width, center.X, center.Y); got != index {
			t.Fatalf("control at %v = %d, want %d", center, got, index)
		}
	}
	if controlAt(width, scene.Min.X+20, 420) != -1 {
		t.Fatal("touch inside the demo selected a side control")
	}
}

func TestNarrowLayoutHidesSideControls(t *testing.T) {
	if got := logicalWidth(640, 480); got != demo.Width {
		t.Fatalf("native layout width = %d", got)
	}
	for _, rect := range controlRects(demo.Width) {
		if !rect.Empty() {
			t.Fatalf("unexpected control over native canvas: %v", rect)
		}
	}
}
