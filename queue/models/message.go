package models

import (
	"api/models"
	"errors"
	"sort"
	"strconv"
	"utils"
)

// QueueMessage is a message that is kept within the queue
type QueueMessage interface {
	AddRecipient(r int64) error
	GetMessage() string
	GetOriginalRecipient() int64
	GetRecipientsAmount() int64
	GetOriginator() string
	GetRecipients() []string
	GetDataCoding() utils.Datacoding
	GetUDH() string
}

type qMessage struct {
	recipients      []int64
	Message         string
	Encoding        utils.Datacoding
	OriginalMessage models.Message
	UDH             string
}

// InitQueueMessage factory method to create QueueMessage
func InitQueueMessage(message string, enc utils.Datacoding, m models.Message, udh string) QueueMessage {
	return &qMessage{[]int64{}, message, enc, m, udh}
}

// GetRecipientsAmount returns the amount of recipients currently added to the message
func (m *qMessage) GetRecipientsAmount() int64 {
	return int64(len(m.recipients))
}

// AddRecipient adds recipient to the list and returns an error such recipient is already added
func (m *qMessage) AddRecipient(r int64) error {
	l := len(m.recipients)

	i := sort.Search(l, func(i int) bool {
		return r == m.recipients[i]
	})

	if i < len(m.recipients) {
		return errors.New("existing recipient")
	}

	m.recipients = append(m.recipients, r)

	return nil
}

// GetMessage returns already encoded message from the message
func (m *qMessage) GetMessage() string {
	return m.Message
}

// GetOriginalRecipient number from the original message
func (m *qMessage) GetOriginalRecipient() int64 {
	return m.OriginalMessage.GetRecipient()
}

// GetOriginator from the original message
func (m *qMessage) GetOriginator() string {
	return m.OriginalMessage.GetOriginator()
}

// GetRecipients immutable collection of added recipients
func (m *qMessage) GetRecipients() []string {
	c := make([]string, len(m.recipients))

	for i, r := range m.recipients {
		c[i] = strconv.FormatInt(r, 10)
	}

	return c
}

// GetDataCoding of the message
func (m *qMessage) GetDataCoding() utils.Datacoding {
	return m.Encoding
}

// GetUDH of the message (if any)
func (m *qMessage) GetUDH() string {
	return m.UDH
}

// ByRecipientsAmount is type for sorting the collection of QueueMessage by recipients amount
type ByRecipientsAmount []QueueMessage

// Len returns length of the collection
func (a ByRecipientsAmount) Len() int {
	return len(a)
}

// Swap returns swapped items
func (a ByRecipientsAmount) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less is used for sorting by amount of recipients amount
func (a ByRecipientsAmount) Less(i, j int) bool {
	return a[i].GetRecipientsAmount() < a[j].GetRecipientsAmount()
}
