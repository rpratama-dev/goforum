package models

type Artist struct {
	BaseModelID	
	Code string `json:"code" validate:"required" gorm:"unique"`
	Name string `json:"name"`
	Poster string `json:"poster"`
	BaseModelAudit
}
