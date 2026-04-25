package cli

import (
	"fmt"

	"github.com/HidemaruOwO/pummit/internal/variable"
	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the version of pummit",
		RunE:  runVersion,
	}
}

func runVersion(cmd *cobra.Command, args []string) error {
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "pummit v%s\n", variable.VERSION)
	return err
}
