package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/brettbuddin/ponyexpress"
	"github.com/brettbuddin/ponyexpress/api"
	"github.com/brettbuddin/ponyexpress/mailbox"

	"gopkg.in/check.v1"
)

var _ = check.Suite(&MessageSuite{})

type MessageSuite struct {
	registry *mailbox.Registry
	server   *httptest.Server
}

func (s *MessageSuite) SetUpTest(c *check.C) {
	s.registry = mailbox.NewRegistry()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "registry", s.registry)
	s.server = httptest.NewServer(ponyexpress.New(ctx))
}

func (s *MessageSuite) TearDownTest(c *check.C) {
	s.server.Close()
	s.registry.Close()
}

type message struct {
	Sender  string `json:"sender"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (s *MessageSuite) TestCreate(c *check.C) {
	mailbox, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	uri, _ := url.Parse(s.server.URL)
	uri.Path = fmt.Sprintf("/mailboxes/%s/messages", mailbox.ID)
	buf, err := json.Marshal(struct {
		Message message `json:"message"`
	}{
		message{
			Sender:  "brett@buddin.us",
			Subject: "subject",
			Body:    "body",
		},
	})

	req, err := http.NewRequest(http.MethodPost, uri.String(), bytes.NewBuffer(buf))
	c.Assert(err, check.IsNil)

	client := http.Client{}
	resp, err := client.Do(req)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 201)
	c.Assert(resp.Header.Get(headerContentType), check.Equals, contentTypeJSON)

	buf, err = ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	validateSchema(c, buf, "../schemas/message.json")

	var content struct {
		Message map[string]interface{} `json:"message"`
	}
	err = json.Unmarshal(buf, &content)
	c.Assert(err, check.IsNil)
	c.Assert(content.Message["sender"], check.Equals, "brett@buddin.us")
}

func (s *MessageSuite) TestCreateMissingSender(c *check.C) {
	mailbox, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	uri, _ := url.Parse(s.server.URL)
	uri.Path = fmt.Sprintf("/mailboxes/%s/messages", mailbox.ID)
	buf, err := json.Marshal(struct {
		Message message `json:"message"`
	}{
		message{
			Sender:  "",
			Subject: "subject",
			Body:    "body",
		},
	})

	req, err := http.NewRequest(http.MethodPost, uri.String(), bytes.NewBuffer(buf))
	c.Assert(err, check.IsNil)

	client := http.Client{}
	resp, err := client.Do(req)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 400)
	c.Assert(resp.Header.Get(headerContentType), check.Equals, contentTypeJSON)

	buf, err = ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	validateSchema(c, buf, "../schemas/error.json")
}

func (s *MessageSuite) TestCreateMissingSubject(c *check.C) {
	mailbox, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	uri, _ := url.Parse(s.server.URL)
	uri.Path = fmt.Sprintf("/mailboxes/%s/messages", mailbox.ID)
	buf, err := json.Marshal(struct {
		Message message `json:"message"`
	}{
		message{
			Sender:  "brett@buddin.us",
			Subject: "",
			Body:    "body",
		},
	})

	req, err := http.NewRequest(http.MethodPost, uri.String(), bytes.NewBuffer(buf))
	c.Assert(err, check.IsNil)

	client := http.Client{}
	resp, err := client.Do(req)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 400)
	c.Assert(resp.Header.Get(headerContentType), check.Equals, contentTypeJSON)

	buf, err = ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	validateSchema(c, buf, "../schemas/error.json")
}

func (s *MessageSuite) TestCreateMissingBody(c *check.C) {
	mailbox, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	uri, _ := url.Parse(s.server.URL)
	uri.Path = fmt.Sprintf("/mailboxes/%s/messages", mailbox.ID)
	buf, err := json.Marshal(struct {
		Message message `json:"message"`
	}{
		message{
			Sender:  "brett@buddin.us",
			Subject: "subject",
			Body:    "",
		},
	})

	req, err := http.NewRequest(http.MethodPost, uri.String(), bytes.NewBuffer(buf))
	c.Assert(err, check.IsNil)

	client := http.Client{}
	resp, err := client.Do(req)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 400)
	c.Assert(resp.Header.Get(headerContentType), check.Equals, contentTypeJSON)

	buf, err = ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	validateSchema(c, buf, "../schemas/error.json")
}

func (s *MessageSuite) TestCreate404Mailbox(c *check.C) {
	_, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	uri, _ := url.Parse(s.server.URL)
	uri.Path = "/mailboxes/b/messages"
	buf, err := json.Marshal(struct {
		Message message `json:"message"`
	}{
		message{
			Sender:  "brett@buddin.us",
			Subject: "subject",
			Body:    "body",
		},
	})

	req, err := http.NewRequest(http.MethodPost, uri.String(), bytes.NewBuffer(buf))
	c.Assert(err, check.IsNil)

	client := http.Client{}
	resp, err := client.Do(req)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 404)
	c.Assert(resp.Header.Get(headerContentType), check.Equals, contentTypeJSON)

	buf, err = ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	validateSchema(c, buf, "../schemas/error.json")
}

func (s *MessageSuite) TestDelete(c *check.C) {
	box, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	message := &mailbox.Message{
		ID:      "b",
		Sender:  "brett@buddin.us",
		Subject: "subject",
		Body:    "body",
	}
	box.Push(message)

	uri, _ := url.Parse(s.server.URL)
	uri.Path = fmt.Sprintf("/mailboxes/%s/messages/%s", box.ID, message.ID)

	req, err := http.NewRequest(http.MethodDelete, uri.String(), nil)
	c.Assert(err, check.IsNil)

	client := http.Client{}
	resp, err := client.Do(req)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 200)
	c.Assert(resp.Header.Get(headerContentType), check.Equals, contentTypeJSON)

	_, err = box.Get("b")
	c.Assert(err, check.NotNil)

	buf, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	validateSchema(c, buf, "../schemas/message.json")
}

func (s *MessageSuite) TestDelete404Mailbox(c *check.C) {
	box, err := s.registry.Create("b")
	c.Assert(err, check.IsNil)

	message := &mailbox.Message{
		ID:      "b",
		Sender:  "brett@buddin.us",
		Subject: "subject",
		Body:    "body",
	}
	box.Push(message)

	uri, _ := url.Parse(s.server.URL)
	uri.Path = fmt.Sprintf("/mailboxes/a/messages/%s", message.ID)

	req, err := http.NewRequest(http.MethodDelete, uri.String(), nil)
	c.Assert(err, check.IsNil)

	client := http.Client{}
	resp, err := client.Do(req)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 404)
	c.Assert(resp.Header.Get(headerContentType), check.Equals, contentTypeJSON)

	buf, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	validateSchema(c, buf, "../schemas/error.json")
}

func (s *MessageSuite) TestDelete404(c *check.C) {
	box, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	message := &mailbox.Message{
		ID:      "b",
		Sender:  "brett@buddin.us",
		Subject: "subject",
		Body:    "body",
	}
	box.Push(message)

	uri, _ := url.Parse(s.server.URL)
	uri.Path = fmt.Sprintf("/mailboxes/%s/messages/%s", box.ID, "c")

	req, err := http.NewRequest(http.MethodDelete, uri.String(), nil)
	c.Assert(err, check.IsNil)

	client := http.Client{}
	resp, err := client.Do(req)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 404)
	c.Assert(resp.Header.Get(headerContentType), check.Equals, contentTypeJSON)

	buf, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	validateSchema(c, buf, "../schemas/error.json")
}

func (s *MessageSuite) TestGet(c *check.C) {
	box, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	message := &mailbox.Message{
		ID:      "b",
		Sender:  "brett@buddin.us",
		Subject: "subject",
		Body:    "body",
	}
	box.Push(message)

	uri, _ := url.Parse(s.server.URL)
	uri.Path = fmt.Sprintf("/mailboxes/%s/messages/%s", box.ID, message.ID)

	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	c.Assert(err, check.IsNil)

	client := http.Client{}
	resp, err := client.Do(req)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 200)
	c.Assert(resp.Header.Get(headerContentType), check.Equals, contentTypeJSON)

	var content struct {
		Message map[string]interface{} `json:"message"`
	}
	err = json.NewDecoder(resp.Body).Decode(&content)
	c.Assert(err, check.IsNil)
	c.Assert(content.Message["id"], check.Equals, "b")
}

func (s *MessageSuite) TestGet404(c *check.C) {
	box, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	message := &mailbox.Message{
		ID:      "b",
		Sender:  "brett@buddin.us",
		Subject: "subject",
		Body:    "body",
	}
	box.Push(message)

	uri, _ := url.Parse(s.server.URL)
	uri.Path = fmt.Sprintf("/mailboxes/%s/messages/%s", box.ID, "c")
	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	c.Assert(err, check.IsNil)

	client := http.Client{}
	resp, err := client.Do(req)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 404)
	c.Assert(resp.Header.Get(headerContentType), check.Equals, contentTypeJSON)

	buf, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	validateSchema(c, buf, "../schemas/error.json")
}

func (s *MessageSuite) TestGet404Mailbox(c *check.C) {
	box, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	message := &mailbox.Message{
		ID:      "b",
		Sender:  "brett@buddin.us",
		Subject: "subject",
		Body:    "body",
	}
	box.Push(message)

	uri, _ := url.Parse(s.server.URL)
	uri.Path = fmt.Sprintf("/mailboxes/b/messages/%s", "b")
	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	c.Assert(err, check.IsNil)

	client := http.Client{}
	resp, err := client.Do(req)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 404)
	c.Assert(resp.Header.Get(headerContentType), check.Equals, contentTypeJSON)

	buf, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	validateSchema(c, buf, "../schemas/error.json")
}

func (s *MessageSuite) TestIndex(c *check.C) {
	box, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	now := time.Now()
	for i := 0; i < 100; i++ {
		box.Push(&mailbox.Message{
			ID:       strconv.Itoa(i),
			Sender:   "brett@buddin.us",
			Subject:  "subject",
			Body:     "body",
			Received: now,
		})
	}

	uri, _ := url.Parse(s.server.URL)
	uri.Path = fmt.Sprintf("/mailboxes/%s/messages", box.ID)
	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	c.Assert(err, check.IsNil)

	client := http.Client{}
	resp, err := client.Do(req)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 200)
	c.Assert(resp.Header.Get(headerContentType), check.Equals, contentTypeJSON)

	buf, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	validateSchema(c, buf, "../schemas/message_index.json")
}

func (s *MessageSuite) TestIndexCursor(c *check.C) {
	box, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	now := time.Now()
	for i := 0; i < 100; i++ {
		box.Push(&mailbox.Message{
			ID:       strconv.Itoa(i),
			Sender:   "brett@buddin.us",
			Subject:  "subject",
			Body:     "body",
			Received: now,
		})
	}

	uri, _ := url.Parse(s.server.URL)
	uri.Path = fmt.Sprintf("/mailboxes/%s/messages", box.ID)
	uri.RawQuery = url.Values{"since_id": []string{"51"}}.Encode()
	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	c.Assert(err, check.IsNil)

	client := http.Client{}
	resp, err := client.Do(req)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 200)
	c.Assert(resp.Header.Get(headerContentType), check.Equals, contentTypeJSON)

	buf, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	validateSchema(c, buf, "../schemas/message_index.json")

	var content struct {
		Messages []interface{} `json:"messages"`
		Meta     api.Meta      `json:"meta"`
	}
	err = json.Unmarshal(buf, &content)
	c.Assert(err, check.IsNil)
	c.Assert(content.Messages, check.HasLen, 48)
	c.Assert(content.Meta.Results, check.Equals, 48)
	c.Assert(content.Meta.SinceID, check.Equals, "51")
	c.Assert(content.Meta.LastID, check.Equals, "99")
}

func (s *MessageSuite) TestIndexCursorWithLimit(c *check.C) {
	box, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	now := time.Now()
	for i := 0; i < 100; i++ {
		box.Push(&mailbox.Message{
			ID:       strconv.Itoa(i),
			Sender:   "brett@buddin.us",
			Subject:  "subject",
			Body:     "body",
			Received: now,
		})
	}

	uri, _ := url.Parse(s.server.URL)
	uri.Path = fmt.Sprintf("/mailboxes/%s/messages", box.ID)
	uri.RawQuery = url.Values{
		"since_id": []string{"51"},
		"limit":    []string{"10"},
	}.Encode()
	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	c.Assert(err, check.IsNil)

	client := http.Client{}
	resp, err := client.Do(req)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 200)
	c.Assert(resp.Header.Get(headerContentType), check.Equals, contentTypeJSON)

	buf, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	validateSchema(c, buf, "../schemas/message_index.json")

	var content struct {
		Messages []interface{} `json:"messages"`
		Meta     api.Meta      `json:"meta"`
	}
	err = json.Unmarshal(buf, &content)
	c.Assert(err, check.IsNil)
	c.Assert(content.Messages, check.HasLen, 10)
	c.Assert(content.Meta.Results, check.Equals, 10)
	c.Assert(content.Meta.SinceID, check.Equals, "51")
	c.Assert(content.Meta.LastID, check.Equals, "61")
}
