package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
)

// TwoIntToFloat 用于aobo机械臂制定的规则，将两个int拼接为float
func TwoIntToFloat(a, b int) float64 {
	// 正负标志
	flag := ""
	if b < 0 {
		flag = "-"
		b = -b
		a = -a
	}
	// 拼接字符串
	str := fmt.Sprintf(flag+"%.4f", float64(a)+float64(b)/100)
	// string 转 float
	floatOutput, err := strconv.ParseFloat(str, 64)
	if err != nil {
		logrus.Errorf("strconv.ParseFloat error: %v", err)
	}
	return floatOutput
}
