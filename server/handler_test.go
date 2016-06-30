package server_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"
	. "gopkg.in/go-playground/assert.v1"

	"github.com/brettbuddin/ponyexpress/server"
)

func TestFormParam(t *testing.T) {
	expected := "fred"

	app := server.New(context.Background())
	app.GET("/", func(c context.Context, w server.ResponseWriter, r *server.Request) {
		Equal(t, r.Param("name"), expected)
	})
	server := httptest.NewServer(app)
	defer server.Close()

	resp, err := http.Get(fmt.Sprintf("%s?%s", server.URL, "name=fred"))
	Equal(t, err, nil)
	Equal(t, resp.StatusCode, 200)
}

func TestURLParam(t *testing.T) {
	expected := "fred"

	app := server.New(context.Background())
	app.GET("/:name", func(c context.Context, w server.ResponseWriter, r *server.Request) {
		Equal(t, r.Param("name"), expected)
	})
	server := httptest.NewServer(app)
	defer server.Close()

	resp, err := http.Get(fmt.Sprintf("%s/%s", server.URL, "fred"))
	Equal(t, err, nil)
	Equal(t, resp.StatusCode, 200)
}

func TestJSONResponder(t *testing.T) {
	type messageBody struct {
		Message string `json:"message"`
	}

	app := server.New(context.Background())
	app.GET("/", func(c context.Context, w server.ResponseWriter, r *server.Request) {
		server.WriteJSON(w, http.StatusOK, messageBody{"Hello, world."})
	})
	server := httptest.NewServer(app)
	defer server.Close()

	resp, err := http.Get(server.URL)
	Equal(t, err, nil)
	Equal(t, resp.StatusCode, http.StatusOK)
	MatchRegex(t, resp.Header.Get("Content-Type"), "application/json")

	var body messageBody
	err = json.NewDecoder(resp.Body).Decode(&body)
	Equal(t, err, nil)
	Equal(t, body.Message, "Hello, world.")
}

func TestRawResponder(t *testing.T) {
	app := server.New(context.Background())
	app.GET("/", func(c context.Context, w server.ResponseWriter, r *server.Request) {
		server.WriteRaw(w, http.StatusOK, []byte("Hello, world."))
	})
	server := httptest.NewServer(app)
	defer server.Close()

	resp, err := http.Get(server.URL)
	Equal(t, err, nil)
	Equal(t, resp.StatusCode, http.StatusOK)
	MatchRegex(t, resp.Header.Get("Content-Type"), "text/plain")

	body, err := ioutil.ReadAll(resp.Body)
	Equal(t, err, nil)
	Equal(t, string(body), "Hello, world.")
}
