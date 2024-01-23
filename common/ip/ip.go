package ip

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"strconv"
	"strings"
)

/*获取包含 内网 IP*/
func GetIntranetIP() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		logrus.Fatalf("GetIntranetIp:panic!err = %v", err)
		os.Exit(1)
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !isPublicIP(ipnet.IP) {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

/*获取外网IP*/
func GetInternetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logrus.Fatalf("GetInternetIP:error!err = %v", err)
		return ""
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && isPublicIP(ipnet.IP) {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	logrus.Errorf("GetInternetIP:cannot get public ip!")
	return ""
}

/*获取内网IP*/
func GetPrivateIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logrus.Errorf("GetInternetIP:InterfaceAddrs error!err = %v", err)
		return ""
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && isPrivateIP(ipnet.IP) {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	logrus.Errorf("GetInternetIP:cannot get private ip!")
	return ""
}

func inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

func inet_aton(ipnr net.IP) int64 {
	bits := strings.Split(ipnr.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

/*将形如192.168.0.64 提取为 []byte*/
func Inet_AtoN(ip string) []byte {
	bits := strings.Split(ip, ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	return []byte{byte(b0), byte(b1), byte(b2), byte(b3)}
}

func ipBetween(from net.IP, to net.IP, test net.IP) bool {
	if from == nil || to == nil || test == nil {
		logrus.Debugf("IPBetween:An IP input is nil")
		return false
	}

	from16 := from.To16()
	to16 := to.To16()
	test16 := test.To16()
	if from16 == nil || to16 == nil || test16 == nil {
		logrus.Debugf("An ip did not convert to a 16 byte")
		return false
	}

	if bytes.Compare(test16, from16) >= 0 && bytes.Compare(test16, to16) <= 0 {
		return true
	}
	return false
}

/*是否为外网IP*/
func isPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}

	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

/*是否为内网IP*/
func isPrivateIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}

	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return true
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return true
		case ip4[0] == 192 && ip4[1] == 168:
			return true
		default:
			return false
		}
	}
	return false
}
