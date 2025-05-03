package cli

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/HidemaruOwO/pummit/internal/cli/alias"
	"github.com/HidemaruOwO/pummit/internal/config"
	// "github.com/HidemaruOwO/pummit/internal/emojis"
	"github.com/HidemaruOwO/pummit/internal/git"
	"github.com/HidemaruOwO/pummit/internal/variable"
	"github.com/spf13/cobra"
)

var (
	version bool
)

var rootCmd = &cobra.Command{
	Use:   "pummit [emoji] [message...]",
	Short: "Make the commit message more beautiful in CLI 🎨",
	// Long: `pummit is a tool that helps create consistent git commit messages using emojis and formatters. For detailed documentation, please refer to:
	// https://github.com/HidemaruOwO/pummit`,
	Long: fmt.Sprintf(`pummit v%s %s
  Make the commit message more beautiful in CLI 🎨`, variable.VERSION, runtime.GOARCH),
	Run: func(cmd *cobra.Command, args []string) {
		// emoji, err := emojis.GetEmojiByName(args[0])
		// if err != nil {
		// 	fmt.Println("Error:", err)
		// 	os.Exit(1)
		// }
		// fmt.Println("Emoji:", emoji)

		if version {
			versionCmd.Run(cmd, args)
			os.Exit(0)
		}

		if len(args) < 2 {
			cmd.Help()
			os.Exit(0)
		}

		cm := git.CommitMessage{
			Emoji:   args[0],
			Message: strings.Join(args[1:], " "),
		}

		git.Commit(cm)
	},
	Args: cobra.ArbitraryArgs,
}

func Execute() error {
	// 初期化
	if err := config.Init(); err != nil {
		return err
	}

	rootCmd.PersistentFlags().BoolVarP(&version, "version", "v", false, "Show the version of pummit")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(alias.AddCmd)
	rootCmd.AddCommand(alias.ListCmd)
	rootCmd.AddCommand(alias.DeleteCmd)
	rootCmd.AddCommand(alias.ResetCmd)

	return rootCmd.Execute()
}
