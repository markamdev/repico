package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/markamdev/repico/gpio"
	v2 "github.com/markamdev/repico/v2"
)

func main() {
	log.Printf("RePiCo server")

	ctrl := gpio.CreateController("/sys/class/gpio")
	mainRouter := v2.CreateHandler(ctrl)
	srv := http.Server{Addr: ":8080", Handler: mainRouter}

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				log.Println("Regular server closing")
			} else {
				log.Fatalln("HTTP server launching error:", err)
			}
		}
	}()

	_ = <-sigChannel
	log.Println("Closing server after signal received")
	srv.Shutdown(context.Background())
}
