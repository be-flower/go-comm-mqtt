package modol

type DigitalTwinHistoryTimes struct {
	// 主键ID
	ID int64 `gorm:"primarykey,autoIncrement,column:'id',type:'int',size:20,comment:'主键ID'" json:"id"`
	// 运行结果 1:OK 2:NG
	Result int `gorm:"column:'result',type:'int',size:1,not null,comment:'运行结果 1:OK 2:NG'" json:"result"`
	// 创建时间
	CreatedAt int64 `gorm:"column:'created_at',type:'int',size:20,autoCreateTime:milli,not null,comment:'创建时间'" json:"created_at"`
	// 更新时间
	UpdatedAt int64 `gorm:"column:'updated_at',type:'int',size:20,autoUpdateTime:milli,not null,comment:'更新时间'" json:"updated_at"`
	// 删除时间
	DeletedAt int64 `gorm:"column:'deleted_at',type:'int',size:20,index,comment:'删除时间'" json:"deleted_at"`
}
