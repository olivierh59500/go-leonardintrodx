package mobile

import (
	"image"
	"image/color"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/olivierh59500/go-leonardintrodx/internal/demo"
)

const (
	controlWireframe = iota
	controlPause
	controlReset
	controlCount
)

var controlLabels = [controlCount]string{"WIRE", "PAUSE", "RESET"}

// logicalWidth preserves the original canvas aspect while reserving wide
// phone sidebars for touch controls.
func logicalWidth(outsideWidth, outsideHeight int) int {
	if outsideWidth <= 0 || outsideHeight <= 0 {
		return demo.Width
	}
	width := (outsideWidth*demo.Height + outsideHeight - 1) / outsideHeight
	return min(1280, max(demo.Width, width))
}

func controlRects(width int) [controlCount]image.Rectangle {
	var rects [controlCount]image.Rectangle
	band := (width - demo.Width) / 2
	if band < 96 {
		return rects
	}
	buttonWidth := min(160, band-30)
	left := (band - buttonWidth) / 2
	right := band + demo.Width + left
	rects[controlWireframe] = image.Rect(left, 326, left+buttonWidth, 380)
	rects[controlPause] = image.Rect(left, 396, left+buttonWidth, 450)
	rects[controlReset] = image.Rect(right, 396, right+buttonWidth, 450)
	return rects
}

func controlAt(width, x, y int) int {
	point := image.Pt(x, y)
	for index, rect := range controlRects(width) {
		if point.In(rect) {
			return index
		}
	}
	return -1
}

func (h *host) readControls() [controlCount]bool {
	var actions [controlCount]bool
	h.held = [controlCount]bool{}
	h.currentTouches = ebiten.AppendTouchIDs(h.currentTouches[:0])
	for _, id := range h.currentTouches {
		x, y := ebiten.TouchPosition(id)
		index := controlAt(h.layoutWidth, x, y)
		if index < 0 {
			continue
		}
		h.held[index] = true
		if !slices.Contains(h.priorTouches, id) {
			actions[index] = true
		}
	}
	h.priorTouches = append(h.priorTouches[:0], h.currentTouches...)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if index := controlAt(h.layoutWidth, x, y); index >= 0 {
			h.held[index] = true
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				actions[index] = true
			}
		}
	}
	return actions
}

func (h *host) drawControls(dst *ebiten.Image) {
	rects := controlRects(h.layoutWidth)
	for index, rect := range rects {
		if rect.Empty() {
			continue
		}
		active := h.held[index]
		if h.scene != nil {
			active = active || (index == controlWireframe && h.scene.Wireframe()) || (index == controlPause && h.scene.Paused())
		}
		fill := color.RGBA{R: 40, G: 26, B: 55, A: 255}
		border := color.RGBA{R: 155, G: 106, B: 177, A: 255}
		if active {
			fill = color.RGBA{R: 105, G: 45, B: 127, A: 255}
			border = color.RGBA{R: 255, G: 203, B: 255, A: 255}
		}
		vector.DrawFilledRect(dst, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Dx()), float32(rect.Dy()), fill, false)
		vector.StrokeRect(dst, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Dx()), float32(rect.Dy()), 2, border, false)
		label := controlLabels[index]
		ebitenutil.DebugPrintAt(dst, label, rect.Min.X+(rect.Dx()-6*len(label))/2, rect.Min.Y+19)
	}
}
