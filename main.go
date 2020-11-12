package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/markamdev/repico/gpio"
	rr "github.com/markamdev/repico/rest"
)

func main() {
	log.Printf("RePiCo server")
	err := initialize()
	if err != nil {
		// handle error and exit
		log.Fatalln("Failed to launch RePiCo:", err)
	}

	interruptCh := make(chan os.Signal)
	signal.Notify(interruptCh, os.Interrupt)

	// wait for app killing
	<-interruptCh
	log.Println("(Killed by Ctrl+C)")
	// TODO Proper HTTP server closing
	//rr.StopServer()
	gpio.DisableGPIO(4)
}

// initialize is a single initialization (and potential failure) point
func initialize() error {
	// TODO change hardcoded pin into value selected by user
	err := gpio.EnableGPIO(4, gpio.Output)
	if err != nil {
		return fmt.Errorf("Failed to set pin 4 as Ouput: %v", err)
	}
	err = rr.LaunchServer()
	if err != nil {
		return fmt.Errorf("Server launching failure: %v", err)
	}
	log.Println("HTTP server started")

	log.Println("Successfully initalized")
	return nil
}
