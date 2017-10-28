package api_test

import (
	"api"
	"mocks"
	"queue"
	"reflect"
	"testing"
	"utils"

	"github.com/stretchr/testify/assert"
)

func TestInitServer(t *testing.T) {
	address := "address"
	v := utils.InitValidator()
	udh := utils.InitEncoder()
	mb := &mocks.ExternalMessageBirdClientMock{}
	q := queue.InitQueue(mb)
	s := api.InitServer(address, v, udh, q)

	e := reflect.ValueOf(s).Elem()

	t.Run("server has provided address", func(t *testing.T) {
		assert.Equal(t, address, e.FieldByName("Address").String())
	})

	t.Run("server has echo instance", func(t *testing.T) {
		assert.Equal(t, "*echo.Echo", e.FieldByName("Instance").Type().String())
	})

	t.Run("server has provided queue instance", func(t *testing.T) {
		assert.Exactly(t, q, e.FieldByName("Queue").Interface())
	})
}

func TestServer_Start(t *testing.T) {
	address := "address"
	v := utils.InitValidator()
	udh := utils.InitEncoder()
	mb := &mocks.ExternalMessageBirdClientMock{}
	q := queue.InitQueue(mb)
	s := api.InitServer(address, v, udh, q)

	t.Run("returns echo error if invalid address provided", func(t *testing.T) {
		assert.Equal(t, "listen tcp: address address: missing port in address", s.Start().Error())
	})
}
