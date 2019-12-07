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
	"flag"
	"fmt"
	"github.com/sbelectronics/streamdeck/pkg/streamdeck"
	//"io/ioutil"
	"bytes"
	"github.com/BurntSushi/graphics-go/graphics"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"time"
)

type MyButtonHandler struct {
	color int
}

func (mbh *MyButtonHandler) OnKeyUp(*streamdeck.Button) {
	mbh.color++
	if mbh.color >= 3 {
		mbh.color = 0
	}

}
func (mbh *MyButtonHandler) OnKeyDown(*streamdeck.Button) {}

var (
	GlobalOptions struct {
		Port          int
		PluginUUID    string
		RegisterEvent string
		Info          string
		Verbose       bool

		DeviceName string
		DeviceId   string
	}

	forceVerbose = false

	color_presets []color.RGBA = []color.RGBA{color.RGBA{200, 0, 0, 255},
		color.RGBA{0, 200, 0, 255},
		color.RGBA{0, 0, 200, 255}}
)

func parseCmdline() {
	log.Printf("Loading command line")

	help := fmt.Sprintf("Port number for websocket")
	flag.IntVar(&(GlobalOptions.Port), "port", 0, help)

	help = fmt.Sprintf("Plugin UUID")
	flag.StringVar(&(GlobalOptions.PluginUUID), "pluginUUID", "", help)

	help = fmt.Sprintf("Register Event")
	flag.StringVar(&(GlobalOptions.RegisterEvent), "registerEvent", "", help)

	help = fmt.Sprintf("Info")
	flag.StringVar(&(GlobalOptions.Info), "info", "", help)

	help = fmt.Sprintf("Verbose mode")
	flag.BoolVar(&(GlobalOptions.Verbose), "v", false, help)

	flag.Parse()

	if forceVerbose {
		GlobalOptions.Verbose = true
	}
}

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
		forceVerbose = true
	}

	parseCmdline()

	sd := streamdeck.StreamDeck{}
	err := sd.Init(GlobalOptions.Port,
		GlobalOptions.PluginUUID,
		GlobalOptions.RegisterEvent,
		GlobalOptions.Info,
		GlobalOptions.Verbose)
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
			// If there's no handler, then attach a handler
			if button.Handler == nil {
				log.Printf("Set handler")
				button.Handler = &MyButtonHandler{}
			}

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
