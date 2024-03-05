package model

import "time"

/*MQTT推送的请求*/
type MQTTCommandRequest struct {
	/*租户，必选*/
	Tenant string `json:"tenant_id"`
	/*门店，可选*/
	StoreID string `json:"store_id,omitempty"`
	/*主题*/
	Topic string `json:"topic,omitemtpy"`
	/*设备类型*/
	DeviceType string `json:"device_type"`
	/*设备，可选*/
	DevId string `json:"dev_id"`
	/*事务ID，必选*/
	Tid string `json:"tid"`
	/*时间，必选,s*/
	Time int64 `json:"time"`
	/*星期几*/
	Weekday time.Weekday `json:"weekday"`
	// 同步事件
	Event string `json:"event"`
	/*时间，YYYY-MM-DD HH:MM:SS*/
	TimeStr string `json:"time_str"`
	/*数据签名，可选，hash(‘sha256’, data+time_str+salt)*/
	Signature string `json:"signature"`
	/*需要同步的数据，必选*/
	Data string `json:"data"`
	// 序列号
	Serial string `json:"serial"`
}
