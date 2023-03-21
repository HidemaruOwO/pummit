package lib

import (
	"os"
	"path/filepath"

	"github.com/HidemaruOwO/nuts/log"
)

func Init(path string, isDebug bool) {
	path = filepath.Join(path)

	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		log.Debugf(isDebug, "File is exists")
	} else {
	}
}
