package cmd

import (
	"github.com/HidemaruOwO/pummit/src/config"
	//  "github.com/fatih/color"
	"fmt"
)

func HelpCmd() {
	fmt.Printf("%s\n", help())
}

func help() string {
	return fmt.Sprintf(`🚛 pummit - %s

Usage:
  pummit [emoji prefix] [subject]
  Create prettier commit message

Subcommand:
  pummit alias [add|del|delete|list|reset]

  add: Add alias
    pummit alias add [alias name] [emoji prefix]

  delete: Delete alias
    pummit alias delete [...alias name]
    Flag:
      --all : Delete all aliases
      pummit alias delete --all

  list: Show alias list
    pummit alias list

  reset: Reset aliases
    pummit alias reset

Flags:
  --help Show help
  --version Show version`, config.Version)
}
