package cli

import (
	"fmt"

	rootconfig "github.com/HidemaruOwO/pummit/internal/infra/config"
	"github.com/HidemaruOwO/pummit/internal/usecase"
	"github.com/spf13/cobra"
)

func newMigrateCommand() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate configuration from JSON to TOML format",
		RunE: func(cmd *cobra.Command, args []string) error {
			service := usecase.NewMigrateService(rootconfig.NewStore(""))
			result, err := service.Migrate(force)
			if err != nil {
				return &ExitError{Code: configErrorCode, Err: err}
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), result.Message)
			return err
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite existing TOML configuration")
	cmd.AddCommand(newMigrateStatusCommand())
	cmd.AddCommand(newMigrateRollbackCommand())
	return cmd
}

func newMigrateStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show configuration file status",
		RunE: func(cmd *cobra.Command, args []string) error {
			service := usecase.NewMigrateService(rootconfig.NewStore(""))
			status, err := service.Status()
			if err != nil {
				return &ExitError{Code: configErrorCode, Err: err}
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "status: %s\nconfig.toml: %s\nconfig.json: %s\n", status.Status, status.ConfigPath, status.JSONPath); err != nil {
				return err
			}
			if status.BackupPath != "" {
				_, err = fmt.Fprintf(cmd.OutOrStdout(), "config.json.bak: %s\n", status.BackupPath)
				return err
			}
			return nil
		},
	}
}

func newMigrateRollbackCommand() *cobra.Command {
	var confirm bool
	cmd := &cobra.Command{
		Use:   "rollback",
		Short: "Rollback configuration to JSON format from backup",
		RunE: func(cmd *cobra.Command, args []string) error {
			service := usecase.NewMigrateService(rootconfig.NewStore(""))
			result, err := service.Rollback(confirm)
			if err != nil {
				return &ExitError{Code: configErrorCode, Err: err}
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), result.Message)
			return err
		},
	}
	cmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm rollback of configuration")
	return cmd
}
