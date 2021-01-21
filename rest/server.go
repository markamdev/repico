package rest

import (
	"context"
	"log"
	"net/http"

	"github.com/markamdev/repico/internal"
)

const (
	defListenAddress = "0.0.0.0:9000"
)

var srv *http.Server

// LaunchServer initializes REST API server
func LaunchServer() error {
	log.Println("Starting HTTP REST server at:", defListenAddress)
	handler := myHandler{}
	// TODO At least listening port should be configurable
	srv = &http.Server{Addr: defListenAddress, Handler: handler}

	// if server will fail to start application will be killed anyway
	go func() {
		err := srv.ListenAndServe()
		// in fact if the ListenAndServe() exits it means that there an error when starting server
		if err != nil {
			internal.MessageBus <- internal.Message{Type: internal.HTTPServerError, Content: err.Error()}
		}
	}()
	// return nil as no errors till now (ex. in config)
	return nil
}

// StopServer sends stop signal to server (ex. on SIGINT)
func StopServer() {
	log.Println("Stopping HTTP server")
	srv.Shutdown(context.Background())
}
