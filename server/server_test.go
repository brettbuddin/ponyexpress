package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"
	. "gopkg.in/go-playground/assert.v1"

	"github.com/brettbuddin/ponyexpress/server"
)

type errorsBody struct {
	Errors interface{} `json:"errors"`
}

func TestPanicHandler(t *testing.T) {
	recovery := func(c context.Context, w server.ResponseWriter, r *server.Request) {
		server.WriteJSON(w, http.StatusInternalServerError, errorsBody{[]string{"Internal server error"}})
	}

	app := server.New(context.Background())
	app.PanicHandler = recovery
	app.GET("/", func(c context.Context, w server.ResponseWriter, r *server.Request) {
		panic("i'm broken.")
	})
	server := httptest.NewServer(app)
	defer server.Close()

	resp, err := http.Get(server.URL)
	Equal(t, err, nil)
	Equal(t, resp.StatusCode, http.StatusInternalServerError)
	MatchRegex(t, resp.Header.Get("Content-Type"), "application/json")

	var body errorsBody
	err = json.NewDecoder(resp.Body).Decode(&body)
	Equal(t, err, nil)
	Equal(t, len(body.Errors.([]interface{})), 1)
}

func TestNotFoundHandler(t *testing.T) {
	notFound := func(c context.Context, w server.ResponseWriter, r *server.Request) {
		server.WriteJSON(w, http.StatusNotFound, errorsBody{[]string{"Not found"}})
	}

	app := server.New(context.Background())
	app.NotFoundHandler = notFound
	server := httptest.NewServer(app)
	defer server.Close()

	resp, err := http.Get(server.URL)
	Equal(t, err, nil)
	Equal(t, resp.StatusCode, http.StatusNotFound)
	MatchRegex(t, resp.Header.Get("Content-Type"), "application/json")

	var body errorsBody
	err = json.NewDecoder(resp.Body).Decode(&body)
	Equal(t, err, nil)
	Equal(t, len(body.Errors.([]interface{})), 1)
}
