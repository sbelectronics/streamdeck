/* streamdeck.go
   (c) Scott M Baker, http://www.smbaker.com/

   This is a golang package for interfacing with the ElGato Streamdeck.
   See examples in the cmd/ directory for example usage.
*/

package streamdeck

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sbelectronics/streamdeck/pkg/util"
	"log"
)

type HandlerInterface interface {
	OnKeyDown(button *Button)
	OnKeyUp(button *Button)
	GetDefaultSettings(button *Button)
}

type Button struct {
	Action     string
	Context    string
	Column     int
	Row        int
	Handler    HandlerInterface
	Settings   map[string]string
	StreamDeck *StreamDeck
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

	ws             *websocket.Conn
	defaultHandler HandlerInterface
}

func (sd *StreamDeck) Init(
	port int,
	pluginUUID string,
	registerEvent string,
	info string,
	verbose bool,
	defaultHandler HandlerInterface) error {
	sd.Port = port
	sd.PluginUUID = pluginUUID
	sd.RegisterEvent = registerEvent
	sd.Info = info
	sd.Verbose = verbose
	sd.Buttons = make(map[string]*Button)
	sd.defaultHandler = defaultHandler

	err := sd.decodeInfo()
	if err != nil {
		return err
	}

	sd.Debugf("Port: %d\n", sd.Port)
	sd.Debugf("PluginUUID: %s", sd.PluginUUID)
	sd.Debugf("RegisterEvent: %s", sd.RegisterEvent)
	sd.Debugf("Info %s", sd.Info)
	sd.Debugf("Device Name %s", sd.DeviceName)
	sd.Debugf("Device Id: %s", sd.DeviceId)

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

func (sd *StreamDeck) Debugf(fmt string, args ...interface{}) {
	if sd.Verbose {
		log.Printf(fmt, args...)
	}
}

func (sd *StreamDeck) Infof(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}

func (sd *StreamDeck) Errorf(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}

func (sd *StreamDeck) Fatalf(fmt string, args ...interface{}) {
	// this will cause an os.exit()
	log.Fatalf(fmt, args...)
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
		sd.Errorf("Failed to unmarshal WillAppear")
		return
	}

	_, ok := sd.Buttons[willAppear.Context]
	if ok {

	} else {
		newButton := Button{Action: willAppear.Action,
			Context:    willAppear.Context,
			Column:     willAppear.Payload.Coordinates.Column,
			Row:        willAppear.Payload.Coordinates.Row,
			Settings:   make(map[string]string),
			Handler:    sd.defaultHandler,
			StreamDeck: sd}

		if newButton.Handler != nil {
			newButton.Handler.GetDefaultSettings(&newButton)
		}

		sd.Buttons[willAppear.Context] = &newButton
		sd.Infof("New button %v", newButton)

		// Try to get settings for this button
		sd.GetSettings(willAppear.Context)
	}
}

func (sd *StreamDeck) onKeyDown(message []byte) {
	var keyDown KeyDownNotification

	err := json.Unmarshal([]byte(message), &keyDown)
	if err != nil {
		sd.Errorf("Failed to unmarshal KeyDown")
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
		sd.Errorf("Failed to unmarshal KeyUp")
		return
	}

	button, ok := sd.Buttons[keyUp.Context]
	if ok {
		if button.Handler != nil {
			button.Handler.OnKeyUp(button)
		}
	}
}

func (sd *StreamDeck) onSendToPlugin(message []byte) {
	var f interface{}

	sd.Debugf("onSendToPlugin")

	err := json.Unmarshal(message, &f)
	if err != nil {
		sd.Errorf("Error unmarshaling SendToPlugin: %v", err)
		return
	}

	root := f.(map[string]interface{})

	context, err := util.StringFromJson(root, "context")
	if err != nil {
		sd.Errorf("failed to get onSendToPlugin.context: %v", err)
		return
	}

	payload, err := util.MapInterfaceFromJson(root, "payload")
	if err != nil {
		sd.Errorf("failed to get onSendToPlugin.context: %v", err)
		return
	}

	button, ok := sd.Buttons[context]
	if !ok {
		sd.Errorf("Failed to find SendToPlugin.button %v", context)
		return
	}

	for k, vInterface := range payload {
		v, ok := vInterface.(string)
		if !ok {
			sd.Errorf("Error converting SendtoPlugin.vInterface %v", vInterface)
			continue
		}
		sd.Infof("Received PropertyInspector k=%s, v=%s", k, v)

		button.Settings[k] = v
	}

	// Tell the streamdeck to store the settings persistently
	sd.SetSettings(context, button.Settings)
}

