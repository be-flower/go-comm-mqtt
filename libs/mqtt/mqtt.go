package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type ConfigurationTopic struct {
	/*topic*/
	Topic string
	QoS   byte
	/*回调函数*/
	Callback func(buff []byte) error
}

/*MQTT配置*/
type Configuration struct {
	/*是否启用SSL*/
	SSL bool
	/*证书文件,crt文件*/
	CertFile string
	/*密钥文件*/
	KeyFile string
	/*客户端ID*/
	ClientId string
	/*用户名*/
	UserName string
	/*密码*/
	Password string
	Retained bool
	/*broker地址*/
	BrokerAddr string
	BrokerPort int
	/*超时时间*/
	Timeout int
	/*回调函数*/
	DefaultCallback func(topic string, buff []byte) error
	/*订阅的主题*/
	SubscriberTopics []ConfigurationTopic
	/*连接的回调*/
	ConnectedCallback func(clientId string) error
	/*连接断开的回调*/
	ConnectionLostCallback func(clientId string) error
}

/*MQTT实例*/
type Client struct {
	Conf Configuration
	cli  mqtt.Client
}

func NewTLSConfig(certFile string, keyFile string) *tls.Config {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(certFile)
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	} else {
		logrus.Errorf("NewTLSConfig:ReadFile error! certFile = %v, err = %v", certFile, err)
	}

	// Import client certificate/key pair
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		logrus.Errorf("NewTLSConfig:LoadX509KeyPair error! err = %v", err)
		panic(err)
	}

	// Just to print out the client certificate..
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		logrus.Errorf("NewTLSConfig:ParseCertificate error! err = %v", err)
		panic(err)
	}

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}
}

/*MQTT日志*/
type MQTTErrorLogger struct{}
type MQTTCriticalLogger struct{}
type MQTTWarnLogger struct{}
type MQTTDebugLogger struct{}

func (MQTTErrorLogger) Println(v ...interface{}) {
	logrus.Errorln(v...)
}

func (MQTTErrorLogger) Printf(format string, v ...interface{}) {
	logrus.Errorf(format, v...)
}

func (MQTTCriticalLogger) Println(v ...interface{}) {
	logrus.Errorln(v...)
}

func (MQTTCriticalLogger) Printf(format string, v ...interface{}) {
	logrus.Errorf(format, v...)
}

func (MQTTWarnLogger) Println(v ...interface{}) {
	logrus.Warnln(v...)
}

func (MQTTWarnLogger) Printf(format string, v ...interface{}) {
	logrus.Warnf(format, v...)
}

func (MQTTDebugLogger) Println(v ...interface{}) {
	logrus.Debugln(v...)
}

func (MQTTDebugLogger) Printf(format string, v ...interface{}) {
	logrus.Debugf(format, v...)
}

/*连接到Mqtt*/
func Connect(conf *Configuration, autoReconnect bool, cleanSession bool, debug bool) *Client {
	inst := new(Client)
	inst.Conf = *conf
	opts := inst.createClientOptions()
	opts.AutoReconnect = autoReconnect
	opts.CleanSession = cleanSession

	if debug {
		mqtt.ERROR = MQTTErrorLogger{}
		mqtt.CRITICAL = MQTTCriticalLogger{}
		mqtt.WARN = MQTTWarnLogger{}
		mqtt.DEBUG = MQTTDebugLogger{}
	} else {
		mqtt.WARN = MQTTWarnLogger{}
		mqtt.ERROR = MQTTErrorLogger{}
		mqtt.CRITICAL = MQTTCriticalLogger{}
	}

	client := mqtt.NewClient(opts)
	inst.cli = client

	token := client.Connect()
	for !token.WaitTimeout(time.Duration(conf.Timeout) * time.Second) {
	}
	if err := token.Error(); err != nil {
		logrus.Errorf("connect:error!err = %v,clientId = %v, BrokerAddr = %v, BrokerPort = %v", err, conf.ClientId, conf.BrokerAddr, conf.BrokerPort)
		return nil
	}

	return inst
}

// IsConnected returns a bool signifying whether
// the client is connected or not.
func (this *Client) IsConnected() bool {
	if this.cli != nil {
		return this.cli.IsConnected()
	}
	return false
}

func (this *Client) IsConnectionOpen() bool {
	if this.cli != nil {
		return this.cli.IsConnectionOpen()
	}
	return false
}

func (this *Client) Disconnect(quiesce uint) {
	if this.cli != nil {
		this.cli.Disconnect(quiesce)
	}
	return
}

func (this *Client) Publish(topic string, qos byte, retained bool, data []byte) error {
	if this.cli != nil {
		token := this.cli.Publish(topic, qos, retained, data)

		for !token.WaitTimeout(time.Duration(this.Conf.Timeout) * time.Second) {
		}

		if err := token.Error(); err != nil {
			logrus.Errorf("Publish:error!err = %v,data = %v", err, data)
			return err
		}

		return nil
	} else {
		return errors.New("not connected!please check connection!")
	}
}

func (this *Client) createClientOptions() *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	if this.Conf.SSL {
		tlsconfig := NewTLSConfig(this.Conf.CertFile, this.Conf.KeyFile)
		opts.AddBroker(fmt.Sprintf("ssl://%s:%d", this.Conf.BrokerAddr, this.Conf.BrokerPort))
		opts.SetClientID(this.Conf.ClientId).SetTLSConfig(tlsconfig)
	} else {
		opts.AddBroker(fmt.Sprintf("tcp://%s:%d", this.Conf.BrokerAddr, this.Conf.BrokerPort))
		opts.SetClientID(this.Conf.ClientId)
	}

	opts.SetUsername(this.Conf.UserName)
	opts.SetPassword(this.Conf.Password)
	opts.WillRetained = this.Conf.Retained

	opts.SetOnConnectHandler(func(cli mqtt.Client) {
		if this.Conf.ConnectedCallback != nil {
			this.Conf.ConnectedCallback(this.Conf.ClientId)
		}
	})
	opts.SetConnectionLostHandler(func(cli mqtt.Client, err error) {
		if this.Conf.ConnectionLostCallback != nil {
			this.Conf.ConnectionLostCallback(this.Conf.ClientId)
		}
	})
	/*收到消息未匹配到任意topic中，*/
	var f mqtt.MessageHandler = func(cli mqtt.Client, msg mqtt.Message) {
		/*2019-10-18:添加默认回调处理*/
		if this.Conf.DefaultCallback != nil {
			this.Conf.DefaultCallback(msg.Topic(), msg.Payload())
		}
	}
	opts.SetDefaultPublishHandler(f)
	return opts
}

func (this *Client) Subscribe(topic string, qos byte, cb func(topic string, buff []byte) error) error {
	token := this.cli.Subscribe(topic, qos, func(client mqtt.Client, msg mqtt.Message) {
		logrus.Infof("Received Message:topic:%v,payload:%v", msg.Topic(), string(msg.Payload()))

		if cb != nil {
			cb(msg.Topic(), msg.Payload())
		}
	})

	if !token.WaitTimeout(3 * time.Second) {
		logrus.Errorf("Subscribe:Error on Subscribe: topic:%v,qos:%v,err:%v", topic, qos, token.Error())
		return token.Error()
	}

	logrus.Infof("Subscribe:topic = %v,qos = %v success", topic, qos)
	return nil
}
