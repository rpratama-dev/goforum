package migration

import (
	models "github.com/rpratama-dev/goforum/src/models/table"
	"github.com/rpratama-dev/goforum/src/services/database"
)

func Migrate() {
	database.Conn.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.Tag{},
		&models.Question{},
		&models.QuestionComment{},
		&models.QuestionVote{},
		&models.Answer{},
		&models.AnswerComment{},
		&models.AnswerVote{},
	)
}
