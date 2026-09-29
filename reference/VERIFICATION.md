# Fidelity and verification

The embedded fonts, indexed artwork, ball texture, icons, scrolling messages
and YM5 soundtrack come from the DirectX 8 executable.
The asset manifest records source and resulting hashes. Source x86 routines
identify the 640 × 480 viewport, effect gates, 16 raster colors, sine
controller parameters, font encodings and the 80-column chip meter.

The [DirectX 8 edition recording](https://x.com/olivierhoute/status/1697733388201853139)
contains 101.546 seconds at 30 frames/s. Comparison frames cover the opening,
each effect entrance, and the 50-, 80- and 100-second marks. A desktop run
with music and all effects completed 500 updates from second 40. The embedded
YM soundtrack is non-silent and loops after its 253.44-second declared
duration.

The presentation MP4 runs for 180 seconds with 9,000 frames at 640 × 480 and
50 frames/s, plus 48 kHz stereo AAC audio. Both streams decode without errors;
the audio has a -18.1 dB mean and -1.5 dB maximum. The poster is a native
canvas capture at second 50.

Original assets, control text, YM data and the principal sine trajectories
come from the executable. DCK's cube projection, shaded bitmap cells and
Direct3D 8 edge coverage are native adaptations, so individual pixels and
some palette transitions differ from the Windows renderer. The public clip
covers the first 101 seconds; later motion follows the same recovered
controllers and the music loops indefinitely.
