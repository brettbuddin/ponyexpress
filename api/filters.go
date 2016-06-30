package api

import (
	"github.com/brettbuddin/ponyexpress/server"
	"golang.org/x/net/context"
)

func SetContentType(next server.ContextHandle) server.ContextHandle {
	return func(ctx context.Context, w server.ResponseWriter, r *server.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(ctx, w, r)
	}
}
