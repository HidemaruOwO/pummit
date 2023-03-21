package cmd

import (
	"fmt"
	"runtime"

	"github.com/HidemaruOwO/pummit/pummit/config"
)

func VersionCmd() {
	fmt.Printf("%s\n", version())
}

func version() string {
	return fmt.Sprintf(`pummit %s %s/%s`, config.Version, runtime.GOOS, runtime.GOARCH)
}
