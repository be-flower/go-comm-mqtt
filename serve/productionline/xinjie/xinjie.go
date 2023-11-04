package xinjie

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/simonvetter/modbus"
	"github.com/sirupsen/logrus"
	"go-comm-mqtt/common/constants"
	"go-comm-mqtt/common/utils"
	"go-comm-mqtt/config"
	"go-comm-mqtt/domains/bos"
	"go-comm-mqtt/domains/vos"
	m "go-comm-mqtt/modbus"
	"go-comm-mqtt/serve/gw"
	"time"
)

/*

 */

func XinJieDealModbus(config config.Config, mqttClient MQTT.Client) {
	modbusClient, err := m.TcpModbusClient(config)
	if err != nil {
		logrus.Error("modbusClient create error!")
	} else {
		//go gw.ReadGatewayModule(mqttClient, config)
		go XinJieReadTcpModbus(mqttClient, config, modbusClient)
	}
}

func XinJieReadTcpModbus(client MQTT.Client, config config.Config, modbusclient *modbus.ModbusClient) {
	logrus.Info("xinjie read tcpmodbus start")
	xinjieJxbDataVo := vos.XinJieJxbDataVo{
		JxbBo: bos.JxbBo{
			Serial:     "JXB001",
			DeviceType: "40001",
		},
		Data: vos.XinJieJxbVo{},
	}
	gwDateVo := vos.GwDataVo{
		GwBo: bos.GwBo{
			Serial:     "TZ001",
			DeviceType: "42001",
		},
		Data: vos.GwModuleVo{},
	}
	logrus.Infof("xinjieJxbDataVo: %+v", xinjieJxbDataVo)
	for {
		time.Sleep(time.Second * time.Duration(config.Tcpmodbus.Interval))
		for _, item := range config.Tcpmodbus.Devices {
			if item.Register == "holding" {
				var sendmsg string
				var sendmodulemsg string
				for _, read := range item.RegisterTable {
					switch read.Type {
					case "aubo-tcp":
						uint16Results, err := modbusclient.ReadRegisters(uint16(read.StartAddr), uint16(read.DataLen), modbus.HOLDING_REGISTER)
						results := utils.Uint16sToInt16s(uint16Results)
						if err != nil {
							logrus.Errorf("××××××××××read holding register xinjie-tcp error: %v", err)
						}
						if len(results) != 12 {
							logrus.Error("failed to read aubo-tcp data")
							continue
						}
						xinjieJxbDataVo.Data.X = fmt.Sprintf("%.2f", utils.TwoIntToFloat(int(results[0]), int(results[1])))
						xinjieJxbDataVo.Data.Y = fmt.Sprintf("%.2f", utils.TwoIntToFloat(int(results[2]), int(results[3])))
						xinjieJxbDataVo.Data.Z = fmt.Sprintf("%.2f", utils.TwoIntToFloat(int(results[4]), int(results[5])))
						xinjieJxbDataVo.Data.Rx = fmt.Sprintf("%.2f", utils.TwoIntToFloat(int(results[6]), int(results[7])))
						xinjieJxbDataVo.Data.Ry = fmt.Sprintf("%.2f", utils.TwoIntToFloat(int(results[8]), int(results[9])))
						xinjieJxbDataVo.Data.Rz = fmt.Sprintf("%.2f", utils.TwoIntToFloat(int(results[10]), int(results[11])))
					case "aubo-joint":
						uint16Results, err := modbusclient.ReadRegisters(uint16(read.StartAddr), uint16(read.DataLen), modbus.HOLDING_REGISTER)
						results := utils.Uint16sToInt16s(uint16Results)
						if err != nil {
							logrus.Errorf("××××××××××read holding register xinjie-joint error: %v", err)
						}
						if len(results) != 12 {
							logrus.Error("failed to read aubo-tcp data")
							continue
						}
						xinjieJxbDataVo.Data.Joint1 = fmt.Sprintf("%.2f", utils.TwoIntToFloat(int(results[0]), int(results[1])))
						xinjieJxbDataVo.Data.Joint2 = fmt.Sprintf("%.2f", utils.TwoIntToFloat(int(results[2]), int(results[3])))
						xinjieJxbDataVo.Data.Joint3 = fmt.Sprintf("%.2f", utils.TwoIntToFloat(int(results[4]), int(results[5])))
						xinjieJxbDataVo.Data.Joint4 = fmt.Sprintf("%.2f", utils.TwoIntToFloat(int(results[6]), int(results[7])))
						xinjieJxbDataVo.Data.Joint5 = fmt.Sprintf("%.2f", utils.TwoIntToFloat(int(results[8]), int(results[9])))
						xinjieJxbDataVo.Data.Joint6 = fmt.Sprintf("%.2f", utils.TwoIntToFloat(int(results[10]), int(results[11])))
					}
				}
				info, err := gw.GetGatewayModuleInfo()
				gwDateVo.Data = info
				if err != nil {
					logrus.Error("failed to get gateway module info")
					continue
				}
				//sendmsg = "{\"serial\":\"" + xinjieJxbDataVo.Serial + "\",\"deviceType\":\"" + xinjieJxbDataVo.DeviceType + "\",\"data\":\"{" +
				//	"\\\"x\\\":" + xinjieJxbDataVo.Data.X + ",\\\"y\\\":" + xinjieJxbDataVo.Data.Y + ",\\\"z\\\":" + xinjieJxbDataVo.Data.Z + ",\\\"rx\\\":" + xinjieJxbDataVo.Data.Rx + ",\\\"ry\\\":" + xinjieJxbDataVo.Data.Ry + ",\\\"rz\\\":" + xinjieJxbDataVo.Data.Rz + "}\"}"
				sendmsg = "{\"serial\":\"" + xinjieJxbDataVo.Serial + "\",\"deviceType\":\"" + xinjieJxbDataVo.DeviceType + "\",\"data\":\"{" +
					"\\\"x\\\":" + xinjieJxbDataVo.Data.X + ",\\\"y\\\":" + xinjieJxbDataVo.Data.Y + ",\\\"z\\\":" + xinjieJxbDataVo.Data.Z + ",\\\"rx\\\":" + xinjieJxbDataVo.Data.Rx + ",\\\"ry\\\":" + xinjieJxbDataVo.Data.Ry + ",\\\"rz\\\":" + xinjieJxbDataVo.Data.Rz + "," +
					"\\\"joint1\\\":" + xinjieJxbDataVo.Data.Joint1 + ",\\\"joint2\\\":" + xinjieJxbDataVo.Data.Joint2 + ",\\\"joint3\\\":" + xinjieJxbDataVo.Data.Joint3 + ",\\\"joint4\\\":" + xinjieJxbDataVo.Data.Joint4 + ",\\\"joint5\\\":" + xinjieJxbDataVo.Data.Joint5 + ",\\\"joint6\\\":" + xinjieJxbDataVo.Data.Joint6 + "}\"}"
				sendmodulemsg = "{\"serial\":\"" + gwDateVo.Serial + "\",\"deviceType\":\"" + gwDateVo.DeviceType + "\",\"data\":\"{" +
					"\\\"imei\\\":\\\"" + gwDateVo.Data.Imei + "\\\",\\\"imsi\\\":\\\"" + gwDateVo.Data.Imsi + "\\\",\\\"ip\\\":\\\"" + gwDateVo.Data.Ip + "\\\",\\\"rssi\\\":" + gwDateVo.Data.Rssi + ",\\\"sinr\\\":" + gwDateVo.Data.Sinr + ",\\\"rsrp\\\":" + gwDateVo.Data.Rsrp + "}\"}"
				point := client.Publish(item.Topic, 1, false, sendmsg)
				point.Wait()
				logrus.Infof("send message on topic: %s ; Message: \x1b[%dm%s\x1b[0m", item.Topic, constants.Cyan, sendmsg)
				module := client.Publish(item.Topic, 1, false, sendmodulemsg)
				module.Wait()
				logrus.Infof("send message on topic: %s ; Message: \x1b[%dm%s\x1b[0m", item.Topic, constants.Cyan, sendmodulemsg)
			}
		}
	}
}
