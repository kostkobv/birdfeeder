package utils_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"utils"
)

func TestInitValidator(t *testing.T) {
	v := utils.InitValidator()
	t.Run("validator should not be empty", func(t *testing.T) {
		assert.NotEmpty(t, v)
	})
}

func TestCValidator_Validate(t *testing.T) {
	t.Run("should validate mdisdn for both int64 and string", func(t *testing.T) {
		t.Run("should be between 6 to 14 symbols long", func(t *testing.T) {
			type vStruct struct {
				A int64  `validate:"msisdn"`
				B string `validate:"msisdn"`
			}

			t.Run("less symbols", func(t *testing.T) {
				v := utils.InitValidator()
				err := v.Validate(&vStruct{12345, "12345"})
				assert.NotNil(t, err)
			})

			t.Run("more symbols", func(t *testing.T) {
				v := utils.InitValidator()
				err := v.Validate(&vStruct{1234567890123456, "1234567890123456"})
				assert.NotNil(t, err)
			})

			t.Run("right amount of symbols", func(t *testing.T) {
				v := utils.InitValidator()
				err := v.Validate(&vStruct{33617129674, "33617129674"})
				assert.Nil(t, err)
			})
		})

		t.Run("should not start with zero", func(t *testing.T) {
			type vStruct struct {
				A string `validate:"msisdn"`
			}

			v := utils.InitValidator()
			err := v.Validate(&vStruct{"02345678901234"})
			assert.NotNil(t, err)
		})
	})

	t.Run("should validate textoriginator", func(t *testing.T) {
		v := utils.InitValidator()

		type vStruct struct {
			A string `validate:"textoriginator"`
		}

		t.Run("should be alphanumeric unicode string between 1 and 11 symbols", func(t *testing.T) {
			t.Run("valid", func(t *testing.T) {
				assert.Nil(t, v.Validate(vStruct{"Hello123"}))
			})

			t.Run("too short", func(t *testing.T) {
				assert.NotNil(t, v.Validate(vStruct{""}))
			})

			t.Run("too long", func(t *testing.T) {
				assert.NotNil(t, v.Validate(vStruct{"Hello12345678"}))
			})

			t.Run("not alphanumeric", func(t *testing.T) {
				assert.NotNil(t, v.Validate(vStruct{"Hello!"}))
			})

			t.Run("unicode", func(t *testing.T) {
				assert.Nil(t, v.Validate(vStruct{"Привет"}))
			})

			t.Run("unicode with not only alphanumeric", func(t *testing.T) {
				assert.NotNil(t, v.Validate(vStruct{"Привет!"}))
			})
		})
	})

	t.Run("msisdn|textoriginator", func(t *testing.T) {
		v := utils.InitValidator()

		type vStruct struct {
			A string `validate:"textoriginator|msisdn"`
		}

		t.Run("valid msisdn", func(t *testing.T) {
			assert.Nil(t, v.Validate(vStruct{"33123123123"}))
		})

		t.Run("valid textoriginator", func(t *testing.T) {
			assert.Nil(t, v.Validate(vStruct{"HelloПривет"}))
		})

		t.Run("not valid textoriginator", func(t *testing.T) {
			assert.NotNil(t, v.Validate(vStruct{"HelloПривет!"}))
		})

		t.Run("not valid msisdn but valid textoriginator", func(t *testing.T) {
			assert.Nil(t, v.Validate(vStruct{"03312312312"}))
		})

		t.Run("not valid msisdn", func(t *testing.T) {
			assert.NotNil(t, v.Validate(vStruct{"033123123121"}))
		})
	})
}

func TestHumaniseValidationErrors(t *testing.T) {
	t.Run("msisdn error", func(t *testing.T) {
		v := utils.InitValidator()

		type vStruct struct {
			A string `validate:"msisdn"`
		}

		err := utils.HumaniseValidationErrors(v.Validate(vStruct{"0"}))

		assert.Equal(t, err, map[string]string{"a": "should be a valid MSISDN"})
	})

	t.Run("textoriginator error", func(t *testing.T) {
		v := utils.InitValidator()

		type vStruct struct {
			A string `validate:"textoriginator"`
		}

		e := v.Validate(vStruct{""})

		err := utils.HumaniseValidationErrors(e)

		assert.Equal(t, err, map[string]string{"a": "use alphanumeric value (max. 11 symbols long)"})
	})

	t.Run("required error", func(t *testing.T) {
		v := utils.InitValidator()

		type vStruct struct {
			A string `validate:"required"`
		}

		e := v.Validate(vStruct{""})

		err := utils.HumaniseValidationErrors(e)

		assert.Equal(t, err, map[string]string{"a": "must have a value"})
	})

	t.Run("required error", func(t *testing.T) {
		v := utils.InitValidator()

		type vStruct struct {
			A string `validate:"required"`
		}

		e := v.Validate(vStruct{""})

		err := utils.HumaniseValidationErrors(e)

		assert.Equal(t, err, map[string]string{"a": "must have a value"})
	})

	t.Run("textoriginator|msisdn error", func(t *testing.T) {
		v := utils.InitValidator()

		type vStruct struct {
			A string `validate:"textoriginator|msisdn"`
		}

		e := v.Validate(vStruct{""})

		err := utils.HumaniseValidationErrors(e)

		assert.Equal(t, err, map[string]string{"a": "use valid MSISDN or alphanumeric value (max. 11 symbols long)"})
	})
}
