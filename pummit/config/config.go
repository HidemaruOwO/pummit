package config

import (
	"os"
	"strconv"

	"github.com/HidemaruOwO/nuts/log"
)

const Version string = "v1.1.2"
const BaseJsonData string = `{
	"writeEmojiPrefix": true,
	"useAlias": true,
  "alias": [  
  ["s", "sparkles", "✨"],
  ["t", "tada", "🎉"],
  ["r", "recycle", "♻️"],
  ["b", "bug", "🐛"],
  ["e", "eyes", "👀"],
  ["d", "books", "📚"],
  ["a", "art", "🎨"],
  ["h", "horse", "🐎"],
  ["w", "wrench", "🔧"],
  ["l", "rotating_light", "🚨"],
  ["p", "hankey", "💩"],
  ["wb", "wastebasket", "🗑️"],
  ["c", "construction", "🚧"],
  ["sm", "snowman", "☃️"]
]
}`

func IsDebug() bool {
	isDebug := false

	envDebug := os.Getenv("DEBUG")
	if envDebug != "" {
		var err error
		isDebug, err = strconv.ParseBool(envDebug)
		if err != nil {
			log.ErrorExit(err)
		}
	}

	return isDebug
}
