package cli

import (
	"fmt"
	"runtime"

	"github.com/HidemaruOwO/pummit/internal/cli/alias"
	"github.com/HidemaruOwO/pummit/internal/config"
	"github.com/HidemaruOwO/pummit/internal/variable"
	"github.com/spf13/cobra"
)

var (
	version bool
)

var rootCmd = &cobra.Command{
	Use:   "pummit",
	Short: "Make the commit message more beautiful in CLI 🎨",
	// Long: `pummit is a tool that helps create consistent git commit messages using emojis and formatters. For detailed documentation, please refer to:
	// https://github.com/HidemaruOwO/pummit`,
	Long: fmt.Sprintf(`pummit v%s %s
  Make the commit message more beautiful in CLI 🎨`, variable.VERSION, runtime.GOARCH),
}

func Execute() error {
	// 初期化
	if err := config.Init(); err != nil {
		return err
	}

	// --versionの時にverionCmdを実行したい
	// fmt.Println(version)

	// rootCmd.PersistentFlags().BoolVarP(&version, "version", "v", false, "Show the version of pummit")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(alias.AddCmd)
	rootCmd.AddCommand(alias.ListCmd)
	rootCmd.AddCommand(alias.DeleteCmd)
	rootCmd.AddCommand(alias.ResetCmd)

	return rootCmd.Execute()
}
