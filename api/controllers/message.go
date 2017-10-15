package controllers

import (
	"api/models"
	"github.com/labstack/echo"
	"net/http"
	"queue"
	"utils"
)

// MessageControllers interface consist all the message endpoints handlers
type MessageControllers interface {
	HandleMessage(c echo.Context) error
}

type mcontroller struct {
	queue queue.MessageQueue
	udh   utils.UDHEncoder
}

// HandleMessage controller
func (mc *mcontroller) HandleMessage(c echo.Context) error {
	var err error

	m := models.InitMessage(mc.udh)

	if err = c.Bind(m); err != nil {
		return err
	}

	if err = c.Validate(m); err != nil {
		t := utils.HumaniseValidationErrors(err)
		return c.JSON(http.StatusUnprocessableEntity, t)
	}

	if err = m.ConvertBody(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	mc.queue.Push(m)

	return c.JSON(http.StatusOK, m)
}

// InitMessageControllers creates the message controller instance
func InitMessageControllers(q queue.MessageQueue, udh utils.UDHEncoder) MessageControllers {
	return &mcontroller{q, udh}
}
