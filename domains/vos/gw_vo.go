package vos

import "go-comm-mqtt/domains/bos"

type GwDataVo struct {
	bos.GwBo
	Data GwModuleVo `json:"data"`
}

// GwModuleVo 网关模组信息
type GwModuleVo struct {
	Imei string `json:"imei"`
	Imsi string `json:"imsi"`
	Ip   string `json:"ip"`
	Rssi string `json:"rssi"`
	Sinr string `json:"sinr"`
	Rsrp string `json:"rsrp"`
}
