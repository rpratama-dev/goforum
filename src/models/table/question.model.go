package models

import "github.com/google/uuid"

type Question struct {
	BaseModelID
	Title				string
	Content			string
	UserID      uuid.UUID	`gorm:"type:uuid;index;not null" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
	BaseModelAudit
}
