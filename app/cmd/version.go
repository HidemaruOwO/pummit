package cmd

import (
	"fmt"
	"runtime"

	"github.com/HidemaruOwO/pummit/app/config"
)

func VersionCmd() {
	fmt.Printf("%s\n", version())
}

func version() string {
	return fmt.Sprintf(`pummit %s %s/%s`, config.VERSION, runtime.GOOS, runtime.GOARCH)
}
