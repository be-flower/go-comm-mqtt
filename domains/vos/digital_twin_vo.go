package vos

type DigitalTwinVo struct {
	//Num int `json:"num"`
	// 上料位前后初始位
	UpInitPos bool `json:"UpInitPos"`
	// 上料位前后取料位
	UpTakePos bool `json:"UpTakePos"`
	// 上料位前后皮带位
	UpBeltPos bool `json:"UpBeltPos"`
	// 皮带线速度（脉冲）（上料）
	BeltSpeed bool `json:"BeltSpeed"`
	// 检测工位抬起高度（上料）
	UpPickHeight bool `json:"UpPickHeight"`
	// 检测工位等料高度（上料）
	UpMateHeight bool `json:"UpMateHeight"`
	// 机械臂位置姿态X
	X float64 `json:"x"`
	// 机械臂位置姿态Y
	Y float64 `json:"y"`
	// 机械臂位置姿态Z
	Z float64 `json:"z"`
	// 皮带线速度（脉冲）（出料）
	OutInitPos bool `json:"OutInitPos"`
	// 检测工位抬起高度（出料）
	OutTakePos bool `json:"OutTakePos"`
	// 检测工位等料高度（出料）
	OutBeltPos bool `json:"OutBeltPos"`
	// 出料位前后皮带线位
	OutBeltLinePos bool `json:"OutBeltLinePos"`
	// 出料位前后Ok位
	OutOk bool `json:"OutOk"`
	// 出料位前后OK抬头位
	OutOkLift bool `json:"OutOkLift"`
	// 检测结果(OK/NG)
	DateResult bool `json:"DateResult"`
	// 启动（1表示1层启动）
	Start int `json:"Start"`
	// 急停（停止）（1表示急停）
	Stop int `json:"Stop"`
	// 复位（1表示复位）
	Reset int `json:"Reset"`
	// 手自动模式（1表示自动，0表示手动模式）
	HandAuto int `json:"HandAuto"`
	// 相机状态
	CameraStatus int `json:"CameraStatus"`
	// 检测工位（1表示下降，2表示上升）
	DelectPos int `json:"DelectPos"`
	// 机械臂位置
	RobotArmPos int `json:"RobotArmPos"`
	// 机器人开始检测标志
	RobotStart bool `json:"RobotStart"`
	// 机器人检测完成标志
	RobotOver bool `json:"RobotOver"`
	// 机械臂位置姿态Rx
	Rx float64 `json:"rx"`
	// 机械臂位置姿态Ry
	Ry float64 `json:"ry"`
	// 机械臂位置姿态Rz
	Rz float64 `json:"rz"`
	//// 关节角度1
	//Joint1 float64 `json:"joint1"`
	//// 关节角度2
	//Joint2 float64 `json:"joint2"`
	//// 关节角度3
	//Joint3 float64 `json:"joint3"`
	//// 关节角度4
	//Joint4 float64 `json:"joint4"`
	//// 关节角度5
	//Joint5 float64 `json:"joint5"`
	//// 关节角度6
	//Joint6 float64 `json:"joint6"`
}

type DigitalTwinJointVo struct {
	// 关节角度1
	Joint1 float64 `json:"joint1"`
	// 关节角度2
	Joint2 float64 `json:"joint2"`
	// 关节角度3
	Joint3 float64 `json:"joint3"`
	// 关节角度4
	Joint4 float64 `json:"joint4"`
	// 关节角度5
	Joint5 float64 `json:"joint5"`
	// 关节角度6
	Joint6 float64 `json:"joint6"`
}
