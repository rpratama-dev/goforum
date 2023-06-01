package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`
	IsActive bool `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uuid.UUID `json:"created_by"`
	CreatedName string `json:"created_name"`
	CreatedFrom string `json:"created_from"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy uuid.UUID `json:"updated_by"`
	UpdatedName string `json:"updated_name"`
	UpdatedFrom string `json:"updated_from"`
}
