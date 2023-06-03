package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rpratama-dev/goforum/src/utils"
)

type AnswerVotePayload struct {
	QuestionID 			string		`json:"question_id" form:"question_id" validate:"uuid,required"`
	AnswerID 				string		`json:"answer_id" form:"answer_id" validate:"uuid,required"`
	VoteType 				string 		`json:"vote" form:"vote" validate:"vote_type,required" gorm:"type:varchar(3);default:'up';check:vote_type IN ('up', 'down')"`
}

type AnswerVote struct {
	BaseModelID
	VoteType 				string 			`gorm:"type:varchar(4);default:'up';check:vote_type IN ('up', 'down')" json:"vote"`
	AnswerID  			uuid.UUID		`gorm:"type:uuid;index;not null" json:"answer_id"`
	Answer    			*Answer			`gorm:"foreignKey:AnswerID" json:"answer,omitempty"`
	UserID      		uuid.UUID		`gorm:"type:uuid;index;not null" json:"user_id"`
	User        		*User      	`gorm:"foreignKey:UserID" json:"user,omitempty"`
	BaseModelAudit
}

func (a *AnswerVotePayload) Validate() []utils.ErrorResponse {
	return voteAnswerValidate(a)
}

func (a *AnswerVote) Append(payload AnswerVotePayload, session Session, apiKey string) {
	answerId, _ := uuid.Parse(payload.AnswerID)
	a.VoteType = payload.VoteType
	a.AnswerID = answerId
	a.UserID = session.User.ID
	a.IsActive = true
	a.CreatedBy = &session.User.ID
	a.CreatedName = session.User.FullName
	a.CreatedFrom = apiKey
}

func voteAnswerValidate(s interface{}) []utils.ErrorResponse {
	validate := validator.New()
	validate.RegisterValidation("uuid", utils.ValidateUUID)
	validate.RegisterValidation("vote_type", utils.ValidateVoteType)
	err := validate.Struct(s)
	if err != nil {
		return utils.ParseErrors(err.(validator.ValidationErrors))
	}
	return nil
}
