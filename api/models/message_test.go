package models_test

import (
	"api/models"
	"testing"

	"reflect"

	"github.com/stretchr/testify/assert"
)

func TestInitMessage(t *testing.T) {
	t.Run("returns instance of message", func(t *testing.T) {
		m := models.InitMessage()
		assert.NotNil(t, m)
	})
}

func TestMes_GetBody(t *testing.T) {
	t.Run("returns message body value", func(t *testing.T) {
		body := "Body"

		m := models.InitMessage()
		reflect.ValueOf(m).Elem().FieldByName("Body").SetString(body)
		assert.Equal(t, body, m.GetBody())
	})
}

func TestMes_GetOriginator(t *testing.T) {
	t.Run("returns message originator value", func(t *testing.T) {
		originator := "originator"

		m := models.InitMessage()
		reflect.ValueOf(m).Elem().FieldByName("Originator").SetString(originator)
		assert.Equal(t, originator, m.GetOriginator())
	})
}

func TestMes_GetRecipient(t *testing.T) {
	t.Run("returns message recipient value", func(t *testing.T) {
		rec := int64(123)

		m := models.InitMessage()
		reflect.ValueOf(m).Elem().FieldByName("Recipient").SetInt(rec)
		assert.Equal(t, rec, m.GetRecipient())
	})
}
