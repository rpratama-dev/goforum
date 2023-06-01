package models

type ProductionCompany struct {
	BaseModelID
	Name string `json:"name"`
	Code string `json:"code" gorm:"unique"`
	LogoPath string `json:"logo_path"`
	BaseModelAudit
}
