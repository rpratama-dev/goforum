package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rpratama-dev/goforum/src/services/database"
	"github.com/rpratama-dev/goforum/src/utils"
	"gorm.io/gorm"
)

type BaseQuestion struct {
	Title				string			`json:"title" form:"title" gorm:"not null" validate:"required,min=10"`
	Content			string			`json:"content" form:"content" gorm:"not null" validate:"required,min=10"`
}

type QuestionPayload struct {
	BaseQuestion
	Tags				[]string	`json:"tags" form:"tags" gorm:"many2many:question_tags" validate:"uuidslices,required,min=1,max=5"`
}

type Question struct {
	BaseModelID
	BaseQuestion
	Tags						*[]Tag							`json:"tags,omitempty" form:"tags" gorm:"many2many:question_tags" validate:"required,min=1,max=5"`
	UserID      		uuid.UUID						`json:"user_id" gorm:"type:uuid;index;not null"`
	User        		*User      					`json:"user,omitempty" gorm:"foreignKey:UserID"`
	Answers 				*[]Answer						`json:"answers,omitempty"`
	Comments 				*[]QuestionComment	`json:"comments,omitempty"`
	Votes 					*[]QuestionVote			`json:"votes,omitempty"`
	BaseModelAudit
}

func (qp *QuestionPayload) Validate() []utils.ErrorResponse {
	validate := validator.New()
	validate.RegisterValidation("uuidslices", utils.ValidateUUIDs)
	err := validate.Struct(qp)
	if err != nil {
		return utils.ParseErrors(err.(validator.ValidationErrors))
	}
	return nil
}

func (s *Question) SoftDelete() error {
	s.IsActive = false
	s.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	return database.Conn.Model(&s).Select(
		"is_active",
		"deleted_by",
		"deleted_at",
		"deleted_name",
		"deleted_from",
	).Updates(s).Error
}
