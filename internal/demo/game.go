package demo

import (
	"errors"
	"fmt"
	"image/color"
	"io"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/effects"
	"github.com/olivierh59500/democonstructionkit/font"
	"github.com/olivierh59500/democonstructionkit/geometry"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/democonstructionkit/sound"
	playback "github.com/olivierh59500/democonstructionkit/sound/ebiten"
	"github.com/olivierh59500/democonstructionkit/sprites"
	"github.com/olivierh59500/go-leonardintrodx/assets"
	"github.com/olivierh59500/go-leonardintrodx/internal/source"
)

const (
	Width          = 640
	Height         = 480
	FPS            = 50
	BallsStart     = 4.0
	BallFieldStart = 12.0
	LogoStart      = 18.0
	SmallTextStart = 26.0
	LargeTextStart = 36.0
)

// Game reproduces the overlapping effect sequence of the DirectX 8 edition.
// Assets and authored control data are local; reusable rendering and music are DCK.
type Game struct {
	art                                       *artwork
	cube                                      *effects.MeshEffect
	balls                                     *sprites.HarmonicField
	bands                                     *composite.HarmonicBands
	logo, large                               *effects.Warp
	small                                     *scrolling.Scrolling
	white                                     *ebiten.Image
	batch                                     *render.Batch
	visual                                    *sound.Stream
	visualPrimed                              bool
	music                                     *playback.Player
	musicData                                 []byte
	samples                                   [960 * 8]byte
	registers                                 [14]uint8
	meter                                     source.Meter
	tick                                      int
	endTick                                   int
	paused, wireframe, showLoad, mute, closed bool
	cubeRealtime                              bool
	lastUpdate                                time.Time
	layer                                     string
}

func NewGame(start int, mute bool) (_ *Game, err error) {
	if start < 0 {
		return nil, fmt.Errorf("demo: negative start tick")
	}
	g := &Game{tick: start, mute: mute, batch: render.NewBatch(4096)}
	defer func() {
		if err != nil {
			g.Close()
		}
	}()
	if g.art, err = loadArtwork(); err != nil {
		return nil, err
	}
	g.musicData, err = assets.Files.ReadFile("original/music.ym")
	if err != nil {
		return nil, err
	}
	g.visual, err = sound.Open("original/music.ym", g.musicData, sound.Options{SampleRate: 48000, Loop: true})
	if err != nil {
		return nil, err
	}
	if err = g.primeMeter(start); err != nil {
		return nil, err
	}
	g.white = ebiten.NewImage(1, 1)
	g.white.Fill(color.White)
	cubeMesh := effects.Cube(source.CubeSide, geometry.Vec2{X: 1, Y: 1}, color.NRGBA{R: 255, G: 200, B: 255, A: 255})
	nativeCubeColors := [6]color.NRGBA{
		{R: 255, G: 240, B: 255, A: 255}, {R: 255, G: 224, B: 255, A: 255},
		{R: 255, G: 208, B: 255, A: 255}, {R: 255, G: 192, B: 255, A: 255},
		{R: 255, G: 176, B: 255, A: 255}, {R: 255, G: 160, B: 255, A: 255},
	}
	// Map the executable's +Z,+X,+Y,-X,-Z,-Y face order to DCK's
	// -Z,+Z,-X,+X,-Y,+Y order after reflecting screen Y.
	for face, sourceFace := range [...]int{4, 0, 3, 1, 2, 5} {
		paint := nativeCubeColors[sourceFace]
		cubeMesh.Triangles[face*2].Color = paint
		cubeMesh.Triangles[face*2+1].Color = paint
	}
	if g.cube, err = effects.NewMesh(cubeMesh, nil, geometry.Camera{
		Center: geometry.Vec2{X: Width / 2, Y: Height / 2}, Focal: Height * math.Sqrt(3) / 2, Near: .1,
	}); err != nil {
		return nil, err
	}
	g.cube.CullBackFaces = true
	if err = g.cube.SetOutline(effects.MeshOutlineConfig{Faces: effects.CubeFaces(), Width: 1.2,
		Color: color.NRGBA{R: 255, G: 198, B: 255, A: 255}, White: g.white}); err != nil {
		return nil, err
	}
	if g.bands, err = newNativeBands(g.white); err != nil {
		return nil, err
	}
	if g.balls, err = sprites.NewHarmonicField(sprites.HarmonicFieldConfig{
		Count: 80, Motion: nativeFormation(ballWaves, 320, 240), ClockScale: [2]float64{85, 0}, PixelSnap: true,
		Style: sprites.FieldStyle{Image: g.art.ball, Appearance: sprites.FieldAppearance{Width: 32, Height: 32, AnchorX: .5, AnchorY: .5}},
	}); err != nil {
		return nil, err
	}
	if g.logo, err = newLogoWarp(g.art.logo); err != nil {
		return nil, err
	}
	if err = g.logo.SetOutline(effects.WarpOutlineConfig{Width: 1.2,
		Color: color.NRGBA{R: 235, G: 135, B: 245, A: 255}, White: g.white}); err != nil {
		return nil, err
	}
	if g.small, err = newSmallScrolling(g.art.smallFace, g.art.smallCells, g.art.smallText, g.white, func() bool { return g.wireframe }); err != nil {
		return nil, err
	}
	if g.large, err = newLargeTextWarp(g.art.largeFace, g.art.largeCells, g.art.largeText, g.white, func() bool { return g.wireframe }); err != nil {
		return nil, err
	}
	if err = g.prepare(g.Seconds()); err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Game) primeMeter(tick int) error {
	g.meter = source.Meter{}
	warmup := min(max(0, tick), 50)
	if _, err := g.visual.Seek(int64(tick-warmup)*960*8, io.SeekStart); err != nil {
		return err
	}
	for i := 0; i <= warmup; i++ {
		if _, err := io.ReadFull(g.visual, g.samples[:]); err != nil {
			return err
		}
		g.registers, _ = g.visual.YMRegisters()
		g.meter.Step(g.registers)
	}
	g.visualPrimed = true
	return nil
}

