package models

type CountryModel struct {
	BaseModelID
	Name string `json:"name" validate:"required"`
	Code string `json:"code" validate:"required" gorm:"unique"`
	BaseModelAudit

}
