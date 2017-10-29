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
	Collection []qModels.QueueMessage // consider it to be a cart with messages putted under the pipe
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
	// start to get messages from the pipe and add it to the collection
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
		q.Collection = n // swap collections (swap carts under the pipe)
		q.Mutex.Unlock()

		// send the processing to the goroutine with passing the reference to our fulled collection (cart)
		go q.sendChanges(c)
	}
}

func (q *queue) sendChanges(c []qModels.QueueMessage) {
	ms := q.getUniqueMessages(c)

	// sort by the biggest amount of recipients
	sort.Sort(sort.Reverse(qModels.ByRecipientsAmount(ms)))

	if len(ms) == 0 {
		return
	}

	// send the message with the biggest amount of recipients
	q.SendMessage(ms[0])

	// put the rest of the messages back to the queue
	q.Push(ms[1:]...)
}

func (q *queue) getUniqueMessages(c []qModels.QueueMessage) []qModels.QueueMessage {
	// set with unique messages bodies as the keys to cache the messages with different recipients and same udh
	mSet := map[string]map[string]qModels.QueueMessage{}

	// iterate through the collection
	for _, m := range c {
		b := m.GetMessage()
		udh := m.GetUDH()

		// check if message body is cached already
		mapItems, ok := mSet[b]

		// if it's not in set - cache it and proceed
		if !ok {
			mSet[b] = map[string]qModels.QueueMessage{udh: m}
			_ = m.AddRecipient(m.GetOriginalRecipient()) // #nosec
			continue
		}

		// check if message with such body and udh is already in the collection
		mItem, ok := mapItems[udh]

		// if it is not - cache it and proceed
		if !ok {
			mapItems[udh] = m
			_ = m.AddRecipient(m.GetOriginalRecipient()) // #nosec
			continue
		}

		// if we had identical message let's try to add recipient to the list of the cached message
		err := mItem.AddRecipient(m.GetOriginalRecipient())

		if err != nil {
			// if it's duplicated identical message with the same recipient - it's intended to be sent twice. Send back to the queue
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

// Push sends all provided messages to the queue (pipe)
func (q *queue) Push(m ...qModels.QueueMessage) {
	for _, mes := range m {
		q.Pipe <- mes
	}
}

// SendMessage sends message to the MessageBird API
func (q *queue) SendMessage(m qModels.QueueMessage) {
	params := external.InitMessageBirdParams(m.GetDataCoding(), m.GetUDH())
	a, err := q.Mb.NewMessage(m.GetOriginator(), m.GetRecipients(), m.GetMessage(), params)

	fmt.Println("----------------")
	fmt.Printf("MB Response:\n %#v;\n\nParams:\n %#v;\n\nOriginal message:\n %#v;\n\nError:\n %#v;\n\n Time:\n %v\n", a, params, m, err, time.Now())
	fmt.Println("----------------")

	if err != nil {
		q.Push(m)
	}
}
