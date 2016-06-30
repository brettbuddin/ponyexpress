package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"
	. "gopkg.in/go-playground/assert.v1"

	"github.com/brettbuddin/ponyexpress/server"
)

func TestCoreFilters(t *testing.T) {
	app := server.New(context.Background())
	app.GET("/", func(c context.Context, w server.ResponseWriter, r *server.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(app)
	defer server.Close()

	resp, err := http.Get(server.URL)
	Equal(t, err, nil)
	Equal(t, resp.StatusCode, 200)
	NotEqual(t, resp.Header.Get("Request-Id"), "")
	NotEqual(t, resp.Header.Get("Runtime"), "")
}

func TestSetResponseHeaderFilter(t *testing.T) {
	app := server.New(context.Background())
	app.AddFilters(
		server.SetResponseHeader("Header1", "Value1"),
		server.SetResponseHeader("Header2", "Value2"),
	)
	app.GET("/", func(c context.Context, w server.ResponseWriter, r *server.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(app)
	defer server.Close()

	resp, err := http.Get(server.URL)
	Equal(t, err, nil)
	Equal(t, resp.StatusCode, 200)
	Equal(t, resp.Header.Get("Header1"), "Value1")
	Equal(t, resp.Header.Get("Header2"), "Value2")
}

func TestSetRequestHeaderFilter(t *testing.T) {
	app := server.New(context.Background())
	app.AddFilters(
		server.SetRequestHeader("HeaderUserRequestedValue", "Value1", false),
		server.SetRequestHeader("HeaderDefaultValue", "Value2", false),
		server.SetRequestHeader("HeaderOvewritten", "Value3", true),
	)
	app.GET("/", func(c context.Context, w server.ResponseWriter, r *server.Request) {
		Equal(t, r.Header.Get("HeaderUserRequestedValue"), "UserValue1")
		Equal(t, r.Header.Get("HeaderDefaultValue"), "Value2")
		Equal(t, r.Header.Get("HeaderOvewritten"), "Value3")
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(app)
	defer server.Close()

	client := &http.Client{}
	req, err := http.NewRequest("GET", server.URL, nil)
	Equal(t, err, nil)
	req.Header.Set("HeaderUserRequestedValue", "UserValue1")
	req.Header.Set("HeaderDefaultValue", "")
	req.Header.Set("HeaderOvewritten", "UserValue3")

	resp, err := client.Do(req)
	Equal(t, err, nil)
	Equal(t, resp.StatusCode, 200)
}
