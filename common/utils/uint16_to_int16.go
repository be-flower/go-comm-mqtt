package utils

/*
	正常来说读出来都是无符号的（uint16），但是有些需求现在是需要读出来有符号的（int16）
	所以添加工具函数，将无符号的uint16转换为有符号的int16
*/

func Uint16sToInt16s(uint16s []uint16) (int16s []int16) {
	for _, v := range uint16s {
		int16s = append(int16s, int16(v))
	}
	return
}

func Uint16ToInt16(uint16s uint16) (int16s int16) {
	return int16(uint16s)
}
