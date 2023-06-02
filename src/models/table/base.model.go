package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModelID struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid,default:uuid_generate_v4(),primaryKey"`
}

type BaseModelAudit struct {
	IsActive bool `json:"is_active"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp without time zone"`
	CreatedBy *uuid.UUID `json:"-" gorm:"type:uuid;default:null"`
	CreatedName string `json:"created_name" gorm:"default:null"`
	CreatedFrom string `json:"-" gorm:"default:null"`
	UpdatedAt time.Time `json:"-" gorm:"type:timestamp without time zone;"`
	UpdatedBy *uuid.UUID `json:"-" gorm:"type:uuid;default:null"`
	UpdatedName string `json:"-" gorm:"default:null"`
	UpdatedFrom string `json:"-" gorm:"default:null"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"type:timestamp without time zone;index;default:null"`
	DeletedBy *uuid.UUID `json:"-" gorm:"type:uuid;default:null"`
	DeletedName string `json:"-" gorm:"default:null"`
	DeletedFrom string `json:"-" gorm:"default:null"`
}
