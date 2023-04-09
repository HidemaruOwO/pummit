package cmd

import (
	"fmt"
)

func ConfigCmd() {
	fmt.Printf("%s\n", configMain())
}

func configMain() string {
	return "Config Cmd"
}
