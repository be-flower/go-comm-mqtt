package main

import (
	"github.com/sirupsen/logrus"
	conf "go-comm-mqtt/config"
	//_ "go-comm-mqtt/data/db"
	_ "go-comm-mqtt/logger"
	"go-comm-mqtt/modbus"
	"go-comm-mqtt/mqtt"
	"go-comm-mqtt/serve/cgxi"
	"go-comm-mqtt/serve/productionline/digitaltwin"
	"go-comm-mqtt/serve/productionline/xinjie"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	catalog, err := os.Getwd()
	if err != nil {
		logrus.Error("get catalog error: %v", err)
	}
	logrus.Infof("catalog: %v", catalog)
	logrus.Info("go-comm-mqtt start")
	config := conf.GetConfig()
	logrus.Info("Using config: %+v\n", config)
	// mqtt
	mqttClient, quit := mqtt.ConnMqtt(config)
	// modbus
	switch config.ServeName {
	case "cgxi":
		cgxi.CgxiDealModbus(config, mqttClient)
	case "test":
		modbus.DealModbus(config, mqttClient)
	case "xinjie":
		xinjie.XinJieDealModbus(config, mqttClient)
	case "digitalTwin":
		digitaltwin.DigitalTwinDealModbus(config, mqttClient)
	default:
		logrus.Error("no config")
	}

	// safe quit
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	cleanup := make(chan bool)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range signalChan {
			quit <- true
			go func() {
				go func() {
					mqttClient.Disconnect(250)
				}()
				time.Sleep(260 * time.Millisecond)
				cleanup <- true
			}()
			<-cleanup
			logrus.Info("safe quit")
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
