package models

import "github.com/google/uuid"

type AnswerVote struct {
	BaseModelID
	AnswerID  			uuid.UUID	`gorm:"type:uuid;index;not null" json:"answer_id"`
	Answer    			Answer		`gorm:"foreignKey:AnswerID" json:"answer"`
	UserID      		uuid.UUID	`gorm:"type:uuid;index;not null" json:"user_id"`
	User        		User      `gorm:"foreignKey:UserID" json:"user"`
	VoteType 				string 			`gorm:"type:varchar(3);default:'up';check:vote_type IN ('up', 'down')"`
	BaseModelAudit
}
