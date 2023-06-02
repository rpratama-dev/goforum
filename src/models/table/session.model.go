package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseSession struct {
	UserID      uuid.UUID	`gorm:"type:uuid;index;not null" json:"user_id"`
	IPAddress   string    `gorm:"not null" json:"ip_address"`
	UserAgent   string    `gorm:"not null" json:"user_agent"`
}

type SessionPayload struct {
	FullName 		string		`gorm:"default:null" json:"full_name"` 
	BaseSession
}

type Session struct {
	BaseModelID
	BaseSession
	ExpiredAt   time.Time `gorm:"type:timestamp without time zone;default:null" json:"expired_at"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
	BaseModelAudit
}

func (s *Session) Append(payload SessionPayload)  {
	s.UserID = payload.UserID;
	s.IPAddress = payload.IPAddress;
	s.UserAgent = payload.UserAgent;
	s.ExpiredAt = time.Now().Add(1 * time.Hour);
	s.CreatedBy = &payload.UserID
	s.CreatedName = payload.FullName
	s.IsActive = true
}
