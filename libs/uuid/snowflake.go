package uuid

import (
	"errors"
	"time"
)

/*
 * 分布式下的UUID生成方式
 * 最高位1比特-0
 * 41比特-毫秒时间
 * 10比特-机器编号ID
 * 索引号12比特
 * 组成：0(1 bit) | timestamp in milli second (41 bit) | machine id (10 bit) | index (12 bit)
 * 每毫秒最多生成4096个id，集群机器最多1024台
 */

type Snowflake struct {
	lastTimestamp int64
	index         int16
	machId        int16
}

func (s *Snowflake) Init(id int16) error {
	if id > 0xff {
		return errors.New("illegal machine id")
	}

	s.machId = id
	s.lastTimestamp = time.Now().UnixNano() / 1e6
	s.index = 0
	return nil
}

func (s *Snowflake) GetIID() (int64, error) {
	curTimestamp := time.Now().UnixNano() / 1e6
	if curTimestamp == s.lastTimestamp {
		s.index++
		if s.index > 0xfff {
			s.index = 0xfff
			return -1, errors.New("out of range")
		}
	} else {
		//fmt.Printf("id/ms:%d -- %d\n", s.lastTimestamp, s.index)
		s.index = 0
		s.lastTimestamp = curTimestamp
	}
	return int64((0x1ffffffffff&s.lastTimestamp)<<22) + int64(0xff<<10) + int64(0xfff&s.index), nil
}
