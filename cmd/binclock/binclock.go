/* binclock.go
   (c) Scott M Baker, http://www.smbaker.com/
*/

package main

import (
	"bytes"
	"github.com/sbelectronics/streamdeck/pkg/binclock"
	"github.com/sbelectronics/streamdeck/pkg/globaloptions"
	"github.com/sbelectronics/streamdeck/pkg/streamdeck"
	"github.com/sbelectronics/streamdeck/pkg/util"
	"image/png"
	"log"
	"os"
	"time"
)

type MyButtonHandler struct {
	color int
}

// On every keyup, change the color
func (mbh *MyButtonHandler) OnKeyUp(*streamdeck.Button)   {}
func (mbh *MyButtonHandler) OnKeyDown(*streamdeck.Button) {}
func (mbh *MyButtonHandler) GetDefaultSettings(button *streamdeck.Button) {
	button.Settings["colorlit"] = binclock.DEFAULT_LIT_COLOR
	button.Settings["colorunlit"] = binclock.DEFAULT_UNLIT_COLOR
	button.Settings["colorback"] = binclock.DEFAULT_BACK_COLOR
}

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
		globaloptions.Verbose,
		&MyButtonHandler{})
	if err != nil {
		log.Fatalf("Error initializing streamdeck plugin %v", err)
	}

	err = sd.Start()
	if err != nil {
		log.Fatalf("Error starting streamdeck plugin %v", err)
	}

	bc := &binclock.BinClock{LitDotColor: binclock.DEFAULT_LIT_COLOR,
		UnlitDotColor: binclock.DEFAULT_UNLIT_COLOR,
		BackColor:     binclock.DEFAULT_BACK_COLOR}
	bc.Create()

	lastSecond := -1
	for {
		// only update if the second has changed
		tNow := time.Now()
		if tNow.Second() == lastSecond {
			continue
		}
		lastSecond = tNow.Second()

		for context, button := range sd.Buttons {
			// override the colors with what might have come from the property inspector
			bc.LitDotColor = util.StringMapGetDefault(button.Settings, "colorlit", binclock.DEFAULT_LIT_COLOR)
			bc.UnlitDotColor = util.StringMapGetDefault(button.Settings, "colorunlit", binclock.DEFAULT_UNLIT_COLOR)
			bc.BackColor = util.StringMapGetDefault(button.Settings, "colorback", binclock.DEFAULT_BACK_COLOR)

			bc.DrawTime(tNow)

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

		time.Sleep(100 * time.Millisecond)
	}
}
