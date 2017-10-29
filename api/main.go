package api

import (
	"queue"
	"utils"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Server interface
type Server interface {
	Start() error
}

type server struct {
	Instance *echo.Echo
	Address  string
	Queue    queue.MessageQueue
}

// InitServer initialize base API server
func InitServer(address string, v echo.Validator, udh utils.UDHEncoder, q queue.MessageQueue) Server {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// assign custom validator
	e.Validator = v

	RegisterEndpoints(e, udh, q)

	return &server{e, address, q}
}

// Start the server
func (s *server) Start() error {
	e := s.Instance.Start(s.Address)
	s.Instance.Logger.Error(e)
	return e
}
