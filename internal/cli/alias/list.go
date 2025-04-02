package alias

import (
	"fmt"
	"sort"

	"github.com/HidemaruOwO/pummit/internal/alias"
	"github.com/HidemaruOwO/pummit/pkg/logger"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "alias:list",
	Short: "Show the list of aliases that have been set",
	// Short: "設定されているエイリアスの一覧を表示します",
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.New()
		aliases := alias.List()

		if len(aliases) == 0 {
			log.Info("Aliases have not been set yet")
			// log.Info("エイリアスはまだ設定されていません")
			return nil
		}

		// ソート
		keys := make([]string, 0, len(aliases))
		for k := range aliases {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		bold := color.New(color.Bold).SprintFunc()
		fmt.Printf("%s\t%s\n", bold("Alias"), bold("Emoji"))
		// fmt.Printf("%s\t%s\n", bold("エイリアス"), bold("絵文字"))
		fmt.Println("----------------------")

		for _, k := range keys {
			fmt.Printf("%s\t%s\n", k, aliases[k])
		}

		return nil
	},
}
