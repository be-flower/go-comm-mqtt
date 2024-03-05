package db

import (
	"go-comm-mqtt/data/model"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
)

var db *gorm.DB = nil
var dbLock sync.Mutex //锁

func Init() {

	openDB()
	if nil == db {
		logrus.Debugln("init db err panic:")
		panic("open db err")
	}
	logrus.Debugln("init db")

	//创建web数据记录表
	logrus.Infof("create table:%v", &model.Device_data{})
	if !db.HasTable(&model.Device_data{}) {
		logrus.Errorf("no table:%v", &model.Device_data{})
		db.CreateTable(&model.Device_data{})
	}
}

func openDB() {
	logrus.Warnln("open sqlite db")
	closeDB()
	var err error
	dbLock.Lock()
	defer dbLock.Unlock()
	db, err = gorm.Open("sqlite3", "go-comm-mqtt.db")
	if nil != err {
		logrus.Warnln("open sqlite db err:", err)
		db = nil
		return
	}
}

func closeDB() {
	logrus.Warnln("close sqlite db")
	dbLock.Lock()
	defer dbLock.Unlock()
	var err error
	if db != nil {
		err = db.Close()
		if nil != err {
			logrus.Warnln("close sqlite db err:", err)
		}
	}
}

func isOpenDB() bool {
	if nil == db {
		closeDB()
		openDB()
	}
	if nil == db {
		logrus.Warnln("sqlite db is nil")
		return false
	}
	return true
}

// AddDeviceData 添加设备数据
func AddDeviceData(deviceData model.Device_data) bool {
	if !isOpenDB() {
		logrus.Warnln("add new device_data err db not open")
		return false
	}
	dbLock.Lock()
	defer dbLock.Unlock()
	//logrus.Info("add new device_data:", deviceData)
	err := db.Create(&deviceData).Error
	if nil != err {
		logrus.Warnln("create new device_data err:", err)
	}
	return true
}

// GetLastDeviceData 获取最后一条设备数据并返回
func GetLastDeviceData() (*model.Device_data, error) {
	deviceData := model.Device_data{}
	if !isOpenDB() {
		logrus.Warnln("get last device_data err db not open")
		return nil, nil
	}
	dbLock.Lock()
	defer dbLock.Unlock()
	err := db.Last(&deviceData).Error
	if nil != err {
		logrus.Warnln("get last device_data err:", err)
		return nil, err
	}
	return &deviceData, err
}
