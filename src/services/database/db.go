package database

import (
	"github.com/rpratama-dev/mymovie/src/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Conn *gorm.DB

func StartDB() (bool) {
	dsn := "host=" + configs.Env.DBHost +
		" port=" + configs.Env.DBPort +
		" dbname=" + configs.Env.DBName +
		" user=" + configs.Env.DBUser +
		" password=" + configs.Env.DBPassword +
		" sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
			return false
	}
	Conn = db
	return true
}
