package streamdeck

type NotificationHeader struct {
	Action  string `json:"action"`
	Event   string `json:"event"`
	Context string `json:"context"`
	Device  string `json:"device"`
}

type Coordinates struct {
	Column int `json:"column"`
	Row    int `json:"row"`
}

type WillAppearPayload struct {
	Coordinates     Coordinates `json:"coordinates"`
	State           int         `json:"state"`
	IsInMultiAction bool        `json:"isInMultiAction"`
}

type WillAppearNotification struct {
	Action  string            `json:"action"`
	Event   string            `json:"event"`
	Context string            `json:"context"`
	Device  string            `json:"device"`
	Payload WillAppearPayload `json:"payload"`
}

type WillDisappearNotification struct {
	Action  string            `json:"action"`
	Event   string            `json:"event"`
	Context string            `json:"context"`
	Device  string            `json:"device"`
	Payload WillAppearPayload `json:"payload"`
}

type KeyPayload struct {
	Coordinates      Coordinates `json:"coordinates"`
	State            int         `json:"state"`
	UserDesiredState bool        `json:"userDesiredState"`
	IsInMultiAction  bool        `json:"isInMultiAction"`
}

type KeyDownNotification struct {
	Action  string     `json:"action"`
	Event   string     `json:"event"`
	Context string     `json:"context"`
	Device  string     `json:"device"`
	Payload KeyPayload `json:"payload"`
}

type KeyUpNotification struct {
	Action  string     `json:"action"`
	Event   string     `json:"event"`
	Context string     `json:"context"`
	Device  string     `json:"device"`
	Payload KeyPayload `json:"payload"`
}
