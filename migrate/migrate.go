package migrate

import (
	"dns-check/database"
	"dns-check/model"
)

func DoMigrate() {
	db := database.GetInstance()
	db.AutoMigrate(&model.Domain{})
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Job{})
}
