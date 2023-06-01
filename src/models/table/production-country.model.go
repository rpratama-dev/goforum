package models

type ProductionCountry struct {
	BaseModelID
	Code string `json:"code" gorm:"unique"`
	Name string `json:"name"`
	BaseModelAudit
}
