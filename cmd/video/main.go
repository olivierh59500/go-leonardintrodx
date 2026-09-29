// Command video exports the native intro canvas with its YM soundtrack.
package main

import (
	"flag"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/video"
	"github.com/olivierh59500/go-leonardintrodx/internal/demo"
)

func main() {
	config := video.Config{
		Output: "recordings/oldskool-directx8-go.mp4", Title: "OldSkool DirectX 8 Go",
		Width: demo.Width, Height: demo.Height, FPS: demo.FPS, TPS: demo.FPS,
		SampleRate: 48000, Duration: 3 * time.Minute, PosterAt: 50 * time.Second,
	}
	config.Flags(flag.CommandLine)
	flag.Parse()
	if config.Duration <= 0 {
		log.Fatal("the original intro loops; recording duration must be positive")
	}
	if err := video.Run(config, func() (ebiten.Game, error) { return demo.NewGame(0, false) }); err != nil {
		log.Fatal(err)
	}
}
