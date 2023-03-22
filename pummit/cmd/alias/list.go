package alias_cmd

import (
	"fmt"

	"github.com/HidemaruOwO/pummit/pummit/lib"
)

func AliasListCmd() {
	fmt.Printf("%s", aliasList())
}

func aliasList() string {
	result := ` 📎 There is aliases
Alias : Prefix : Emoji
----------------------
`
	for _, value := range lib.GetAliasList() {
		emojiList := fmt.Sprintf("  %s : %s : %s\n", value[0], value[1], value[2])
		result = result + emojiList
	}

	return result
}
