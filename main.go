package main

import (
	"bytes"
	"flag"
	"image"
	"image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	capture "github.com/olivierh59500/democonstructionkit/fidelity/ebiten"
	"github.com/olivierh59500/go-leonardintrodx/assets"
	"github.com/olivierh59500/go-leonardintrodx/internal/demo"
)

func main() {
	start := flag.Float64("start", 0, "production position in seconds")
	mute := flag.Bool("mute", false, "disable device audio")
	directory := flag.String("capture", "", "write one native frame to a directory")
	limit := flag.Int("ticks", 0, "optional number of 50 Hz updates before exit")
	layer := flag.String("layer", "", "optional isolated layer for capture")
	wireframe := flag.Bool("wireframe", false, "start in outlined geometry mode")
	flag.Parse()
	if *start < 0 || *limit < 0 || math.IsNaN(*start) || math.IsInf(*start, 0) {
		log.Fatal("invalid start time")
	}
	first := int(math.Round(*start * demo.FPS))
	if *directory != "" {
		if err := capture.Run(capture.Config{Directory: *directory, Frames: []int{0}, Width: demo.Width, Height: demo.Height}, func() (ebiten.Game, error) {
			game, err := demo.NewGame(first, true)
			if err != nil {
				return nil, err
			}
			game.SetWireframe(*wireframe)
			if err = game.SetLayer(*layer); err != nil {
				game.Close()
				return nil, err
			}
			return game, nil
		}); err != nil {
			log.Fatal(err)
		}
		return
	}
	game, err := demo.NewGame(first, *mute)
	if err != nil {
		log.Fatal(err)
	}
	game.SetWireframe(*wireframe)
	if err = game.SetLayer(*layer); err != nil {
		log.Fatal(err)
	}
	game.SetTickLimit(*limit)
	defer game.Close()
	ebiten.SetTPS(demo.FPS)
	ebiten.SetWindowSize(demo.Width*3/2, demo.Height*3/2)
	ebiten.SetWindowTitle("OldSkool DirectX 8 Go - Leonard / Oxygene")
	icons := make([]image.Image, 0, 2)
	for _, name := range []string{"original/icon-16.png", "original/icon-32.png"} {
		data, readErr := assets.Files.ReadFile(name)
		if readErr != nil {
			log.Fatal(readErr)
		}
		icon, decodeErr := png.Decode(bytes.NewReader(data))
		if decodeErr != nil {
			log.Fatal(decodeErr)
		}
		icons = append(icons, icon)
	}
	ebiten.SetWindowIcon(icons)
	ebiten.SetRunnableOnUnfocused(true)
	if err = ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
