package models

// Message interface
type Message interface {
	GetBody() string
	GetRecipient() int64
	GetOriginator() string
}

type mes struct {
	Recipient  int64  `json:"recipient" validate:"required,msisdn"`
	Originator string `json:"originator" validate:"required,textoriginator|msisdn"`
	Body       string `json:"message" validate:"required,max=1377"`
}

// InitMessage is a Message factory method
func InitMessage() Message {
	return &mes{}
}

// GetBody returns message body
func (m *mes) GetBody() string {
	return m.Body
}

// GetRecipient returns message recipient
func (m *mes) GetRecipient() int64 {
	return m.Recipient
}

// GetOriginator returns message originator
func (m *mes) GetOriginator() string {
	return m.Originator
}
