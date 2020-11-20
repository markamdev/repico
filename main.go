package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/markamdev/repico/config"
	"github.com/markamdev/repico/gpio"
	rr "github.com/markamdev/repico/rest"
)

var currentConfig config.Params

func main() {
	log.Printf("RePiCo server")
	err := initialize()
	if err != nil {
		// handle error and exit
		log.Fatalln("Failed to launch RePiCo:", err)
	}

	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)

	sg := <-interruptCh
	log.Println("Closed by singal:", sg)
	deinit()

	log.Println("Exiting...")
}

// initialize is a single initialization (and potential failure) point
func initialize() error {
	cfg, err := config.ReadConfig()
	if err != nil {
		return fmt.Errorf("Configuration reading error: %v", err.Error())
	}
	currentConfig = cfg
	if len(currentConfig.Pins) == 0 {
		return fmt.Errorf("No GPIO to manage - exiting")
	}
	for _, pin := range currentConfig.Pins {
		log.Printf("Enabling output GPIO %d", pin)
		err := gpio.EnableGPIO(pin, gpio.Output)
		if err != nil {
			return fmt.Errorf("Failed to set pin %d as Ouput: %v", pin, err)
		}
	}

	err = rr.LaunchServer()
	if err != nil {
		return fmt.Errorf("Server launching failure: %v", err)
	}
	log.Println("HTTP server started")

	log.Println("Successfully initalized")
	return nil
}

func deinit() {
	for _, pin := range currentConfig.Pins {
		log.Printf("Disabling output GPIO %d", pin)
		gpio.DisableGPIO(pin)
	}
}
