package alias

import (
	"github.com/HidemaruOwO/pummit/internal/alias"
	"github.com/HidemaruOwO/pummit/pkg/logger"
	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "alias:delete [name]",
	Short: "Delete the specified alias",
	// Short: "指定したエイリアスを削除します",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.New()
		name := args[0]

		if err := alias.Delete(name); err != nil {
			if err == alias.ErrAliasNotFound {
				log.Errorf("The alias '%s' does not exist", name)
				// log.Errorf("エイリアス '%s' は存在しません", name)
				return nil
			}
			return err
		}

		log.Errorf("Deleted alias '%s'", name)
		// log.Infof("エイリアス '%s' を削除しました", name)
		return nil
	},
}
