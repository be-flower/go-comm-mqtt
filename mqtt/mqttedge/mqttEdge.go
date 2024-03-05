package mqttedge

import (
	"errors"
	"go-comm-mqtt/conf"
	"go-comm-mqtt/libs/constants"
	"go-comm-mqtt/libs/uuid"

	"github.com/sirupsen/logrus"

	"go-comm-mqtt/libs/mqtt"

	"sync"
	"time"
)

var (
	MQTTEdgeManager *MQTTEdgesManager
)

// MQTTEdgesManager 边缘设备mqtt
type MQTTEdgesManager struct {
	/*接收设备的 mqtt 实例*/
	deviceCli *mqtt.Client
	init      bool
	lock      sync.Mutex
}

func NewMQTTEdgeManager() *MQTTEdgesManager {
	return &MQTTEdgesManager{}
}

// 初始化
func (dm *MQTTEdgesManager) Init() {
	if dm.init {
		return
	}
	InitEdgeTopics()

	dm.init = true
	dm.deviceCli = nil
	dm.mqttReload()

	go dm.CheckDeviceConnection()
	//初始化mqtt
}

type Topic struct {
	topic string
	qos   byte
	cb    func(topic string, buff []byte) error
}

var edgeTopics []Topic

func InitEdgeTopics() {
	edgeTopics = []Topic{}
}

func (dm *MQTTEdgesManager) mqttReload() {
	//mqtt 服务改变
	if dm.deviceCli != nil {
		dm.deviceCli.Disconnect(0)
	}

	logrus.Infof("preparing EdgesConnect!")
	dm.lock.Lock()
	defer dm.lock.Unlock()
	if conf.Conf.MqttEdge.ClientId == "" {
		conf.Conf.MqttEdge.ClientId = string(uuid.NewV4().String())
	}
	dm.deviceCli = mqtt.Connect(&mqtt.Configuration{
		ClientId:               conf.Conf.MqttEdge.ClientId,
		UserName:               conf.Conf.MqttEdge.UserName,
		Password:               conf.Conf.MqttEdge.PassWord,
		BrokerAddr:             conf.Conf.MqttEdge.Host,
		BrokerPort:             conf.Conf.MqttEdge.Port,
		Timeout:                10,
		DefaultCallback:        dm.deviceCmdRequestProcess,
		ConnectedCallback:      dm.deviceConnected,
		ConnectionLostCallback: dm.deviceConnectionLost,
	}, false, false, false)

	if dm.deviceCli != nil {
		logrus.Infof("EdgesConnect success!")
	} else {
		logrus.Errorf("EdgesConnect error!")
	}

}

func (dm *MQTTEdgesManager) deviceConnected(clientId string) error {
	logrus.Infof("cbEdgesConnected!clientId = %v", clientId)

	// 订阅topic
	logrus.Infof("cbEdgeConnected!clientId = %v,edgeTopics len = %v", clientId, len(edgeTopics))
	dm.lock.Lock()
	defer dm.lock.Unlock()
	// 订阅topic
	if dm.deviceCli != nil {
		for _, t := range edgeTopics {
			go func(topic string, qos byte, cb func(topic string, buff []byte) error) {
				logrus.Infof("cbEdgeConnected:Subscribing topic:%v!", topic)

				for {
					if dm.deviceCli != nil && dm.deviceCli.IsConnected() {
						err := dm.deviceCli.Subscribe(topic, qos, cb)
						if err != nil {
							logrus.Errorf("cbEdgeConnected:Subscribe topic:%v Failed!", topic)
							time.Sleep(time.Second * time.Duration(10))
							continue
						} else {
							logrus.Infof("cbEdgeConnected:Subscribe topic:%v success!", topic)
						}
					}
					break
				}
			}(t.topic, t.qos, t.cb)
		}
	}

	return nil
}

func (dm *MQTTEdgesManager) deviceConnectionLost(clientId string) error {
	logrus.Infof("cbEdgesConnectionLost!clientId = %v", clientId)

	return nil
}

func (dm *MQTTEdgesManager) CheckDeviceConnection() {
	for {
		if dm.deviceCli == nil || !dm.deviceCli.IsConnected() {
			dm.mqttReload()
			logrus.Infof("CheckEdgesConnection:CloudReconnect")
		}

		time.Sleep(time.Second * time.Duration(10))
	}
}

/*注册TOPIC的回调*/
func (dm *MQTTEdgesManager) deviceCmdRequestProcess(topic string, msg []byte) error {
	return nil
}

func Publish2Edge(topic string, qos byte, retained bool, data []byte) error {
	logrus.Infof("Publish2Edge:topic = %v, qos = %v, retained = %v, data = \x1b[%dm%v\x1b[0m", topic, qos, retained, constants.Cyan, string(data))

	if MQTTEdgeManager.deviceCli != nil {
		return MQTTEdgeManager.deviceCli.Publish(topic, qos, retained, data)
	} else {
		return errors.New("not connected!please check connection!")
	}
}
