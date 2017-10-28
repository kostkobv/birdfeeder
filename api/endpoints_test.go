package api_test

import (
	"api"
	"mocks"
	"testing"

	"github.com/labstack/echo"
)

func TestRegisterEndpoints(t *testing.T) {
	t.Run("registered POST /message", func(t *testing.T) {
		e := echo.New()
		udh := &mocks.UDHEncoderMock{}
		q := &mocks.MessageQueue{}
		api.RegisterEndpoints(e, udh, q)

		for _, r := range e.Routes() {
			if r.Path == "/message" && r.Method == "POST" {
				return
			}
		}

		t.Fail()
	})
}
