package main

import (
	"os"

	"github.com/HidemaruOwO/pummit/internal/app"
	"github.com/HidemaruOwO/pummit/legacy/logger"
)

func main() {
	log := logger.New()

	if err := app.Run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		log.Error(err.Error())
		os.Exit(app.ExitCode(err))
	}
}
