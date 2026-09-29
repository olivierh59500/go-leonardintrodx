package source

const (
	CubeSide  = 60
	CubeViewZ = 200
)

// CubePosition evaluates the six translation oscillators stored in the
// original executable. Their phases advance at 85 source units per second.
func CubePosition(seconds float64) (x, y, z float64) {
	x = cubeWaves[0].At(seconds, 0) + cubeWaves[1].At(seconds, 0)
	y = cubeWaves[2].At(seconds, 0) + cubeWaves[3].At(seconds, 0)
	z = 100 + cubeWaves[4].At(seconds, 0) + cubeWaves[5].At(seconds, 0)
	return
}

var cubeWaves = [6]Oscillator{
	{Amplitude: 60, Rate: .0323}, {Amplitude: 60, Rate: .02},
	{Amplitude: 50, Rate: .0321}, {Amplitude: 50, Rate: .0257},
	{Amplitude: 50, Rate: .0357}, {Amplitude: 50, Rate: .0279},
}
