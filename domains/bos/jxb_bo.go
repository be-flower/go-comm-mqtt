package bos

type JxbBo struct {
	Serial     string `json:"serial"`      //TZ:探针 JXB:机械臂
	DeviceType string `json:"device_type"` //设备类型
}

type CmdBo struct {
	Enable *int `json:"enable"` //使能
}
