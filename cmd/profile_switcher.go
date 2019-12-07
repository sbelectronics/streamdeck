/* profile_switcher.go
   (c) Scott M Baker, http://www.smbaker.com/
*/

package main

import (
	"encoding/json"
	"fmt"
	"github.com/sbelectronics/streamdeck/pkg/globaloptions"
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
		globaloptions.ForceVerbose = true
	}

	globaloptions.ParseCmdline()

	patterns, err := loadPatterns()
	if err != nil {
		log.Fatalf("Error loading patterns: %v", err)
	}

	sd := streamdeck.StreamDeck{}
	err = sd.Init(globaloptions.Port,
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

	lastTitle := ""
	for {
		title, err := windows.GetTopmostWindowTitle()
		if err != nil {
			log.Printf("Error %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if lastTitle != title {
			if globaloptions.Verbose {
				log.Printf("Title: %s", title)
			}

			lastTitle = title

			matchedSomething := false
			for _, p := range patterns {
				if matched, _ := regexp.MatchString(p.Regex, title); matched {
					if globaloptions.Verbose {
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
