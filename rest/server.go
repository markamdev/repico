package rest

import (
	"log"
	"net/http"
)

const (
	defListenAddress = "localhost:9000"
)

var srv *http.Server

// LaunchServer initializes REST API server
func LaunchServer() error {
	log.Println("Starting HTTP REST server at:", defListenAddress)
	handler := myHandler{}
	// TODO At least listening port should be configurable
	srv = &http.Server{Addr: defListenAddress, Handler: handler}

	// if server will fail to start application will be killed anyway
	go log.Fatal(srv.ListenAndServe())
	// return nil as no errors till now (ex. in config)
	return nil
}

// StopServer sends stop signal to server (ex. on SIGINT)
func StopServer() {
	srv.Shutdown(nil)
}
