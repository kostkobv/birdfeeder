package mocks

import (
	"reflect"

	"github.com/go-playground/universal-translator"
	"github.com/stretchr/testify/mock"
)

// FieldErrorMock contains all functions to get error details
type FieldErrorMock struct {
	mock.Mock
}

// Tag returns the validation tag that failed. if the
// validation was an alias, this will return the
// alias name and not the underlying tag that failed.
//
// eg. alias "iscolor": "hexcolor|rgb|rgba|hsl|hsla"
// will return "iscolor"
func (e *FieldErrorMock) Tag() string {
	args := e.Called()
	return args.String(0)
}

// ActualTag returns the validation tag that failed, even if an
// alias the actual tag within the alias will be returned.
// If an 'or' validation fails the entire or will be returned.
//
// eg. alias "iscolor": "hexcolor|rgb|rgba|hsl|hsla"
// will return "hexcolor|rgb|rgba|hsl|hsla"
func (e *FieldErrorMock) ActualTag() string {
	args := e.Called()
	return args.String(0)
}

// Namespace returns the namespace for the field error, with the tag
// name taking precedence over the fields actual name.
//
// eg. JSON name "User.fname"
//
// See StructNamespace() for a version that returns actual names.
//
// NOTE: this field can be blank when validating a single primitive field
// using validate.Field(...) as there is no way to extract it's name
func (e *FieldErrorMock) Namespace() string {
	args := e.Called()
	return args.String(0)
}

// StructNamespace returns the namespace for the field error, with the fields
// actual name.
//
// eq. "User.FirstName" see Namespace for comparison
//
// NOTE: this field can be blank when validating a single primitive field
// using validate.Field(...) as there is no way to extract it's name
func (e *FieldErrorMock) StructNamespace() string {
	args := e.Called()
	return args.String(0)
}

// Field returns the fields name with the tag name taking precedence over the
// fields actual name.
//
// eq. JSON name "fname"
// see ActualField for comparison
func (e *FieldErrorMock) Field() string {
	args := e.Called()
	return args.String(0)
}

// StructField returns the fields actual name from the struct, when able to determine.
//
// eq.  "FirstName"
// see Field for comparison
func (e *FieldErrorMock) StructField() string {
	args := e.Called()
	return args.String(0)
}

// Value returns the actual fields value in case needed for creating the error
// message
func (e *FieldErrorMock) Value() interface{} {
	args := e.Called()
	return args.String(0)
}

// Param returns the param value, in string form for comparison; this will also
// help with generating an error message
func (e *FieldErrorMock) Param() string {
	args := e.Called()
	return args.String(0)
}

// Kind returns the Field's reflect Kind
//
// eg. time.Time's kind is a struct
func (e *FieldErrorMock) Kind() reflect.Kind {
	args := e.Called()
	return args.Get(0).(reflect.Kind)
}

// Type returns the Field's reflect Type
//
// // eg. time.Time's type is time.Time
func (e *FieldErrorMock) Type() reflect.Type {
	args := e.Called()
	return args.Get(0).(reflect.Type)
}

// Translate returns the FieldError's translated error
// from the provided 'ut.Translator' and registered 'TranslationFunc'
//
// NOTE: is not registered translation can be found it returns the same
// as calling fe.Error()
func (e *FieldErrorMock) Translate(ut ut.Translator) string {
	args := e.Called(ut)
	return args.String(0)
}

// Error returns the fieldError's error message
func (e *FieldErrorMock) Error() string {
	args := e.Called()
	return args.String(0)
}
