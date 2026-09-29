# Native composition

The original PE32 program renders at 640 × 480 through DirectX 8. Its frame
delta is measured from a performance counter and capped at 1/30 second. Most
effect controllers advance from a common `delta × 85` clock; soundtrack data
is YM5 at 50 frames/s. This Go version uses 50 logical updates/s. The
reference recording runs for 101.546 seconds at 30 video frames/s, while the
original chip track contains 253.44 seconds before looping.

| Source routine | First cue | Native composition |
| --- | ---: | --- |
| `004010c0` | 0 s | Rotating and moving colored cube, built with DCK's solid-cube geometry. |
| `00401570` | 0 s | Eighty music-register level columns, decoded through DCK's YM stream. |
| `00402b30` | 4 s | Sixteen pastel ribbon strips, batched through DCK. |
| `00402d40` | 12 s | Eighty original ball images placed by four spaced source sine banks and drawn by DCK's sprite renderer. |
| `00402660` | 18 s | The 180 × 80 Oxygene logo, projected in ten strips by DCK's image warp. |
| `00402110` | 26 s | The original 8 × 8 font, advanced by DCK's scrolling transport and painted as separate shaded cells. |
| `00401a40` | 36 s | The original 32 × 30 font, advanced by a second DCK scroll and painted in bounded cells. |

The raster colors come from the 16 original vertex-color words. The original
oscillator pairs set the ribbon endpoints and all 80 ball positions; the logo
uses four additional source oscillators for its vertical and row motion. The
big and small text scroll at 510 and 595 source units/s. Their bitmap character
orders and messages are recovered from the executable. RGB-tinted font cells,
3D perspective and clipping are native approximations of Direct3D's raster
rules, using bounded GPU batches and persistent DCK surfaces.

The chip meter uses the original sixteen-word volume table, tone bins chosen
by `round(10000/period)`, 8-pixel spacing and a 3-pixel decay per update.

**P** freezes the graphics while the music continues, as in the original
message loop. **W** outlines the cube, the sixteen ribbons, the ball quads,
the warped logo strips and both fonts' individual cells. The outline mode
uses bounded DCK batches; polygon-edge coverage can differ by a pixel from
the Direct3D 8 renderer.