var logoWaves = [4]source.Oscillator{
	{Amplitude: 101.75, Rate: .0342, Spacing: .2736},
	{Amplitude: 101.75, Rate: .0143, Spacing: .1144},
	{Amplitude: 94.25, Rate: .0342},
	{Amplitude: 94.25, Rate: .0443},
}

func newLogoWarp(logo *ebiten.Image) (*effects.Warp, error) {
	layer := kit.Func{OnDraw: func(dst *ebiten.Image) {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(233.0/180.0, 103.0/80.0)
		dst.DrawImage(logo, op)
	}}
	warp, err := effects.NewWarp(layer, 233, 103, 1, 10)
	if err != nil {
		return nil, err
	}
	warp.Map = func(x, y, t float64) geometry.Vec2 {
		strip := int(math.Round(y * 10 / 103))
		return geometry.Vec2{
			X: 203.5 + x + logoWaves[0].At(t, strip) + logoWaves[1].At(t, strip),
			Y: 188.5 + y + logoWaves[2].At(t, 0) + logoWaves[3].At(t, 0),
		}
	}
	return warp, nil
}

func newLargeTextWarp(face scrolling.Face, bank *font.CellBank, text string, white *ebiten.Image, wireframe func() bool) (*effects.Warp, error) {
	const width, height, columns, rows = 720, 480, 40, 30
	cells := largeCells(bank, white, wireframe)
	scroll, err := scrolling.New(scrolling.Config{
		Text: text, Fonts: map[string]scrolling.Face{"original": face}, Font: "original",
		Speed: 510, X: 720, Y: 0, Repeat: true, Gap: 0,
		Map: func(sample scrolling.Sample, op *ebiten.DrawImageOptions) bool {
			if sample.Glyph.Image == nil {
				return false
			}
			x0, y0 := op.GeoM.Apply(0, 0)
			x1, y1 := op.GeoM.Apply(float64(sample.Glyph.Image.Bounds().Dx()), float64(sample.Glyph.Image.Bounds().Dy()))
			// The per-cell motion can extend beyond the glyph's ordinary bounds.
			const margin = 64.0
			return x1 > -margin && x0 < float64(width)+margin && y1 > -margin && y0 < float64(height)+margin
		},
		Shape: "authored", Modes: map[string]scrolling.Mode{"authored": {Cells: &cells}},
	})
	if err != nil {
		return nil, err
	}
	warp, err := effects.NewWarp(scroll, width, height, columns, rows)
	if err != nil {
		return nil, err
	}
	warp.Map = func(x, y, t float64) geometry.Vec2 {
		return geometry.Vec2{X: x - 20, Y: y + 20}
	}
	warp.Tint = func(x, y, t float64) color.Color {
		position := max(0, min(1, y/480+0.07*math.Sin(t*0.11)))
		warm := max(0, min(1, (t-24)/24))
		return color.RGBA{
			R: uint8(20 + 205*position*position),
			G: uint8(240 - 215*position),
			B: uint8(245 - 20*position - 165*warm*position), A: 255,
		}
	}
	return warp, nil
}

func (g *Game) Update() error {
	if g.endTick > 0 && g.tick >= g.endTick {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.ToggleWireframe()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		g.showLoad = !g.showLoad
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if err := g.TogglePause(); err != nil {
			return err
		}
	}
	if !g.mute && g.music == nil {
		var err error
		g.music, err = playback.Open(nil, "original/music.ym", g.musicData, sound.Options{SampleRate: 48000, Loop: true})
		if err != nil {
			return err
		}
		if g.tick > 0 {
			if err = g.music.Seek(time.Duration(g.tick) * time.Second / FPS); err != nil {
				return err
			}
		}
		g.music.Play()
	}
	if g.paused {
		return nil
	}
	if g.visual != nil {
		if g.visualPrimed {
			g.visualPrimed = false
		} else {
			if _, err := io.ReadFull(g.visual, g.samples[:]); err != nil {
				return err
			}
			g.registers, _ = g.visual.YMRegisters()
			g.meter.Step(g.registers)
		}
	}
	g.tick++
	if err := g.prepare(g.Seconds()); err != nil {
		return err
	}
	g.lastUpdate = time.Now()
	return nil
}

