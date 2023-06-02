package models

type TagPayload struct {
	Name						string `json:"name" form:"name" gorm:"unique" validate:"required"`
}

type Tag struct {
	BaseModelID
	TagPayload
	BaseModelAudit
}
