// Command checkframes compares complete filled or wireframe production frames.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/go-leonardintrodx/internal/demo"
)

type sample struct {
	Tick   int    `json:"tick"`
	SHA256 string `json:"sha256"`
}

type report struct {
	Width        int      `json:"width"`
	Height       int      `json:"height"`
	Rate         int      `json:"ticks_per_second"`
	Frames       int      `json:"drawn_frames"`
	Wireframe    bool     `json:"wireframe"`
	Samples      []sample `json:"samples"`
	UpdateMeanUS float64  `json:"update_mean_us,omitempty"`
	DrawMeanUS   float64  `json:"draw_submission_mean_us,omitempty"`
}

type probe struct {
	game         *demo.Game
	surface      *ebiten.Image
	pixels       []byte
	tick, limit  int
	measure      bool
	update, draw time.Duration
	report       report
	done         bool
	err          error
}

func (*probe) Layout(int, int) (int, int) { return demo.Width, demo.Height }
func (p *probe) Update() error {
	if p.done {
		return ebiten.Termination
	}
	return nil
}

func (p *probe) Draw(dst *ebiten.Image) {
	if p.done {
		return
	}
	if p.game == nil {
		p.game, p.err = demo.NewGame(0, true)
		if p.err != nil {
			p.done = true
			return
		}
		p.game.SetWireframe(p.report.Wireframe)
		p.surface = render.NewSurface(demo.Width, demo.Height)
		p.pixels = make([]byte, demo.Width*demo.Height*4)
	}
	for step := 0; step < 12 && !p.done; step++ {
		var started time.Time
		if p.measure {
			started = time.Now()
		}
		if p.tick > 0 {
			if p.err = p.game.Update(); p.err != nil {
				p.done = true
				return
			}
		}
		if p.measure {
			p.update += time.Since(started)
			started = time.Now()
		}
		p.game.Draw(p.surface)
		if p.measure {
			p.draw += time.Since(started)
		}
		p.report.Frames++
		capture := p.tick < 12 || p.tick%31 == 0 || p.tick == p.limit
		for _, cue := range []int{200, 600, 900, 1300, 1800} {
			capture = capture || p.tick >= cue-2 && p.tick <= cue+12
		}
		if capture {
			p.surface.ReadPixels(p.pixels)
			digest := sha256.Sum256(p.pixels)
			p.report.Samples = append(p.report.Samples, sample{p.tick, hex.EncodeToString(digest[:])})
		}
		p.done = p.tick >= p.limit
		p.tick++
	}
	dst.DrawImage(p.surface, nil)
}

func (p *probe) Close() {
	if p.game != nil {
		p.game.Close()
	}
	if p.surface != nil {
		p.surface.Deallocate()
	}
}

func main() {
	output := flag.String("output", "", "write complete-frame fingerprints as JSON")
	limit := flag.Int("last-tick", 9000, "last 50 Hz update tick to render")
	outline := flag.Bool("wireframe", false, "capture the original outlined geometry mode")
	measure := flag.Bool("timing", false, "measure muted CPU update/draw submission, excluding readback and GPU completion")
	flag.Parse()
	if *output == "" || *limit < 0 || *limit > 60000 {
		log.Fatal("-output and a last tick within 0..60000 are required")
	}
	p := &probe{limit: *limit, measure: *measure, report: report{Wireframe: *outline,
		Width: demo.Width, Height: demo.Height, Rate: demo.FPS}}
	defer p.Close()
	ebiten.SetWindowSize(demo.Width, demo.Height)
	ebiten.SetWindowTitle("OldSkool / frame verification")
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetVsyncEnabled(false)
	if err := ebiten.RunGame(p); err != nil {
		log.Fatal(err)
	}
	if p.err != nil {
		log.Fatal(p.err)
	}
	if p.measure {
		p.report.UpdateMeanUS = float64(p.update) / float64(p.report.Frames) / 1000
		p.report.DrawMeanUS = float64(p.draw) / float64(p.report.Frames) / 1000
	}
	data, err := json.MarshalIndent(p.report, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(*output, append(data, '\n'), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d frames drawn; %d fingerprints recorded; wireframe=%t\n", p.report.Frames, len(p.report.Samples), p.report.Wireframe)
}
