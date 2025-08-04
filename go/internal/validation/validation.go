package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

func regexpValidation(fl validator.FieldLevel) bool {
	pattern := fl.Param()
	field := fl.Field().String()

	matched, err := regexp.MatchString(pattern, field)
	if err != nil {
		return false
	}
	return matched
}

func NewValidator() (*validator.Validate, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.RegisterValidation("regexp", regexpValidation)
	if err != nil {
		return nil, err
	}
	return validate, nil
}

func MustNewValidator() *validator.Validate {
	validate, err := NewValidator()
	if err != nil {
		panic(fmt.Sprintf("MustNewValidator: %v", err))
	}
	return validate
}
