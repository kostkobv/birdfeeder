package models

import "utils"

// Message interface
type Message interface {
	ConvertBody() error
}

type mes struct {
	Recipient       int64  `json:"recipient" validate:"required,msisdn"`
	Originator      string `json:"originator" validate:"required,textoriginator|msisdn"`
	Body            string `json:"message" validate:"required"`
	encoder         utils.UDHEncoder
	EncodedMessages []string
	Encoding        string
}

// InitMessage is a Message factory method
func InitMessage(e utils.UDHEncoder) Message {
	return &mes{encoder: e}
}

func (m *mes) ConvertBody() error {
	encoded, err := m.encoder.Encode(m.Body)

	if err != nil {
		return err
	}

	m.EncodedMessages = encoded.Messages
	m.Encoding = string(encoded.Encoding)

	return nil
}
