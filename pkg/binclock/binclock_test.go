package binclock

import (
	"image/png"
	"os"
	"testing"
	"time"
)

func TestBinClock(t *testing.T) {
	bc := &BinClock{LitDotColor: DEFAULT_LIT_COLOR,
		UnlitDotColor: DEFAULT_UNLIT_COLOR,
		BackColor:     DEFAULT_BACK_COLOR}
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
