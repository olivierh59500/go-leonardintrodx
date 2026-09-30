package source

import (
	"math"
	"testing"

	"github.com/olivierh59500/democonstructionkit/motion"
)

// SmallScrollRows retains the independent source equation for pose verification.
func SmallScrollRows(seconds float64) [SmallGlyphHeight]SmallScrollRow {
	globalY := smallGlobalY[0].At(seconds, 0) + smallGlobalY[1].At(seconds, 0)
	globalZ := smallGlobalZ[0].At(seconds, 0) + smallGlobalZ[1].At(seconds, 0)
	var rows [SmallGlyphHeight]SmallScrollRow
	for row := range rows {
		rows[row] = SmallScrollRow{
			Y: 120 + globalY - float64(row*SmallRowPitch) - smallRowY.At(seconds, row),
			Z: 140 + globalZ + smallRowZ.At(seconds, row),
		}
	}
	return rows
}

func TestSharedSmallScrollProfileMatchesEveryNativePose(t *testing.T) {
	profile, err := NewSmallScrollProfile()
	if err != nil {
		t.Fatal(err)
	}
	for tick := 0; tick <= 9000; tick++ {
		seconds := float64(tick) / 50
		expected := SmallScrollRows(seconds)
		for row, original := range expected {
			got, ok := profile.Apply(row, seconds, motion.Point{X: 120, Y: 140})
			if !ok || got.X != original.Y || got.Y != original.Z {
				t.Fatalf("tick %d row %d: %v; want (%v,%v)", tick, row, got, original.Y, original.Z)
			}
		}
	}
}

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
