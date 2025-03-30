package database

import (
	"fmt"
	"information-release/util"
	"log"
	"os"
	"path"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	globalDB *gorm.DB
)

func ConnectDB(confDir, confFile, fileType, logDir string) {
	viper := util.InitViper(confDir, confFile, fileType)

	db := viper.GetString("mysql.user")
	pass := viper.GetString("mysql.pass")
	host := viper.GetString("mysql.host")
	port := viper.GetString("mysql.port")
	dbname := viper.GetString("mysql.dbname")
	logFileName := viper.GetString("mysql.log")
	logFile, _ := os.OpenFile(path.Join(logDir, logFileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)

	DataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)

	newLogger := logger.New(
		log.New(logFile, "\r\n", log.LstdFlags), // io writer，可以输出到文件，也可以输出到os.Stdout
		logger.Config{
			SlowThreshold:             100 * time.Millisecond, //耗时超过此值认定为慢查询
			LogLevel:                  logger.Info,            // LogLevel的最低阈值，Silent为不输出日志
			IgnoreRecordNotFoundError: true,                   // 忽略RecordNotFound这种错误日志
			Colorful:                  false,                  // 禁用颜色
		},
	)

	db, err := gorm.Open(mysql.Open(DataSourceName), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	globalDB = db
}
