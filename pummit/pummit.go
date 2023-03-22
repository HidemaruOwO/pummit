package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/HidemaruOwO/nuts/log"
	"github.com/HidemaruOwO/pummit/pummit/cmd"
	"github.com/HidemaruOwO/pummit/pummit/cmd/alias"
	"github.com/HidemaruOwO/pummit/pummit/lib"
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

	lib.Init(lib.PlatformPath("pummit"), isDebug)

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
		log.Warnf("Not enough arguments\n")
		os.Exit(0)
	}

	if args[0] == "alias" {
		if len(args) == 1 {
			cmd.HelpCmd()
			fmt.Printf("\n")
			log.Warnf("The alias command requires a second argument\n")
			os.Exit(0)
		} else {
			if args[1] == "add" {
				alias_cmd.AliasAddCmd()
			} else if args[1] == "delete" || args[1] == "del" {
				alias_cmd.AliasDeleteCmd()
			} else if args[1] == "list" {
				alias_cmd.AliasListCmd()
			} else {
				cmd.HelpCmd()
				fmt.Printf("\n")
				log.Warnf("The second argument of the alias command is wrong. Please check the help command\n")
				os.Exit(0)
			}
		}
		os.Exit(0)
	}

	cmd.RootCmd()
}
