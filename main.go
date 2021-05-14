package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/markamdev/repico/gpio"
	v2 "github.com/markamdev/repico/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	initLogger()

	logrus.Debugln("RePiCo server")

	ctrl := gpio.CreateController("/sys/class/gpio")
	mainRouter := v2.CreateHandler(ctrl)
	srv := http.Server{Addr: ":8080", Handler: mainRouter}

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
