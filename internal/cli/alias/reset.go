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
  RunE: func(cmd *cobra.Command, args []string) error {
    log := logger.New()

    confirmFlag, err := cmd.Flags().GetBool("confirm")
    if err != nil {
      return fmt.Errorf("failed to get confirm flag: %w", err)
    }

    confirmed := confirmFlag
    if !confirmed {
      result, err := prompt.Run(
        "Are you sure you want to reset all alias settings?",
      )
      if err != nil {
        return fmt.Errorf("failed to show confirmation prompt: %w", err)
      }
      confirmed = result.Confirmed
    }

    if !confirmed {
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

func init() {
  ResetCmd.Flags().Bool("confirm", false, "Skip confirmation prompt")
}
