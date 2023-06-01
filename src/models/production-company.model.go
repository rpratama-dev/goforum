package models

type ProductionCompany struct {
	BaseModel
	Code string `json:"code" gorm:"unique"`
	LogoPath string `json:"logo_path"`
	Name string `json:"name"`
}
