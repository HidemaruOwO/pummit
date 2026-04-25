package cli

import (
	"fmt"

	rootconfig "github.com/HidemaruOwO/pummit/internal/infra/config"
	"github.com/HidemaruOwO/pummit/internal/usecase"
	"github.com/spf13/cobra"
)

func newDoctorCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Diagnose and validate your pummit environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			service := usecase.NewDoctorService(usecase.NewConfigService(rootconfig.NewStore("")), rootconfig.NewStore(""))
			report := service.Run()
			for _, item := range report.System {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", item.Name, item.Message); err != nil {
					return err
				}
			}
			for _, item := range report.Checks {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s [%s]: %s\n", item.Name, item.Status, item.Message); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
