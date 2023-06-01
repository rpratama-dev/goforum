package models

type GenreModel struct {
	BaseModelID
	Name string `json:"name" gorm:"unique" validate:"required"`
	BaseModelAudit

}
