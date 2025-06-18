package cli

import (
	"fmt"
	"os"

	"github.com/HidemaruOwO/pummit/internal/config"
	"github.com/spf13/cobra"
)

var (
	forceFlag  bool
	dryRunFlag bool
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate configuration from JSON to TOML format",
	Long: `Migrate existing JSON configuration to the new TOML format.

This command will:
1. Check for existing config.json file
2. Convert it to TOML format
3. Create a backup of the original JSON file
4. Save the new config.toml file

Examples:
  pummit migrate                  # Standard migration
  pummit migrate --force          # Overwrite existing TOML config
  pummit migrate --dry-run        # Preview changes without applying`,
	Run: func(cmd *cobra.Command, args []string) {
		// マイグレーション実行
		result, err := config.MigrateConfig(forceFlag, dryRunFlag)
		if err != nil {
			fmt.Printf("Migration failed: %v\n", err)
			os.Exit(1)
		}

		// 結果を表示
		fmt.Println(result.Message)

		if dryRunFlag {
			// ドライラン時の詳細情報
			fmt.Printf("\nDry run results:\n")
			fmt.Printf("  JSON config: %s\n", result.JSONPath)
			fmt.Printf("  TOML config: %s\n", result.TOMLPath)
			if result.ConvertedFrom != "" {
				fmt.Printf("  Source: %s\n", result.ConvertedFrom)
			}
			return
		}

		if result.Success {
			fmt.Printf("\n✅ Migration completed successfully!\n")
			fmt.Printf("   New TOML config: %s\n", result.TOMLPath)

			if result.BackupPath != "" {
				fmt.Printf("   JSON backup: %s\n", result.BackupPath)
			}

			if result.ConvertedFrom == "json" {
				fmt.Printf("\n💡 Your configuration has been converted from JSON to TOML format.\n")
				fmt.Printf("   The application will now use the TOML configuration by default.\n")
			}
		} else {
			fmt.Printf("\n⚠️  %s\n", result.Message)
		}
	},
}

// ロールバックコマンド
var rollbackCmd = &cobra.Command{
	Use:   "rollback [backup-file]",
	Short: "Rollback configuration to JSON format from backup",
	Long: `Rollback configuration from a backup file created during migration.

This will:
1. Restore the JSON configuration from the backup
2. Remove the TOML configuration file
3. Reset the application to use JSON format

Example:
  pummit rollback ~/.config/pummit/config.20240118_143022.bak`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		backupPath := args[0]

		err := config.RollbackConfig(backupPath)
		if err != nil {
			fmt.Printf("Rollback failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Configuration rolled back successfully!\n")
		fmt.Printf("   Restored from: %s\n", backupPath)
		fmt.Printf("   Application will now use JSON configuration.\n")
	},
}

// ステータス確認コマンド
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show configuration file status",
	Long:  `Display the current status of configuration files (JSON/TOML).`,
	Run: func(cmd *cobra.Command, args []string) {
		status, err := config.CheckConfigStatus()
		if err != nil {
			fmt.Printf("Failed to check status: %v\n", err)
			os.Exit(1)
		}

		configDir, _ := config.GetConfigDir()
		fmt.Printf("Configuration directory: %s\n\n", configDir)

		switch status {
		case "both":
			fmt.Printf("📁 Configuration files found:\n")
			fmt.Printf("   ✅ config.toml (active)\n")
			fmt.Printf("   ⚠️  config.json (legacy)\n")
			fmt.Printf("\n💡 TOML format is being used. Consider removing the old JSON file.\n")
		case "toml_only":
			fmt.Printf("📁 Configuration files found:\n")
			fmt.Printf("   ✅ config.toml (active)\n")
			fmt.Printf("\n✨ Using modern TOML configuration format.\n")
		case "json_only":
			fmt.Printf("📁 Configuration files found:\n")
			fmt.Printf("   ⚠️  config.json (legacy)\n")
			fmt.Printf("\n💡 Consider migrating to TOML format: pummit migrate\n")
		case "none":
			fmt.Printf("📁 No configuration files found.\n")
			fmt.Printf("\n💡 Default configuration will be created on next run.\n")
		}
	},
}

func init() {
	// フラグを追加
	migrateCmd.Flags().BoolVar(&forceFlag, "force", false, "Overwrite existing TOML configuration")
	migrateCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "Show what would be done without making changes")

	// サブコマンドを追加
	migrateCmd.AddCommand(rollbackCmd)
	migrateCmd.AddCommand(statusCmd)
}
