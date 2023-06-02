package models

import "github.com/google/uuid"

type AnswerComment struct {
	BaseModelID
	Content					string
	AnswerID  			uuid.UUID	`gorm:"type:uuid;index;not null" json:"answer_id"`
	Answer    			Answer	`gorm:"foreignKey:AnswerID" json:"answer"`
	UserID      		uuid.UUID	`gorm:"type:uuid;index;not null" json:"user_id"`
	User        		User      `gorm:"foreignKey:UserID" json:"user"`
	BaseModelAudit
}
