/* binclock.go
   (c) Scott M Baker, http://www.smbaker.com/
*/

package main

import (
	"bytes"
	"github.com/sbelectronics/streamdeck/pkg/binclock"
	"github.com/sbelectronics/streamdeck/pkg/globaloptions"
	"github.com/sbelectronics/streamdeck/pkg/streamdeck"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"
)

func main() {
	// Create the file c:\junk\binclock_plugin.log and we will append
	// log messages to it.
	if _, err := os.Stat("c:\\junk\\binclock_plugin.log"); err == nil {
		logf, err := os.OpenFile("c:\\junk\\binclock_plugin.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
		globaloptions.Verbose)
	if err != nil {
		log.Fatalf("Error initializing streamdeck plugin %v", err)
	}

	err = sd.Start()
	if err != nil {
		log.Fatalf("Error starting streamdeck plugin %v", err)
	}

	bc := &binclock.BinClock{LitDotColor: color.RGBA{0, 0, 200, 255},
		UnlitDotColor: color.RGBA{80, 80, 80, 255},
		BackColor:     color.RGBA{0, 0, 0, 255}}
	bc.Create()

	for {
		for context, _ := range sd.Buttons {
			bc.DrawTime(time.Now())

			// Encode the image into a png and set it on the StreamDeck
			buf := new(bytes.Buffer)
			err := png.Encode(buf, bc.Image())
			if err != nil {
				log.Printf("Error encoding png %v", err)
			}
			err = sd.SetImage(context, buf.Bytes(), streamdeck.TYPE_PNG, streamdeck.TARGET_BOTH)
			if err != nil {
				log.Printf("Error setimage: %v", err)
			}
		}

		time.Sleep(1000 * time.Millisecond)
	}
}
