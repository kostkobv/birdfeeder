package api

import (
	"api/controllers"
	"github.com/labstack/echo"
	"queue"
	"utils"
)

// RegisterEndpoints for API server
func RegisterEndpoints(e *echo.Echo, udh utils.UDHEncoder, q queue.MessageQueue) {
	mControllers := controllers.InitMessageControllers(q, udh)

	e.POST("/message", mControllers.HandleMessage)
}
