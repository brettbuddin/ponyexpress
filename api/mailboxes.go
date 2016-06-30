package api

import (
	"encoding/json"
	"net/http"

	"github.com/satori/go.uuid"
	"golang.org/x/net/context"

	"github.com/brettbuddin/ponyexpress/mailbox"
	"github.com/brettbuddin/ponyexpress/server"
)

const RegistryKey = "registry"

type MailboxResponse struct {
	Mailbox *mailbox.Mailbox `json:"mailbox"`
}

func MailboxCreate(ctx context.Context, w server.ResponseWriter, r *server.Request) {
	registry := ctx.Value(RegistryKey).(*mailbox.Registry)
	box, err := registry.Create(uuid.NewV4().String())
	if err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(MailboxResponse{box}); err != nil {
		writeError(w, http.StatusInternalServerError, errInternalServerError)
		return
	}
}

func MailboxDelete(ctx context.Context, w server.ResponseWriter, r *server.Request) {
	registry := ctx.Value(RegistryKey).(*mailbox.Registry)
	address := r.URLParams.ByName(ParamAddress)
	box, err := registry.Remove(address)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(MailboxResponse{box}); err != nil {
		writeError(w, http.StatusInternalServerError, errInternalServerError)
		return
	}
}
