package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

var instance *validator.Validate

func init() {
	instance = mustNewValidator()
}

func Get() *validator.Validate {
	return instance
}

func regexpValidation(fl validator.FieldLevel) bool {
	pattern := fl.Param()
	field := fl.Field().String()
	matched, err := regexp.MatchString(pattern, field)
	if err != nil {
		return false
	}
	return matched
}

func phoneNumberValidation(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	matched, err := regexp.MatchString(`^\+[1-9]\d{1,14}$`, field)
	if err != nil {
		return false
	}
	return matched
}

func userIDValidation(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	matched, err := regexp.MatchString(`^usr-[A-Za-z0-9]+$`, field)
	if err != nil {
		return false
	}
	return matched
}

func accountNumberValidation(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	matched, err := regexp.MatchString(`^01\d{6}$`, field)
	if err != nil {
		return false
	}
	return matched
}

func transactionIDValidation(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	matched, err := regexp.MatchString(`^tan-[A-Za-z0-9]$`, field)
	if err != nil {
		return false
	}
	return matched
}

func newValidator() (*validator.Validate, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.RegisterValidation("regexp", regexpValidation)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterValidation("phone", phoneNumberValidation)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterValidation("userID", userIDValidation)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterValidation("acctNum", accountNumberValidation)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterValidation("tanID", transactionIDValidation)
	if err != nil {
		return nil, err
	}
	return validate, nil
}

func mustNewValidator() *validator.Validate {
	validate, err := newValidator()
	if err != nil {
		panic(fmt.Sprintf("MustNewValidator: %v", err))
	}
	return validate
}
