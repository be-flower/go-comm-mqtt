package config

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type MqttInfo struct {
	Host     string
	Port     int
	ClientId string
	UserName string
	PassWord string
	SubList  []string
	PubList  []string
	Qos      int
}

type RegisterTable struct {
	StartAddr int
	DataLen   int
	Type      string
	Name      string
}

type Device struct {
	Register      string
	Topic         string
	RegisterTable []RegisterTable
}

type TcpModbus struct {
	Enable   bool
	Host     string
	Port     int
	SlaveID  int
	Interval int
	Devices  []Device
}

type RtuModbus struct {
	Enable   bool
	Device   string
	BaudRate int
	DataBits int
	Parity   string
	StopBits int
	SlaveID  int
	Interval int
	Devices  []Device
}

type Config struct {
	ServeName string
	Mqttinfo  MqttInfo
	Tcpmodbus TcpModbus
	Rtumodbus RtuModbus
}

func GetConfig() Config {

	config := Config{}
	config.Mqttinfo.Host = "127.0.0.1"
	config.Mqttinfo.Port = 1883
	config.Mqttinfo.Qos = 0

	//modbusTCP配置
	config.Tcpmodbus.Enable = false
	config.Tcpmodbus.Host = "127.0.0.1"
	config.Tcpmodbus.Port = 502
	config.Tcpmodbus.SlaveID = 1
	config.Tcpmodbus.Interval = 3

	//modbusRTU配置
	config.Rtumodbus.Enable = false
	config.Rtumodbus.BaudRate = 9600
	config.Rtumodbus.DataBits = 8
	config.Rtumodbus.Parity = "N"
	config.Rtumodbus.StopBits = 1
	config.Rtumodbus.SlaveID = 1
	config.Rtumodbus.Interval = 3

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Fatal("Config file not found")
		} else {
			fmt.Println(err.Error())
		}
		return config
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		logrus.Fatalf("unable to decode into struct, %v", err)
	}

	return config
}
