package source

const (
	SmallCellSide    = 20
	SmallColumnPitch = 28
	SmallRowPitch    = 30
	SmallWorldOffset = 470
)

// SmallScrollRow is the world-space center of one row of the original
// eight-row cubelet scroller. The column transport supplies X separately.
type SmallScrollRow struct{ Y, Z float64 }

var smallGlobalY = [2]Oscillator{
	{Amplitude: 70, Rate: .0067}, {Amplitude: 70, Rate: .0239},
}
var smallGlobalZ = [2]Oscillator{
	{Amplitude: 30, Rate: .0097}, {Amplitude: 30, Rate: .0339},
}
var smallRowY = Oscillator{Amplitude: 30, Rate: .103, Spacing: .222}
var smallRowZ = Oscillator{Amplitude: 10, Rate: .1, Spacing: .17}

// SmallScrollRows evaluates the source oscillators after the effect cue.
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
