package modbus

import (
	"encoding/json"
	"errors"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	modbus "github.com/simonvetter/modbus"
	"github.com/sirupsen/logrus"
	"go-comm-mqtt/common/constants"
	"go-comm-mqtt/config"
	"go-comm-mqtt/domains/bos"
	"strconv"
	"time"
)

func TcpModbusClient(config config.Config) (*modbus.ModbusClient, error) {
	logrus.Info("TCPModbus start")
	// for a TCP endpoint
	// (see examples/tls_client.go for TLS usage and options)
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     "tcp://" + config.Tcpmodbus.Host + ":" + strconv.Itoa(config.Tcpmodbus.Port),
		Timeout: 1 * time.Second,
	})
	if err != nil {
		// error out if client creation failed
		logrus.Error("tcpmodbus connect error!")
		return nil, errors.New("tcpmodbus connect error!")
	}
	err = client.SetEncoding(modbus.BIG_ENDIAN, modbus.LOW_WORD_FIRST)
	if err != nil {
		logrus.Error("tcpmodbus SetEncoding error!")
		return nil, errors.New("tcpmodbus SetEncoding error!")
	}
	// now that the client is created and configured, attempt to connect
	err = client.Open()
	if err != nil {
		// error out if we failed to connect/open the device
		// note: multiple Open() attempts can be made on the same client until
		// the connection succeeds (i.e. err == nil), calling the constructor again
		// is unnecessary.
		// likewise, a client can be opened and closed as many times as needed.
		logrus.Error("tcpmodbus connect error!")
		return nil, errors.New("tcpmodbus connect error!")
	}

	logrus.Info("tcpmodbus connect to " + config.Tcpmodbus.Host + " successful")

	return client, nil
}

/* 使用modbus slave模拟器测试 */

func ReadTcpModbus(client MQTT.Client, config config.Config, modbusclient *modbus.ModbusClient) {
	logrus.Info("read tcpmodbus start")
	for {
		time.Sleep(time.Second * time.Duration(config.Tcpmodbus.Interval))
		for _, item := range config.Tcpmodbus.Devices {
			if item.Register == "holding" {
				var sendmsg string
				for _, read := range item.RegisterTable {
					var val string
					switch read.Type {
					case "int":
						results, err := modbusclient.ReadRegister(uint16(read.StartAddr), modbus.HOLDING_REGISTER)
						if err != nil {
							logrus.Errorf("read holding register int error: %v", err)
						}
						val = strconv.Itoa(int(results))
						sendmsg = "{\"key\":" + "\"" + read.Name + "\"," + "\"val\":" + val + "}"
						//sendmsg = "{\"serial\":1,\"deviceType\":1,\"data\":\"{\\\"x\\\":1,\\\"y\\\":1,\\\"z\\\":1,\\\"rx\\\":1,\\\"ry\\\":1,\\\"rz\\\":1}\"}"
					case "float":
						results, err := modbusclient.ReadFloat32(uint16(read.StartAddr), modbus.HOLDING_REGISTER)
						if err != nil {
							fmt.Println(err.Error())
						}
						val = strconv.FormatFloat(float64(results), 'f', 2, 64)
						sendmsg = "{\"key\":" + "\"" + read.Name + "\"," + "\"val\":" + val + "}"
					}
					publish := client.Publish(item.Topic, 1, false, sendmsg)
					publish.Wait()
					logrus.Infof("send message on topic: %s ; Message: \x1b[%dm%s\x1b[0m", item.Topic, constants.Cyan, sendmsg)
				}
			}
		}
	}
}

// 读取mqtt信息并根据配置文件写入到指定的modbus地址中
func WriteTcpModbus(client MQTT.Client, config config.Config, modbusclient *modbus.ModbusClient) {
	logrus.Info("write tcpmodbus start")
	for {
		time.Sleep(time.Second * time.Duration(config.Tcpmodbus.Interval))
		// 读取mqtt消息
		for _, sub := range config.Mqttinfo.SubList {
			token := client.Subscribe(sub, 1, func(client MQTT.Client, msg MQTT.Message) {
				logrus.Infof("Received message on topic: %s\nMessage: \n%s", msg.Topic(), msg.Payload())
				enable := bos.CmdBo{}
				// 根据配置文件写入到modbus地址中
				for _, item := range config.Tcpmodbus.Devices {
					switch item.Register {
					case "write":
						if item.Topic == msg.Topic() {
							err := json.Unmarshal(msg.Payload(), &enable)
							if err != nil {
								logrus.Errorf("enable json unmarshal error: %v", err)
								return
							}
							if enable.Enable == nil {
								logrus.Errorf("enable json unmarshal nil")
								return
							}
							logrus.Infof("enable: %d", *enable.Enable)
							for _, write := range item.RegisterTable {
								switch *enable.Enable {
								case 1: // 启动
									if write.Name == "start" {
										err := startTcp(modbusclient, uint16(write.StartAddr))
										if err != nil {
											logrus.Errorf("start cgxi error: %v", err)
											return
										}
										logrus.Info("----------start cgxi successful----------")
									}
								case 0: // 停止
									if write.Name == "stop" {
										err := stopTcp(modbusclient, uint16(write.StartAddr))
										if err != nil {
											logrus.Errorf("stop cgxi error: %v", err)
											return
										}
										logrus.Info("----------stop cgxi successful----------")
									}
								}

							}
						}
					}
				}
			})
			token.Wait()
		}
	}
}

func startTcp(modbusclient *modbus.ModbusClient, addr uint16) error {
	// 启动
	err := modbusclient.WriteRegister(addr, uint16(1))
	if err != nil {
		err := fmt.Errorf("write start register error: %v", err)
		return err
	}
	time.Sleep(time.Second * time.Duration(1))
	// 回写
	err = modbusclient.WriteRegister(addr, uint16(0))
	if err != nil {
		err := fmt.Errorf("write start register recover error: %v", err)
		return err
	}

	return nil
}

func stopTcp(modbusclient *modbus.ModbusClient, addr uint16) error {
	// 停止
	err := modbusclient.WriteRegister(addr, uint16(1))
	if err != nil {
		err := fmt.Errorf("write stop register error: %v", err)
		return err
	}
	time.Sleep(time.Second * time.Duration(1))
	// 回写
	err = modbusclient.WriteRegister(addr, uint16(0))
	if err != nil {
		err := fmt.Errorf("write stop register recover error: %v", err)
		return err
	}

	return nil
}
