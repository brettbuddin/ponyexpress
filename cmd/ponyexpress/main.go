package main

import (
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/brettbuddin/ponyexpress"
	"github.com/brettbuddin/ponyexpress/logger"
	"github.com/brettbuddin/ponyexpress/mailbox"
)

const timeout = 5 * time.Second

func main() {
	registry := mailbox.NewRegistry()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "registry", registry)
	app := ponyexpress.New(ctx)

	// TODO
	// go serveSMTP(registry)

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":3000"
	}
	logger.Infof("Listening at http://localhost%s", addr)
	server := &http.Server{
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		Addr:         addr,
		Handler:      app,
	}
	if err := server.ListenAndServe(); err != nil {
		logger.Errorf(err.Error())
	}
}
