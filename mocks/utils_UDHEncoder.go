package mocks

import (
	"utils"

	"github.com/stretchr/testify/mock"
)

type UDHEncoderMock struct {
	mock.Mock
}

func (um *UDHEncoderMock) Encode(m string) *utils.Encoded {
	args := um.Called(m)
	return args.Get(0).(*utils.Encoded)
}

func (um *UDHEncoderMock) GenerateUDH(p uint8, parts uint8, mesHash uint32) string {
	args := um.Called(p, parts, mesHash)
	return args.String(0)
}

func (um *UDHEncoderMock) SplitTextMessage(m string) *utils.Encoded {
	args := um.Called(m)
	return args.Get(0).(*utils.Encoded)
}
