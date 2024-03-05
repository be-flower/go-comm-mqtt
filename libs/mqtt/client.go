package mqtt

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Subscribers struct {
	Topic    string `json:"topic"`
	Qos      byte   `json:"qos"`
	Callback func(topic string, msg []byte) error
}

type MQTTClient struct {
	conf          *Configuration
	autoReconnect bool
	cleanSession  bool
	debug         bool

	subs []Subscribers

	cli *Client
}

func NewMQTTClient(conf *Configuration, autoReconnect bool, cleanSession bool, debug bool, subs []Subscribers) (*MQTTClient, error) {
	if conf.ClientId == "" {
		logrus.Errorf("NewMQTTClient:缺少参数ClientId!")
		return nil, fmt.Errorf("NewMQTTClient:缺少参数ClientId!")
	}

	this := &MQTTClient{
		conf:          conf,
		autoReconnect: autoReconnect,
		cleanSession:  cleanSession,
		debug:         debug,

		subs: subs,
	}

	if conf.Timeout == 0 {
		conf.Timeout = 10
	}

	if conf.DefaultCallback == nil {
		conf.DefaultCallback = this.Callback
	}

	if conf.ConnectedCallback == nil {
		conf.ConnectedCallback = this.Connected
	}

	if conf.ConnectionLostCallback == nil {
		conf.ConnectionLostCallback = this.ConnectionLost
	}

	cli := Connect(conf, autoReconnect, cleanSession, debug)
	if cli != nil {
		logrus.Infof("NewMQTTClient success!")
	} else {
		logrus.Errorf("NewMQTTClient error!")
		return nil, fmt.Errorf("NewMQTTClient error!")
	}
	this.cli = cli

	go this.CheckConnection()

	return this, nil
}

func (this *MQTTClient) CheckConnection() {
	for {
		if this.cli == nil || !this.cli.IsConnected() {
			this.Reconnect()
		}

		time.Sleep(time.Second * time.Duration(30))
	}
}

func (this *MQTTClient) Reconnect() {
	cli := Connect(this.conf, this.autoReconnect, this.cleanSession, this.debug)
	if cli != nil {
		this.cli = cli
		logrus.Infof("Reconnect:Connect success!")
	} else {
		logrus.Errorf("Reconnect:Connect error!")
	}
}

func (this *MQTTClient) Disconnect() {
	if this.cli != nil {
		this.cli.Disconnect(1000)
	}
}

func (this *MQTTClient) Connected(clientID string) error {
	logrus.Infof("Connected!clientId = %v", clientID)

	if this.cli != nil {
		for i, _ := range this.subs {
			sub := this.subs[i]
			go this.cli.Subscribe(sub.Topic, sub.Qos, sub.Callback)
		}
	}

	return nil
}

func (this *MQTTClient) ConnectionLost(clientID string) error {
	logrus.Infof("ConnectionLost!clientID = %v", clientID)
	return nil
}

func (this *MQTTClient) Callback(topic string, buff []byte) error {
	logrus.Infof("Callback!buff = %v", string(buff))
	return nil
}

func (this *MQTTClient) Publish(topic string, qos byte, retained bool, data []byte) error {
	logrus.Infof("Publish:topic = %v, qos = %v, retained = %v, data = %v", topic, qos, retained, string(data))

	if this.cli != nil && this.cli.IsConnected() {
		return this.cli.Publish(topic, qos, retained, data)
	} else {
		return fmt.Errorf("not connected!please check connection!")
	}
}
