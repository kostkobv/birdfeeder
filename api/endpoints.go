package api

import (
	"api/controllers"
	"github.com/labstack/echo"
	"queue"
)

// RegisterEndpoints for API server
func RegisterEndpoints(e *echo.Echo, q queue.MessageQueue) {
	mControllers := controllers.InitMessageControllers(q)

	e.POST("/message", mControllers.HandleMessage)
}
