package mqttcloud

//连接 云端 mqtt
import (
	"encoding/json"
	"errors"
	"go-comm-mqtt/conf"
	"go-comm-mqtt/data/model"
	"go-comm-mqtt/libs/constants"
	"go-comm-mqtt/libs/mqtt"
	"go-comm-mqtt/libs/uuid"
	"go-comm-mqtt/serve/productionline/web"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	MQTTCloudManager *MQTTCloudsManager
)

// MQTTCloudsManager 云端mqtt
type MQTTCloudsManager struct {
	/*接收设备的 mqtt 实例*/
	deviceCli *mqtt.Client
	init      bool
	lock      sync.Mutex
}

func NewMQTTCloudManager() *MQTTCloudsManager {
	return &MQTTCloudsManager{}
}

// 初始化
func (dm *MQTTCloudsManager) Init() {
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
var CloudRequestProcessFirst = true

func InitEdgeTopics() {
	edgeTopics = []Topic{
		// 接收云平台消息
		Topic{"/platform_command_push_request/", 0, CloudRequestProcess},
		// 测试
		//Topic{"hello", 0, hello},
	}
}

func (dm *MQTTCloudsManager) mqttReload() {
	//mqtt 服务改变
	if dm.deviceCli != nil {
		dm.deviceCli.Disconnect(0)
	}

	logrus.Infof("preparing EdgesConnect!")
	dm.lock.Lock()
	defer dm.lock.Unlock()
	if conf.Conf.MqttCloud.ClientId == "" {
		conf.Conf.MqttCloud.ClientId = string(uuid.NewV4().String())
	}
	dm.deviceCli = mqtt.Connect(&mqtt.Configuration{
		ClientId:               conf.Conf.MqttCloud.ClientId,
		UserName:               conf.Conf.MqttCloud.UserName,
		Password:               conf.Conf.MqttCloud.PassWord,
		BrokerAddr:             conf.Conf.MqttCloud.Host,
		BrokerPort:             conf.Conf.MqttCloud.Port,
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

func (dm *MQTTCloudsManager) deviceConnected(clientId string) error {
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

func (dm *MQTTCloudsManager) deviceConnectionLost(clientId string) error {
	logrus.Infof("cbEdgesConnectionLost!clientId = %v", clientId)

	return nil
}

func (dm *MQTTCloudsManager) CheckDeviceConnection() {
	for {
		if dm.deviceCli == nil || !dm.deviceCli.IsConnected() {
			dm.mqttReload()
			logrus.Infof("CheckEdgesConnection:CloudReconnect")
		}

		time.Sleep(time.Second * time.Duration(10))
	}
}

/*注册TOPIC的回调*/
func (dm *MQTTCloudsManager) deviceCmdRequestProcess(topic string, msg []byte) error {
	return nil
}

// CloudRequestProcess 处理云平台消息
func CloudRequestProcess(topic string, msg []byte) error {

	// 第一次不处理
	if CloudRequestProcessFirst {
		CloudRequestProcessFirst = false
		return nil
	}

	MQTTCloudManager.lock.Lock()
	defer MQTTCloudManager.lock.Unlock()

	var req model.MQTTCommandRequest

	if err := json.Unmarshal(msg, &req); err != nil {
		logrus.Errorf("CloudCmdRequestProcess:Unmarshal error! err = %v,msg = %v", err, string(msg))
		return err
	}

	if req.Serial == "CX001" {
		err := web.ProductionLineDealCmd(req.Event)
		if err != nil {
			logrus.Errorf("CloudCmdRequestProcess:ProductionLineDealCmd error! err = %v", err)
			return err
		}
	}

	return nil
}

func hello(topic string, msg []byte) error {

	MQTTCloudManager.lock.Lock()
	defer MQTTCloudManager.lock.Unlock()

	logrus.Infof(string(msg))

	return nil
}

func Publish2Cloud(topic string, qos byte, retained bool, data []byte) error {
	logrus.Infof("Publish2Cloud:topic = %v, qos = %v, retained = %v, data = \x1b[%dm%v\x1b[0m", topic, qos, retained, constants.Cyan, string(data))

	if MQTTCloudManager.deviceCli != nil {
		return MQTTCloudManager.deviceCli.Publish(topic, qos, retained, data)
	} else {
		return errors.New("not connected!please check connection!")
	}
}
