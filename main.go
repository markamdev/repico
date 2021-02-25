package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/markamdev/repico/config"
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
	currConfig, err := config.LoadDefault()
	if err != nil {
		// configuration reading error is not an fatal app error
		log.Println("Configuration reading error:", err.Error())
	}

	// apply fetched configuration
	config.ApplyConfig(currConfig)

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
	// TODO configuration storage closing should be added here when implemented
	log.Println("Deinitialization completed")
}
