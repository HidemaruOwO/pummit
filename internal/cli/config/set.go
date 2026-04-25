package config

import (
	"fmt"

	infraconfig "github.com/HidemaruOwO/pummit/internal/infra/config"
	"github.com/HidemaruOwO/pummit/internal/usecase"
	"github.com/spf13/cobra"
)

func newSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Update a configuration value and save it",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			service := usecase.NewConfigService(infraconfig.NewStore(""))
			if _, err := service.Set(args[0], args[1]); err != nil {
				return &exitError{code: configErrorCode, err: err}
			}
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "Updated %s\n", args[0])
			return err
		},
	}
}
