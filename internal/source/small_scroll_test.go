package source

import (
	"math"
	"testing"
)

func TestSmallScrollRowsMatchSourceOscillators(t *testing.T) {
	for _, tt := range []struct {
		seconds float64
		row     int
		y, z    float64
	}{
		{0, 0, 120, 140},
		{0, 7, -119.99576835057925, 149.28368967249168},
		{1, 0, 201.8222058720541, 177.72622619369278},
		{1, 7, 33.64922620706474, 167.120119566579},
	} {
		row := SmallScrollRows(tt.seconds)[tt.row]
		if math.Abs(row.Y-tt.y) > 1e-9 || math.Abs(row.Z-tt.z) > 1e-9 {
			t.Fatalf("t=%g row=%d: got (%g, %g), want (%g, %g)", tt.seconds, tt.row, row.Y, row.Z, tt.y, tt.z)
		}
	}
}
