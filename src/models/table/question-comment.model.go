package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rpratama-dev/mymovie/src/utils"
)

type QuestionCommentPayload struct {
	Content					string	`json:"content" form:"content" validate:"required,min=10"`
	QuestionID  		string	`json:"question_id" form:"question_id" validate:"uuid,required"`
}

type QuestionCommentPayloadUpdate struct {
	QuestionCommentPayload
	CommentID  			string	`json:"comment_id" form:"comment_id" validate:"uuid,required"`
}

type QuestionComment struct {
	BaseModelID
	Content					string		`gorm:"not null" json:"content"`
	QuestionID  		uuid.UUID	`gorm:"type:uuid;index;not null" json:"question_id"`
	Question    		Question	`gorm:"foreignKey:QuestionID" json:"question"`
	UserID      		uuid.UUID	`gorm:"type:uuid;index;not null" json:"user_id"`
	User        		User      `gorm:"foreignKey:UserID" json:"user"`
	BaseModelAudit
}

func (a *QuestionCommentPayload) Validate() []utils.ErrorResponse {
	return questionCommentValidate(a)
}

func (a *QuestionCommentPayloadUpdate) Validate() []utils.ErrorResponse {
	return questionCommentValidate(a)
}

func (a *QuestionComment) Append(payload QuestionCommentPayload, session Session, apiKey string) {
	questionId, _ := uuid.Parse(payload.QuestionID)
	a.Content = payload.Content
	a.QuestionID = questionId
	a.UserID = session.User.ID
	a.IsActive = true
	a.CreatedBy = &session.User.ID
	a.CreatedName = session.User.FullName
	a.CreatedFrom = apiKey
}

func questionCommentValidate(s interface{}) []utils.ErrorResponse {
	validate := validator.New()
	validate.RegisterValidation("uuid", utils.ValidateUUID)
	err := validate.Struct(s)
	if err != nil {
		return utils.ParseErrors(err.(validator.ValidationErrors))
	}
	return nil
}
