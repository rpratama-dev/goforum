package utils

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rpratama-dev/mymovie/src/variables"
)


type ErrorResponse struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func ValidatePhoneNumber(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()
	match, _ := regexp.MatchString(variables.REGEXP_PHONE_NUMBER, phoneNumber)
	return match
}

func ValidateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	match, _ := regexp.MatchString(`[a-z]`, password)
	if !match {
		return false
	}

	match, _ = regexp.MatchString(`[A-Z]`, password)
	if !match {
		return false
	}

	match, _ = regexp.MatchString(`\d`, password)
	if !match {
		return false
	}

	match, _ = regexp.MatchString(`[@#$%^&+=]`, password)
	if !match {
		return false
	}

	if len(password) < 8 {
		return false
	}

	return true
}

func ValidateUUID(fl validator.FieldLevel) bool {
	str := fl.Field().Interface().(string)
	_, err := uuid.Parse(str)
	return err == nil
}

func ValidateUUIDs(fl validator.FieldLevel) bool {
	slice := fl.Field().Interface().([]string)
	for _, str := range slice {
		_, err := uuid.Parse(str)
		if err != nil {
			return false
		}
	}
	return true
}

func ParseErrors(validationErrors validator.ValidationErrors) []ErrorResponse {
	errorResponses := make([]ErrorResponse, len(validationErrors))
	for i, ve := range validationErrors {
		errorResponses[i] = ErrorResponse{
			Field: ve.StructField(),
			Error: strings.Split(ve.Error(), "Error:")[1],
		}
	}
	return errorResponses
}
