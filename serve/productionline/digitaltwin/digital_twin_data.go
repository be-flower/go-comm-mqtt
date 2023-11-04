package digitaltwin

import (
	"github.com/sirupsen/logrus"
	"go-comm-mqtt/data/db"
	"go-comm-mqtt/data/modol"
)

func CreateHistoryTable() {
	table := db.DB.Migrator().HasTable(&modol.DigitalTwinHistoryTimes{})
	if !table {
		err := db.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&modol.DigitalTwinHistoryTimes{})
		if err != nil {
			logrus.Fatal("failed to create table ")
		}
	}
}
