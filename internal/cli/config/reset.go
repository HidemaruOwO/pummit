package config

import (
	"fmt"

	cfg "github.com/HidemaruOwO/pummit/internal/config"
	"github.com/HidemaruOwO/pummit/internal/prompt"
	"github.com/HidemaruOwO/pummit/pkg/logger"
	"github.com/spf13/cobra"
)

var ResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to default values",
	RunE: func(cmd *cobra.Command, args []string) error {
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return fmt.Errorf("failed to get force flag: %w", err)
		}

		confirmed := force
		if !confirmed {
			result, err := prompt.Run("Are you sure you want to reset configuration to defaults?")
			if err != nil {
				return fmt.Errorf("failed to show confirmation prompt: %w", err)
			}
			confirmed = result.Confirmed
		}

		if !confirmed {
			logger.New().Info("Reset canceled")
			return nil
		}

		cfg.CurrentTOMLConfig = cfg.GetDefaultTOMLConfig()
		if err := cfg.SaveTOMLConfig(); err != nil {
			return err
		}

		logger.New().Info("Configuration has been reset to defaults")
		return nil
	},
}

func init() {
	ResetCmd.Flags().Bool("force", false, "Skip confirmation prompt")
}
