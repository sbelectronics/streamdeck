/* http.go
   (c) Scott M Baker, http://www.smbaker.com/

   This is a streamdeck plugin that will perform an HTTP operation
   against a given URL. You can specify the type of HTTP operation
   to perform (GET, POST, PUT, PATCH, DELETE) and for operations
   that send data (POST, PUT, etc) you can specify what data to
   send.
*/

package main

import (
	"github.com/sbelectronics/streamdeck/pkg/globaloptions"
	"github.com/sbelectronics/streamdeck/pkg/streamdeck"
	"github.com/sbelectronics/streamdeck/pkg/util"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type MyButtonHandler struct {
}

// On every keyup, change the color
func (mbh *MyButtonHandler) OnKeyUp(button *streamdeck.Button) {
	url, ok := button.Settings["url"]
	if !ok || (url == "") {
		button.StreamDeck.Errorf("Error - no url defined for button")
		return
	}

	url, err := util.SanitizeUrl(url)
	if err != nil {
		button.StreamDeck.Errorf("Error sanitizing url: %v", err)
	}
	button.StreamDeck.Debugf("%s", url)

	var resp *http.Response
	var req *http.Request

	operation := util.StringMapGetDefault(button.Settings, "operation", "GET")
	data := util.StringMapGetDefault(button.Settings, "data", "")
	mimetype := util.StringMapGetDefault(button.Settings, "mimetype", "text/plain")

	switch operation {
	case "POST":
		resp, err = http.Post(url, "text/plain", strings.NewReader(data))
	case "PUT":
		req, err = http.NewRequest("PUT", url, strings.NewReader(data))
		if err == nil {
			req.Header.Set("Content-Type", mimetype)
			resp, err = http.DefaultClient.Do(req)
		}
	case "PATCH":
		req, err = http.NewRequest("PATCH", url, strings.NewReader(data))
		if err == nil {
			req.Header.Set("Content-Type", mimetype)
			resp, err = http.DefaultClient.Do(req)
		}
	case "DELETE":
		req, err = http.NewRequest("DELETE", url, strings.NewReader(data))
		if err == nil {
			req.Header.Set("Content-Type", mimetype)
			resp, err = http.DefaultClient.Do(req)
		}
	default:
		resp, err = http.Get(url)
	}

	if err != nil {
		button.StreamDeck.Errorf("Error while performing operation %s on url %s: %v", operation, url, err)
	} else {
		button.StreamDeck.Debugf("Http response %v", resp)
	}

	// suppress error
	_ = resp
}
func (mbh *MyButtonHandler) OnKeyDown(*streamdeck.Button)          {}
func (mbh *MyButtonHandler) GetDefaultSettings(*streamdeck.Button) {}

func main() {
	// Create the file c:\junk\demo_plugin.log and we will append
	// log messages to it.
	if _, err := os.Stat("c:\\junk\\http_plugin.log"); err == nil {
		logf, err := os.OpenFile("c:\\junk\\http_plugin.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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

	for {
		time.Sleep(1000 * time.Millisecond)
	}
}