func (g *Game) prepare(t float64) error {
	g.cube.Transform = cubeTransform(t)
	if err := g.cube.Update(kit.Frame{Time: t}); err != nil {
		return err
	}
	if t >= BallsStart {
		if err := g.bands.Update(kit.Frame{Time: t - BallsStart}); err != nil {
			return err
		}
	}
	if t >= BallFieldStart {
		if err := g.balls.Update(kit.Frame{Time: t - BallFieldStart}); err != nil {
			return err
		}
	}
	if t >= LogoStart {
		if err := g.logo.Update(kit.Frame{Time: t - LogoStart}); err != nil {
			return err
		}
	}
	if t >= SmallTextStart {
		if err := g.small.Update(kit.Frame{Time: t - SmallTextStart}); err != nil {
			return err
		}
	}
	if t >= LargeTextStart {
		if err := g.large.Update(kit.Frame{Time: t - LargeTextStart}); err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) Seconds() float64         { return float64(g.tick) / FPS }
func (*Game) Layout(int, int) (int, int) { return Width, Height }

// SetTickLimit optionally ends a short interactive device check.
func (g *Game) SetTickLimit(count int) {
	if count > 0 {
		g.endTick = g.tick + count
	}
}

// SetWireframe selects outlined geometry for the composed scene.
func (g *Game) SetWireframe(enabled bool) { g.wireframe = enabled }

// ToggleWireframe switches between filled and outlined geometry.
func (g *Game) ToggleWireframe() { g.wireframe = !g.wireframe }

// Wireframe reports the current geometry mode for external controls.
func (g *Game) Wireframe() bool { return g.wireframe }

// TogglePause freezes or resumes the picture while the YM music continues.
func (g *Game) TogglePause() error {
	g.paused = !g.paused
	if !g.paused && g.music != nil {
		return g.primeMeter(int(math.Round(g.music.Position().Seconds() * FPS)))
	}
	return nil
}

// Paused reports whether visual updates are frozen.
func (g *Game) Paused() bool { return g.paused }

// EnableRealtimeCube interpolates the native cube pose between logical updates.
// Offline captures leave this disabled so their frames stay deterministic.
func (g *Game) EnableRealtimeCube(enabled bool) {
	g.cubeRealtime = enabled
	g.lastUpdate = time.Now()
}

func (g *Game) Draw(dst *ebiten.Image) {
	dst.Fill(color.Black)
	t := g.Seconds()
	if g.cubeRealtime && !g.paused {
		cubeTime := t + min(time.Since(g.lastUpdate).Seconds(), 1.0/FPS)
		g.cube.Transform = cubeTransform(cubeTime)
		_ = g.cube.Update(kit.Frame{Time: cubeTime})
	}
	switch g.layer {
	case "atlas-large":
		dst.DrawImage(g.art.largeAtlas, nil)
		return
	case "atlas-small":
		dst.DrawImage(g.art.smallAtlas, nil)
		return
	case "raw-large":
		g.large.Source.Draw(dst)
		return
	case "raw-small":
		g.small.Draw(dst)
		return
	case "large":
		g.large.Draw(dst)
		return
	case "small":
		g.small.Draw(dst)
		return
	case "logo":
		g.logo.Draw(dst)
		return
	}
	if g.wireframe {
		g.drawWireframe(dst, t)
		return
	}
	if t >= BallsStart {
		g.bands.Draw(dst)
	}
	if t >= BallFieldStart {
		g.balls.Draw(dst)
	}
	if t >= LogoStart {
		g.logo.Draw(dst)
	}
	if t >= LargeTextStart {
		g.large.Draw(dst)
	}
	g.cube.Draw(dst)
	if t >= SmallTextStart {
		g.small.Draw(dst)
	}
	g.drawMeters(dst)
	if g.showLoad {
		g.drawLoad(dst)
	}
}

// SetLayer selects one composition for a native diagnostic capture.
func (g *Game) SetLayer(layer string) error {
	switch layer {
	case "", "atlas-large", "atlas-small", "raw-large", "raw-small", "large", "small", "logo":
		g.layer = layer
		return nil
	default:
		return fmt.Errorf("demo: unknown layer %q", layer)
	}
}

func (g *Game) Close() error {
	if g.closed {
		return nil
	}
	g.closed = true
	var err error
	for _, e := range []kit.Effect{g.logo, g.small, g.large} {
		if e != nil {
			err = errors.Join(err, kit.Close(e))
		}
	}
	if g.cube != nil {
		err = errors.Join(err, g.cube.Close())
	}
	if g.balls != nil {
		err = errors.Join(err, g.balls.Close())
	}
	if g.bands != nil {
		err = errors.Join(err, g.bands.Close())
	}
	if g.music != nil {
		err = errors.Join(err, g.music.Close())
	}
	if g.visual != nil {
		err = errors.Join(err, g.visual.Close())
	}
	if g.art != nil {
		g.art.close()
	}
	if g.white != nil {
		g.white.Deallocate()
	}
	return err
}
