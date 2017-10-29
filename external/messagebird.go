package external

import (
	"time"
	"utils"

	mb "github.com/messagebird/go-rest-api"
)

// MessageBirdClient is an overlay for MessageBird Client
type MessageBirdClient interface {
	// NewMessage submits message to the MessageBird API
	NewMessage(originator string, recipients []string, body string, msgParams *mb.MessageParams) (*mb.Message, error)
}

// InitMessageBirdClient is a factory method for MessageBirdClient
func InitMessageBirdClient(key string) MessageBirdClient {
	return mb.New(key)
}

// InitMessageBirdParams is a factory method for MessageBird MessageParams
func InitMessageBirdParams(dc utils.Datacoding, udh string) *mb.MessageParams {
	td := mb.TypeDetails{}

	if udh != "" {
		td["udh"] = udh
	}

	return &mb.MessageParams{
		Type:              "binary",
		Reference:         "",
		Validity:          0,
		Gateway:           0,
		TypeDetails:       td,
		DataCoding:        string(dc),
		ScheduledDatetime: time.Time{},
	}
}
