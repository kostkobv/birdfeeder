package controllers_test

import (
	"api/controllers"
	"testing"

	"errors"

	"mocks"

	"net/http"

	"utils"

	"time"

	apiModels "api/models"
	"queue/models"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInitMessageControllers(t *testing.T) {
	qMock := &mocks.MessageQueue{}
	udhMock := &mocks.UDHEncoderMock{}
	c := controllers.InitMessageControllers(qMock, udhMock)

	t.Run("initialize message controller", func(t *testing.T) {
		assert.NotNil(t, c)
	})
}

func TestMcontroller_HandleMessage(t *testing.T) {
	qMock := &mocks.MessageQueue{}
	udhMock := &mocks.UDHEncoderMock{}
	c := controllers.InitMessageControllers(qMock, udhMock)

	t.Run("returns error if didn't manage to bind the request", func(t *testing.T) {
		cm := new(mocks.EchoContextMock)
		e := errors.New("Couldn't bind")
		cm.On("Bind", mock.Anything).Return(e)

		err := c.HandleMessage(cm)
		assert.Equal(t, e, err)
	})

	t.Run("returns error if didn't manage to validate the request", func(t *testing.T) {
		cm := new(mocks.EchoContextMock)

		et := &mocks.FieldErrorMock{}
		et.On("Field").Return("test")
		et.On("Error").Return("test")
		et.On("Translate", mock.Anything).Return("Test")
		e := validator.ValidationErrors{et}

		err := map[string]string{"test": "Test"}

		cm.On("Bind", mock.Anything).Return(nil)
		cm.On("Validate", mock.Anything).Return(e)
		cm.On("JSON", http.StatusUnprocessableEntity, err).Return(nil)

		returnedError := c.HandleMessage(cm)
		assert.Nil(t, returnedError)
		cm.AssertCalled(t, "JSON", http.StatusUnprocessableEntity, err)
	})

	t.Run("sends message to queue and renders passed message object", func(t *testing.T) {
		cm := new(mocks.EchoContextMock)
		cm.On("Bind", mock.Anything).Return(nil)
		cm.On("Validate", mock.Anything).Return(nil)

		chanWait := make(chan time.Time)

		cm.On("JSON", http.StatusOK, mock.Anything).Return(nil).WaitFor = chanWait

		mes := []string{"a", "b"}
		enc := &utils.Encoded{
			Encoding: "plain",
			Messages: mes,
		}

		udhMock.On("SplitTextMessage", mock.Anything).Return(enc)
		udhMock.On("GenerateUDH", mock.Anything, mock.Anything, mock.Anything).Return("")

		qMock.On("Push", mock.Anything).Return(nil).RunFn = func(arguments mock.Arguments) {
			m := arguments.Get(0).(models.QueueMessage)

			if m.GetMessage() == mes[len(mes)-1] {
				chanWait <- time.Now()
			}
		}

		om := apiModels.InitMessage()
		m1 := models.InitQueueMessage("a", "plain", om, "")
		m2 := models.InitQueueMessage("b", "plain", om, "")

		returnedError := c.HandleMessage(cm)
		assert.Nil(t, returnedError)
		qMock.AssertCalled(t, "Push", m1)
		qMock.AssertCalled(t, "Push", m2)
		qMock.AssertNumberOfCalls(t, "Push", 2)
	})
}
