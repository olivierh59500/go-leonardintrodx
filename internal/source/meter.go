package source

import "github.com/olivierh59500/democonstructionkit/sound"

// YMVolumeTable is the sixteen-word level table in the original executable.
var YMVolumeTable = [16]uint16{
	62, 161, 265, 377, 580, 774, 1155, 1575,
	2260, 3088, 4570, 6233, 9330, 13187, 21220, 32767,
}

// Meter configures the production's frequency columns and level curve. DCK
// retains peaks, applies YM gating and charges them before the per-tick decay.
// The zero value initializes its controller on the first Step.
type Meter struct{ controller *sound.YMPeriodMeter }

func (m *Meter) Step(registers [14]uint8) {
	if m.controller == nil {
		config := sound.YMPeriodMeterConfig{
			Columns: 80, Gain: .003, Decay: 3, PeriodMin: 1, PeriodMaxExclusive: 4095,
			FrequencyScale: 10000, Rounding: sound.YMPeriodRoundNearest, Gating: sound.YMPeriodGateFixed,
			Envelope: true, EnvelopeShift: 6, EnvelopeLevel: 15, EnvelopeGain: 1,
		}
		for index, level := range YMVolumeTable {
			config.Levels[index] = float64(level)
		}
		var err error
		m.controller, err = sound.NewYMPeriodMeter(config)
		if err != nil {
			panic(err)
		}
	}
	_ = m.controller.Step(registers)
}

func (m *Meter) Level(column int) float64 {
	if m == nil {
		return 0
	}
	return m.controller.Level(column)
}
