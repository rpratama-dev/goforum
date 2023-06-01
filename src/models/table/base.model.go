package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModelID struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}

type BaseModelAudit struct {
	IsActive bool `json:"is_active"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp without time zone;"`
	CreatedBy *uuid.UUID `json:"created_by" gorm:"type:uuid;default:null"`
	CreatedName string `json:"created_name" gorm:"default:null"`
	CreatedFrom string `json:"created_from" gorm:"default:null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp without time zone;"`
	UpdatedBy *uuid.UUID `json:"updated_by" gorm:"type:uuid;default:null"`
	UpdatedName string `json:"updated_name" gorm:"default:null"`
	UpdatedFrom string `json:"updated_from" gorm:"default:null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp without time zone;index;default:null"`
	DeletedBy *uuid.UUID `json:"deleted_by" gorm:"type:uuid;default:null"`
	DeletedName string `json:"deleted_name" gorm:"default:null"`
	DeletedFrom string `json:"deleted_from" gorm:"default:null"`
}
