package server

import (
	"net/http"
)

// ResponseWriter tracks status code and bytes written to a response.
type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher

	Status() int
	Written() bool
	Size() int
	BeforeWrite(func(ResponseWriter))
}

// NewResponseWriter wraps an http.ResponseWriter for tracking.
func NewResponseWriter(w http.ResponseWriter) ResponseWriter {
	return &responseWriter{ResponseWriter: w}
}

type responseWriter struct {
	http.ResponseWriter
	status, size int
	beforeWrites []func(ResponseWriter)
}

func (w *responseWriter) Status() int {
	return w.status
}

func (w *responseWriter) Size() int {
	return w.size
}

func (w *responseWriter) Written() bool {
	return w.status != 0
}

func (w *responseWriter) WriteHeader(s int) {
	w.status = s
	w.callBefores()
	w.ResponseWriter.WriteHeader(s)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if !w.Written() {
		w.WriteHeader(http.StatusOK)
	}
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

// BeforeWrite registers a callback to be called just before writing the response.
func (w *responseWriter) BeforeWrite(before func(ResponseWriter)) {
	w.beforeWrites = append(w.beforeWrites, before)
}

func (w *responseWriter) callBefores() {
	for i := len(w.beforeWrites) - 1; i >= 0; i-- {
		w.beforeWrites[i](w)
	}
}

func (w *responseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
