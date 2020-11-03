package main

import (
	"errors"
	"fmt"
	"log"

	rr "github.com/markamdev/repico/rest"
)

func main() {
	log.Printf("RePiCo server")
	err := initialize()
	if err != nil {
		// handle error and exit
		log.Fatalln("Failed to launch RePiCo:", errors.Unwrap(err))
	}
}

// initialize is a single initialization (and potential failure) point
func initialize() error {
	err := rr.LaunchServer()
	if err != nil {
		return fmt.Errorf("Server launching failure: %v", err)
	}

	log.Println("Successfully initalized")
	return nil
}
