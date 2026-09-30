# OldSkool DirectX 8 Go

A native Go/Ebitengine conversion of Leonard/Oxygene’s *Old skool Demo second
edition*, using Demo Construction Kit v1.0.7. The original Oxygene logo,
ball artwork, two bitmap fonts, two scrolling messages and YM5 music are
embedded in this version. DCK supplies the cube, batched sprite rendering,
font-aware scrolling, image warps, bounded rendering and audio playback.
The smaller scrolling text uses the original rainbow cubelets and moves across
the full screen without an intermediate clipping surface.

Both scrollings use DCK's owned `Mode.Cells`. Their native bitmaps are prepared
once in immutable `font.CellBank` caches. The intro supplies its original row
oscillators, cell dimensions, camera and colors; DCK owns flat rectangles,
cuboid faces, depth/culling and filled/wireframe rendering. All 750 sampled
complete frames match the preceding renderer across two 9,001-frame traversals.
One muted native replay per implementation/mode measured mean CPU draw
submissions of 737.34 to 715.49 microseconds filled and 1,118.58 to 1,029.59
wireframe; those timings exclude GPU completion and readback.
`go run ./cmd/checkframes -output captures/frames.json -timing` records complete
frame fingerprints; add `-wireframe` for the outlined composition.

The sixteen colored strips use DCK's `composite.HarmonicBands`, including their
filled and outlined materials. The eighty ball sprites use `sprites.HarmonicField`.
The intro supplies its four oscillator values, 85-unit clock, colors, integer
rounding and clamps; DCK samples and caches the poses once per update. All 750
complete-frame samples still match, and every applicable band/sprite position
matches the original oscillator routines through tick 9,000. No working image
is added. Absolute clocks preserve those poses when opening at a later time.

Cube, logo-grid and sprite-box outlines also use DCK materials. They reuse the
mesh's transformed points, the warp's current map/time and the sprite field's
cached samples; the intro no longer projects its own cube edges or loops over
logo/sprite borders. All 750 complete filled/wireframe samples remain identical.
The outline width, colors, native quad boundaries and layout remain parameters.
Optional `-cpuprofile /path/to/profile.pprof` on `cmd/checkframes` records a Go
CPU profile of the traversal; profiler timings include their own sampling cost.

```sh
GOWORK=off go run .
```

The 640 × 480 intro runs at 50 logical updates/s. On desktop, the cube also
interpolates its original motion between updates for smooth display refresh.
Effects enter at 4, 12, 18, 26 and 36 seconds, then run together. **Escape**
exits, **W** switches to wireframe, **P** pauses the image while music continues,
and **Tab** toggles the load bar.
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

## Android

`./scripts/run-android.sh` builds the ARM64 APK, installs it on the single
connected Android device and launches **OldSkool DirectX 8 Go**. Add
`--build-only` to produce the APK without installing it. The debug APK is at
`android/app/build/outputs/apk/debug/app-debug.apk`.

The original 640 × 480 canvas stays centered on a landscape display. Wide
sidebars contain **WIRE**, **PAUSE** and **RESET** touch buttons. Pause freezes
the picture while the YM music continues; Reset restarts both. The app keeps
the screen awake while it is in the foreground.

The 30 September 2026 ARM64 package embeds DCK v1.0.7 after the shared font-cell,
harmonic and outline migrations. Its native ELF load segments and APK placement
pass 16 KiB alignment. Desktop verification reproduces 750 complete filled and
wireframe captures on the published module.

An explicit verification launch can render over a locked screen without
dismissing the keyguard or changing the device's security settings:

```sh
adb shell am start -S -W -n com.olivierh.leonardintrodx/.MainActivity --ez oldskool_verify true
adb logcat -s GoLog:I AndroidRuntime:E
```

The host logs simulation/display cadence every ten seconds. The flag changes
only lock-screen visibility; clocks, effects, controls and music remain normal.

The DCK 1.0.6 package was installed and exercised on the Pixel 10a on
30 September 2026. After every effect had entered, 37 ten-second samples
recorded 49.2–50.9 logical updates/s and 59.7–60.2 displayed frames/s, with no
observed crash. This filled-mode runtime check is separate from the desktop
filled/wireframe image comparisons; it does not measure native/GPU memory.

The subsequent DCK 1.0.7 package was also installed and run on the Pixel.
After the final effect entrance, 35 ten-second samples recorded 49.2–50.9
updates/s and 59.8–60.1 displayed frames/s without an observed crash. The
published-module desktop replay still matches all 750 filled/wireframe samples.
