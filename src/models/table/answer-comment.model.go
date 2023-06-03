package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rpratama-dev/goforum/src/utils"
)

type AnswerCommentPayload struct {
	Content					string	`json:"content" form:"content" validate:"required,min=10"`
	QuestionID  		string	`json:"question_id" form:"question_id" validate:"uuid,required"`
	AnswerID  			string	`json:"answer_id" form:"answer_id" validate:"uuid,required"`
}

type AnswerCommentPayloadUpdate struct {
	AnswerCommentPayload
	CommentID  			string	`json:"comment_id" form:"comment_id" validate:"uuid,required"`
}

type AnswerComment struct {
	BaseModelID
	Content					string
	AnswerID  			uuid.UUID	`gorm:"type:uuid;index;not null" json:"answer_id"`
	Answer    			*Answer		`gorm:"foreignKey:AnswerID" json:"answer,omitempty"`
	UserID      		uuid.UUID	`gorm:"type:uuid;index;not null" json:"user_id"`
	User        		*User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BaseModelAudit
}

func (a *AnswerCommentPayload) Validate() []utils.ErrorResponse {
	return answerCommentValidate(a)
}

func (a *AnswerCommentPayloadUpdate) Validate() []utils.ErrorResponse {
	return answerCommentValidate(a)
}

func (a *AnswerComment) Append(payload AnswerCommentPayload, session Session, apiKey string) {
	answerId, _ := uuid.Parse(payload.AnswerID)
	a.Content = payload.Content
	a.AnswerID = answerId
	a.UserID = session.User.ID
	a.IsActive = true
	a.CreatedBy = &session.User.ID
	a.CreatedName = session.User.FullName
	a.CreatedFrom = apiKey
}

func answerCommentValidate(s interface{}) []utils.ErrorResponse {
	validate := validator.New()
	validate.RegisterValidation("uuid", utils.ValidateUUID)
	err := validate.Struct(s)
	if err != nil {
		return utils.ParseErrors(err.(validator.ValidationErrors))
	}
	return nil
}
