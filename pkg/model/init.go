package model

import (
	"github.com/sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"dice/pkg/global"
	"dice/pkg/utils"
	"time"
)

var Gdb *gorm.DB

func InitDatabase() {
	config := global.Config()
	var err error
	for Gdb, err = gorm.Open("sqlite3", utils.GetRunPath()+config.DB.Dsn); err != nil; {
		logrus.Fatalln("DB Connect Error:" + err.Error())
		logrus.Warningln("Ready to retry in 5 seconds")
		time.Sleep(5 * time.Second)
		Gdb, err = gorm.Open("mysql", config.DB.Dsn)
	}
	Gdb.SingularTable(true)
	Gdb.DB().SetMaxOpenConns(config.DB.MaxOpen)
	Gdb.DB().SetMaxIdleConns(config.DB.MaxIdle)
	if global.Debug {
		Gdb.LogMode(true)
	}
	logrus.Infoln("db connected")
}
