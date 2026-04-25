package alias

import (
	"fmt"
	"sort"
	"strings"

	"github.com/HidemaruOwO/pummit/legacy/config"
	"github.com/HidemaruOwO/pummit/legacy/logger"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show the list of aliases that have been set",
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.New()
		aliases := config.CurrentTOMLConfig.Alias.Entries

		if len(aliases) == 0 {
			fmt.Println()
			log.Info("🔍 No aliases have been configured yet")
			fmt.Println("💡 Use 'pummit alias add <shortcut> <name> <emoji>' to add your first alias")
			fmt.Println()
			return nil
		}

		// Sort by name for better organization
		sort.SliceStable(aliases, func(i, j int) bool {
			return aliases[i].Name < aliases[j].Name
		})

		// Color functions for visual hierarchy
		title := color.New(color.Bold, color.FgCyan).SprintFunc()
		header := color.New(color.Bold, color.FgWhite).SprintFunc()
		emoji := color.New(color.FgYellow).SprintFunc()
		name := color.New(color.FgGreen).SprintFunc()
		shortcut := color.New(color.FgBlue).SprintFunc()
		separator := color.New(color.FgHiBlack).SprintFunc()

		fmt.Println()
		fmt.Printf("  %s\n", title("📋 Configured Aliases"))
		fmt.Println()

		// Calculate column widths for proper alignment
		maxNameWidth := 4     // "Name"
		maxShortcutWidth := 9 // "Shortcuts"
		for _, entry := range aliases {
			if len(entry.Name) > maxNameWidth {
				maxNameWidth = len(entry.Name)
			}
			shortcutStr := strings.Join(entry.Shortcuts, ", ")
			if len(shortcutStr) > maxShortcutWidth {
				maxShortcutWidth = len(shortcutStr)
			}
		}

		// Add padding
		nameWidth := maxNameWidth + 2
		shortcutWidth := maxShortcutWidth + 2

		// Print header with proper spacing
		fmt.Printf("  %s  %-*s  %s\n",
			header("Emoji"),
			nameWidth, header("Name"),
			header("Shortcuts"))

		// Print separator line with exact alignment
		fmt.Printf("  %s%s%s\n",
			separator("─────"),
			separator(strings.Repeat("─", nameWidth+2)),
			separator(strings.Repeat("─", shortcutWidth)))

		// Print data rows with consistent formatting
		for _, entry := range aliases {
			shortcuts := strings.Join(entry.Shortcuts, ", ")
			fmt.Printf("  %s  %-*s  %s\n",
				emoji(entry.Emoji),
				nameWidth, name(entry.Name),
				shortcut(shortcuts))
		}

		fmt.Println()
		fmt.Printf("  %s %d aliases configured\n",
			separator("✓"), len(aliases))
		fmt.Println()

		return nil
	},
}
