package models_test

import (
	apiModels "api/models"
	"queue/models"
	"reflect"
	"testing"
	"utils"

	"strconv"

	"github.com/stretchr/testify/assert"
)

func TestInitQueueMessage(t *testing.T) {
	body := "body"
	dc := "plain"
	rm := apiModels.InitMessage()
	udh := "udh"

	m := models.InitQueueMessage(body, utils.Datacoding(dc), rm, udh)
	reflected := reflect.ValueOf(m).Elem()

	t.Run("new message keeps reference to api model", func(t *testing.T) {
		assert.Equal(t, rm, reflected.FieldByName("OriginalMessage").Interface())
	})

	t.Run("new message keeps the same passed datacoding", func(t *testing.T) {
		assert.Equal(t, dc, reflected.FieldByName("Encoding").String())
	})

	t.Run("new message keeps the same passed udh", func(t *testing.T) {
		assert.Equal(t, udh, reflected.FieldByName("UDH").String())
	})

	t.Run("new message keeps the same passed message body", func(t *testing.T) {
		assert.Equal(t, body, reflected.FieldByName("Message").String())
	})
}

func TestByRecipientsAmount_Len(t *testing.T) {
	body := "body"
	dc := "plain"
	rm := apiModels.InitMessage()
	udh := "udh"

	m1 := models.InitQueueMessage(body, utils.Datacoding(dc), rm, udh)
	m2 := models.InitQueueMessage(body, utils.Datacoding(dc), rm, udh)

	c := models.ByRecipientsAmount{m1}
	assert.Equal(t, 1, c.Len())
	c = append(c, m2)
	assert.Equal(t, 2, c.Len())
}

func TestByRecipientsAmount_Less(t *testing.T) {
	body := "body"
	dc := "plain"
	rm := apiModels.InitMessage()
	udh := "udh"

	m1 := models.InitQueueMessage(body, utils.Datacoding(dc), rm, udh)
	m2 := models.InitQueueMessage(body, utils.Datacoding(dc), rm, udh)

	m1.AddRecipient(0)

	m2.AddRecipient(0)
	m2.AddRecipient(1)

	c := models.ByRecipientsAmount{m1, m2}

	t.Run("true if second provided message has less recipients", func(t *testing.T) {

		assert.Equal(t, true, c.Less(0, 1))
	})

	t.Run("false if first provided message has less recipients", func(t *testing.T) {
		c := models.ByRecipientsAmount{m1, m2}

		assert.Equal(t, false, c.Less(1, 0))
	})
}

func TestByRecipientsAmount_Swap(t *testing.T) {
	body := "body"
	dc := "plain"
	rm := apiModels.InitMessage()
	udh := "udh"

	m1 := models.InitQueueMessage(body, utils.Datacoding(dc), rm, udh)
	m2 := models.InitQueueMessage(body, utils.Datacoding(dc), rm, udh)

	t.Run("swaps elements in collection", func(t *testing.T) {
		c := models.ByRecipientsAmount{m1, m2}
		expected := models.ByRecipientsAmount{m2, m1}
		c.Swap(0, 1)

		assert.Equal(t, expected, c)
	})
}

func TestQMessage_AddRecipient(t *testing.T) {
	body := "body"
	dc := "plain"
	rm := apiModels.InitMessage()
	udh := "udh"

	m := models.InitQueueMessage(body, utils.Datacoding(dc), rm, udh)

	t.Run("adds recipient to the message", func(t *testing.T) {
		assert.Equal(t, m.GetRecipientsAmount(), int64(0))
		m.AddRecipient(0)
		assert.Equal(t, int64(1), m.GetRecipientsAmount())
	})

	t.Run("doesn't add the same recipients twice", func(t *testing.T) {
		assert.Equal(t, m.GetRecipientsAmount(), int64(1))
		m.AddRecipient(0)
		assert.Equal(t, int64(1), m.GetRecipientsAmount())
	})
}

func TestQMessage_GetDataCoding(t *testing.T) {
	body := "body"
	dc := "plain"
	rm := apiModels.InitMessage()
	udh := "udh"

	m := models.InitQueueMessage(body, utils.Datacoding(dc), rm, udh)

	t.Run("returns provided datacoding", func(t *testing.T) {
		assert.Equal(t, utils.Datacoding(dc), m.GetDataCoding())
	})
}

func TestQMessage_GetMessage(t *testing.T) {
	body := "body"
	dc := "plain"
	rm := apiModels.InitMessage()
	udh := "udh"

	m := models.InitQueueMessage(body, utils.Datacoding(dc), rm, udh)

	t.Run("returns provided message", func(t *testing.T) {
		assert.Equal(t, body, m.GetMessage())
	})
}

func TestQMessage_GetOriginalRecipient(t *testing.T) {
	body := "body"
	dc := "plain"
	rm := apiModels.InitMessage()
	udh := "udh"

	m := models.InitQueueMessage(body, utils.Datacoding(dc), rm, udh)

	t.Run("returns originally provided recipient", func(t *testing.T) {
		originalRecipient := int64(123123123)
		reflect.ValueOf(rm).Elem().FieldByName("Recipient").SetInt(originalRecipient)
		m.AddRecipient(123)
		assert.Equal(t, originalRecipient, m.GetOriginalRecipient())
	})
}

func TestQMessage_GetOriginator(t *testing.T) {
	body := "body"
	dc := "plain"
	rm := apiModels.InitMessage()
	udh := "udh"

	m := models.InitQueueMessage(body, utils.Datacoding(dc), rm, udh)

	t.Run("returns originally provided recipient", func(t *testing.T) {
		originator := "originator"
		reflect.ValueOf(rm).Elem().FieldByName("Originator").SetString(originator)
		m.AddRecipient(123)
		assert.Equal(t, originator, m.GetOriginator())
	})
}

func TestQMessage_GetRecipients(t *testing.T) {
	body := "body"
	dc := "plain"
	rm := apiModels.InitMessage()
	udh := "udh"

	m := models.InitQueueMessage(body, utils.Datacoding(dc), rm, udh)

	t.Run("returns provided recipients", func(t *testing.T) {
		original := []int64{31, 32, 33, 34}
		recipients := make([]string, len(original))

		for i, r := range []int64{31, 32, 33, 34} {
			m.AddRecipient(r)
			recipients[i] = strconv.FormatInt(r, 10)
		}

		assert.Equal(t, recipients, m.GetRecipients())
	})
}

func TestQMessage_GetRecipientsAmount(t *testing.T) {
	body := "body"
	dc := "plain"
	rm := apiModels.InitMessage()
	udh := "udh"

	m := models.InitQueueMessage(body, utils.Datacoding(dc), rm, udh)

	t.Run("returns provided recipients", func(t *testing.T) {
		original := []int64{31, 32, 33, 34}

		for _, r := range []int64{31, 32, 33, 34} {
			m.AddRecipient(r)
		}

		assert.Equal(t, int64(len(original)), m.GetRecipientsAmount())
	})
}

func TestQMessage_GetUDH(t *testing.T) {
	body := "body"
	dc := "plain"
	rm := apiModels.InitMessage()
	udh := "udh"

	m := models.InitQueueMessage(body, utils.Datacoding(dc), rm, udh)

	t.Run("returns originally provided udh", func(t *testing.T) {
		assert.Equal(t, udh, m.GetUDH())
	})
}
