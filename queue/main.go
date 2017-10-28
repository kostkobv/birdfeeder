package queue

import (
	"external"
	"fmt"
	qModels "queue/models"
	"sort"
	"sync"
	"time"
)

// MessageQueue for sending data to the third parties
type MessageQueue interface {
	Push(m ...qModels.QueueMessage)
}

type queue struct {
	Pipe       chan qModels.QueueMessage
	Mutex      *sync.Mutex
	Collection []qModels.QueueMessage
	Mb         external.MessageBirdClient
}

// InitQueue for sending messages to third-parties
func InitQueue(mb external.MessageBirdClient) MessageQueue {
	c := []qModels.QueueMessage{}
	q := &queue{make(chan qModels.QueueMessage), &sync.Mutex{}, c, mb}

	go q.listenForChanges()

	return q
}

func (q *queue) startCollectingChanges() {
	// start to get messages from the queue and add it to the collection
	for v := range q.Pipe {
		// append is not thread-safe
		q.Mutex.Lock()
		q.Collection = append(q.Collection, v)
		q.Mutex.Unlock()
	}
}

func (q *queue) listenForChanges() {
	go q.startCollectingChanges()

	for {
		time.Sleep(time.Second)

		// prevent data race
		q.Mutex.Lock()
		l := len(q.Collection)
		q.Mutex.Unlock()

		// do nothing if we have nothing in the queue
		if l == 0 {
			continue
		}

		c := q.Collection             // old collection
		n := []qModels.QueueMessage{} // new empty collection

		// not thread-safe
		q.Mutex.Lock()
		// meanwhile put new empty collection here (old collection would still be accessible within goroutine)
		q.Collection = n // swap collections
		q.Mutex.Unlock()

		// send the processing to the goroutine with passing the reference to our collection
		go q.sendChanges(c)
	}
}

func (q *queue) sendChanges(c []qModels.QueueMessage) {
	ms := q.getUniqueMessages(c)

	sort.Sort(sort.Reverse(qModels.ByRecipientsAmount(ms)))

	if len(ms) == 0 {
		return
	}

	q.SendMessage(ms[0])
	q.Push(ms[1:]...)
}

func (q *queue) getUniqueMessages(c []qModels.QueueMessage) []qModels.QueueMessage {
	// set with unique messages bodies as the keys
	mSet := map[string]map[string]qModels.QueueMessage{}

	// iterate through the collection
	for _, m := range c {
		b := m.GetMessage()
		udh := m.GetUDH()
		mapItems, ok := mSet[b]

		// if it's not in set - add it and proceed
		if !ok {
			mSet[b] = map[string]qModels.QueueMessage{udh: m}
			_ = m.AddRecipient(m.GetOriginalRecipient()) // #nosec
			continue
		}

		mItem, ok := mapItems[udh]

		if !ok {
			mapItems[udh] = m
			_ = m.AddRecipient(m.GetOriginalRecipient()) // #nosec
			continue
		}

		// add to the list of recipients
		err := mItem.AddRecipient(m.GetOriginalRecipient())

		if err != nil {
			// if there was an error - send it back to the queue
			q.Push(m)
		}
	}

	var result []qModels.QueueMessage

	for _, mi := range mSet {
		for _, m := range mi {
			result = append(result, m)
		}
	}

	return result
}

func (q *queue) Push(m ...qModels.QueueMessage) {
	for _, mes := range m {
		q.Pipe <- mes
	}
}

func (q *queue) SendMessage(m qModels.QueueMessage) {
	params := external.InitMessageBirdParams(m.GetDataCoding(), m.GetUDH())
	a, err := q.Mb.NewMessage(m.GetOriginator(), m.GetRecipients(), m.GetMessage(), params)

	fmt.Printf("Sent message: %#v; Error: %#v; Time: %v", a, err, time.Now())

	if err != nil {
		q.Push(m)
	}
}
