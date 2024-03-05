package main

import (
	"encoding/json"
	conf "go-comm-mqtt/conf"
	"go-comm-mqtt/data/db"
	_ "go-comm-mqtt/logger"
	"go-comm-mqtt/mqtt/mqttcloud"
	"go-comm-mqtt/mqtt/mqttedge"
	"go-comm-mqtt/serve/productionline/web"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("go-comm-mqtt start")

	// 初始化配置
	conf.InitConfig()

	// 上研需求，获取mqtt配置信息
	//if conf.Conf.ServeName == "iotmqtt" {
	//	err := Iot_mqtt.GetIotMqttConfig(conf.Conf)
	//	if err != nil {
	//		logrus.Error("GetIotMqttConfig error:", err)
	//		return
	//	}
	//}

	// 打印配置信息
	configuration, _ := json.Marshal(conf.Conf)
	logrus.Infof("Using conf: %v", string(configuration))

	//初始化db
	db.Init()

	/** mqtt **/
	// 初始化mqtt edge
	mqttedge.MQTTEdgeManager = mqttedge.NewMQTTEdgeManager()
	mqttedge.MQTTEdgeManager.Init()

	// 初始化mqtt cloud
	mqttcloud.MQTTCloudManager = mqttcloud.NewMQTTCloudManager()
	mqttcloud.MQTTCloudManager.Init()

	// serve
	switch conf.Conf.ServeName {
	//case "cgxi":
	//	cgxi.CgxiDealModbus(config, mqttClient)
	//case "test":
	//	modbus.DealModbus(config, mqttClient)
	//case "xinjie":
	//	xinjie.XinJieDealModbus(config, mqttClient)
	//case "digitalTwin":
	//	digitaltwin.DigitalTwinDealModbus(config, mqttClient)
	//case "iotmqtt":
	//	Iot_mqtt.IotMqttImpl(config, mqttClient)
	case "productionline_web":
		web.ProductionLineWebDeal()
	default:
		logrus.Error("config error, no serve name")
	}

	//保持存活
	c := initSignal()
	handleSignal(c)

}

// initSignal register signals handler.
func initSignal() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	return c
}

// handleSignal fetch signal from chan then do exit or reload.
func handleSignal(c chan os.Signal) {
	// Block until a signal is received.
	for {
		s := <-c
		logrus.Infof("get a signal %s", s.String())
		switch s {
		case os.Interrupt:
			return
		case syscall.SIGHUP:
			// TODO reload
			//return
		default:
			return
		}
	}
}
