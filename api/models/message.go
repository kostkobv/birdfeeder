package models

// Message interface
type Message interface {
}

type mes struct {
	Recipient  int64  `json:"recipient" validate:"required,msisdn"`
	Originator string `json:"originator" validate:"required,textoriginator|msisdn"`
	Message    string `json:"message" validate:"required"`
}

// InitMessage is a Message factory method
func InitMessage() Message {
	return &mes{}
}
