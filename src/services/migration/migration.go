package migration

import (
	models "github.com/rpratama-dev/mymovie/src/models/table"
	"github.com/rpratama-dev/mymovie/src/services/database"
)

func Migrate() {
	database.Conn.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.Tag{},
		&models.Question{},
		&models.QuestionTag{},
		&models.QuestionComment{},
		&models.QuestionVote{},
		&models.Answer{},
		&models.AnswerComment{},
		&models.AnswerVote{},
	)
}
