package bos

import "encoding/json"

type HubCloudMqttReport struct {
	Serial     string          `json:"serial"`
	DeviceType string          `json:"deviceType"`
	Data       json.RawMessage `json:"data"`
}
