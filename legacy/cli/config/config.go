package config

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Manage pummit configuration",
	Long: `Manage pummit configuration.

Available subcommands:
  list      Show all configuration values in TOML format
  get       Get a configuration value by key
  set       Update a configuration value and save it
  edit      Open the configuration file in your editor
  validate  Validate the configuration file
  reset     Reset configuration to default values

Examples:
  pummit config list
  pummit config get base.emoji
  pummit config set base.emoji true
  pummit config edit
  pummit config validate
  pummit config reset --force`,
	Run: func(cmd *cobra.Command, args []string) {
		// デフォルト動作: ヘルプを表示
		_ = cmd.Help()
	},
}

func init() {
	Cmd.AddCommand(ListCmd)
	Cmd.AddCommand(GetCmd)
	Cmd.AddCommand(SetCmd)
	Cmd.AddCommand(EditCmd)
	Cmd.AddCommand(ValidateCmd)
	Cmd.AddCommand(ResetCmd)
}
