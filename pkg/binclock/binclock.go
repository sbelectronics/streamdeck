package binclock

import (
	"github.com/fogleman/gg"
	"image"
	"time"
)

const (
	dotSize = 3
	width   = 72
	height  = 72
	xOffset = 8
	yOffset = 22
	vOffset = 8
	hOffset = 8
	space   = 24

	DEFAULT_LIT_COLOR   = "#0000C8"
	DEFAULT_UNLIT_COLOR = "#505050"
	DEFAULT_BACK_COLOR  = "#000000"
)

type BinClock struct {
	Img           *gg.Context
	LitDotColor   string
	UnlitDotColor string
	BackColor     string
}

func (bc *BinClock) Create() {
	bc.Img = gg.NewContext(width, height)
}

func (bc *BinClock) DrawDot(x int, y int, lit bool) {
	bc.Img.Push()
	bc.Img.DrawCircle(float64(x), float64(height-y), float64(dotSize))
	if lit {
		bc.Img.SetHexColor(bc.LitDotColor)
	} else {
		bc.Img.SetHexColor(bc.UnlitDotColor)
	}
	bc.Img.Fill()
	bc.Img.Pop()
}

func (bc *BinClock) DrawNibble(x int, y int, bits int, v int) {
	for i := 0; i < bits; i++ {
		bc.DrawDot(x, y, (v&1) == 1)
		v = v >> 1
		y = y + vOffset
	}
}

func (bc *BinClock) DrawDigit(x int, y int, bits int, v int) {
	bc.DrawNibble(x, y, bits-4, v/10)
	bc.DrawNibble(x+hOffset, y, 4, v%10)
}

func (bc *BinClock) DrawTime(t time.Time) {
	bc.Img.Push()
	bc.Img.DrawRectangle(0, 0, width, height)
	bc.Img.SetHexColor(bc.BackColor)
	bc.Img.Fill()
	bc.Img.Pop()

	bc.DrawDigit(xOffset, yOffset, 6, t.Hour())
	bc.DrawDigit(xOffset+space, yOffset, 7, t.Minute())
	bc.DrawDigit(xOffset+space*2, yOffset, 7, t.Second())
}

func (bc *BinClock) Image() image.Image {
	return bc.Img.Image()
}
