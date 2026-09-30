package demo

import (
	"math"
	"testing"

	"github.com/olivierh59500/democonstructionkit/motion"
)

func TestSharedLogoProfileMatchesEveryNativeVertex(t *testing.T) {
	profile, err := logoRowProfile()
	if err != nil {
		t.Fatal(err)
	}
	for tick := 0; tick <= 9000; tick++ {
		seconds := float64(tick) / FPS
		for row := 0; row < 10; row++ {
			for _, y := range []float64{float64(row*103) / 10, float64(row*103)/10 + 103.0/10} {
				strip := int(math.Round(y * 10 / 103))
				for _, x := range []float64{0, 233} {
					want := motion.Point{
						X: 203.5 + x + logoWaves[0].At(seconds, strip) + logoWaves[1].At(seconds, strip),
						Y: 188.5 + y + logoWaves[2].At(seconds, 0) + logoWaves[3].At(seconds, 0),
					}
					got, ok := profile.Apply(strip, seconds, motion.Point{X: 203.5 + x, Y: 188.5 + y})
					if !ok || got != want {
						t.Fatalf("tick %d vertex (%g,%g): %v; want %v", tick, x, y, got, want)
					}
				}
			}
		}
	}
}

func TestSharedLargeCellProfileMatchesEveryNativeRow(t *testing.T) {
	profile, err := largeRowProfile()
	if err != nil {
		t.Fatal(err)
	}
	for tick := 0; tick <= 9000; tick++ {
		seconds := float64(tick) / FPS
		for row := 0; row < profile.Len(); row++ {
			want := motion.Point{X: largeRowX.At(seconds, row), Y: largeRowY.At(seconds, row)}
			got, ok := profile.Apply(row, seconds, motion.Point{})
			if !ok || got != want {
				t.Fatalf("tick %d row %d: %v; want %v", tick, row, got, want)
			}
		}
	}
}
