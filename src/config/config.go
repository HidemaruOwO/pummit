package config

import (
	"os"
	"strconv"

	"github.com/HidemaruOwO/nuts/log"
)

const VERSION string = "v1.2.3"

// const BaseJsonData string = `{
const BASE_JSON_DATA string = `{
	"writeEmoji": true,
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

var IS_DEBUG bool = isDebug()

func isDebug() bool {
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
