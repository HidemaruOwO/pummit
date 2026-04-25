package cli

import (
	"fmt"
	"io"
	"runtime"

	aliascmd "github.com/HidemaruOwO/pummit/internal/cli/alias"
	configcmd "github.com/HidemaruOwO/pummit/internal/cli/config"
	"github.com/HidemaruOwO/pummit/internal/variable"
	"github.com/spf13/cobra"
)

func Execute(args []string, stdout, stderr io.Writer) error {
	cmd := NewRootCommand()
	cmd.SetArgs(args)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	return cmd.Execute()
}

func NewRootCommand() *cobra.Command {
	var versionFlag bool

	cmd := &cobra.Command{
		Use:   "pummit",
		Short: "Create consistent git commits with emoji support",
		Long: fmt.Sprintf(
			"pummit v%s %s\n  Rewrite bootstrap for the next implementation phases.",
			variable.VERSION,
			runtime.GOARCH,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			if versionFlag {
				return runVersion(cmd, args)
			}

			if len(args) >= 2 {
				return runCompatCommit(cmd, args)
			}

			return cmd.Help()
		},
		Args: cobra.ArbitraryArgs,
	}

	cmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Show the version of pummit")
	cmd.AddCommand(newVersionCommand())
	cmd.AddCommand(newCommitCommand())
	cmd.AddCommand(configcmd.NewCommand())
	cmd.AddCommand(aliascmd.NewCommand())
	cmd.AddCommand(newDoctorCommand())
	cmd.AddCommand(newMigrateCommand())
	cmd.AddCommand(newMCPCommand())

	return cmd
}
