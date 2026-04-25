package config

import (
	"github.com/BurntSushi/toml"
	infraconfig "github.com/HidemaruOwO/pummit/internal/infra/config"
	"github.com/HidemaruOwO/pummit/internal/usecase"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Show all configuration values in TOML format",
		RunE: func(cmd *cobra.Command, args []string) error {
			service := usecase.NewConfigService(infraconfig.NewStore(""))
			cfg, err := service.List()
			if err != nil {
				return &exitError{code: configErrorCode, err: err}
			}
			return toml.NewEncoder(cmd.OutOrStdout()).Encode(cfg)
		},
	}
}
