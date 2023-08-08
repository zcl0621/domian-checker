package database

import (
	"dns-check/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

var DB *gorm.DB

func GetInstance() *gorm.DB {
	if DB == nil {
		connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", config.Conf.DataBase.Host, config.Conf.DataBase.Port, config.Conf.DataBase.User, config.Conf.DataBase.DBName, config.Conf.DataBase.Password, "disable")
		DB, _ = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
		sqlDB, err := DB.DB()
		if err != nil {
			log.Panic(err.Error())
			return nil
		}
		sqlDB.SetMaxOpenConns(1024)
		sqlDB.SetMaxIdleConns(512)
		sqlDB.SetConnMaxIdleTime(time.Minute)
		sqlDB.SetConnMaxLifetime(time.Minute)
	}
	if config.RunMode == "debug" || config.RunMode == "dev" {
		DB = DB.Debug()
	}
	return DB
}
