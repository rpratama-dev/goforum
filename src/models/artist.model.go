package models

type Artist struct {
	BaseModel
	Code string `json:"code" validate:"required" gorm:"unique"`
	Name string `json:"name"`
	Poster string `json:"poster"`
}
