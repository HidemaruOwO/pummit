package config

import (
	"fmt"
	"os"

	"github.com/HidemaruOwO/pummit/legacy/config"
	"github.com/HidemaruOwO/pummit/legacy/logger"
	"github.com/HidemaruOwO/pummit/legacy/prompt"
	"github.com/spf13/cobra"
)

var EditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open the configuration file in your editor",
	Long: `Open the configuration file in your editor.

If $EDITOR environment variable is set, it will be used.
Otherwise, an interactive editor selector will be shown.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEdit()
	},
}

func runEdit() error {
	log := logger.New()
	configPath := config.TOMLConfigPath

	// $EDITOR 環境変数をチェック
	if editor := os.Getenv("EDITOR"); editor != "" {
		log.Debugf("Using $EDITOR: %s", editor)
		fmt.Printf("Opening %s with %s...\n", configPath, editor)
		return prompt.OpenWithEditor(editor, configPath)
	}

	// $EDITOR が未設定の場合はエディター選択UIを表示
	log.Debug("$EDITOR not set, showing editor selector")
	result, err := prompt.RunEditorSelector(configPath)
	if err != nil {
		return fmt.Errorf("failed to run editor selector: %w", err)
	}

	if result.Canceled {
		log.Info("Edit canceled")
		return nil
	}

	fmt.Printf("Opening %s with %s...\n", configPath, result.Command)
	return prompt.OpenWithEditor(result.Command, configPath)
}
