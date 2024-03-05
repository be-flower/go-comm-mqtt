package bos

type IotMqttGetMqttInfoReq struct {
	DeviceId string `json:"deviceId"`
	Password string `json:"password"`
}

type IotMqttGetMqttInfoResp struct {
	Code    string                       `json:"code"`
	Message string                       `json:"message"`
	Result  IotMqttGetMqttInfoRespResult `json:"result"`
}

type IotMqttGetMqttInfoRespResult struct {
	ClientId string `json:"clientId"`
	UserName string `json:"userName"`
	Password string `json:"password"`
	MqttHost string `json:"mqttHost"`
	MqttPort int    `json:"mqttPort"`
}
