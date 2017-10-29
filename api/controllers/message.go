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

// MessageControllers interface consists all the message endpoints handlers
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

	// create new message instance
	m := models.InitMessage()

	// bind request data into message
	if err = c.Bind(m); err != nil {
		return err
	}

	// validate data
	if err = c.Validate(m); err != nil {
		t := utils.HumaniseValidationErrors(err)
		return c.JSON(http.StatusUnprocessableEntity, t)
	}

	// send message to the subroutine for processing
	go mc.SendMessageToQueue(m)

	return c.JSON(http.StatusOK, m)
}

// SendMessageToQueue splits the submitted message, generated UHD and pushes it to the queue. In fact is not a controller method but rather a helper function
func (mc *mcontroller) SendMessageToQueue(m models.Message) {
	body := m.GetBody()

	// split the message
	mes := mc.Udh.SplitTextMessage(body)

	parts := len(mes.Messages)

	var udh string

	// unique hash based on message body
	hash := mc.generateMessageHash(body)

	for p, encoded := range mes.Messages {
		if parts > 1 {
			// generates udh for provided message part if needed
			udh = mc.Udh.GenerateUDH(uint8(p+1), uint8(parts), hash)
		}

		// create QueueMessage instance based on the message part
		qm := qModels.InitQueueMessage(encoded, mes.Encoding, m, udh)

		// push it to the queue
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
