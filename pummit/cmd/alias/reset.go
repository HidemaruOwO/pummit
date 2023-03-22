package alias_cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/HidemaruOwO/nuts/log"
	"github.com/HidemaruOwO/pummit/pummit/config"
	"github.com/HidemaruOwO/pummit/pummit/lib"
	"github.com/fatih/color"
)

func AliasResetCmd() {
	aliasReset()
}

func aliasReset() {
	fmt.Printf("%s May I reset the aliases? :(Y/n) ", color.HiRedString(">"))
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.ErrorExit(err)
	}

	input = strings.TrimSpace(input)
	if input == "y" || input == "Y" {
		var alias lib.Alias
		if err := json.Unmarshal([]byte(config.BaseJsonData), &alias); err != nil {
			log.Criticalf("JSON encoding failed\n")
			log.ErrorExit(err)
		}
		lib.WriteConfig(alias)
		log.Infof("Alias reseted\n")
		fmt.Printf("\n")
		AliasListCmd()
	} else if input == "n" || input == "N" {
		log.Infof("The alias was not reset because N was selected\n")
	} else {
		log.Warnf("Please enter Y on n\n")
	}
}
