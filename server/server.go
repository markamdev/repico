package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type muxWrapper struct {
	router *mux.Router
	server http.Server
}

func NewHandler() Handler {
	result := muxWrapper{router: mux.NewRouter(), server: http.Server{}}
	result.router.StrictSlash(true)
	result.server.Handler = result.router
	return &result
}

func (mw *muxWrapper) GetSubRouter(path string) *mux.Router {
	return mw.router.PathPrefix(path).Subrouter()
}

func (mw *muxWrapper) ServeHTTP(port int) error {
	if port < 1 {
		return errors.New("invalid port number")
	}
	mw.server.Addr = ":" + strconv.Itoa(port)
	return mw.server.ListenAndServe()
}

func (mw *muxWrapper) Shutdown(ctx context.Context) {
	mw.server.Shutdown(ctx)
}
