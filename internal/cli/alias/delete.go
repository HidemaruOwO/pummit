package alias

import (
	"fmt"

	"github.com/HidemaruOwO/pummit/internal/alias"
	"github.com/HidemaruOwO/pummit/internal/prompt"
	"github.com/HidemaruOwO/pummit/pkg/logger"
	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete the specified alias",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.New()
		name := args[0]

		confirmFlag, err := cmd.Flags().GetBool("confirm")
		if err != nil {
			return fmt.Errorf("failed to get confirm flag: %w", err)
		}

		if !confirmFlag {
			question := fmt.Sprintf("Are you sure you want to delete the alias '%s'?", name)
			result, err := prompt.Run(question)
			if err != nil {
				return fmt.Errorf("failed to show confirmation prompt: %w", err)
			}

			if !result.Confirmed {
				log.Info("Delete canceled")
				return nil
			}
		}

		if err := alias.Delete(name); err != nil {
			if err == alias.ErrAliasNotFound {
				log.Errorf("The alias '%s' does not exist", name)
				return nil
			}
			return err
		}

		log.Infof("Deleted alias '%s'", name)
		return nil
	},
}

func init() {
	DeleteCmd.Flags().Bool("confirm", false, "Skip confirmation prompt")
}
