package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/HidemaruOwO/nuts/log"
	"github.com/HidemaruOwO/pummit/src/cmd"
	alias_cmd "github.com/HidemaruOwO/pummit/src/cmd/alias"
	"github.com/HidemaruOwO/pummit/src/config"
	"github.com/HidemaruOwO/pummit/src/lib"
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

	lib.Init(lib.PlatformPath("pummit"))

	app := &cli.App{
		EnableBashCompletion: true,
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Present() {
				cmd.RootCmd(ctx.Args().Get(0), ctx.Args().Get(1))
			} else {
				cli.ShowAppHelp(ctx)
			}

			return nil
		},
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
						if ctx.Args().Present() {
							alias_cmd.AliasAddCmd(ctx.Args().Get(0), ctx.Args().Get(1))
						} else {
							return fmt.Errorf(
								"Error: %s",
								"Arguments must have both a prefix and an emoji",
							)
						}
						return nil
					},
				},
				{
					Name:  "delete",
					Usage: "Remove aliases for git characters.",
					Action: func(ctx *cli.Context) error {
						if ctx.Args().Present() {
							alias_cmd.AliasDeleteCmd(ctx.Args().Slice())
						} else {
							return fmt.Errorf("Error: %s",
								"delete requires one or more pairs of prefix and emoji in arguments")
						}

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
