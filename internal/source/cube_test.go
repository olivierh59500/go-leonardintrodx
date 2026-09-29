package source

import (
	"math"
	"testing"
)

func TestCubePositionMatchesOriginalOscillators(t *testing.T) {
	for _, tt := range []struct {
		seconds, x, y, z float64
	}{
		{0, 0, 0, 100},
		{1, 82.64887709158766, 60.94822207117116, 140.1544909447826},
		{5, 102.94443869598712, -5.8602193800396165, 92.9631323554035},
	} {
		x, y, z := CubePosition(tt.seconds)
		if math.Abs(x-tt.x) > 1e-9 || math.Abs(y-tt.y) > 1e-9 || math.Abs(z-tt.z) > 1e-9 {
			t.Fatalf("t=%g: got (%g, %g, %g), want (%g, %g, %g)", tt.seconds, x, y, z, tt.x, tt.y, tt.z)
		}
	}
}
