package server

import (
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

// New creates a new Server.
func New(ctx context.Context) *Server {
	return &Server{
		router:  httprouter.New(),
		context: ctx,
		filters: []Filter{
			setRequestIDHeader,
			setRuntimeHeader,
		},
	}
}

// Server is a HTTP multiplexer.
type Server struct {
	router  *httprouter.Router
	context context.Context
	filters []Filter
	once    sync.Once

	NotFoundHandler ContextHandle
	PanicHandler    ContextHandle
}

// AddFilters adds Filter functions to the chain. These functions are composed in the in-order and wrap *all* handlers
// registered with the Server.
func (r *Server) AddFilters(filters ...Filter) {
	r.filters = append(r.filters, filters...)
}

// HTTP methods.
const (
	MethodDELETE  = "DELETE"
	MethodGET     = "GET"
	MethodHEAD    = "HEAD"
	MethodOPTIONS = "OPTIONS"
	MethodPATCH   = "PATCH"
	MethodPOST    = "POST"
	MethodPUT     = "PUT"
)

// DELETE registers a DELETE handler at a path.
func (r *Server) DELETE(path string, h ContextHandle) {
	r.Handle(MethodDELETE, path, h)
}

// GET registers a GET handler at path.
func (r *Server) GET(path string, h ContextHandle) {
	r.Handle(MethodGET, path, h)
}

// HEAD registers a HEAD handler at a path.
func (r *Server) HEAD(path string, h ContextHandle) {
	r.Handle(MethodHEAD, path, h)
}

// OPTIONS registers an OPTIONS handler at a path.
func (r *Server) OPTIONS(path string, h ContextHandle) {
	r.Handle(MethodOPTIONS, path, h)
}

// PATCH registers a PATCH handler at a path.
func (r *Server) PATCH(path string, h ContextHandle) {
	r.Handle(MethodPATCH, path, h)
}

// POST registers an POST handler at a path.
func (r *Server) POST(path string, h ContextHandle) {
	r.Handle(MethodPOST, path, h)
}

// PUT registers an PUT handler at a path.
func (r *Server) PUT(path string, h ContextHandle) {
	r.Handle(MethodPUT, path, h)
}

// Handle registers a handler for a particular HTTP method at a path.
func (r *Server) Handle(method, path string, h ContextHandle) {
	r.handle(method, path, r.applyFilters(h))
}

func (r *Server) applyFilters(h ContextHandle) ContextHandle {
	if len(r.filters) == 0 {
		return h
	}
	chain := h
	for i := len(r.filters) - 1; i >= 0; i-- {
		chain = r.filters[i](chain)
	}
	return chain
}

func (r *Server) handle(method, path string, h ContextHandle) {
	r.router.Handle(method, path, func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		h(r.context, NewResponseWriter(w), &Request{req, p})
	})
}

// ServeHTTP is the entry-point for HTTP routing. It conforms to `http.Handler` interface. The first time this method is
// called (by an inbound request) error handlers are registered with the underlying HTTP router.
func (r *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.once.Do(r.registerErrorHandlers)
	r.router.ServeHTTP(w, req)
}

func (r *Server) registerErrorHandlers() {
	if r.NotFoundHandler != nil {
		r.router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			r.applyFilters(r.NotFoundHandler)(r.context, NewResponseWriter(w), &Request{req, nil})
		})
	}

	if r.PanicHandler != nil {
		r.router.PanicHandler = func(w http.ResponseWriter, req *http.Request, _ interface{}) {
			r.applyFilters(r.PanicHandler)(r.context, NewResponseWriter(w), &Request{req, nil})
		}
	}
}
