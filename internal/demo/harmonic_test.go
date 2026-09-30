package demo

import (
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/sprites"
)

// The source oscillator remains an independent oracle for every production
// tick and item, including poses between complete-frame fingerprint samples.
func TestSharedHarmonicPosesMatchNativeOscillators(t *testing.T) {
	white := ebiten.NewImage(1, 1)
	defer white.Deallocate()
	bands, err := newNativeBands(white)
	if err != nil {
		t.Fatal(err)
	}
	defer bands.Close()
	balls, err := sprites.NewHarmonicField(sprites.HarmonicFieldConfig{
		Count: 80, Motion: nativeFormation(ballWaves, 320, 240), ClockScale: [2]float64{85}, PixelSnap: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer balls.Close()
	for tick := 0; tick <= 9000; tick++ {
		time := float64(tick) / FPS
		if time >= BallsStart {
			local := time - BallsStart
			if err := bands.Update(kit.Frame{Time: local}); err != nil {
				t.Fatal(err)
			}
			for i, p := range bands.Poses() {
				left := max(0, min(463, math.Round(232+ribbonWaves[0].At(local, i)+ribbonWaves[1].At(local, i))))
				right := max(0, min(463, math.Round(232+ribbonWaves[2].At(local, i)+ribbonWaves[3].At(local, i))))
				if p.X != left || p.Y != right {
					t.Fatalf("band %d tick %d: got %+v, want (%v,%v)", i, tick, p, left, right)
				}
			}
		}
		if time >= BallFieldStart {
			local := time - BallFieldStart
			if err := balls.Update(kit.Frame{Time: local}); err != nil {
				t.Fatal(err)
			}
			for i, p := range balls.Samples() {
				x := math.Round(320 + ballWaves[0].At(local, i) + ballWaves[1].At(local, i))
				y := math.Round(240 + ballWaves[2].At(local, i) + ballWaves[3].At(local, i))
				if p.X != x || p.Y != y {
					t.Fatalf("ball %d tick %d: got (%v,%v), want (%v,%v)", i, tick, p.X, p.Y, x, y)
				}
			}
		}
	}
}
