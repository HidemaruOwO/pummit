package config

import (
	"os"

	"github.com/BurntSushi/toml"
	cfg "github.com/HidemaruOwO/pummit/internal/config"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all configuration values in TOML format",
	RunE: func(cmd *cobra.Command, args []string) error {
		encoder := toml.NewEncoder(os.Stdout)
		return encoder.Encode(cfg.CurrentTOMLConfig)
	},
}
