package modbus

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	config "go-comm-mqtt/config"
)

func DealModbus(config config.Config, mqttClient MQTT.Client) {
	modbusClient, err := TcpModbusClient(config)
	if err != nil {
		logrus.Error("modbusClient create error!")
	} else {
		go ReadTcpModbus(mqttClient, config, modbusClient)
		go WriteTcpModbus(mqttClient, config, modbusClient)
	}
}
