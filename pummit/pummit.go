package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/HidemaruOwO/nuts/log"
	"github.com/HidemaruOwO/pummit/src/cmd"
)

var isDebug bool = false

func main() {
	envDebug := os.Getenv("DEBUG")

	if envDebug != "" {
		var err error
		isDebug, err = strconv.ParseBool(envDebug)
		if err != nil {
			log.ErrorExit(err)
		}
	}

	//set flags
	version := flag.Bool("version", false, "Print version information")
	help := flag.Bool("help", false, "Print help")
	flag.Parse()

	args := flag.Args()

	log.Debugf(isDebug, "args: %s\n", args)

	// use flags
	if *help {
		cmd.HelpCmd()
		os.Exit(0)
	}
	if *version {
		cmd.VersionCmd()
		os.Exit(0)
	}

	if len(args) == 0 {
		cmd.HelpCmd()
		fmt.Printf("\n")
		log.Warnf("Not enough arguments")
		os.Exit(0)
	}

	cmd.RootCmd()
}
