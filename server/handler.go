package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

// ContextHandle is a function that handles an inbound request.
type ContextHandle func(context.Context, ResponseWriter, *Request)

// Filter is a middleware.
type Filter func(ContextHandle) ContextHandle

// NewRequest creates a new Request.
func NewRequest(method, url string, body io.Reader, params httprouter.Params) (*Request, error) {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return &Request{r, params}, nil
}

// Request represents an inbound request.
type Request struct {
	*http.Request
	URLParams httprouter.Params
}

// Param looks for "parameters":
// - In the URL parameters (`:key`-style slugs in the URL)
// - In form data
func (r *Request) Param(key string) string {
	if r.URLParams != nil {
		val := r.URLParams.ByName(key)
		if val != "" {
			return val
		}
	}

	return r.Request.FormValue(key)
}

// Responder commits a status code and content for a request to a ResponseWriter.
type Responder interface {
	Respond(ResponseWriter) error
}

// WriteRaw writes plain text to a ResponseWriter.
func WriteRaw(w ResponseWriter, status int, content []byte) error {
	return RawResponder{status, content}.Respond(w)
}

// RawResponder writes plain text to a ResponseWriter.
type RawResponder struct {
	StatusCode int
	Content    []byte
}

func (r RawResponder) Respond(w ResponseWriter) error {
	w.WriteHeader(r.StatusCode)
	w.Write(r.Content)
	return nil
}

// WriteJSON writes JSON to a ResponseWriter.
func WriteJSON(w ResponseWriter, status int, content interface{}) error {
	return JSONResponder{status, content}.Respond(w)
}

// JSONResponder writes JSON to a ResponseWriter.
type JSONResponder struct {
	StatusCode int
	Content    interface{}
}

func (r JSONResponder) Respond(w ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode)
	return json.NewEncoder(w).Encode(r.Content)
}
