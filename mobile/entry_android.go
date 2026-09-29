//go:build android

package mobile

import (
	"github.com/hajimehoshi/ebiten/v2"
	enginemobile "github.com/hajimehoshi/ebiten/v2/mobile"
	"github.com/olivierh59500/go-leonardintrodx/internal/demo"
)

func init() {
	ebiten.SetTPS(demo.FPS)
	enginemobile.SetGame(&host{layoutWidth: demo.Width})
}

// Dummy gives the Android binding a stable exported Go symbol.
func Dummy() {}
