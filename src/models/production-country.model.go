package models

type ProductionCountry struct {
	BaseModel
	Code string `json:"code" gorm:"unique"`
	Name string `json:"name"`
}
