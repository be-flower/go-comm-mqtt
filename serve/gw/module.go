package gw

import (
	"fmt"
	"go-comm-mqtt/conf"
	"go-comm-mqtt/domains/bos"
	"go-comm-mqtt/domains/vos"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

func ReadGatewayModule(client MQTT.Client, config conf.Config) {
	for {
		var sendmodulemsg string
		time.Sleep(time.Second * time.Duration(config.Tcpmodbus.Interval))
		gwDateVo := vos.GwDataVo{
			GwBo: bos.GwBo{
				Serial:     "TZ001",
				DeviceType: "42001",
			},
			Data: vos.GwModuleVo{},
		}
		info, err := GetGatewayModuleInfo()
		gwDateVo.Data = info
		if err != nil {
			logrus.Error("failed to get gateway module info")
			continue
		}
		sendmodulemsg = "{\"serial\":\"" + gwDateVo.Serial + "\",\"deviceType\":\"" + gwDateVo.DeviceType + "\",\"data\":\"{" +
			"\\\"imei\\\":\\\"" + gwDateVo.Data.Imei + "\\\",\\\"imsi\\\":\\\"" + gwDateVo.Data.Imsi + "\\\",\\\"ip\\\":\\\"" + gwDateVo.Data.Ip + "\\\",\\\"rssi\\\":" + gwDateVo.Data.Rssi + ",\\\"sinr\\\":" + gwDateVo.Data.Sinr + ",\\\"rsrp\\\":" + gwDateVo.Data.Rsrp + "}\"}"
		client.Publish("hello", 1, false, sendmodulemsg)
	}
}

func GetGatewayModuleInfo() (vos.GwModuleVo, error) {

	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("GetGatewayModuleInfo panic error: %v", err)
		}
	}()

	gwModuleVo := vos.GwModuleVo{}

	// Imei
	for i := 0; i < 3; i++ {
		imeiCmd := "ate.sh ati | awk '{if (NR==6) {print $2}}'"
		imei, err := startCmd(imeiCmd)
		if err != nil {
			logrus.Errorf("GetGatewayModuleInfo imeiCmd error: %v", err)
			return vos.GwModuleVo{}, err
		}
		if imei != "" {
			gwModuleVo.Imei = imei
			break
		}
	}

	// Imsi
	for i := 0; i < 3; i++ {
		imsiCmd := "ate.sh at+cimi | awk '{if (NR==2) {print $1}}'"
		imsi, err := startCmd(imsiCmd)
		if err != nil {
			logrus.Errorf("GetGatewayModuleInfo imsiCmd error: %v", err)
			return vos.GwModuleVo{}, err
		}
		if imsi != "" {
			gwModuleVo.Imsi = imsi
			break
		}
	}

	// Ip
	ipCmd := "ate.sh at+cgdcont? | grep +CGDCONT | awk '{if (NR==1) {print $0}}' | awk -F, '{print $4}' | sed 's/\\\"//g'"
	ip, err := startCmd(ipCmd)
	if err != nil {
		logrus.Errorf("GetGatewayModuleInfo ipCmd error: %v", err)
		return vos.GwModuleVo{}, err
	}
	gwModuleVo.Ip = ip

	// Rssi
	rssiCmd := "ate.sh at+csq? | grep CSQ: | awk '{print $2}' | awk -F, '{print $1}'"
	rssiStr, err := startCmd(rssiCmd)
	if err != nil {
		logrus.Errorf("GetGatewayModuleInfo rssiCmd error: %v", err)
		return vos.GwModuleVo{}, err
	}
	var rssi float64 = 0
	if rssiStr != "" {
		rssi, err = strconv.ParseFloat(rssiStr, 64)
		if err != nil {
			logrus.Errorf("GetGatewayModuleInfo rssiStr to float64 error: %v", err)
			return vos.GwModuleVo{}, err
		}
	}
	if rssi <= 0 {
		rssi = -113
	} else if rssi > 32 {
		rssi = 255
	} else {
		rssi = (2 * rssi) - 113
	}
	gwModuleVo.Rssi = fmt.Sprintf("%.2f", rssi)

	// Sinr
	sinrCmd := "ate.sh at+gtccinfo? | awk '{if (NR==4) {print $0}}' | awk -F, '{print $11}'"
	sinrStr, err := startCmd(sinrCmd)
	if err != nil {
		logrus.Errorf("GetGatewayModuleInfo sinrCmd error: %v", err)
		return vos.GwModuleVo{}, err
	}
	var sinr float64 = 0
	if sinrStr != "" {
		sinr, err = strconv.ParseFloat(sinrStr, 64)
		if err != nil {
			logrus.Errorf("GetGatewayModuleInfo sinrStr to float64 error: %v", err)
			return vos.GwModuleVo{}, err
		}
	}
	if sinr > 127 || sinr < 0 {
		sinr = 255
	} else {
		sinr = (sinr * 0.5) - 23.0
	}
	gwModuleVo.Sinr = fmt.Sprintf("%.2f", sinr)

	// Rsrp
	rsrpCmdGen := "ate.sh at+gtccinfo?"
	generation := 0
	gen, err := startCmd(rsrpCmdGen)
	if err != nil {
		logrus.Errorf("GetGatewayModuleInfo rsrpCmdGen error: %v", err)
		return vos.GwModuleVo{}, err
	}
	if strings.Contains(gen, "NR service cell") {
		generation = 5
	} else if strings.Contains(gen, "LTE-NR EN-DC service cell") {
		generation = 5
	} else if strings.Contains(gen, "LTE") {
		generation = 4
	}
	rsrpCmd := "ate.sh at+gtccinfo? | awk '{if (NR==4) {print $0}}' | awk -F, '{print $13}'"
	rsrpStr, err := startCmd(rsrpCmd)
	if err != nil {
		logrus.Errorf("GetGatewayModuleInfo rsrpCmd error: %v", err)
		return vos.GwModuleVo{}, err
	}
	var rsrp float64 = 0
	if rsrpStr != "" {
		rsrp, err = strconv.ParseFloat(rsrpStr, 64)
		if err != nil {
			logrus.Errorf("GetGatewayModuleInfo rsrpStr to float64 error: %v", err)
			return vos.GwModuleVo{}, err
		}
	}
	if generation == 5 {
		if rsrp > 126 || rsrp < 0 {
			rsrp = 255
		} else {
			rsrp = (rsrp * 0.5) - 43.0
		}
	} else if generation == 4 {
		if rsrp > 34 || rsrp < 0 {
			rsrp = 255
		} else {
			rsrp = (rsrp * 0.5) - 19.5
		}
	}
	gwModuleVo.Rsrp = fmt.Sprintf("%.2f", rsrp)

	return gwModuleVo, nil
}

func startCmd(command string) (string, error) {

	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("startCmd panic: %v", err)
		}
	}()

	cmd := exec.Command("/bin/bash", "-c", command)
	// 获取管道输入
	output, err := cmd.StdoutPipe()
	if err != nil {
		err := fmt.Errorf("无法获取命令的标准输出管道: %v", err)
		return "", err
	}

	// 执行Linux命令
	if err = cmd.Start(); err != nil {
		err := fmt.Errorf("linux命令执行失败，请检查命令输入是否有误: %v", err)
		return "", err
	}

	// 读取所有输出
	bytes, err := ioutil.ReadAll(output)
	if err != nil {
		err := fmt.Errorf("打印异常，请检查")
		return "", err
	}

	if err = cmd.Wait(); err != nil {
		err := fmt.Errorf("wait: %v", err.Error())
		return "", err
	}

	// 读出来后面多了一个换行符
	var str string
	if len(bytes) != 0 {
		str = string(bytes[:len(bytes)-1])
	}

	return str, nil
}
