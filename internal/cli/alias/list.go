package alias

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/HidemaruOwO/pummit/internal/config"
	"github.com/HidemaruOwO/pummit/pkg/logger"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "alias:list",
	Short: "Show the list of aliases that have been set",
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.New()
		rawAliases := config.CurrentConfig.Aliases

		if len(rawAliases) == 0 {
			log.Info("Aliases have not been set yet")
			return nil
		}

		sort.SliceStable(rawAliases, func(i, j int) bool {
			aliasI := rawAliases[i][0]
			if idx := strings.Index(aliasI, ","); idx != -1 {
				aliasI = aliasI[:idx]
			}
			aliasJ := rawAliases[j][0]
			if idx := strings.Index(aliasJ, ","); idx != -1 {
				aliasJ = aliasJ[:idx]
			}
			return aliasI < aliasJ
		})

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		bold := color.New(color.Bold).SprintFunc()

		fmt.Fprintf(w, "%s | %s | %s\n", bold("Emoji"), bold("Prefix"), bold("Alias"))
		fmt.Fprintln(w, "-------------------")

		for _, item := range rawAliases {
			aliasesStr := item[0]
			prefix := ""
			emoji := ""

			if len(item) >= 3 {
				prefix = item[1]
				emoji = item[2]
			} else if len(item) == 2 {
				emoji = item[1]
			} else if len(item) == 1 {
				continue
			}

			fmt.Fprintf(w, "%s | %s | %s\n", emoji, prefix, aliasesStr)
		}

		w.Flush()

		return nil
	},
}
