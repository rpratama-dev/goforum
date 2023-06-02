package models

import "github.com/google/uuid"

type QuestionVote struct {
	BaseModelID
	QuestionID  		uuid.UUID	`gorm:"type:uuid;index;not null" json:"question_id"`
	Question    		Question	`gorm:"foreignKey:QuestionID" json:"question"`
	UserID      		uuid.UUID	`gorm:"type:uuid;index;not null" json:"user_id"`
	User        		User      `gorm:"foreignKey:UserID" json:"user"`
	VoteType 				string 		`gorm:"type:varchar(3);default:'up';check:vote_type IN ('up', 'down')"`
	BaseModelAudit
}
