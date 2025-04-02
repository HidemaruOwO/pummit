package alias

import (
	"github.com/HidemaruOwO/pummit/internal/alias"
	"github.com/HidemaruOwO/pummit/pkg/logger"
	"github.com/spf13/cobra"
)

var ResetCmd = &cobra.Command{
	Use:   "alias:reset",
	Short: "Reset all alias settings",
	// Short: "すべてのエイリアス設定をリセットします",
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.New()

		if err := alias.Reset(); err != nil {
			return err
		}

		log.Info("All alias settings have been reset")
		// log.Info("すべてのエイリアス設定をリセットしました")
		return nil
	},
}
