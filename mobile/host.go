// Package mobile hosts the original 640 by 480 production on Android.
package mobile

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/go-leonardintrodx/internal/demo"
)

type host struct {
	scene          *demo.Game
	canvas         *ebiten.Image
	layoutWidth    int
	currentTouches []ebiten.TouchID
	priorTouches   []ebiten.TouchID
	held           [controlCount]bool
	lastReport     time.Time
}

func (h *host) start() error {
	if h.scene != nil {
		if err := h.scene.Close(); err != nil {
			return err
		}
	}
	game, err := demo.NewGame(0, false)
	if err != nil {
		return err
	}
	game.EnableRealtimeCube(true)
	h.scene = game
	if h.canvas == nil {
		h.canvas = ebiten.NewImage(demo.Width, demo.Height)
	}
	h.lastReport = time.Now()
	return nil
}

func (h *host) Update() error {
	// Graphics and the device audio stream start after EbitenView is ready.
	if h.scene == nil {
		if err := h.start(); err != nil {
			return err
		}
	}
	actions := h.readControls()
	if actions[controlReset] {
		if err := h.start(); err != nil {
			return err
		}
	}
	if actions[controlWireframe] {
		h.scene.ToggleWireframe()
	}
	if actions[controlPause] {
		if err := h.scene.TogglePause(); err != nil {
			return err
		}
	}
	if err := h.scene.Update(); err != nil {
		return err
	}
	if time.Since(h.lastReport) >= 10*time.Second {
		log.Printf("oldskool_mobile seconds=%.1f tps=%.1f fps=%.1f", h.scene.Seconds(), ebiten.ActualTPS(), ebiten.ActualFPS())
		h.lastReport = time.Now()
	}
	return nil
}

func (h *host) Draw(dst *ebiten.Image) {
	dst.Fill(color.Black)
	if h.scene == nil || h.canvas == nil {
		return
	}
	h.scene.Draw(h.canvas)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64((h.layoutWidth-demo.Width)/2), 0)
	dst.DrawImage(h.canvas, op)
	h.drawControls(dst)
}

func (h *host) Layout(outsideWidth, outsideHeight int) (int, int) {
	h.layoutWidth = logicalWidth(outsideWidth, outsideHeight)
	return h.layoutWidth, demo.Height
}
