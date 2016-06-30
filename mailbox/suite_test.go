package mailbox

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"gopkg.in/check.v1"
)

var _ = check.Suite(&Suite{})

func Test(t *testing.T) {
	check.TestingT(t)
}

type Suite struct {
	registry *Registry
}

func (s *Suite) SetUpTest(c *check.C) {
	s.registry = NewRegistry()
}

func (s *Suite) TearDownTest(c *check.C) {
	s.registry.Close()
}

func (s Suite) BenchmarkMailboxCreate(c *check.C) {
	for i := 0; i < c.N; i++ {
		s.registry.Create(strconv.Itoa(i))
	}
}

func (s Suite) BenchmarkMessageCreate(c *check.C) {
	box, _ := s.registry.Create("mailbox")
	now := time.Now()
	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		box.Push(&Message{
			ID:       strconv.Itoa(i),
			Received: now,
		})
	}
}

func (s Suite) BenchmarkMessageGetNearEndElement(c *check.C) {
	box, _ := s.registry.Create("mailbox")
	now := time.Now()

	for i := 0; i < 400; i++ {
		box.Push(&Message{
			ID:       strconv.Itoa(i),
			Received: now,
		})
	}

	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		box.Get("399")
	}
}

func (s Suite) BenchmarkMessageGetNearFirstElement(c *check.C) {
	box, _ := s.registry.Create("mailbox")
	now := time.Now()

	for i := 0; i < 400; i++ {
		box.Push(&Message{
			ID:       strconv.Itoa(i),
			Received: now,
		})
	}

	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		box.Get("1")
	}
}

func (s Suite) TestBasicOperations(c *check.C) {
	// Create a mailbox
	box, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	// Add a message
	box.Push(&Message{
		ID:       "message-id",
		Sender:   "brett@buddin.us",
		Received: time.Now(),
	})

	// Get the message
	msg, err := box.Get("message-id")
	c.Assert(err, check.IsNil)
	c.Assert(msg.ID, check.Equals, "message-id")
	c.Assert(msg.Sender, check.Equals, "brett@buddin.us")

	// Get a non-existant message
	_, err = box.Get("does-not-exist")
	c.Assert(err, check.NotNil)

	// Remove the message
	_, err = box.Remove("message-id")
	c.Assert(err, check.IsNil)

	// Remove a non-existant message
	_, err = box.Remove("does-not-exist")
	c.Assert(err, check.NotNil)

	// Get the same mailbox from the registry
	gottenBox, err := s.registry.Get("a")
	c.Assert(err, check.IsNil)
	c.Assert(box, check.Equals, gottenBox)

	// Attempt to create the same mailbox
	_, err = s.registry.Create("a")
	c.Assert(err, check.NotNil)

	// Remove the mailbox
	_, err = s.registry.Remove("a")
	c.Assert(err, check.IsNil)

	// Remove an unknown mailbox
	_, err = s.registry.Remove("b")
	c.Assert(err, check.NotNil)

	// Make sure the mailbox is gone
	_, err = s.registry.Get("a")
	c.Assert(err, check.NotNil)
}

func (s Suite) TestListOrder(c *check.C) {
	b, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	for i := 0; i < 10; i++ {
		m := &Message{
			ID:       fmt.Sprintf("id-%d", i),
			Sender:   "brett@buddin.us",
			Received: time.Now().Add(time.Duration(i) * time.Second),
		}
		b.Push(m)
	}

	messages := b.List("", 100)
	c.Assert(messages, check.HasLen, 10)
	for i, m := range messages {
		c.Assert(m.ID, check.Equals, fmt.Sprintf("id-%d", 9-i))
	}
}

func (s Suite) TestListLimit(c *check.C) {
	b, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	for i := 0; i < 10; i++ {
		m := &Message{
			ID:       fmt.Sprintf("id-%d", i),
			Sender:   "brett@buddin.us",
			Received: time.Now().Add(time.Duration(i) * time.Second),
		}
		b.Push(m)
	}

	messages := b.List("", 1)
	c.Assert(messages, check.HasLen, 1)
}

func (s Suite) TestListCursor(c *check.C) {
	b, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	for i := 0; i < 10; i++ {
		m := &Message{
			ID:       fmt.Sprintf("id-%d", i),
			Sender:   "brett@buddin.us",
			Received: time.Now().Add(time.Duration(i) * time.Second),
		}
		b.Push(m)
	}

	messages := b.List("id-1", 2)
	c.Assert(messages, check.HasLen, 2)
	c.Assert(messages[0].ID, check.Equals, "id-3")
	c.Assert(messages[1].ID, check.Equals, "id-2")

	messages = b.List("id-3", 2)
	c.Assert(messages, check.HasLen, 2)
	c.Assert(messages[0].ID, check.Equals, "id-5")
	c.Assert(messages[1].ID, check.Equals, "id-4")

	messages = b.List("id-8", 2)
	c.Assert(messages, check.HasLen, 1)
	c.Assert(messages[0].ID, check.Equals, "id-9")
}

func (s Suite) TestEviction(c *check.C) {
	b, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	now := time.Now()
	for i := 0; i < 10; i++ {
		m := &Message{
			ID:       fmt.Sprintf("id-%d", i),
			Sender:   "brett@buddin.us",
			Received: now.Add(time.Duration(i) * time.Minute),
		}
		b.Push(m)
	}

	c.Assert(b.Evict(now.Add(1*time.Minute)), check.Equals, 1)
	c.Assert(b.Evict(now.Add(4*time.Minute)), check.Equals, 3)
}
