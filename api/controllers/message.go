package controllers

import (
	"api/models"
	"github.com/labstack/echo"
	"hash/fnv"
	"net/http"
	"queue"
	qModels "queue/models"
	"strconv"
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

	m := models.InitMessage()

	if err = c.Bind(m); err != nil {
		return err
	}

	if err = c.Validate(m); err != nil {
		t := utils.HumaniseValidationErrors(err)
		return c.JSON(http.StatusUnprocessableEntity, t)
	}

	go mc.sendMessageToQueue(m)

	return c.JSON(http.StatusOK, m)
}

func (mc *mcontroller) sendMessageToQueue(m models.Message) {
	body := m.GetBody()
	orig := m.GetOriginator()
	mes := mc.udh.Encode(body)

	parts := len(mes.Messages)

	var udh string

	for p, encoded := range mes.Messages {
		if parts > 1 {
			pStringified := strconv.FormatInt(int64(p+1), 16)
			partsStringified := strconv.FormatInt(int64(parts), 16)
			hash := mc.generateMessageHash(pStringified, partsStringified, body, orig)
			udh = mc.udh.GenerateUDH(uint8(p+1), uint8(parts), hash)
		}

		qm := qModels.InitQueueMessage(encoded, mes.Encoding, m, udh)
		mc.queue.Push(qm)
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
