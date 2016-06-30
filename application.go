package ponyexpress

import (
	"golang.org/x/net/context"

	"github.com/brettbuddin/ponyexpress/api"
	"github.com/brettbuddin/ponyexpress/server"
)

type Application struct {
	*server.Server
}

func New(ctx context.Context) *Application {
	server := server.New(ctx)
	server.AddFilters(api.SetContentType)
	server.PanicHandler = api.PanicRecovery
	server.NotFoundHandler = api.NotFound

	// Mailboxes
	server.POST("/mailboxes", api.MailboxCreate)
	server.DELETE("/mailboxes/:address", api.MailboxDelete)

	// Messages
	server.GET("/mailboxes/:address/messages", api.MessageIndex)
	server.POST("/mailboxes/:address/messages", api.MessageCreate)
	server.GET("/mailboxes/:address/messages/:message_id", api.MessageShow)
	server.DELETE("/mailboxes/:address/messages/:message_id", api.MessageDelete)

	return &Application{server}
}
