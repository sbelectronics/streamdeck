package binclock

import (
	//"github.com/stretchr/testify/assert"
	"image/color"
	"image/png"
	"os"
	"testing"
	"time"
)

func TestBinClock(t *testing.T) {
	bc := &BinClock{LitDotColor: color.RGBA{0, 0, 200, 255},
		UnlitDotColor: color.RGBA{80, 80, 80, 255},
		BackColor:     color.RGBA{0, 0, 0, 255}}
	bc.Create()

	ts := time.Now()

	bc.DrawTime(ts)

	f, err := os.Create("binclock_test.png")
	if err != nil {
		t.Fatalf("Failed to create binclock_test.png: %v", err)
	}
	defer f.Close()

	err = png.Encode(f, bc.Image())
	if err != nil {
		t.Fatalf("Failed to encode png: %v", err)
	}
}
