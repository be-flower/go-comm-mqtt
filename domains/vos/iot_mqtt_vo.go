package vos

// 上报的数据
type IotMqttDataVo struct {
	// 上料位前后初始位
	UpInitPos int `json:"QD1Aix"`
	// 上料位前后取料位
	UpTakePos int `json:"QD1E9m"`
	// 上料位前后皮带位
	UpBeltPos int `json:"QD18xp"`
	// 皮带线速度（脉冲）（上料）
	BeltSpeed int `json:"QD1CD7"`
	// 检测工位抬起高度（上料）
	UpPickHeight int `json:"QD19ba"`
	// 检测工位等料高度（上料）
	UpMateHeight int `json:"QD1D5A"`
	// 皮带线速度（脉冲）（出料）
	OutInitPos int `json:"QD1Cq9"`
	// 检测工位抬起高度（出料）
	OutTakePos int `json:"QD1OJK"`
	// 检测工位等料高度（出料）
	OutBeltPos int `json:"QD1LIR"`
	// 出料位前后皮带线位
	OutBeltLinePos int `json:"QD1D3k"`
	// 出料位前后Ok位
	OutOk int `json:"QD17wK"`
	// 出料位前后OK抬头位
	OutOkLift int `json:"QD14Ji"`
	// 检测结果(OK/NG)
	DateResult int `json:"QD14BU"`
	// 启动（1表示1层启动）
	Start int `json:"QD19HC"`
	// 急停（停止）（1表示急停）
	Stop int `json:"QD1B36"`
	// 复位（1表示复位）
	Reset int `json:"QD1DKs"`
	// 手自动模式（1表示自动，0表示手动模式）
	HandAuto int `json:"QD1K3O"`
	// 相机状态
	CameraStatus int `json:"QD1MTh"`
	// 检测工位（1表示下降，2表示上升）
	DelectPos int `json:"QD1CWf"`
	// 机械臂位置
	RobotArmPos int `json:"QD1DPE"`
	// 机器人开始检测标志
	RobotStart int `json:"QD1NHt"`
	// 机器人检测完成标志
	RobotOver int `json:"QD1KVu"`
}

// 上研 采集数据
type IotMqttGatherVo struct {
	Et string                      `json:"et"`
	Da []IotMqttGatherDeviceDataVo `json:"da"`
}

type IotMqttGatherDeviceDataVo struct {
	Id string        `json:"id"`
	Da IotMqttDataVo `json:"da"`
}

// 上研 上电
type IotMqttPowerOnVo struct {
	Et string `json:"et"`
	Ip string `json:"ip"`
}

// 上研 心跳状态
type IotMqttHeartbeatVo struct {
	Et string                         `json:"et"`
	Ip string                         `json:"ip"`
	Da []IotMqttHeartbeatDeviceDataVo `json:"da"`
}

type IotMqttHeartbeatDeviceDataVo struct {
	Id string `json:"id"`
	Ds int    `json:"ds"` // 0:在线 1:离线
}
