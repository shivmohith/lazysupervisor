package main

import (
	"fmt"
	"io"
	"os"

	"github.com/shivmohith/tui-supervisor/supervisord"
	"github.com/shivmohith/tui-supervisor/tui"
	log "github.com/sirupsen/logrus"
)

func main() {
	// nolint:gomnd
	logFileHandler, err := os.OpenFile(
		"/home/shivmohith/go_personal_projects/tui-supervisor/tui.log",
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		0666,
	)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	initializeLogger(logFileHandler)

	client, err := supervisord.NewClient()
	if err != nil {
		log.Fatalf("failed to get new supervisord client because %v", err)
	}

	app := tui.New(client)

	if err := app.Start(); err != nil {
		log.Fatalf("failed to start the tui app because %v", err)
	}
}

func initializeLogger(w io.Writer) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(w)
}
