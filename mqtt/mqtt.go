package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"go-comm-mqtt/config"
	"os"
	"strconv"
	"time"
)

func ConnMqtt(config config.Config) (MQTT.Client, chan bool) {

	hostname, _ := os.Hostname()

	server := "tcp://" + config.Mqttinfo.Host + ":" + strconv.Itoa(config.Mqttinfo.Port)
	//subtopic := config.Mqttinfo.SubList
	//pubtopic := config.Mqttinfo.PubList
	//qos := config.Mqttinfo.Qos
	clientid := hostname + strconv.Itoa(time.Now().Second())
	username := config.Mqttinfo.UserName
	password := config.Mqttinfo.PassWord

	connOpts := MQTT.NewClientOptions().AddBroker(server).SetClientID(clientid).SetCleanSession(true)
	if username != "" {
		connOpts.SetUsername(username)
		if password != "" {
			connOpts.SetPassword(password)
		}
	}
	// tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	// connOpts.SetTLSConfig(tlsConfig)

	// 自动重连机制，如网络不稳定可开启
	connOpts.SetAutoReconnect(true)      //启用自动重连功能
	connOpts.SetMaxReconnectInterval(30) //每30秒尝试重连
	connOpts.SetResumeSubs(true)         //设置为true，客户端会在重新连接后重新订阅之前的topic

	quit := make(chan bool)
	recmsg := make(chan [2]string, 300)

	// 设置当接收到不匹配任何已知订阅的消息时将调用的MessageHandler
	//connOpts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
	//	recmsg <- [2]string{msg.Topic(), string(msg.Payload())}
	//})

	// 现在作用为保活和退出
	go dealMqttMsg(recmsg, quit)

	client := MQTT.NewClient(connOpts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		logrus.Infof("Connected to %s", server)
	}

	/*
		下面有说到类似回显，但是重连之后虽然通过 SetResumeSubs 方法重新订阅了主题，
		但是不能自动设置 SetDefaultPublishHandler 方法来接收消息，导致 go dealMqttMsg(recmsg, quit) 无法将日志打出来
		修改整体逻辑，将日志打印放到发送的时候打印，并等待发送成功打印日志
	*/
	// 其实不用订阅这个也没事，你在读的时候那里发布已经发布到mqtt了，这里订阅只是类似于回显
	//for _, item := range pubtopic {
	//	if token := client.Subscribe(item, byte(qos), nil); token.Wait() && token.Error() != nil {
	//		panic(token.Error())
	//	} else {
	//		logrus.Infof("Subscribe topic  %s  successful", item)
	//	}
	//}

	return client, quit
}

// 处理订阅到的MQTT消息
func dealMqttMsg(msg chan [2]string, exit chan bool) {
	for {
		select {
		//case incoming := <-msg:
		//	logrus.Infof("send message on topic: %s ; Message: \x1b[%dm%s\x1b[0m", incoming[0], constants.Cyan, incoming[1])
		case <-exit:
			return
		default:
			//logrus.Printf("empty\n")
			time.Sleep(time.Millisecond * 10)
		}
	}
}
