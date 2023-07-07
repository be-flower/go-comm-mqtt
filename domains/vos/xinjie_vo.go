package vos

import "go-comm-mqtt/domains/bos"

type XinJieJxbDataVo struct {
	bos.JxbBo
	Data XinJieJxbVo `json:"data"`
}

type XinJieJxbVo struct {
	XinJieJxbTcp
	XinJieJoint
}

type XinJieJxbTcp struct {
	X  string `json:"x"`
	Y  string `json:"y"`
	Z  string `json:"z"`
	Rx string `json:"rx"`
	Ry string `json:"ry"`
	Rz string `json:"rz"`
}

type XinJieJoint struct {
	Joint1 string `json:"joint1"`
	Joint2 string `json:"joint2"`
	Joint3 string `json:"joint3"`
	Joint4 string `json:"joint4"`
	Joint5 string `json:"joint5"`
	Joint6 string `json:"joint6"`
}
