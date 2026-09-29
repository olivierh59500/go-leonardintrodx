package assets

import (
	"image/png"
	"io"
	"testing"
	"time"

	"github.com/olivierh59500/democonstructionkit/sound"
)

func TestEmbeddedArtworkAndSoundtrack(t *testing.T) {
	for name, want := range map[string][2]int{
		"original/ball.png": {32, 32},
		"original/logo.png": {180, 80},
	} {
		file, err := Files.Open(name)
		if err != nil {
			t.Fatal(err)
		}
		image, err := png.Decode(file)
		file.Close()
		if err != nil {
			t.Fatal(err)
		}
		if bounds := image.Bounds(); bounds.Dx() != want[0] || bounds.Dy() != want[1] {
			t.Fatalf("%s: unexpected artwork size %v", name, bounds)
		}
	}
	data, err := Files.ReadFile("original/music.ym")
	if err != nil {
		t.Fatal(err)
	}
	stream, err := sound.Open("original/music.ym", data, sound.Options{SampleRate: 48000})
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()
	info := stream.Metadata()
	if info.Format != sound.FormatYM || info.Title != "Bankok Knights" ||
		info.Duration < 253*time.Second || info.Duration > 254*time.Second {
		t.Fatalf("unexpected decoded soundtrack: %+v", info)
	}
	pcm := make([]byte, 48000*8*5)
	n, err := io.ReadFull(stream, pcm)
	if err != nil || n != len(pcm) {
		t.Fatalf("soundtrack did not decode: %d samples, %v", n, err)
	}
	nonzero := false
	for _, sample := range pcm {
		if sample != 0 {
			nonzero = true
			break
		}
	}
	if !nonzero {
		t.Fatal("the decoded soundtrack is silent")
	}
}

func TestOriginalYMLoopContinuesPastSongDuration(t *testing.T) {
	data, err := Files.ReadFile("original/music.ym")
	if err != nil {
		t.Fatal(err)
	}
	stream, err := sound.Open("original/music.ym", data, sound.Options{SampleRate: 8000, Loop: true})
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()
	const duration = 256
	if count, err := io.CopyN(io.Discard, stream, duration*8000*8); err != nil {
		t.Fatalf("loop stopped after %d stereo bytes: %v", count, err)
	}
	var frame [8]byte
	if _, err := io.ReadFull(stream, frame[:]); err != nil {
		t.Fatalf("loop did not continue after %d seconds: %v", duration, err)
	}
}
