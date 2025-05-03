package cli

import (
	"fmt"
	"github.com/HidemaruOwO/pummit/internal/variable"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of pummit",
	// Short: "pummitのバージョンを表示します",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("pummit v%s\n", variable.VERSION)
	},
}
