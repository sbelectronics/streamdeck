/* messages.go
   (c) Scott M Baker, http://www.smbaker.com/

   JSON message structures, sent from Plugin to Streamdeck
*/

package streamdeck

const (
	TARGET_BOTH     = 0
	TARGET_HARDWARE = 1
	TARGET_SOFTWARE = 2

	TYPE_JPG = "image/jpg"
	TYPE_PNG = "image/png"
	TYPE_BMP = "image/bmp"
)

type RegisterMessage struct {
	Event string `json:"event"`
	Uuid  string `json:"uuid"`
}

type ProfilePayload struct {
	Profile string `json:"profile"`
}

type SwitchProfileMessage struct {
	Event   string         `json:"event"`
	Context string         `json:"context"`
	Device  string         `json:"device"`
	Payload ProfilePayload `json:"payload"`
}

type SetImagePayload struct {
	Image  string `json:"image"`
	Target int    `json:"target"`
}

type SetImageMessage struct {
	Event   string          `json:"event"`
	Context string          `json:"context"`
	Payload SetImagePayload `json:"payload"`
}

type SetTitlePayload struct {
	Title  string `json:"title"`
	Target int    `json:"target"`
}

type SetTitleMessage struct {
	Event   string          `json:"event"`
	Context string          `json:"context"`
	Payload SetTitlePayload `json:"payload"`
}

type ShowAlertMessage struct {
	Event   string `json:"event"`
	Context string `json:"context"`
}

type ShowOkMessage struct {
	Event   string `json:"event"`
	Context string `json:"context"`
}

type GetSettingsMessage struct {
	Event   string `json:"event"`
	Context string `json:"context"`
}
