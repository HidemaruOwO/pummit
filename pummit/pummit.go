package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/HidemaruOwO/nuts/log"
	alias_cmd "github.com/HidemaruOwO/pummit/pummit/cmd/alias"
	"github.com/HidemaruOwO/pummit/pummit/config"
	"github.com/urfave/cli/v2"
)

var isDebug bool = config.IsDebug()

func main() {
	envDebug := os.Getenv("DEBUG")

	if envDebug != "" {
		var err error
		isDebug, err = strconv.ParseBool(envDebug)
		if err != nil {
			log.ErrorExit(err)
		}
	}

	app := &cli.App{
		EnableBashCompletion: true,
	}

	app.Name = "pummit"
	app.Usage = "pummit <emoji prefix> <subject>"
	app.Description = "Easily create nicely formatted commit messages "

	app.Version = fmt.Sprintf("%s %s", config.Version, runtime.GOARCH)
	app.Commands = []*cli.Command{
		{
			Name:  "alias",
			Usage: "Subcommand to manage aliases for git characters.",
			Subcommands: []*cli.Command{
				{
					Name:  "add",
					Usage: "Add an alias for the git character.",
					Action: func(ctx *cli.Context) error {
						alias_cmd.AliasAddCmd(ctx.Args().Slice())
						return nil
					},
				},
				{
					Name:  "delete",
					Usage: "Remove aliases for git characters.",
					Action: func(ctx *cli.Context) error {
						alias_cmd.AliasDeleteCmd(ctx.Args().Slice())
						return nil
					},
				},
				{
					Name:  "list",
					Usage: "Displays a list of aliases for git characters.",
					Action: func(ctx *cli.Context) error {
						alias_cmd.AliasListCmd()
						return nil
					},
				},
				{
					Name:  "reset",
					Usage: "Reset the list of aliases for git characters.",
					Action: func(ctx *cli.Context) error {
						alias_cmd.AliasResetCmd()
						return nil
					},
				},
			},
		},
		{
			Name:  "complete",
			Usage: "Outputs a completion script.",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "bash",
					Usage: "Outputs a completion script for bash.",
				},
				&cli.BoolFlag{
					Name:  "zsh",
					Usage: "Outputs a completion script for zsh.",
				},
				&cli.BoolFlag{
					Name:  "fish",
					Usage: "Outputs a completion script for fish.",
				},
			},
			Action: func(ctx *cli.Context) error {
				if ctx.Bool("bash") {
					bashComplete()
				}

				if ctx.Bool("zsh") {
					zshComplete()
				}

				if ctx.Bool("fish") {
					fishComplete()
				}

				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
