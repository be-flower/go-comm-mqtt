package model

// 设备列表返回
type Device_data struct {
	Id     int    `gorm:"column:id;type:integer;primary_key;autoincrement" json:"id"` //主键ID
	Serial string `gorm:"column:serial;type:varchar(128)" json:"serial"`              //物联网终端唯一识别码
	Type   string `gorm:"column:type;type:varchar(128)" json:"type"`                  //设备类型
	Data   string `gorm:"column:data;type:varchar(128)" json:"data"`                  //数据
	Time   string `gorm:"column:time;type:varchar(128)" json:"time"`                  //时间
}
