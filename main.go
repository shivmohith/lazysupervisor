package main

import (
	"io"
	"os"

	"github.com/shivmohith/lazysupervisor/supervisord"
	"github.com/shivmohith/lazysupervisor/tui"
	log "github.com/sirupsen/logrus"
)

func main() {
	initializeLogger(os.Stdout)

	client, err := supervisord.NewClient()
	if err != nil {
		log.Fatalf("getting new supervisord client: %v", err)
	}

	app := tui.New(client)

	app.BuildLayout()

	if err := app.Start(); err != nil {
		log.Fatalf("starting the tui app: %v", err)
	}
}

func initializeLogger(w io.Writer) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(w)
}
