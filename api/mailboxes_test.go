package api_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"golang.org/x/net/context"

	"github.com/brettbuddin/ponyexpress"
	"github.com/brettbuddin/ponyexpress/mailbox"

	"gopkg.in/check.v1"
)

var _ = check.Suite(&MailboxSuite{})

func Test(t *testing.T) {
	check.TestingT(t)
}

type MailboxSuite struct {
	registry *mailbox.Registry
	server   *httptest.Server
}

func (s *MailboxSuite) SetUpTest(c *check.C) {
	s.registry = mailbox.NewRegistry()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "registry", s.registry)
	s.server = httptest.NewServer(ponyexpress.New(ctx))
}

func (s *MailboxSuite) TearDownTest(c *check.C) {
	s.server.Close()
	s.registry.Close()
}

func (s *MailboxSuite) TestMailboxCreate(c *check.C) {
	uri, _ := url.Parse(s.server.URL)
	uri.Path = "/mailboxes"

	req, err := http.NewRequest(http.MethodPost, uri.String(), nil)
	c.Assert(err, check.IsNil)

	client := http.Client{}
	resp, err := client.Do(req)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 201)
	c.Assert(resp.Header.Get(headerContentType), check.Equals, contentTypeJSON)

	buf, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	validateSchema(c, buf, "../schemas/mailbox.json")
}

func (s *MailboxSuite) TestMailboxDelete(c *check.C) {
	mailbox, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	uri, _ := url.Parse(s.server.URL)
	uri.Path = fmt.Sprintf("/mailboxes/%s", mailbox.ID)

	req, err := http.NewRequest(http.MethodDelete, uri.String(), nil)
	c.Assert(err, check.IsNil)

	client := http.Client{}
	resp, err := client.Do(req)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 200)
	c.Assert(resp.Header.Get(headerContentType), check.Equals, contentTypeJSON)

	buf, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	validateSchema(c, buf, "../schemas/mailbox.json")

	_, err = s.registry.Get("a")
	c.Assert(err, check.NotNil)
}

func (s *MailboxSuite) TestMailboxDelete404(c *check.C) {
	_, err := s.registry.Create("a")
	c.Assert(err, check.IsNil)

	uri, _ := url.Parse(s.server.URL)
	uri.Path = fmt.Sprintf("/mailboxes/%s", "b")

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
