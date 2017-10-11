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
}

// HandleMessage controller
func (mc *mcontroller) HandleMessage(c echo.Context) error {
	var err error

	m := models.InitMessage()

	if err = c.Bind(m); err != nil {
		return err
	}

	if err = c.Validate(m); err != nil {
		t := utils.HumaniseValidationErrors(err)
		return c.JSON(http.StatusUnprocessableEntity, t)
	}

	mc.queue.Push(m)

	return c.JSON(http.StatusOK, m)
}

// InitMessageControllers creates the message controller instance
func InitMessageControllers(q queue.MessageQueue) MessageControllers {
	return &mcontroller{q}
}
