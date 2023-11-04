package db

import (
	"gorm.io/gorm"
)

var DB *gorm.DB

//func init() {
//	var err error
//	DB, err = gorm.Open(sqlite.Open("comm.db"), &gorm.Config{})
//	if err != nil {
//		logrus.Fatal("failed to connect database")
//	}
//}
