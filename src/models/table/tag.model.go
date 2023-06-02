package models

type Tag struct {
	BaseModelID
	Name						string `json:"name" form:"name" gorm:"unique"`
	BaseModelAudit
}
