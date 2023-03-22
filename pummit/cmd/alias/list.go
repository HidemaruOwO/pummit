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
	ch := make(chan string)
	for _, value := range lib.GetAliasList() {
		go func(value []string) {
			emojiList := fmt.Sprintf("  %s : %s : %s\n", value[0], value[1], value[2])
			ch <- emojiList
		}(value)
	}

	for range lib.GetAliasList() {
		result = result + <-ch
	}

	return result
}
