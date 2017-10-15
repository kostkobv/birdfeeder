package utils

import (
	"config"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
	"regexp"
	"strconv"
	"strings"
)

var trans ut.Translator

var msisdnRegex = regexp.MustCompile(`^[1-9]\d{5,14}$`)              // first symbol is number between 1 and 9; 6 to 15 digits
var textoriginatorRegex = regexp.MustCompile(`^[\p{L}\p{N}]{1,11}$`) // alphanumeric unicode string between 1 and 11 symbols

// msisdnValidator checks if passed data is valid msisdn
func msisdnValidator(fl validator.FieldLevel) bool {
	v := fl.Field()
	t := v.Type().Name()

	switch t {
	case "int64":
		return msisdnRegex.MatchString(strconv.FormatInt(v.Int(), 10))
	case "string":
		return msisdnRegex.MatchString(v.String())
	default:
		return false
	}
}

// textoriginatorValidator checks if value is valid alphanumeric originator
func textoriginatorValidator(fl validator.FieldLevel) bool {
	v := fl.Field()

	if v.Type().Name() != "string" {
		return false
	}

	return textoriginatorRegex.MatchString(v.String())
}

// CustomValidator is interface for validator that matches echo.Validator
type CustomValidator interface {
	Validate(i interface{}) error
}

type cValidator struct {
	validator *validator.Validate
	trans     ut.Translator
}

// Validate the provided struct
func (v *cValidator) Validate(i interface{}) error {
	err := v.validator.Struct(i)

	if err != nil {
		return err
	}

	return nil
}

// HumaniseValidationErrors makes errors less detailed
func HumaniseValidationErrors(err error) map[string]string {
	e := map[string]string{}

	errs := err.(validator.ValidationErrors)
	for _, val := range errs {
		key := strings.ToLower(val.Field())
		e[key] = val.Translate(trans)
	}

	return e
}

// RegisterCustomTranslations which would be readable for end-users
func (v *cValidator) RegisterCustomTranslations() {
	for key, text := range config.ValidationMessages {
		v.registerTranslation(key, text)
	}
}

func (v *cValidator) registerTranslation(tag string, text string) {
	v.validator.RegisterTranslation(tag, v.trans, func(ut ut.Translator) error {
		return ut.Add(tag, text, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())

		return t
	})
}

// InitValidator is the CustomValidator factory method
func InitValidator() CustomValidator {
	en := en.New()
	uni := ut.New(en, en)

	trans, _ = uni.GetTranslator("en")

	v := validator.New()
	v.RegisterValidation("msisdn", msisdnValidator)
	v.RegisterValidation("textoriginator", textoriginatorValidator)

	val := &cValidator{v, trans}
	val.RegisterCustomTranslations()

	return val
}
