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
	"github.com/olivierh59500/democonstructionkit/effects"
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
	cube                                      *effects.SolidCube
	ballRenderer                              *sprites.FieldRenderer
	ballSamples                               [80]sprites.FieldSample
	logo, small, large                        *effects.Warp
	white                                     *ebiten.Image
	batch                                     *render.Batch
	visual                                    *sound.Stream
	visualPrimed                              bool
	music                                     *playback.Player
	musicData                                 []byte
	samples                                   [960 * 8]byte
	registers                                 [14]uint8
	tick                                      int
	endTick                                   int
	paused, wireframe, showLoad, mute, closed bool
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
	if start > 0 {
		_, err = g.visual.Seek(int64(start)*960*8, io.SeekStart)
		if err != nil {
			return nil, err
		}
	}
	if _, err = io.ReadFull(g.visual, g.samples[:]); err != nil {
		return nil, err
	}
	g.registers, _ = g.visual.YMRegisters()
	g.visualPrimed = true
	g.white = ebiten.NewImage(1, 1)
	g.white.Fill(color.White)
	cubeConfig := effects.DefaultSolidCubeConfig(74)
	cubeConfig.Perspective = 310
	cubeConfig.EdgeWidth = 0
	cubeConfig.CullBackFaces = true
	cubeConfig.FaceColors = [6]color.RGBA{
		{250, 170, 253, 255}, {255, 224, 255, 255}, {255, 194, 255, 255},
		{246, 156, 245, 255}, {245, 205, 254, 255}, {255, 187, 249, 255},
	}
	if g.cube, err = effects.NewSolidCube(cubeConfig); err != nil {
		return nil, err
	}
	g.ballRenderer = sprites.NewFieldRenderer(80)
	if g.logo, err = newLogoWarp(g.art.logo); err != nil {
		return nil, err
	}
	if g.small, err = newTextWarp(g.art.smallFace, g.art.smallText, g.white, false, func() bool { return g.wireframe }); err != nil {
		return nil, err
	}
	if g.large, err = newTextWarp(g.art.largeFace, g.art.largeText, g.white, true, func() bool { return g.wireframe }); err != nil {
		return nil, err
	}
	if err = g.prepare(g.Seconds()); err != nil {
		return nil, err
	}
	return g, nil
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

func newTextWarp(face scrolling.Face, text string, white *ebiten.Image, large bool, wireframe func() bool) (*effects.Warp, error) {
	width, height, columns, rows := 1120, 280, 40, 8
	speed, entry := 595.0, 1120.0
	var painter scrolling.Painter
	var err error
	if large {
		width, height, columns, rows = 720, 480, 40, 30
		speed, entry = 510, 720
		painter, err = newLargePainter(white, wireframe)
	} else {
		painter, err = newSmallPainter(white, wireframe)
	}
	if err != nil {
		return nil, err
	}
	scroll, err := scrolling.New(scrolling.Config{
		Text: text, Fonts: map[string]scrolling.Face{"original": face}, Font: "original",
		Speed: speed, X: entry, Y: 0, Repeat: true, Gap: 0,
		Shape: "authored", Modes: map[string]scrolling.Mode{"authored": {Paint: painter}},
	})
	if err != nil {
		return nil, err
	}
	warp, err := effects.NewWarp(scroll, width, height, columns, rows)
	if err != nil {
		return nil, err
	}
	if large {
		warp.Map = func(x, y, t float64) geometry.Vec2 {
			return geometry.Vec2{X: x - 230, Y: y + 20}
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
	} else {
		warp.Map = func(x, y, t float64) geometry.Vec2 {
			return geometry.Vec2{X: x - 160, Y: y + 60}
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
		g.wireframe = !g.wireframe
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		g.showLoad = !g.showLoad
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.paused = !g.paused
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
		}
	}
	g.tick++
	return g.prepare(g.Seconds())
}

func (g *Game) prepare(t float64) error {
	g.cube.Rotation = geometry.Vec3{X: 0.65 + t*1.57, Y: 0.4 + t*1.72, Z: 0.38 + t*0.56}
	if t >= BallFieldStart {
		g.updateBallSamples(t - BallFieldStart)
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
func (g *Game) Draw(dst *ebiten.Image) {
	dst.Fill(color.Black)
	t := g.Seconds()
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
		g.small.Source.Draw(dst)
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
		g.drawBands(dst, t)
	}
	if t >= BallFieldStart {
		g.ballRenderer.Draw(dst, g.ballSamples[:], sprites.FieldStyle{
			Image: g.art.ball, Appearance: sprites.FieldAppearance{Width: 32, Height: 32, AnchorX: .5, AnchorY: .5},
		})
	}
	if t >= LogoStart {
		g.logo.Draw(dst)
	}
	if t >= LargeTextStart {
		g.large.Draw(dst)
	}
	centerX, centerY := cubePosition(t)
	if g.wireframe {
		g.drawCubeWire(dst, centerX, centerY)
	} else {
		g.cube.DrawAt(dst, centerX, centerY)
	}
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
	if g.ballRenderer != nil {
		err = errors.Join(err, g.ballRenderer.Close())
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
