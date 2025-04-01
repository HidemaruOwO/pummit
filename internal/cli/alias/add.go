package alias

import (
	"github.com/HidemaruOwO/pummit/internal/alias"
	"github.com/HidemaruOwO/pummit/pkg/logger"
	"github.com/spf13/cobra"
)

// AddCmd はエイリアス追加コマンドです
var AddCmd = &cobra.Command{
	Use:   "alias:add [name] [emoji]",
	Short: "Add a new alias",
	// Short: "新しいエイリアスを追加します",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.New()
		name := args[0]
		emoji := args[1]

		if err := alias.Add(name, emoji); err != nil {
			if err == alias.ErrAliasExists {
				log.Errorf("The alias '%s' already exists", name)
				// log.Errorf("エイリアス '%s' は既に存在します", name)
				return nil
			}
			return err
		}

		log.Infof("Added alias '%s' as '%s'", name, emoji)
		// log.Infof("エイリアス '%s' を '%s' に設定しました", name, emoji)
		return nil
	},
}
