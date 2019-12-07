/* globaloptions.go
   (c) Scott M Baker, http://www.smbaker.com/
*/

package globaloptions

import (
	"flag"
	"fmt"
	"log"
)

var (
	Port          int
	PluginUUID    string
	RegisterEvent string
	Info          string
	Verbose       bool

	ForceVerbose = false
)

func ParseCmdline() {
	log.Printf("Loading command line")

	help := fmt.Sprintf("Port number for websocket")
	flag.IntVar(&(Port), "port", 0, help)

	help = fmt.Sprintf("Plugin UUID")
	flag.StringVar(&(PluginUUID), "pluginUUID", "", help)

	help = fmt.Sprintf("Register Event")
	flag.StringVar(&(RegisterEvent), "registerEvent", "", help)

	help = fmt.Sprintf("Info")
	flag.StringVar(&(Info), "info", "", help)

	help = fmt.Sprintf("Verbose mode")
	flag.BoolVar(&(Verbose), "v", false, help)

	flag.Parse()

	if ForceVerbose {
		Verbose = true
	}
}
