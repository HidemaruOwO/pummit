package emojis

import (
	"encoding/json"
	"fmt"

	"github.com/HidemaruOwO/pummit/internal/variable"
)

type Emoji struct {
	Emoji string `json:"emoji"`
	Name  string `json:"name"`
}

var (
	mapping   map[string]string
	allEmojis []Emoji
)

func init() {
	mapping = make(map[string]string)

	err := json.Unmarshal([]byte(variable.EMOJIS_JSON), &allEmojis)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse emojis.json: %v", err))
	}

	for _, e := range allEmojis {
		mapping[e.Name] = e.Emoji
	}
}

func GetEmojiByName(name string) (string, error) {
	emoji, ok := mapping[name]
	if !ok {
		return "", fmt.Errorf("emoji not found for name: %s", name)
	}
	return emoji, nil
}

func GetAllEmojis() map[string]string {
	copiedMap := make(map[string]string, len(mapping))
	for k, v := range mapping {
		copiedMap[k] = v
	}
	return copiedMap
}

func GetAllEmojiStructs() []Emoji {
	return allEmojis
}
