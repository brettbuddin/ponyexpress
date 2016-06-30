package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/context"

	"github.com/brettbuddin/ponyexpress/server"
)

var (
	errMethodNotAllowed    = fmt.Errorf("method not allowed")
	errInternalServerError = fmt.Errorf("internal server error")
	errBadRequest          = fmt.Errorf("bad request")
	errNotFound            = fmt.Errorf("not found")
)

type errResp struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(errResp{err.Error()}); err != nil {
		log.Printf("failure to write error: %s (%d)\n", err, status)
	}
}

func PanicRecovery(ctx context.Context, w server.ResponseWriter, r *server.Request) {
	writeError(w, http.StatusInternalServerError, errInternalServerError)
}

func NotFound(ctx context.Context, w server.ResponseWriter, r *server.Request) {
	writeError(w, http.StatusInternalServerError, errNotFound)
}
