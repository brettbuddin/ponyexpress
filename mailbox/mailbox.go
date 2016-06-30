package mailbox

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type Message struct {
	ID       string    `json:"id"`
	Sender   string    `json:"sender"`
	Subject  string    `json:"subject"`
	Body     string    `json:"body"`
	Received time.Time `json:"received"`
}

func (m *Message) Key() string {
	return m.ID
}

func NewMailbox(id string, dirty chan *Mailbox) *Mailbox {
	return &Mailbox{
		ID:    id,
		list:  newIndexedList(),
		dirty: dirty,
	}
}

type Mailbox struct {
	sync.RWMutex
	ID    string `json:"id"`
	list  *indexedList
	dirty chan *Mailbox
}

func (b *Mailbox) Push(m *Message) {
	b.Lock()
	defer b.Unlock()
	if b.list.Len() > SizeLimit {
		b.list.Remove(b.list.Front())
	}
	b.list.PushBack(m)
	b.dirty <- b
}

func (b *Mailbox) Get(id string) (*Message, error) {
	b.RLock()
	defer b.RUnlock()
	if e, ok := b.list.GetKey(id); ok {
		return e.Value.(*Message), nil
	}
	return nil, fmt.Errorf("unknown message: %s", id)
}

func (b *Mailbox) Remove(id string) (*Message, error) {
	b.Lock()
	defer b.Unlock()
	if e, ok := b.list.GetKey(id); ok {
		return b.list.Remove(e).(*Message), nil
	}
	return nil, fmt.Errorf("unknown message: %s", id)
}

func (b *Mailbox) List(sinceID string, limit int) []*Message {
	b.Lock()
	defer b.Unlock()
	messages := []*Message{}

	var e *list.Element
	if sinceID != "" {
		var ok bool
		e, ok = b.list.GetKey(sinceID)
		if ok {
			e = e.Next()
		}
	} else {
		e = b.list.Front()
	}

	count := 0
	for ; e != nil; e = e.Next() {
		if count >= limit {
			break
		}
		msg := e.Value.(*Message)
		messages = append([]*Message{msg}, messages...)
		count++
	}

	return messages
}

func (b *Mailbox) Evict(cutoff time.Time) int {
	b.Lock()
	defer b.Unlock()
	evicted := 0
	var next *list.Element
	for e := b.list.Front(); e != nil; e = next {
		next = e.Next()
		msg := e.Value.(*Message)
		if msg.Received.Before(cutoff) {
			b.list.Remove(e)
			evicted++
		}
	}
	return evicted
}
