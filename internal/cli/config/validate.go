package config

import (
	"fmt"

	infraconfig "github.com/HidemaruOwO/pummit/internal/infra/config"
	"github.com/HidemaruOwO/pummit/internal/usecase"
	"github.com/spf13/cobra"
)

const configErrorCode = 2

type exitError struct {
	code int
	err  error
}

func (e *exitError) Error() string {
	if e == nil || e.err == nil {
		return ""
	}

	return e.err.Error()
}

func (e *exitError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.err
}

func (e *exitError) ExitCode() int {
	if e == nil {
		return 0
	}

	return e.code
}

func newValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate the configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			service := usecase.NewConfigService(infraconfig.NewStore(""))
			if _, err := service.Validate(); err != nil {
				return &exitError{code: configErrorCode, err: err}
			}

			_, err := fmt.Fprintln(cmd.OutOrStdout(), "Configuration is valid")
			return err
		},
	}
}
