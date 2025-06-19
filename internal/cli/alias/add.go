package alias

import (
	"errors"

	"github.com/HidemaruOwO/pummit/internal/alias"
	"github.com/HidemaruOwO/pummit/internal/emojis"
	"github.com/HidemaruOwO/pummit/pkg/logger"
	"github.com/spf13/cobra"
)

var emojiFlag string

var AddCmd = &cobra.Command{
	Use:   "add [name] [prefix]",
	Short: "Add a new alias. Optionally specify the exact emoji with --emoji.",
	Long: `Add a new alias for an emoji prefix.
If the --emoji flag is provided, the specified emoji will be used directly.
Otherwise, the emoji will be inferred from the [prefix] argument.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.New()
		name := args[0]
		prefix := args[1]
		var actualEmoji string

		if emojiFlag != "" {
			actualEmoji = emojiFlag
			log.Debugf("Using specified emoji from --emoji flag: %s", actualEmoji)
		} else {
			log.Debugf("Inferring emoji from prefix: %s", prefix)
			inferredEmoji, err := emojis.GetEmojiByName(prefix)
			if err != nil {
				log.Errorf("Failed to find emoji for prefix '%s': %v. Use --emoji flag to specify it directly.", prefix, err)
				return nil
			}
			actualEmoji = inferredEmoji
			log.Debugf("Inferred emoji: %s", actualEmoji)
		}

		if err := alias.Add(name, prefix, actualEmoji); err != nil {
			if errors.Is(err, alias.ErrAliasExists) {
				log.Errorf("Failed to add alias: %v. The alias '%s' might already exist, or the emoji '%s' might be associated with a different prefix.", err, name, actualEmoji)
				return nil
			}
			log.Errorf("Failed to add alias: %v", err)
			return nil
		}

		log.Infof("Added alias '%s' for prefix '%s' with emoji '%s'", name, prefix, actualEmoji)
		return nil
	},
}

func init() {
	AddCmd.Flags().StringVar(&emojiFlag, "emoji", "", "Specify the exact emoji to use, overriding inference from prefix.")
}
