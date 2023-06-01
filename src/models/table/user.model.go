package models

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rpratama-dev/mymovie/src/services/database"
	"github.com/rpratama-dev/mymovie/src/utils"
	"gorm.io/gorm"
)

type BaseUser struct {
	FullName     string `json:"full_name" form:"full_name" validate:"required"`
	BirthDate    string `json:"birth_date" form:"birth_date" validate:"required,datetime=2006-01-02"`
	Email        string `json:"email" form:"email" validate:"required,email" gorm:"unique"`
	PhoneNumber  string `json:"phone_number" form:"phone_number" validate:"required,min=10,phone_number" gorm:"unique"`
	Password     string `json:"password" form:"password" validate:"required,min=8,strong_password"`
}

type UserLogin struct {
	Email string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type UserPayload struct {
	BaseUser
}

type User struct {
	BaseModelID
	BaseUser
	VerifiedToken string `json:"verified_token"` 
	IsVerified bool `json:"is_verified"`
	VerifiedAt *time.Time `json:"verified_at" gorm:"type:timestamp without time zone;default:null"` 
	BaseModelAudit
}

func (u *UserLogin) Validate() []utils.ErrorResponse {
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		return utils.ParseErrors(err.(validator.ValidationErrors))
	}
	return nil
}

func (u *UserLogin) IsPasswordMatch(hashedPassword string) bool  {
	return utils.IsPasswordMatch(hashedPassword, u.Password)
}

func (u *UserPayload) Validate() []utils.ErrorResponse {
	validate := validator.New()
	validate.RegisterValidation("phone_number", utils.ValidatePhoneNumber)
	validate.RegisterValidation("strong_password", utils.ValidateStrongPassword)

	err := validate.Struct(u)
	if err != nil {
		return utils.ParseErrors(err.(validator.ValidationErrors))
	}
	return nil
}

func (u *User) Append(user UserPayload) {
	u.FullName = user.FullName;
	u.BirthDate = user.BirthDate;
	u.Email = user.Email;
	u.PhoneNumber = user.PhoneNumber;
	u.Password = user.Password;
}

func (u *User) FindByEmail() error  {
	err := database.Conn.First(&u, "email = ?", strings.ToLower(u.Email)).Error
	return err

}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate a hashed password
	u.CreatedBy = &u.ID
	u.CreatedName = u.FullName
	u.CreatedFrom = "System Registration"
	u.Email = strings.ToLower(u.Email)
	u.Password = utils.HashPassword(u.Password)
	u.VerifiedToken = generateVerifiedToken()
  return
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
  if u.ID != uuid.Nil {
    tx.Model(u).Update("created_by", u.ID)
  }
	// Remove password
	u.Password = ""
  return
}

func generateVerifiedToken() string {
	// Generate random bytes
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	// Encode bytes to Base64 without slash
	randomString := base64.RawURLEncoding.EncodeToString(bytes)
	return randomString
}
