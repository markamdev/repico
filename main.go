package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/markamdev/repico/gpio"
	"github.com/markamdev/repico/server"
	v2 "github.com/markamdev/repico/v2"
	"github.com/namsral/flag"
	"github.com/sirupsen/logrus"
)

var (
	port = flag.Int("repico-port", 8080, "Repico listening port")
)

func main() {
	initFlags()
	initLogger()

	logrus.Debugln("RePiCo server")

	ctrl := gpio.CreateController("/sys/class/gpio")
	mainRouter := v2.CreateHandler(ctrl)

	srv, err := server.Create(*port, mainRouter)
	if err != nil {
		logrus.Fatalln("Server cration error:", err)
	}

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				logrus.Debugln("Regular server closing")
			} else {
				logrus.Fatalln("HTTP server launching error:", err)
			}
		}
	}()

	_ = <-sigChannel

	logrus.Debugln("Closing server after signal received")
	srv.Shutdown(context.Background())
}

func initLogger() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}

func initFlags() {
	flag.Parse()
}
