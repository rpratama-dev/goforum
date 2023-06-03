package models

import (
	"fmt"

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
	Question    	*Question					`json:"question,omitempty" gorm:"foreignKey:QuestionID"`
	UserID      	uuid.UUID					`json:"user_id" gorm:"type:uuid;index;not null"`
	User        	*User      				`json:"user,omitempty" gorm:"foreignKey:UserID"`
	Comments 			*[]AnswerComment	`json:"comments,omitempty"`
	Votes 				*[]AnswerVote			`json:"votes,omitempty"`
	Score					int32						`json:"score" gorm:"-"`
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

func (q *Answer) CalculateScore() {
	var upVotes, downVotes int32
	if q.Votes != nil {
		for _, vote := range *q.Votes {
			fmt.Println("VoteType", vote.VoteType)
			if vote.VoteType == "up" {
				upVotes++
			} else if vote.VoteType == "down" {
				downVotes++
			}
		}
	}
	// Calculate score based on the number of upVotes and downVotes
	score := (upVotes * 5) - (downVotes * 2)
	if (score < 0) {
		score = 0
	}
	fmt.Println("VoteType score", score)
	q.Score = score
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
