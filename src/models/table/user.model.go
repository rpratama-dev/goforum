package models

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rpratama-dev/mymovie/src/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type BaseUser struct {
	FullName     string `json:"full_name" form:"full_name" validate:"required"`
	BirthDate    string `json:"birth_date" form:"birth_date" validate:"required,datetime=2006-01-02"`
	Email        string `json:"email" form:"email" validate:"required,email" gorm:"unique"`
	PhoneNumber  string `json:"phone_number" form:"phone_number" validate:"required,min=10,phone_number" gorm:"unique"`
	Password     string `json:"password" form:"password" validate:"required,min=8,strong_password"`
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

type ErrorResponse struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func (u *UserPayload) Validate() []ErrorResponse {
	validate := validator.New()
	validate.RegisterValidation("phone_number", utils.ValidatePhoneNumber)
	validate.RegisterValidation("strong_password", utils.ValidateStrongPassword)

	err := validate.Struct(u)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorResponses := make([]ErrorResponse, len(validationErrors))

		for i, ve := range validationErrors {
			errorResponses[i] = ErrorResponse{
				Field: ve.StructField(),
				Error: strings.Split(ve.Error(), "Error:")[1],
			}
		}

		return errorResponses
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

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate a hashed password	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	u.CreatedBy = &u.ID;
	u.CreatedName = u.FullName;
	u.CreatedFrom = "System Registration";
	u.Password = string(hashedPassword)
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
