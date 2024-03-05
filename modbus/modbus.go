package modbus

import (
	config "go-comm-mqtt/conf"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

func DealModbus(config config.Config, mqttClient MQTT.Client) {
	modbusClient, err := TcpModbusClient()
	if err != nil {
		logrus.Error("modbusClient create error!")
	} else {
		go ReadTcpModbus(mqttClient, config, modbusClient)
		go WriteTcpModbus(mqttClient, config, modbusClient)
	}
}
