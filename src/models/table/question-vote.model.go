package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rpratama-dev/goforum/src/utils"
)

type QuestionVotePayload struct {
	QuestionID 			string		`validate:"uuid,required"`
	VoteType 				string 		`json:"vote" form:"vote" validate:"vote_type,required" gorm:"type:varchar(3);default:'up';check:vote_type IN ('up', 'down')"`
}

type QuestionVote struct {
	BaseModelID
	VoteType 				string 		`json:"vote" gorm:"type:varchar(4);default:'up';check:vote_type IN ('up', 'down')"`
	QuestionID  		uuid.UUID	`gorm:"type:uuid;index;not null" json:"question_id"`
	Question    		*Question	`gorm:"foreignKey:QuestionID" json:"question,omitempty"`
	UserID      		uuid.UUID	`gorm:"type:uuid;index;not null" json:"user_id"`
	User        		*User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BaseModelAudit
}

func (a *QuestionVotePayload) Validate() []utils.ErrorResponse {
	return voteValidate(a)
}

func (a *QuestionVote) Append(payload QuestionVotePayload, session Session, apiKey string) {
	questionId, _ := uuid.Parse(payload.QuestionID)
	a.VoteType = payload.VoteType
	a.QuestionID = questionId
	a.UserID = session.User.ID
	a.IsActive = true
	a.CreatedBy = &session.User.ID
	a.CreatedName = session.User.FullName
	a.CreatedFrom = apiKey
}

func voteValidate(s interface{}) []utils.ErrorResponse {
	validate := validator.New()
	validate.RegisterValidation("uuid", utils.ValidateUUID)
	validate.RegisterValidation("vote_type", utils.ValidateVoteType)
	err := validate.Struct(s)
	if err != nil {
		return utils.ParseErrors(err.(validator.ValidationErrors))
	}
	return nil
}
