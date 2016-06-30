package server

import (
	"fmt"
	"time"

	"github.com/nu7hatch/gouuid"
	"golang.org/x/net/context"

	"github.com/brettbuddin/ponyexpress/logger"
)

// Keys set in the Context
const (
	ContextRequestID = "request_id"

	requestIDHeader = "Request-Id"
	runtimeHeader   = "Runtime"
)

// SetRequestHeader sets a request header to a specified value.
func SetRequestHeader(key string, val string, overwrite bool) Filter {
	return func(next ContextHandle) ContextHandle {
		return func(c context.Context, w ResponseWriter, r *Request) {
			inVal := r.Header.Get(key)
			if overwrite || inVal == "" {
				r.Header.Set(key, val)
			}
			next(c, w, r)
		}
	}
}

// SetResponseHeader sets a response header to a specified value.
func SetResponseHeader(key string, val string) Filter {
	return func(next ContextHandle) ContextHandle {
		return func(c context.Context, w ResponseWriter, r *Request) {
			w.Header().Set(key, val)
			next(c, w, r)
		}
	}
}

func setRequestIDHeader(next ContextHandle) ContextHandle {
	return func(c context.Context, w ResponseWriter, r *Request) {
		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}
		w.Header().Set(requestIDHeader, id.String())
		next(context.WithValue(c, ContextRequestID, id.String()), w, r)
	}
}

func setRuntimeHeader(next ContextHandle) ContextHandle {
	return func(c context.Context, w ResponseWriter, r *Request) {
		var start time.Time
		requestID := c.Value(ContextRequestID)

		w.BeforeWrite(func(w ResponseWriter) {
			elapsed := time.Since(start).Seconds()
			w.Header().Set(runtimeHeader, fmt.Sprintf("%f", elapsed))
			logger.Infof("finished: %s", logger.Fields{
				"request_id": requestID,
				"method":     r.Method,
				"uri":        r.RequestURI,
				"status":     w.Status(),
				"elapsed":    fmt.Sprintf("%f", elapsed),
			})
		})

		logger.Debugf("started: %s", logger.Fields{
			"id":     requestID,
			"method": r.Method,
			"uri":    r.RequestURI,
		})
		start = time.Now()
		next(c, w, r)
	}
}
