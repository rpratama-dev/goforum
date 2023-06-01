package models

type GenreModel struct {
	BaseModel
	Name string `json:"name" gorm:"unique" validate:"required"`
}
