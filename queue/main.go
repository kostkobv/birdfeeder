package queue

import (
	"api/models"
	"fmt"
	"sync"
	"time"
)

// MessageQueue for sending data to the third parties
type MessageQueue interface {
	Push(m models.Message)
}

type queue struct {
	pipe       chan models.Message
	mutex      *sync.Mutex
	collection []models.Message
}

// InitQueue for sending messages to third-parties
func InitQueue() MessageQueue {
	c := []models.Message{}
	q := &queue{make(chan models.Message), &sync.Mutex{}, c}

	go q.listenForChanges()

	return q
}

func (q *queue) startCollectingChanges() {
	for v := range q.pipe {
		q.mutex.Lock()
		q.collection = append(q.collection, v)
		q.mutex.Unlock()
	}
}

func (q *queue) listenForChanges() {
	go q.startCollectingChanges()

	for {
		time.Sleep(time.Second)

		// do nothing if we have nothing in the queue
		if len(q.collection) == 0 {
			continue
		}

		c := q.collection

		// thread-safe lock
		q.mutex.Lock()

		// meanwhile put new empty collection here (old collection would still be accessible within goroutine)
		q.collection = []models.Message{}

		// unlock
		q.mutex.Unlock()

		// send the processing to the goroutine with passing the reference to our collection
		go q.sendChanges(c)

	}
}

func (q *queue) sendChanges(c []models.Message) {
	for _, m := range c {
		fmt.Println(m)
	}
}

// Push message to the queue
func (q *queue) Push(m models.Message) {
	q.pipe <- m
}
