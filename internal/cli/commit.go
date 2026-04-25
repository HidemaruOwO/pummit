package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	rootconfig "github.com/HidemaruOwO/pummit/internal/infra/config"
	rootemoji "github.com/HidemaruOwO/pummit/internal/infra/emoji"
	rootgit "github.com/HidemaruOwO/pummit/internal/infra/git"
	"github.com/HidemaruOwO/pummit/internal/usecase"
	"github.com/spf13/cobra"
)

const (
	configErrorCode = 2
	gitErrorCode    = 3
	userErrorCode   = 4
)

type commitResult interface {
	GetCommitted() bool
	GetMessage() string
}

func newCommitCommand() *cobra.Command {
	var emojiInput string
	var autoEmoji bool

	cmd := &cobra.Command{
		Use:   "commit [message...]",
		Short: "Create a git commit with emoji support",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			message := strings.TrimSpace(strings.Join(args, " "))
			if message == "" {
				return &ExitError{Code: userErrorCode, Err: errors.New("commit message is required")}
			}

			if autoEmoji && emojiInput != "" {
				return &ExitError{Code: userErrorCode, Err: errors.New("--emoji and --auto-emoji cannot be used together")}
			}

			if !autoEmoji && emojiInput == "" {
				return &ExitError{Code: userErrorCode, Err: errors.New("either --emoji or --auto-emoji is required")}
			}

			return runCommit(cmd, emojiInput, autoEmoji, message)
		},
	}

	cmd.Flags().StringVar(&emojiInput, "emoji", "", "Emoji name or alias to use for the commit")
	cmd.Flags().BoolVar(&autoEmoji, "auto-emoji", false, "Resolve the emoji from branchMapping")
	return cmd
}

func runCompatCommit(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return cmd.Help()
	}

	return runCommit(cmd, args[0], false, strings.Join(args[1:], " "))
}

func runCommit(cmd *cobra.Command, emojiInput string, autoEmoji bool, message string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	service := usecase.NewCommitService(
		rootconfig.NewStore(""),
		rootgit.NewClient(wd),
		rootemoji.NewCatalog(),
		rootemoji.NewGitmojiClient(),
	)

	var result commitResult
	if autoEmoji {
		autoResult, err := service.CommitAuto(message)
		if err != nil {
			return classifyCommitError(err)
		}
		result = autoResult
	} else {
		explicitResult, err := service.CommitExplicit(emojiInput, message)
		if err != nil {
			return classifyCommitError(err)
		}
		result = explicitResult
	}

	if !result.GetCommitted() {
		_, err := fmt.Fprintln(cmd.OutOrStdout(), "No staged files to commit")
		return err
	}

	_, err = fmt.Fprintln(cmd.OutOrStdout(), result.GetMessage())
	return err
}

func classifyCommitError(err error) error {
	msg := err.Error()
	if strings.Contains(msg, "meta.") || strings.Contains(msg, "base.") || strings.Contains(msg, "unknown keys") || strings.Contains(msg, "templates.") || strings.Contains(msg, "branchMapping.") || strings.Contains(msg, "alias ") || strings.Contains(msg, "locale.") || strings.Contains(msg, "interactive.") || strings.Contains(msg, "scope.") {
		return &ExitError{Code: configErrorCode, Err: err}
	}

	if strings.Contains(msg, "git ") {
		return &ExitError{Code: gitErrorCode, Err: err}
	}

	return err
}
