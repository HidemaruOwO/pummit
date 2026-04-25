package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	cfg "github.com/HidemaruOwO/pummit/legacy/config"
	"github.com/spf13/cobra"
)

var ValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(cfg.TOMLConfigPath)
		if err != nil {
			return fmt.Errorf("failed to read config: %w", err)
		}

		var parsed cfg.TOMLConfig
		if err := toml.Unmarshal(data, &parsed); err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}

		if err := ValidateRequiredFields(parsed); err != nil {
			return err
		}

		fmt.Println("✅ Configuration is valid")
		return nil
	},
}

// ValidateRequiredFields validates the configuration for required fields and constraints.
func ValidateRequiredFields(c cfg.TOMLConfig) error {
	if c.Meta.Version == "" {
		return fmt.Errorf("meta.version is required")
	}
	if c.Templates.Enabled {
		if c.Templates.DefaultTemplate == "" {
			return fmt.Errorf("templates.defaultTemplate is required when templates.enabled is true")
		}
		if _, ok := c.Templates.Definitions[c.Templates.DefaultTemplate]; !ok {
			return fmt.Errorf("templates.defaultTemplate '%s' is not defined", c.Templates.DefaultTemplate)
		}
	}
	if c.Scope.HistoryLimit < 0 {
		return fmt.Errorf("scope.historyLimit must be >= 0")
	}
	if c.Base.FilesLength < 0 {
		return fmt.Errorf("base.filesLength must be >= 0")
	}
	return nil
}
