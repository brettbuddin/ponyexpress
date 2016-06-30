package mailbox

import (
	"fmt"
	"sync"
	"time"

	"github.com/brettbuddin/ponyexpress/logger"
)

var (
	SizeLimit   = 500
	ExpireAfter = time.Hour
)

func NewRegistry() *Registry {
	r := &Registry{
		boxes: map[string]*Mailbox{},
		dirty: make(chan *Mailbox),
	}
	go r.eviction()
	return r
}

type Registry struct {
	sync.RWMutex
	boxes map[string]*Mailbox
	dirty chan *Mailbox
}

func (r *Registry) Close() {
	close(r.dirty)
}

func (r *Registry) Create(id string) (*Mailbox, error) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.boxes[id]; ok {
		return nil, fmt.Errorf("mailbox already exists: %s", id)
	}
	b := NewMailbox(id, r.dirty)
	r.boxes[id] = b
	return b, nil
}

func (r *Registry) Get(id string) (*Mailbox, error) {
	r.RLock()
	defer r.RUnlock()
	b, ok := r.boxes[id]
	if !ok {
		return nil, fmt.Errorf("unknown mailbox: %s", id)
	}
	return b, nil
}

func (r *Registry) Remove(id string) (*Mailbox, error) {
	r.Lock()
	defer r.Unlock()
	box, ok := r.boxes[id]
	if !ok {
		return nil, fmt.Errorf("unknown mailbox: %s", id)
	}
	delete(r.boxes, id)
	return box, nil
}

var (
	DirtyMax   = 100
	EvictEvery = 30 * time.Second
)

func (r *Registry) eviction() {
	var (
		dirty = map[*Mailbox]struct{}{}
		evict = func() {
			logger.Debugf("gc: started")
			expire := time.Now().Add(-ExpireAfter)
			logger.Debugf("gc: evicting messages older than %s", expire)
			for mb := range dirty {
				evicted := mb.Evict(expire)
				logger.Debugf("gc: evicted %d messages from %s\n", evicted, mb.ID)
			}
			dirty = map[*Mailbox]struct{}{}
			logger.Debugf("gc: completed")
		}
		tick = time.Tick(EvictEvery)
	)

	for {
		select {
		case mailbox, ok := <-r.dirty:
			if !ok {
				return
			}
			if _, ok := dirty[mailbox]; !ok {
				dirty[mailbox] = struct{}{}
			}
			if len(dirty) > DirtyMax {
				evict()
			}
		case <-tick:
			if len(dirty) == 0 {
				continue
			}
			evict()
		}
	}
}
