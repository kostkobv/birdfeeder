package external_test

import (
	"external"
	"reflect"
	"testing"

	"time"

	"utils"

	mb "github.com/messagebird/go-rest-api"
	"github.com/stretchr/testify/assert"
)

func TestInitMessageBirdClient(t *testing.T) {
	t.Run("returns new MessageBird client instance with provided key", func(t *testing.T) {
		key := "key"
		mb := external.InitMessageBirdClient(key)
		assert.Equal(t, key, reflect.ValueOf(mb).Elem().FieldByName("AccessKey").String())
	})
}

func TestInitMessageBirdParams(t *testing.T) {
	t.Run("returns struct with empty udh if udh is empty string", func(t *testing.T) {
		dc := "plain"

		params := external.InitMessageBirdParams(utils.Datacoding(dc), "")
		expected := &mb.MessageParams{
			Type:              "binary",
			Reference:         "",
			Validity:          0,
			Gateway:           0,
			TypeDetails:       mb.TypeDetails{},
			DataCoding:        dc,
			ScheduledDatetime: time.Time{},
		}

		assert.Equal(t, expected, params)
	})

	t.Run("returns struct with udh if udh not empty string", func(t *testing.T) {
		dc := "plain"
		udh := "udh"

		params := external.InitMessageBirdParams(utils.Datacoding(dc), udh)
		expected := &mb.MessageParams{
			Type:              "binary",
			Reference:         "",
			Validity:          0,
			Gateway:           0,
			TypeDetails:       mb.TypeDetails{"udh": udh},
			DataCoding:        dc,
			ScheduledDatetime: time.Time{},
		}

		assert.Equal(t, expected, params)
	})
}
