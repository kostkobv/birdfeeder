package api

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"queue"
	"utils"
)

// Server interface
type Server interface {
	Start()
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

	e.Validator = v

	RegisterEndpoints(e, udh, q)

	return &server{e, address, q}
}

// Start the server
func (s *server) Start() {
	s.Instance.Logger.Fatal(s.Instance.Start(s.Address))
}
