package migration

import (
	models "github.com/rpratama-dev/mymovie/src/models/table"
	"github.com/rpratama-dev/mymovie/src/services/database"
)

func Migrate() {
	database.Conn.AutoMigrate(&models.User{})
}
