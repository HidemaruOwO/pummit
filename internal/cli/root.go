package cli

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/HidemaruOwO/pummit/internal/cli/alias"
	"github.com/HidemaruOwO/pummit/internal/config"
	"github.com/HidemaruOwO/pummit/internal/git"
	"github.com/HidemaruOwO/pummit/internal/variable"
	"github.com/HidemaruOwO/pummit/pkg/gitmoji"
	"github.com/HidemaruOwO/pummit/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	version     bool
	offlineMode bool
)

var rootCmd = &cobra.Command{
	Use:   "pummit [emoji] [message...]",
	Short: "Make the commit message more beautiful in the CLI 🎨",
	// Long: `pummit is a tool that helps create consistent git commit messages using emojis and formatters. For detailed documentation, please refer to:
	// https://github.com/HidemaruOwO/pummit`,
	Long: fmt.Sprintf(`pummit v%s %s
  Make your commit messages beautiful, consistent, and meaningful with emoji support and smart automation 🎨`, variable.VERSION, runtime.GOARCH),
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

		log := logger.New()

		// グローバル--offlineフラグまたは自動検出でオフラインモードを制御
		if offlineMode {
			log.Info("Offline mode enabled via --offline flag")
		} else if !gitmoji.IsOnline() {
			// ネットワークが利用できない場合、自動的にオフラインモードを有効化
			log.Info("Network appears to be offline, enabling offline mode automatically")
			offlineMode = true
		}

		cm := git.CommitMessage{
			Emoji:   args[0],
			Message: strings.Join(args[1:], " "),
		}

		git.CommitWithOfflineMode(cm, offlineMode)
	},
	Args: cobra.ArbitraryArgs,
}

func Execute() error {
	// 初期化
	if err := config.Init(); err != nil {
		return err
	}

	rootCmd.PersistentFlags().BoolVarP(&version, "version", "v", false, "Show the version of pummit")
	rootCmd.PersistentFlags().BoolVar(&offlineMode, "offline", false, "Run with offline mode（disable calling Gitmoji API）")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(alias.AddCmd)
	rootCmd.AddCommand(alias.ListCmd)
	rootCmd.AddCommand(alias.DeleteCmd)
	rootCmd.AddCommand(alias.ResetCmd)

	return rootCmd.Execute()
}
