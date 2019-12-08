/* demo.go
   (c) Scott M Baker, http://www.smbaker.com/

   NOTE: Make sure to leave the title blank when adding the button

   This is a simple demo that will show a rotating square with a counter
   that increments in the middle of it. Pressing the button will change
   the color red -> Green -> Blue, and back to red again. If the counter
   doesn't count, then you didn't heed the above note about making sure
   the title is blank.
*/

package main

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/graphics-go/graphics"
	"github.com/sbelectronics/streamdeck/pkg/globaloptions"
	"github.com/sbelectronics/streamdeck/pkg/streamdeck"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"time"
)

var (
	color_presets []color.RGBA = []color.RGBA{color.RGBA{200, 0, 0, 255},
		color.RGBA{0, 200, 0, 255},
		color.RGBA{0, 0, 200, 255}}
)

type MyButtonHandler struct {
	color int
}

// On every keyup, change the color
func (mbh *MyButtonHandler) OnKeyUp(*streamdeck.Button) {
	mbh.color++
	if mbh.color >= 3 {
		mbh.color = 0
	}

}
func (mbh *MyButtonHandler) OnKeyDown(*streamdeck.Button)          {}
func (mbh *MyButtonHandler) GetDefaultSettings(*streamdeck.Button) {}

func main() {
	// Create the file c:\junk\demo_plugin.log and we will append
	// log messages to it.
	if _, err := os.Stat("c:\\junk\\demo_plugin.log"); err == nil {
		logf, err := os.OpenFile("c:\\junk\\demo_plugin.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer logf.Close()
		log.SetOutput(logf)
		globaloptions.ForceVerbose = true
	}

	globaloptions.ParseCmdline()

	sd := streamdeck.StreamDeck{}
	err := sd.Init(globaloptions.Port,
		globaloptions.PluginUUID,
		globaloptions.RegisterEvent,
		globaloptions.Info,
		globaloptions.Verbose,
		&MyButtonHandler{})
	if err != nil {
		log.Fatalf("Error initializing streamdeck plugin %v", err)
	}

	err = sd.Start()
	if err != nil {
		log.Fatalf("Error starting streamdeck plugin %v", err)
	}

	counter := 0
	angle := math.Pi * 0.01
	for {
		for context, button := range sd.Buttons {
			// Draw a colored box and rotate it
			thisColor := color_presets[button.Handler.(*MyButtonHandler).color]
			img := image.NewRGBA(image.Rect(0, 0, 72, 72))
			draw.Draw(img, image.Rect(12, 12, 60, 60), &image.Uniform{thisColor}, image.ZP, draw.Src)
			srcDim := img.Bounds()
			rotatedImg := image.NewRGBA(image.Rect(0, 0, srcDim.Dy(), srcDim.Dx()))
			graphics.Rotate(rotatedImg, img, &graphics.RotateOptions{angle})
			angle += math.Pi * 0.01

			// Encode the image into a png and set it on the StreamDeck
			buf := new(bytes.Buffer)
			err := png.Encode(buf, rotatedImg)
			if err != nil {
				log.Printf("Error encoding png %v", err)
			}
			err = sd.SetImage(context, buf.Bytes(), streamdeck.TYPE_PNG, streamdeck.TARGET_BOTH)
			if err != nil {
				log.Printf("Error setimage: %v", err)
			}

			// Set the title to the number
			err = sd.SetTitle(context, fmt.Sprintf("%d", counter%1000), streamdeck.TARGET_BOTH)
			if err != nil {
				log.Printf("Error settitle: %v", err)
			}
		}
		counter += 1

		time.Sleep(100 * time.Millisecond)
	}
}
