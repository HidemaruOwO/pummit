package alias

import (
	"errors"
	"fmt"
	"strings"

	rootconfig "github.com/HidemaruOwO/pummit/internal/infra/config"
	rootemoji "github.com/HidemaruOwO/pummit/internal/infra/emoji"
	"github.com/HidemaruOwO/pummit/internal/usecase"
	"github.com/spf13/cobra"
)

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

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alias",
		Short: "Manage emoji aliases for commit prefixes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newAddCommand())
	cmd.AddCommand(newListCommand())
	cmd.AddCommand(newDeleteCommand())
	cmd.AddCommand(newResetCommand())
	return cmd
}

func newService() *usecase.AliasService {
	configService := usecase.NewConfigService(rootconfig.NewStore(""))
	return usecase.NewAliasService(configService, rootemoji.NewCatalog(), rootemoji.NewGitmojiClient())
}

func newAddCommand() *cobra.Command {
	var emojiValue string
	cmd := &cobra.Command{
		Use:   "add [shortcut] [name]",
		Short: "Add a new alias",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := newService().Add(args[0], args[1], emojiValue); err != nil {
				return &exitError{code: 2, err: err}
			}
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "Added alias %q for %q\n", args[0], args[1])
			return err
		},
	}
	cmd.Flags().StringVar(&emojiValue, "emoji", "", "Specify the exact emoji to use")
	return cmd
}

func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured aliases",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := newService().List()
			if err != nil {
				return &exitError{code: 2, err: err}
			}
			for _, entry := range entries {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s %s => %s\n", entry.Emoji, entry.Name, strings.Join(entry.Shortcuts, ",")); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func newDeleteCommand() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "delete [shortcut]",
		Short: "Delete an alias shortcut",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !force {
				return &exitError{code: 4, err: errors.New("alias delete requires --force")}
			}
			if err := newService().Delete(args[0]); err != nil {
				return &exitError{code: 2, err: err}
			}
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "Deleted alias %q\n", args[0])
			return err
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")
	return cmd
}

func newResetCommand() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset aliases to default values",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !force {
				return &exitError{code: 4, err: errors.New("alias reset requires --force")}
			}
			if err := newService().Reset(); err != nil {
				return &exitError{code: 2, err: err}
			}
			_, err := fmt.Fprintln(cmd.OutOrStdout(), "Aliases have been reset to defaults")
			return err
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")
	return cmd
}
