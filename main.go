package main

import (
	"dns-check/Job"
	"dns-check/config"
	"dns-check/database"
	"dns-check/migrate"
	"dns-check/redisUtils"
	"dns-check/server"
)

func main() {
	config.InitConf()
	database.GetInstance()
	redisUtils.InitRedis()
	migrate.DoMigrate()
	Job.HandlerJob()
	err := server.StartGinServer()
	if err != nil {
		panic(err)
	}
}
