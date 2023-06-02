package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rpratama-dev/mymovie/src/utils"
)

type AnswerPayload struct {
	Content				string `json:"content" form:"content" validate:"required,min=10"`
	QuestionID  	uuid.UUID	`json:"question_id" form:"question_id" gorm:"type:uuid,index,not null"`
}

type Answer struct {
	BaseModelID
	AnswerPayload
	IsTheBest			bool
	Question    	Question	`gorm:"foreignKey:QuestionID" json:"question"`
	UserID      	uuid.UUID	`gorm:"type:uuid;index;not null" json:"user_id"`
	User        	User      `gorm:"foreignKey:UserID" json:"user"`
	BaseModelAudit
}

func (a *AnswerPayload) Validate() []utils.ErrorResponse {
	validate := validator.New()
	err := validate.Struct(a)
	if err != nil {
		return utils.ParseErrors(err.(validator.ValidationErrors))
	}
	return nil
}
