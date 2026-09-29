package source

import (
	"math"
	"testing"
)

func TestOriginalYMRegistersChargeFrequencyColumns(t *testing.T) {
	// Decoder registers from the first frame of the embedded Bankok Knights YM5.
	first := [14]uint8{145, 1, 134, 3, 35, 3, 9, 248, 14, 15, 14, 0, 0, 0}
	var meter Meter
	meter.Step(first)
	for column, want := range map[int]float64{11: 95.301, 12: 60.66, 25: 60.66} {
		if math.Abs(meter.Level(column)-want) > 0.001 {
			t.Fatalf("column %d: level %.3f, want %.3f", column, meter.Level(column), want)
		}
	}
	first[7] = 255
	meter.Step(first)
	if math.Abs(meter.Level(25)-57.66) > 0.001 {
		t.Fatal("inactive tone did not decay by three pixels")
	}
}
