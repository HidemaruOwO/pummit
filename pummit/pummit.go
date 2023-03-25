package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/HidemaruOwO/nuts/log"
	"github.com/HidemaruOwO/pummit/pummit/cmd"
	alias_cmd "github.com/HidemaruOwO/pummit/pummit/cmd/alias"
	"github.com/HidemaruOwO/pummit/pummit/config"
	"github.com/HidemaruOwO/pummit/pummit/lib"
)

var isDebug bool = config.IsDebug()

func main() {
	// set flags
	version := flag.Bool("version", false, "Print version information")
	help := flag.Bool("help", false, "Print help")
	flag.Parse()

	args := flag.Args()

	log.Debugf(isDebug, "args: %s\n", args)

	lib.Init(lib.PlatformPath("pummit"))

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
		}

		if len(args) != 1 {
			if args[1] == "add" {
				alias_cmd.AliasAddCmd(args)
			} else if args[1] == "delete" || args[1] == "del" {
				alias_cmd.AliasDeleteCmd(args)
			} else if args[1] == "list" {
				alias_cmd.AliasListCmd()
			} else if args[1] == "reset" {
				alias_cmd.AliasResetCmd()
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
