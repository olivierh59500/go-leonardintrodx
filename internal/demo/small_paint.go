package demo

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/geometry"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/go-leonardintrodx/internal/source"
)

const smallHalfCell = source.SmallCellSide / 2

// The source's first four vertices are white and its lower four are the
// yellow, orange, blue and cyan corners of every twenty-unit cubelet.
var smallCellCorners = [8]geometry.Vec3{
	{X: -smallHalfCell, Y: smallHalfCell, Z: -smallHalfCell},
	{X: smallHalfCell, Y: smallHalfCell, Z: -smallHalfCell},
	{X: smallHalfCell, Y: smallHalfCell, Z: smallHalfCell},
	{X: -smallHalfCell, Y: smallHalfCell, Z: smallHalfCell},
	{X: -smallHalfCell, Y: -smallHalfCell, Z: -smallHalfCell},
	{X: smallHalfCell, Y: -smallHalfCell, Z: -smallHalfCell},
	{X: smallHalfCell, Y: -smallHalfCell, Z: smallHalfCell},
	{X: -smallHalfCell, Y: -smallHalfCell, Z: smallHalfCell},
}

var smallCellColors = [8]color.RGBA{
	{255, 255, 255, 255}, {255, 255, 255, 255},
	{255, 255, 255, 255}, {255, 255, 255, 255},
	{255, 255, 0, 255}, {255, 128, 0, 255},
	{0, 0, 255, 255}, {0, 255, 255, 255},
}

var smallCellFaces = [6][6]int{
	{0, 3, 1, 1, 3, 2},
	{0, 1, 4, 1, 5, 4},
	{1, 2, 5, 2, 6, 5},
	{7, 2, 3, 7, 6, 2},
	{0, 7, 3, 4, 7, 0},
	{4, 6, 7, 4, 5, 6},
}

var smallCellFaceCorners = [6][4]int{
	{0, 3, 2, 1}, {0, 1, 5, 4}, {1, 2, 6, 5},
	{7, 6, 2, 3}, {0, 4, 7, 3}, {4, 5, 6, 7},
}

type smallVisibleFace struct {
	index int
	depth float64
}

func newSmallScrolling(face scrolling.Face, text string, white *ebiten.Image, wireframe func() bool) (*scrolling.Scrolling, error) {
	painter, err := newSmallPainter(white, wireframe)
	if err != nil {
		return nil, err
	}
	return scrolling.New(scrolling.Config{
		Text: text, Fonts: map[string]scrolling.Face{"original": face}, Font: "original",
		Speed: 85 * 7, X: 39 * source.SmallColumnPitch, Repeat: true,
		Map: func(sample scrolling.Sample, op *ebiten.DrawImageOptions) bool {
			if sample.Glyph.Image == nil {
				return false
			}
			left, _ := op.GeoM.Apply(0, 0)
			right, _ := op.GeoM.Apply(float64(sample.Glyph.Image.Bounds().Dx()), 0)
			// The source's 3D projection can bring a whole glyph into view.
			return right > source.SmallWorldOffset-400 && left < source.SmallWorldOffset+400
		},
		Shape: "source-3d", Modes: map[string]scrolling.Mode{"source-3d": {Paint: painter}},
	})
}

// newSmallPainter projects the executable's cubelet vertices directly onto
// the complete game surface. It avoids clipping the moving rows to an atlas.
func newSmallPainter(white *ebiten.Image, wireframe func() bool) (scrolling.Painter, error) {
	font, err := readSmallFont()
	if err != nil {
		return nil, err
	}
	batch := render.NewBatch(4096)
	camera := geometry.Camera{
		Center: geometry.Vec2{X: Width / 2, Y: Height / 2},
		Focal:  Height * math.Sqrt(3) / 2, Near: .1,
	}
	return func(dst *ebiten.Image, sample scrolling.Sample, op ebiten.DrawImageOptions) {
		if sample.Glyph.Rune > 255 {
			return
		}
		rows := source.SmallScrollRows(sample.Time)
		outlined := wireframe()
		batch.Begin(dst, white)
		for row := 0; row < source.SmallGlyphHeight; row++ {
			for column := 0; column < source.SmallGlyphWidth; column++ {
				if !font.Pixel(byte(sample.Glyph.Rune), column, row) {
					continue
				}
				worldColumn, _ := op.GeoM.Apply(float64(column), float64(row))
				center := geometry.Vec3{
					X: worldColumn - source.SmallWorldOffset,
					Y: rows[row].Y,
					Z: source.CubeViewZ + rows[row].Z,
				}
				var world [8]geometry.Vec3
				var screen [8]geometry.Vec2
				minX, minY, maxX, maxY := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
				for i, corner := range smallCellCorners {
					world[i] = center.Add(corner)
					// The source uses upward-positive Y; DCK's camera uses downward-positive Y.
					screen[i], _, _ = camera.Project(geometry.Vec3{X: world[i].X, Y: -world[i].Y, Z: world[i].Z})
					minX, minY = min(minX, screen[i].X), min(minY, screen[i].Y)
					maxX, maxY = max(maxX, screen[i].X), max(maxY, screen[i].Y)
				}
				if maxX < 0 || minX >= Width || maxY < 0 || minY >= Height {
					continue
				}
				var faces [6]smallVisibleFace
				count := 0
				for face, indices := range smallCellFaces {
					a, b, c := world[indices[0]], world[indices[1]], world[indices[2]]
					if b.Sub(a).Cross(c.Sub(a)).Dot(a) >= 0 {
						continue
					}
					depth := (world[smallCellFaceCorners[face][0]].Z + world[smallCellFaceCorners[face][1]].Z +
						world[smallCellFaceCorners[face][2]].Z + world[smallCellFaceCorners[face][3]].Z) / 4
					insert := count
					for insert > 0 && faces[insert-1].depth < depth {
						faces[insert] = faces[insert-1]
						insert--
					}
					faces[insert] = smallVisibleFace{index: face, depth: depth}
					count++
				}
				for _, visible := range faces[:count] {
					face := visible.index
					if outlined {
						corners := smallCellFaceCorners[face]
						for edge := 0; edge < 4; edge++ {
							strokeSmallEdge(batch, screen[corners[edge]], screen[corners[(edge+1)%4]])
						}
						continue
					}
					indices := smallCellFaces[face]
					for triangle := 0; triangle < 6; triangle += 3 {
						var vertices [3]ebiten.Vertex
						for corner := 0; corner < 3; corner++ {
							index := indices[triangle+corner]
							point := screen[index]
							vertices[corner] = render.Vertex(point.X, point.Y, 0, 0, smallCellColors[index])
						}
						batch.Triangle(vertices[0], vertices[1], vertices[2])
					}
				}
			}
		}
		batch.Flush()
	}, nil
}

func strokeSmallEdge(batch *render.Batch, a, b geometry.Vec2) {
	dx, dy := b.X-a.X, b.Y-a.Y
	length := math.Hypot(dx, dy)
	if length == 0 {
		return
	}
	nx, ny := -dy/length*.6, dx/length*.6
	paint := color.RGBA{255, 255, 255, 255}
	batch.Quad([4]ebiten.Vertex{
		render.Vertex(a.X+nx, a.Y+ny, 0, 0, paint),
		render.Vertex(b.X+nx, b.Y+ny, 1, 0, paint),
		render.Vertex(b.X-nx, b.Y-ny, 1, 1, paint),
		render.Vertex(a.X-nx, a.Y-ny, 0, 1, paint),
	})
}
