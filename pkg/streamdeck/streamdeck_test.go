package streamdeck

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeviceInfo(t *testing.T) {
	info := `{
		"application": {
		  "language": "en",
		  "platform": "mac",
		  "version": "4.1.0"
		},
		"plugin": {
		  "version": "1.1"
		},
		"devicePixelRatio": 2,
		"devices": [
		  {
			"id": "55F16B35884A859CCE4FFA1FC8D3DE5B",
			"name": "Device Name",
			"size": {
			  "columns": 5,
			  "rows": 3
			},
			"type": 0
		  },
		  {
			"id": "B8F04425B95855CF417199BCB97CD2BB",
			"name": "Another Device",
			"size": {
			  "columns": 3,
			  "rows": 2
			},
			"type": 1
		  }
		]
		}`

	sd := StreamDeck{}
	sd.Init(0, "uuid", "reg", info, false, nil)

	err := sd.decodeInfo()

	assert.Nil(t, err)
	assert.Equal(t, "55F16B35884A859CCE4FFA1FC8D3DE5B", sd.DeviceId)
	assert.Equal(t, "Device Name", sd.DeviceName)
}
