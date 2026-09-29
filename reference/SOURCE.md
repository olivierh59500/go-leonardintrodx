# Original program and assets

The 24,576-byte Windows PE32 release expands with UPX. The packed executable
has SHA-256
`e266c7202b5c5b2f7a389af11901f9ebc92befcec416ea3e82e8f410932913b5`.
The expanded image has SHA-256
`f4bb66e03c3a3b8247d96d1975ab95a1dee365cbec82ea451eca42ef117a62b1`.

The executable carries six numbered resources. Their embedded counterparts
and checksums are listed in `assets/manifest.json`. The 16- and 32-pixel
application icons preserve the original ball artwork in the desktop window.

| ID | Native content | Embedded asset |
| --- | --- | --- |
| 106 | 32 × 32 indexed ball bitmap | `original/ball.png` |
| 107 | Forty-eight 32 × 30 monochrome glyphs | `original/font-large.bin` |
| 108 | 256-entry character-to-glyph table | `original/font-map.bin` |
| 109 | One hundred 8 × 8 monochrome glyphs | `original/font-small.bin` |
| 110 | 180 × 80 indexed Oxygene logo | `original/logo.png` |
| 111 | YM5 soundtrack | `original/music.ym` |

The YM5 declares 12,672 frames at 50 Hz, a 2 MHz sound-chip clock, title
“Bankok Knights” and author Big Alec (Matt Gray). The original notice credits
Big Alec/Delta Force for this arrangement, Matt Gray for the C64 composition,
Mon/Oxygene for the logo and Leonard/Oxygene for the code and YM playback.
The scrolling message names Escape, W (wireframe), P (pause) and Tab (CPU load)
as interactive keys.

The [DirectX 8 edition recording](https://x.com/olivierhoute/status/1697733388201853139)
shows the second edition. The archived
[first-edition screenshot](https://www.pouet.net/prod.php?which=735) shows an
earlier version with different layout and timing.

The program opens a 640 × 480 window. A rotating cube and eighty sound-chip
level columns appear from the start. Sixteen pastel ribbons enter at 4
seconds; eighty textured balls enter at 12 seconds; the logo appears in ten
animated strips at 18 seconds. Independent bitmap-font compositions enter at
26 and 36 seconds. The recording covers 101.546 seconds of the continuous
intro; the YM file contains 253.44 seconds of music before its loop point.

The large font has four stored bytes per 30-pixel row. Bytes for columns 0–15
and 16–31 occupy separate 60-byte halves of each glyph. The 512-byte lookup
contains big-endian 16-bit entries; the original renderer reads the low byte
at `2*character+1`. The small font stores eight row bytes for each character
from ASCII space onward. Its corresponding composition starts at 26 seconds.
Both scrolling messages are embedded as `original/large-message.txt` and
`original/small-message.txt`.
The large composition uses 40 visible columns by 30 rows of 16-pixel blocks,
with a column pitch of 18 and a transport rate of 510 source units/second.
The small composition uses 40 visible columns by eight rows, a pitch of 28
and a transport rate of 595 source units/second. Each program begins advancing
when its effect gate opens, at 36 and 26 seconds respectively.
