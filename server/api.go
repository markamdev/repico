package server

import (
	"context"

	"github.com/gorilla/mux"
)

type Handler interface {
	GetSubRouter(path string) *mux.Router
	ServeHTTP(port int) error
	Shutdown(ctx context.Context)
}
