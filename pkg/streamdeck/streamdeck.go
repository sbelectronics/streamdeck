/* streamdeck.go
   (c) Scott M Baker, http://www.smbaker.com/
*/

package streamdeck

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type HandlerInterface interface {
	OnKeyDown(button *Button)
	OnKeyUp(button *Button)
}

type Button struct {
	Context string
	Column  int
	Row     int
	Handler HandlerInterface
}

type StreamDeck struct {
	Port          int
	PluginUUID    string
	RegisterEvent string
	Info          string
	Verbose       bool

	DeviceName string
	DeviceId   string

	lastProfile string

	Buttons map[string]*Button

	ws *websocket.Conn
}

func (sd *StreamDeck) Init(
	port int,
	pluginUUID string,
	registerEvent string,
	info string,
	verbose bool) error {
	sd.Port = port
	sd.PluginUUID = pluginUUID
	sd.RegisterEvent = registerEvent
	sd.Info = info
	sd.Verbose = verbose
	sd.Buttons = make(map[string]*Button)

	err := sd.decodeInfo()
	if err != nil {
		return err
	}

	if sd.Verbose {
		log.Printf("Port: %d\n", sd.Port)
		log.Printf("PluginUUID: %s", sd.PluginUUID)
		log.Printf("RegisterEvent: %s", sd.RegisterEvent)
		log.Printf("Info %s", sd.Info)
		log.Printf("Device Name %s", sd.DeviceName)
		log.Printf("Device Id: %s", sd.DeviceId)
	}

	return nil
}

func (sd *StreamDeck) Start() error {
	err := sd.initWebsocket()
	if err != nil {
		return err
	}
	go sd.processWebsocketIncoming()

	err = sd.register()
	if err != nil {
		return err
	}

	return nil
}

// Extract the device id from the Info JSON
func (sd *StreamDeck) decodeInfo() error {
	var f interface{}

	if sd.Info != "" {
		err := json.Unmarshal([]byte(sd.Info), &f)
		if err != nil {
			return err
		}

		root := f.(map[string]interface{})
		devices := root["devices"].([]interface{})
		device := devices[0].(map[string]interface{})

		sd.DeviceName = device["name"].(string)
		sd.DeviceId = device["id"].(string)
	}

	return nil
}

func (sd *StreamDeck) initWebsocket() error {
	url := fmt.Sprintf("ws://127.0.0.1:%d", sd.Port)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	sd.ws = c
	return nil
}

func (sd *StreamDeck) onWillAppear(message []byte) {
	var willAppear WillAppearNotification

	err := json.Unmarshal([]byte(message), &willAppear)
	if err != nil {
		log.Printf("Failed to unmarshal WillAppear")
		return
	}

	_, ok := sd.Buttons[willAppear.Context]
	if ok {

	} else {
		newButton := Button{Context: willAppear.Context,
			Column: willAppear.Payload.Coordinates.Column,
			Row:    willAppear.Payload.Coordinates.Row}
		sd.Buttons[willAppear.Context] = &newButton
		log.Printf("New button %v", newButton)
	}
}

func (sd *StreamDeck) onKeyDown(message []byte) {
	var keyDown KeyDownNotification

	err := json.Unmarshal([]byte(message), &keyDown)
	if err != nil {
		log.Printf("Failed to unmarshal KeyDown")
		return
	}

	button, ok := sd.Buttons[keyDown.Context]
	if ok {
		if button.Handler != nil {
			button.Handler.OnKeyDown(button)
		}
	}
}

func (sd *StreamDeck) onKeyUp(message []byte) {
	var keyUp KeyUpNotification

	err := json.Unmarshal([]byte(message), &keyUp)
	if err != nil {
		log.Printf("Failed to unmarshal KeyUp")
		return
	}

	button, ok := sd.Buttons[keyUp.Context]
	if ok {
		if button.Handler != nil {
			button.Handler.OnKeyUp(button)
		}
	}
}

func (sd *StreamDeck) processWebsocketIncoming() {
	for {
		_, message, err := sd.ws.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			return
		}

		log.Printf("Receive: %s", message)

		var header NotificationHeader
		err = json.Unmarshal([]byte(message), &header)
		if err != nil {
			log.Printf("Error unmarshaling header %v", err)
			continue
		}
		log.Printf("event: %s", header.Event)
		switch header.Event {
		case "willAppear":
			sd.onWillAppear(message)
		case "willDisappear":
		case "keyDown":
			sd.onKeyDown(message)
		case "keyUp":
			sd.onKeyUp(message)
		}
	}
}

func (sd *StreamDeck) SwitchProfile(profile string) error {
	msg := SwitchProfileMessage{"switchToProfile",
		sd.PluginUUID,
		sd.DeviceId,
		ProfilePayload{profile}}

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if sd.Verbose {
		log.Printf("Send JSON %s", string(b))
	}

	err = sd.ws.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return err
	}

	return nil
}

func (sd *StreamDeck) SwitchProfileIfChanged(profile string) error {
	if sd.lastProfile != profile {
		sd.lastProfile = profile
		return sd.SwitchProfile(profile)
	}
	return nil
}

// untested
func (sd *StreamDeck) SetImage(context string, image []byte, mimeType string, target int) error {
	sEnc := b64.StdEncoding.EncodeToString(image)
	imageStr := "data:" + mimeType + ";base64," + sEnc

	msg := SetImageMessage{"setImage",
		context, //sd.PluginUUID,
		SetImagePayload{imageStr, target}}

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if sd.Verbose {
		log.Printf("Send JSON %s", string(b))
	}

	err = sd.ws.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return err
	}

	return nil
}

// untested
func (sd *StreamDeck) SetTitle(context string, title string, target int) error {
	msg := SetTitleMessage{"setTitle",
		context, //sd.PluginUUID,
		SetTitlePayload{title, target}}

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if sd.Verbose {
		log.Printf("Send JSON %s", string(b))
	}

	err = sd.ws.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return err
	}

	return nil
}

func (sd *StreamDeck) ShowAlert(context string) error {
	msg := ShowAlertMessage{"showAlert",
		context, //sd.PluginUUID,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if sd.Verbose {
		log.Printf("Send JSON %s", string(b))
	}

	err = sd.ws.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return err
	}

	return nil
}

func (sd *StreamDeck) register() error {
	msg := RegisterMessage{sd.RegisterEvent, sd.PluginUUID}

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if sd.Verbose {
		log.Printf("Send JSON %s", string(b))
	}

	err = sd.ws.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return err
	}

	return nil
}
