# OldSkool DirectX 8 Go

A native Go/Ebitengine conversion of Leonard/Oxygene’s *Old skool Demo second
edition*, using Demo Construction Kit v1.0.0. The original Oxygene logo,
ball artwork, two bitmap fonts, two scrolling messages and YM5 music are
embedded in this version. DCK supplies the cube, batched sprite rendering,
font-aware scrolling, image warps, bounded rendering and audio playback.

```sh
GOWORK=off go run .
```

The 640 × 480 intro runs at 50 updates/s. Effects enter at 4, 12, 18, 26 and
36 seconds, then run together. **Escape** exits, **W** switches to wireframe,
**P** pauses the image while music continues, and **Tab** toggles the load bar.
`-start 45` opens the intro at 45 seconds; `-mute` silences playback and
`-wireframe` starts in outline mode. To save a still frame, use
`-start 45 -capture captures/check`.

The intro remains active while the original soundtrack loops. The
[source inventory](reference/SOURCE.md) covers the artwork, music and resource
formats. The [rendering notes](reference/RENDERING.md) describe the effects and
their DCK composition. [Fidelity and verification](reference/VERIFICATION.md)
covers the visual, audio and video checks.

`GOWORK=off go run ./cmd/video` exports a three-minute 640 × 480 MP4 at 50
frames/s with the original YM music and a native poster frame. `-duration`
chooses a different finite span of the looping intro. The capture contains
only the Ebitengine canvas and DCK audio.
