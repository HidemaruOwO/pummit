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

Flags:
  --help Show help
  --version Show version`, config.Version)
}
