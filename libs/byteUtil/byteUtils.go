package byteUtils

import "strings"

func TrimEscapeString(src []byte) []byte {
	dst := src
	if len(src) > 0 {
		// ""开始
		if src[0] == '"' {
			// 如果紧跟的是转义字符
			data := src[1 : len(src)-1]
			data2 := strings.Replace(string(data), `\`, ``, -1)

			if data2[0] == '"' {
				dst = []byte(data2[1 : len(data2)-1])
			} else {
				dst = []byte(data2)
			}
		}
	}

	return dst
}
