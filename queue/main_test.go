package queue_test

import (
	"mocks"
	"queue"
	"reflect"
	"testing"

	"queue/models"

	apiModels "api/models"

	"time"

	"errors"

	"github.com/messagebird/go-rest-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInitQueue(t *testing.T) {
	mb := &mocks.ExternalMessageBirdClientMock{}
	q := queue.InitQueue(mb)

	rq := reflect.ValueOf(q).Elem()
	t.Run("inits queue with provided messagebird client", func(t *testing.T) {
		assert.Equal(t, mb, rq.FieldByName("Mb").Interface())
	})

	t.Run("inits queue with mutex", func(t *testing.T) {
		assert.Equal(t, "*sync.Mutex", rq.FieldByName("Mutex").Type().String())
	})

	t.Run("inits queue with pipe", func(t *testing.T) {
		assert.Equal(t, "chan models.QueueMessage", rq.FieldByName("Pipe").Type().String())
	})

	t.Run("inits queue with empty working messages collection", func(t *testing.T) {
		assert.Equal(t, "[]models.QueueMessage", rq.FieldByName("Collection").Type().String())
	})
}

func TestQueue_Push(t *testing.T) {
	t.Run("pushed messages are sent to messagebird every second", func(t *testing.T) {
		t.Run("two identical messages should be sent separately twice", func(t *testing.T) {
			t.Parallel()
			mb := &mocks.ExternalMessageBirdClientMock{}
			q := queue.InitQueue(mb)

			m1 := models.InitQueueMessage("m1", "", apiModels.InitMessage(), "")
			m2 := models.InitQueueMessage("m2", "", apiModels.InitMessage(), "")

			mbMes := &messagebird.Message{}
			mb.On("NewMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mbMes, nil)

			q.Push(m1, m2)

			// 1 sec per message + threshold
			time.Sleep(2*time.Second + 100*time.Millisecond)

			mb.AssertNumberOfCalls(t, "NewMessage", 2)
		})

		t.Run("identical messages with different recipients sent as one message", func(t *testing.T) {
			t.Parallel()
			mb := &mocks.ExternalMessageBirdClientMock{}
			q := queue.InitQueue(mb)

			rm1 := apiModels.InitMessage()
			reflect.ValueOf(rm1).Elem().FieldByName("Recipient").SetInt(123123)
			rm2 := apiModels.InitMessage()
			reflect.ValueOf(rm2).Elem().FieldByName("Recipient").SetInt(123)

			m1 := models.InitQueueMessage("m1", "", rm1, "")
			m2 := models.InitQueueMessage("m1", "", rm2, "")

			mbMes := &messagebird.Message{}
			mb.On("NewMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mbMes, nil)

			q.Push(m1, m2)

			// 1 sec per message + threshold
			time.Sleep(2*time.Second + 100*time.Millisecond)

			mb.AssertNumberOfCalls(t, "NewMessage", 1)
		})

		t.Run("if there was an error - add it back to the queue", func(t *testing.T) {
			t.Parallel()
			mb := &mocks.ExternalMessageBirdClientMock{}
			q := queue.InitQueue(mb)
			m := models.InitQueueMessage("m1", "", apiModels.InitMessage(), "")

			mbMes := &messagebird.Message{}
			err := errors.New("err")
			mb.On("NewMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mbMes, err)

			q.Push(m)

			// 1 sec per message + threshold
			time.Sleep(2*time.Second + 100*time.Millisecond)

			mb.AssertNumberOfCalls(t, "NewMessage", 2)
		})

		t.Run("message with bigger amount of recipients should be a priority", func(t *testing.T) {
			t.Parallel()
			mb := &mocks.ExternalMessageBirdClientMock{}
			q := queue.InitQueue(mb)

			rm2 := apiModels.InitMessage()
			reflect.ValueOf(rm2).Elem().FieldByName("Recipient").SetInt(123)

			rm3 := apiModels.InitMessage()
			reflect.ValueOf(rm3).Elem().FieldByName("Recipient").SetInt(123123123)

			m1 := models.InitQueueMessage("m1", "", apiModels.InitMessage(), "1")
			m2 := models.InitQueueMessage("m2", "", rm2, "2")
			m3 := models.InitQueueMessage("m2", "", rm3, "2")

			mbMes := &messagebird.Message{}
			err := errors.New("err")
			mb.On("NewMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mbMes, err)

			q.Push(m1, m2, m3)

			// 1 sec per message + threshold
			time.Sleep(1*time.Second + 100*time.Millisecond)

			mb.AssertNumberOfCalls(t, "NewMessage", 1)
			mb.AssertCalled(t, "NewMessage", mock.Anything, []string{"123", "123123123"}, "m2", mock.Anything)
		})
	})
}
