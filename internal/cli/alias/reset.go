package alias

import (
	"fmt"

	"github.com/HidemaruOwO/pummit/internal/alias"
	"github.com/HidemaruOwO/pummit/internal/prompt"
	"github.com/HidemaruOwO/pummit/pkg/logger"
	"github.com/spf13/cobra"
)

var ResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset all alias settings",
	// Short: "すべてのエイリアス設定をリセットします",
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.New()

		result, err := prompt.Run("Are you sure you want to reset all alias settings?")
		if err != nil {
			return fmt.Errorf("failed to show confirmation prompt: %w", err)
		}

		if !result.Confirmed {
			log.Info("Reset canceled")
			return nil
		}

		if err := alias.Reset(); err != nil {
			return err
		}

		log.Info("All alias settings have been reset")
		return nil
	},
}
