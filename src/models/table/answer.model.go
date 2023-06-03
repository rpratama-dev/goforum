package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rpratama-dev/goforum/src/utils"
)

type AnswerPayload struct {
	Content				string 		`json:"content" form:"content" validate:"required,min=10"`
	QuestionID  	string		`json:"question_id" form:"question_id" validate:"uuid,required"`
}

type AnswerPayloadPatch struct {
	QuestionID  	string		`json:"question_id" form:"question_id" validate:"uuid,required"`
	AnswerID  		string		`json:"answer_id" form:"answer_id" validate:"uuid,required"`
}

type AnswerPayloadUpdate struct {
	Content				string 		`json:"content" form:"content" validate:"required,min=10" gorm:"not null"`
	AnswerPayloadPatch
}

type Answer struct {
	BaseModelID
	Content				string 						`json:"content" form:"content" gorm:"not null"`
	QuestionID  	uuid.UUID					`json:"question_id" form:"question_id" gorm:"type:uuid,index,not null"`
	IsTheBest			bool							`json:"is_the_best"`
	Question    	*Question					`gorm:"foreignKey:QuestionID" json:"question,omitempty"`
	UserID      	uuid.UUID					`gorm:"type:uuid;index;not null" json:"user_id"`
	User        	*User      				`gorm:"foreignKey:UserID" json:"user,omitempty"`
	Comments 			*[]AnswerComment	`json:"comments,omitempty"`
	Votes 				*[]AnswerVote			`json:"votes,omitempty"`
	BaseModelAudit
}

func (a *AnswerPayload) Validate() []utils.ErrorResponse {
	return answerValidate(a)
}

func (a *AnswerPayloadUpdate) Validate() []utils.ErrorResponse {
	return answerValidate(a)
}

func (a *AnswerPayloadPatch) Validate() []utils.ErrorResponse {
	return answerValidate(a)
}

func (a *Answer) Append(payload AnswerPayload, session Session, apiKey string) {
	questionId, _ := uuid.Parse(payload.QuestionID)
	a.Content = payload.Content
	a.QuestionID = questionId
	a.UserID = session.User.ID
	a.IsActive = true
	a.CreatedBy = &session.User.ID
	a.CreatedName = session.User.FullName
	a.CreatedFrom = apiKey
}

func answerValidate(s interface{}) []utils.ErrorResponse {
	validate := validator.New()
	validate.RegisterValidation("uuid", utils.ValidateUUID)
	err := validate.Struct(s)
	if err != nil {
		return utils.ParseErrors(err.(validator.ValidationErrors))
	}
	return nil
}
