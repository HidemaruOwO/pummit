package config

import (
	"errors"
	"fmt"

	infraconfig "github.com/HidemaruOwO/pummit/internal/infra/config"
	"github.com/HidemaruOwO/pummit/internal/usecase"
	"github.com/spf13/cobra"
)

func newResetCommand() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset configuration to default values",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !force {
				return &exitError{code: 4, err: errors.New("config reset requires --force")}
			}
			service := usecase.NewConfigService(infraconfig.NewStore(""))
			if _, err := service.Reset(); err != nil {
				return &exitError{code: configErrorCode, err: err}
			}
			_, err := fmt.Fprintln(cmd.OutOrStdout(), "Configuration has been reset to defaults")
			return err
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")
	return cmd
}
