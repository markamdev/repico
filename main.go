package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/markamdev/repico/gpio"
	"github.com/markamdev/repico/internal"
	rr "github.com/markamdev/repico/rest"
)

func main() {
	log.Printf("RePiCo server")
	err := initialize()
	if err != nil {
		// handle error and exit
		log.Fatalln("Failed to launch RePiCo:", err)
	}

	setSighandler()

	for {
		msg := <-internal.MessageBus
		switch msg.Type {
		case internal.ClosedBySignal:
			log.Println("Exiting due to received signal")
			deinit()
			return
		case internal.HTTPServerError:
			log.Println("Exiting due to HTTP Server error")
			deinit()
			return
		case internal.Info:
			log.Println("Info:", msg.Content)
		}
	}
}

// initialize is a single initialization (and potential failure) point
func initialize() error {
	// TODO add initial config reading and setting here
	_, err := internal.ReadConfig()
	if err != nil {
		return fmt.Errorf("Configuration reading error: %v", err.Error())
	}

	err = rr.LaunchServer()
	if err != nil {
		return fmt.Errorf("Server launching failure: %v", err)
	}

	log.Println("Initialization completed")
	return nil
}

func setSighandler() {
	log.Println("Attaching signal handler")
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		_ = <-interruptCh
		internal.MessageBus <- internal.Message{Type: internal.ClosedBySignal, Content: "Closing signal received"}
	}()
}

func deinit() {
	gpio.ClearPins()
	rr.StopServer()
	log.Println("Deinitialization completed")
}
