# Native composition

The original PE32 program renders at 640 × 480 through DirectX 8. Its frame
delta is measured from a performance counter and capped at 1/30 second. Most
effect controllers advance from a common `delta × 85` clock; soundtrack data
is YM5 at 50 frames/s. This Go version uses 50 logical updates/s and
interpolates the cube pose between desktop updates. Offline captures keep a
fixed clock. The reference recording runs for 101.546 seconds at 30 video
frames/s, while the original chip track contains 253.44 seconds before looping.

| Source routine | First cue | Native composition |
| --- | ---: | --- |
| `004010c0` | 0 s | Rotating and moving colored cube, built with DCK's mesh renderer and the original six translation oscillators. |
| `00401570` | 0 s | Eighty music-register level columns, decoded through DCK's YM stream. |
| `00402b30` | 4 s | Sixteen pastel ribbon strips, batched through DCK. |
| `00402d40` | 12 s | Eighty original ball images placed by four spaced source sine banks and drawn by DCK's sprite renderer. |
| `00402660` | 18 s | The 180 × 80 Oxygene logo, projected in ten strips by DCK's image warp. |
| `00402110` | 26 s | The original 8 × 8 font, advanced by DCK's scrolling transport and drawn as source-colored, projected 3D cubelets. |
| `00401a40` | 36 s | The original 32 × 30 font, advanced by a second DCK scroll and painted in bounded cells with the original row oscillators. |

The cube uses the executable's 60-unit vertices, six face colors, 200-unit view
translation and projection matrix. At 640 × 480, the matrix gives a focal
length of 415.69 pixels. The X rotation advances at π/2 radians/s; the Y
rotation is 1.1 times that angle. A DCK mesh applies the original depth change
to both its center and its size.

The raster colors come from the 16 original vertex-color words. The original
oscillator pairs set the ribbon endpoints and all 80 ball positions; the logo
uses four additional source oscillators for its vertical and row motion. The
big and small text scroll at 510 and 595 source units/s. Their bitmap character
orders and messages are recovered from the executable. Large-text tinting and
Direct3D raster-edge coverage remain native adaptations, using bounded GPU
batches and persistent DCK surfaces.

The large scroller enters at the right edge. Its 40-column source surface uses
an 18-pixel column pitch and a 20-pixel left offset. Each of its 30 rows has
two 32-pixel sine motions with source rates 0.08 and 0.103 per speed-clock
unit. Its painter skips glyphs outside the source surface, allowing for cell
displacement.

The smaller scroller uses 40 columns spaced 28 world units and eight rows
spaced 30 units. Each lit font bit becomes a 20-unit cube with the original
white, yellow, orange, blue and cyan vertex colors. Four global oscillators
and two row oscillators move those cubes through the same source projection
as the main cube. DCK controls the message and transport; the cubelets draw
directly to the 640 × 480 canvas, with conservative 3D visibility culling.
They can move below the former 280-pixel intermediate surface without being
cut off.

The chip meter uses the original sixteen-word volume table, tone bins chosen
by `round(10000/period)`, 8-pixel spacing and a 3-pixel decay per update.

**P** freezes the graphics while the music continues, as in the original
message loop. **W** outlines the cube, the sixteen ribbons, the ball quads,
the warped logo strips and both fonts' individual cells. The outline mode
uses bounded DCK batches; polygon-edge coverage can differ by a pixel from
the Direct3D 8 renderer.
