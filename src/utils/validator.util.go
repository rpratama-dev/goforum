package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/rpratama-dev/mymovie/src/variables"
)

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
