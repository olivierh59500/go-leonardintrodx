package source

import "math"

// YMVolumeTable is the sixteen-word level table in the original executable.
var YMVolumeTable = [16]uint16{
	62, 161, 265, 377, 580, 774, 1155, 1575,
	2260, 3088, 4570, 6233, 9330, 13187, 21220, 32767,
}

// Meter retains eighty 8-pixel columns. Three YM tone periods select columns;
// their original volume levels charge them before a three-pixel decay.
type Meter struct{ levels [80]float64 }

func (m *Meter) Step(registers [14]uint8) {
	envelope := (int(registers[12])*256 + int(registers[11]&0xc0)) >> 6
	for channel := 0; channel < 3; channel++ {
		volume := registers[8+channel] & 0x1f
		period := int(registers[channel*2]) + int(registers[channel*2+1]&0x0f)*256
		strength := int(volume & 0x0f)
		if volume&0x10 != 0 {
			strength = 15
		} else if registers[7]&(1<<channel) != 0 {
			strength = 0
		}
		m.charge(period, strength)
		if volume&0x10 != 0 {
			m.charge(envelope, 15)
		}
	}
	for index := range m.levels {
		m.levels[index] = max(0, m.levels[index]-3)
	}
}

func (m *Meter) charge(period, strength int) {
	if period <= 0 || period >= 4095 || strength < 0 || strength > 15 {
		return
	}
	column := int(math.Round(10000 / float64(period)))
	if column < 0 || column >= len(m.levels) {
		return
	}
	height := float64(YMVolumeTable[strength]) * .003
	m.levels[column] = max(m.levels[column], height)
}

func (m *Meter) Level(column int) float64 {
	if column < 0 || column >= len(m.levels) {
		return 0
	}
	return m.levels[column]
}
