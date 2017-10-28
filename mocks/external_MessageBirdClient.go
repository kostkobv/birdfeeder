package mocks

import "github.com/stretchr/testify/mock"
import mb "github.com/messagebird/go-rest-api"

type ExternalMessageBirdClientMock struct {
	mock.Mock
}

func (m *ExternalMessageBirdClientMock) NewMessage(originator string, recipients []string, body string, msgParams *mb.MessageParams) (*mb.Message, error) {
	args := m.Called(originator, recipients, body, msgParams)
	return args.Get(0).(*mb.Message), args.Error(1)
}
