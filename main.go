package main

import (
	"os"

	"github.com/HidemaruOwO/pummit/legacy/cli"
	"github.com/HidemaruOwO/pummit/legacy/logger"
)

func main() {
	log := logger.New()

	if err := cli.Execute(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
