package models

import "github.com/google/uuid"

type QuestionTag struct {
	BaseModelID
	QuestionID  		uuid.UUID	`gorm:"type:uuid;index;not null" json:"question_id"`
	Question    		Question	`gorm:"foreignKey:QuestionID" json:"question"`
	TagID  					uuid.UUID	`gorm:"type:uuid;index;not null" json:"tag_id"`
	Tag    					Tag	`gorm:"foreignKey:TagID" json:"tag"`
	BaseModelAudit
}
