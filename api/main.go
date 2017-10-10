package api

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Server interface
type Server interface {
	Start()
}

type server struct {
	Instance *echo.Echo
	Address  string
}

// InitServer initialize base API server
func InitServer(address string, v echo.Validator) Server {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Validator = v

	RegisterEndpoints(e)

	return &server{e, address}
}

// Start the server
func (s *server) Start() {
	s.Instance.Logger.Fatal(s.Instance.Start(s.Address))
}
