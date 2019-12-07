/* profile_switcher.go
   (c) Scott M Baker, http://www.smbaker.com/
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sbelectronics/streamdeck/pkg/platform/windows"
	"github.com/sbelectronics/streamdeck/pkg/streamdeck"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"
)

type Pattern struct {
	Regex   string
	Profile string
}

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

func loadPatterns() ([]Pattern, error) {
	var patterns []Pattern

	jsonFile, err := os.Open("patterns.json")
	if err != nil {
		return nil, fmt.Errorf("Failed to open patterns.json: %v", err)
	}
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to open patterns.json: %v", err)
	}
	err = json.Unmarshal(jsonData, &patterns)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal patterns.json: %v", err)
	}

	return patterns, nil
}

func main() {
	// Create the file c:\junk\profile_switcher_plugin.log and we will append
	// log messages to it.
	if _, err := os.Stat("c:\\junk\\profile_switcher_plugin.log"); err == nil {
		logf, err := os.OpenFile("c:\\junk\\profile_switcher_plugin.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer logf.Close()
		log.SetOutput(logf)
		forceVerbose = true
	}

	parseCmdline()

	patterns, err := loadPatterns()
	if err != nil {
		log.Fatalf("Error loading patterns: %v", err)
	}

	sd := streamdeck.StreamDeck{}
	err = sd.Init(GlobalOptions.Port,
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

	lastTitle := ""
	for {
		title, err := windows.GetTopmostWindowTitle()
		if err != nil {
			log.Printf("Error %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if lastTitle != title {
			if GlobalOptions.Verbose {
				log.Printf("Title: %s", title)
			}

			lastTitle = title

			matchedSomething := false
			for _, p := range patterns {
				if matched, _ := regexp.MatchString(p.Regex, title); matched {
					if GlobalOptions.Verbose {
						log.Printf("Matched %v", p)
					}
					sd.SwitchProfileIfChanged(p.Profile)
					matchedSomething = true
				}
			}
			if !matchedSomething {
				sd.SwitchProfileIfChanged("") // this will switch to the default
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}