func (sd *StreamDeck) onDidReceiveSettings(message []byte) {
	var f interface{}

	sd.Debugf("onDidReceiveSettings")

	err := json.Unmarshal(message, &f)
	if err != nil {
		sd.Errorf("Error unmarshaling DidReceiveSettings: %v", err)
		return
	}

	root := f.(map[string]interface{})

	context, err := util.StringFromJson(root, "context")
	if err != nil {
		sd.Errorf("failed to get onDidReceiveSettings.context: %v", err)
		return
	}

	payload, err := util.MapInterfaceFromJson(root, "payload")
	if err != nil {
		sd.Errorf("failed to get onDidReceiveSettings.payload: %v", err)
		return
	}

	settings, err := util.MapInterfaceFromJson(payload, "settings")
	if err != nil {
		sd.Errorf("failed to get onDidReceiveSettings.setings: %v", err)
		return
	}

	button, ok := sd.Buttons[context]
	if !ok {
		sd.Errorf("Failed to find DidReceiveSettings.button %v", context)
		return
	}

	for k, vInterface := range settings {
		v, ok := vInterface.(string)
		if !ok {
			sd.Errorf("Error converting DidReceiveSettings.vInterface %v", vInterface)
			continue
		}
		sd.Infof("Received DidReceiveSettings k=%s, v=%s", k, v)

		button.Settings[k] = v
	}

	// Tell the streamdeck to store the settings persistently
	//    (did this in onPropertyInspectorDidAppear instead)
	//sd.SendToPropertyInspector(action, context, button.Settings)
}

func (sd *StreamDeck) onPropertyInspectorDidAppear(message []byte) {
	var pda PropertyInspectorDidAppearNotification

	err := json.Unmarshal([]byte(message), &pda)
	if err != nil {
		sd.Errorf("Failed to unmarshal PropertyInspectorDidAppear")
		return
	}

	button, ok := sd.Buttons[pda.Context]
	if ok {
		sd.SendToPropertyInspector(button.Action, button.Context, button.Settings)
	}
}

func (sd *StreamDeck) processWebsocketIncoming() {
	for {
		_, message, err := sd.ws.ReadMessage()
		if err != nil {
			sd.Errorf("Read error: %v", err)
			sd.Fatalf("FATAL: Websocket has failed -- exiting plugin")
			return // this will probably not be executed
		}

		sd.Debugf("Websocket Receive: %s", message)

		var header NotificationHeader
		err = json.Unmarshal([]byte(message), &header)
		if err != nil {
			sd.Errorf("Error unmarshaling header %v", err)
			continue
		}

		sd.Debugf("Websocket Event: %s", header.Event)

		switch header.Event {
		case "willAppear":
			sd.onWillAppear(message)
		case "willDisappear":
		case "keyDown":
			sd.onKeyDown(message)
		case "keyUp":
			sd.onKeyUp(message)
		case "sendToPlugin":
			sd.onSendToPlugin(message)
		case "didReceiveSettings":
			sd.onDidReceiveSettings(message)
		case "propertyInspectorDidAppear":
			sd.onPropertyInspectorDidAppear(message)
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

	sd.Debugf("Send JSON %s", string(b))

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

	sd.Debugf("Send JSON %s", string(b))

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

	sd.Debugf("Send JSON %s", string(b))

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

	sd.Debugf("Send JSON %s", string(b))

	err = sd.ws.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return err
	}

	return nil
}

func (sd *StreamDeck) GetSettings(context string) error {
	msg := GetSettingsMessage{"getSettings",
		context,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	sd.Debugf("Send JSON %s", string(b))

	err = sd.ws.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return err
	}

	return nil
}

func (sd *StreamDeck) SetSettings(context string, settings map[string]string) error {
	msg := make(map[string]interface{})
	msg["event"] = "setSettings"
	msg["context"] = context
	msg["payload"] = settings

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	sd.Debugf("Send JSON %s", string(b))

	err = sd.ws.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return err
	}

	return nil
}

func (sd *StreamDeck) SendToPropertyInspector(action string, context string, settings map[string]string) error {
	msg := make(map[string]interface{})
	msg["action"] = action
	msg["event"] = "sendToPropertyInspector"
	msg["context"] = context
	msg["payload"] = settings

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	sd.Debugf("Send JSON %s", string(b))

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

	sd.Debugf("Send JSON %s", string(b))

	err = sd.ws.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return err
	}

	return nil
}
