package models

type CountryModel struct {
	BaseModel
	Name string `json:"name" validate:"required"`
	Code string `json:"code" validate:"required" gorm:"unique"`
}
