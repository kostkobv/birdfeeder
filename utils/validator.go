package utils

import (
	"github.com/go-playground/validator"
	"regexp"
	"strconv"
)

// ^[1-9]\d{6,14}$
func msisdnValidator(fl validator.FieldLevel) bool {
	msisdnRegex := regexp.MustCompile(`^[1-9]\d{6,14}$`)
	s := strconv.FormatInt(fl.Field().Int(), 10)
	return msisdnRegex.MatchString(s)
}

type CustomValidator interface {
	Validate(i interface{}) error
}

type cValidator struct {
	validator *validator.Validate
}

func (v *cValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

// Init the validator instance
func Init() CustomValidator {
	v := validator.New()

	v.RegisterValidation("msisdn", msisdnValidator)

	return &cValidator{v}
}
