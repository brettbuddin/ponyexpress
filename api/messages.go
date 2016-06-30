package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/satori/go.uuid"
	"golang.org/x/net/context"

	"github.com/brettbuddin/ponyexpress/mailbox"
	"github.com/brettbuddin/ponyexpress/server"
)

const (
	ParamMessageID = "message_id"
	ParamAddress   = "address"
	ParamLimit     = "limit"
	ParamSinceID   = "since_id"
)

type MessageResponse struct {
	Message *mailbox.Message `json:"message"`
}

type MessageListResponse struct {
	Messages []*mailbox.Message `json:"messages"`
	Meta     Meta               `json:"meta"`
}

type Meta struct {
	Results int    `json:"results"`
	Limit   int    `json:"limit"`
	SinceID string `json:"since_id"`
	LastID  string `json:"last_id"`
}

func MessageIndex(ctx context.Context, w server.ResponseWriter, r *server.Request) {
	registry := ctx.Value(RegistryKey).(*mailbox.Registry)
	box, err := registry.Get(r.URLParams.ByName(ParamAddress))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	params, err := extractListParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	messages := box.List(params.SinceID, params.Limit)

	var lastID string
	if len(messages) > 0 {
		lastID = messages[0].ID
	}

	resp := MessageListResponse{
		Messages: messages,
		Meta: Meta{
			Results: len(messages),
			Limit:   params.Limit,
			SinceID: params.SinceID,
			LastID:  lastID,
		},
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeError(w, http.StatusInternalServerError, errInternalServerError)
		return
	}
}

func MessageShow(ctx context.Context, w server.ResponseWriter, r *server.Request) {
	registry := ctx.Value(RegistryKey).(*mailbox.Registry)
	box, err := registry.Get(r.URLParams.ByName(ParamAddress))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	msg, err := box.Get(r.URLParams.ByName(ParamMessageID))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(MessageResponse{msg}); err != nil {
		writeError(w, http.StatusInternalServerError, errInternalServerError)
		return
	}
}

type MessagePayload struct {
	Message struct {
		Sender  string `json:"sender"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	} `json:"message"`
}

func MessageCreate(ctx context.Context, w server.ResponseWriter, r *server.Request) {
	registry := ctx.Value(RegistryKey).(*mailbox.Registry)
	box, err := registry.Get(r.URLParams.ByName(ParamAddress))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	var in MessagePayload
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, errBadRequest)
		return
	}

	fields := map[string]string{
		"sender":  in.Message.Sender,
		"subject": in.Message.Subject,
		"body":    in.Message.Body,
	}
	for k, v := range fields {
		if v == "" {
			writeError(w, http.StatusBadRequest, fmt.Errorf("%s is required", k))
			return
		}
	}

	msg := &mailbox.Message{
		ID:       uuid.NewV4().String(),
		Sender:   in.Message.Sender,
		Subject:  in.Message.Subject,
		Body:     in.Message.Body,
		Received: time.Now(),
	}
	box.Push(msg)

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(MessageResponse{msg}); err != nil {
		writeError(w, http.StatusInternalServerError, errInternalServerError)
		return
	}
}

func MessageDelete(ctx context.Context, w server.ResponseWriter, r *server.Request) {
	registry := ctx.Value(RegistryKey).(*mailbox.Registry)

	box, err := registry.Get(r.URLParams.ByName(ParamAddress))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	msg, err := box.Remove(r.URLParams.ByName(ParamMessageID))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(MessageResponse{msg}); err != nil {
		writeError(w, http.StatusInternalServerError, errInternalServerError)
		return
	}
}

type ListParams struct {
	Limit   int
	SinceID string
}

func extractListParams(r *server.Request) (*ListParams, error) {
	limit := 100
	if v := r.FormValue(ParamLimit); v != "" {
		var err error
		limit, err = strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
	}

	params := &ListParams{
		Limit:   limit,
		SinceID: r.FormValue(ParamSinceID),
	}

	return params, nil
}
