package cmd

import (
	"github.com/HidemaruOwO/pummit/pummit/config"
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
  pummit alias [add|del|delete|list]
  
  add: Add alias
    pummit alias add [alias name] [emoji prefix]

  delete: Delete alias
    pummit alias delete [alias name]

  list: Show alias list
    pummit alias list

Flags:
  --help Show help
  --version Show version`, config.Version)
}
