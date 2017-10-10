package api

import (
	"api/controllers"
	"github.com/labstack/echo"
)

// RegisterEndpoints for API server
func RegisterEndpoints(e *echo.Echo) {
	e.POST("/message", controllers.HandleMessage)
}
