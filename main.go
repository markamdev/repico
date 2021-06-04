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
	port  = flag.Int("repico-port", 8080, "Repico listening port")
	level = flag.String("log-level", "ERROR", "Log level: ERROR, DEBUG or VERBOSE")
)

func main() {
	initLogger()
	initFlags()

	logrus.Debugln("RePiCo starts listening on port", *port)

	ctrl := gpio.CreateController("/sys/class/gpio")

	newRouter := server.NewHandler()
	gpioSubRouter := newRouter.GetSubRouter("/v2")
	v2.AttachHandlers(gpioSubRouter, ctrl)

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)

	go func() {
		err := newRouter.ServeHTTP(*port)
		if err != nil {
			if err == http.ErrServerClosed {
				logrus.Debugln("Regular server closing")
			} else {
				logrus.Fatalln("HTTP server launching error:", err)
			}
		}
	}()

	<-sigChannel

	logrus.Debugln("Closing server after signal received")
	newRouter.Shutdown(context.Background())
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

	switch *level {
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "VERBOSE":
		logrus.SetLevel(logrus.TraceLevel)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}
