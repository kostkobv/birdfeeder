package controllers

import (
	"api/models"
	"hash/fnv"
	"net/http"
	"queue"
	qModels "queue/models"
	"utils"

	"github.com/labstack/echo"
)

// MessageControllers interface consist all the message endpoints handlers
type MessageControllers interface {
	HandleMessage(c echo.Context) error
	SendMessageToQueue(m models.Message)
}

type mcontroller struct {
	Queue queue.MessageQueue
	Udh   utils.UDHEncoder
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

	go mc.SendMessageToQueue(m)

	return c.JSON(http.StatusOK, m)
}

func (mc *mcontroller) SendMessageToQueue(m models.Message) {
	body := m.GetBody()
	mes := mc.Udh.SplitTextMessage(body)

	parts := len(mes.Messages)

	var udh string

	hash := mc.generateMessageHash(body)

	for p, encoded := range mes.Messages {
		if parts > 1 {
			udh = mc.Udh.GenerateUDH(uint8(p+1), uint8(parts), hash)
		}

		qm := qModels.InitQueueMessage(encoded, mes.Encoding, m, udh)
		mc.Queue.Push(qm)
	}
}

func (mc *mcontroller) generateMessageHash(s ...string) uint32 {
	h := fnv.New32a()

	for _, i := range s {
		h.Write([]byte(i))
	}

	return h.Sum32()
}

// InitMessageControllers creates the message controller instance
func InitMessageControllers(q queue.MessageQueue, udh utils.UDHEncoder) MessageControllers {
	return &mcontroller{q, udh}
}
