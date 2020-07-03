package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"opsHeart/conf"
	"opsHeart/logger"
)

var DB *gorm.DB

func InitDB() {
	dbType := conf.GetDbType()
	dbUrl := conf.GetDbUrl()

	var err error
	DB, err = gorm.Open(dbType, dbUrl)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	} else {
		logger.ServerLog.Infof("action=init db;status=success")
		DB.SingularTable(true)
	}
}
