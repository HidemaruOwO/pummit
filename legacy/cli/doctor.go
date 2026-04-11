package cli

import (
	"fmt"
	"os"

	"github.com/HidemaruOwO/pummit/legacy/doctor"
	"github.com/HidemaruOwO/pummit/legacy/logger"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose and validate your pummit environment",
	Long: `The doctor command performs comprehensive diagnostics of your pummit environment.

It checks:
- Git configuration (user.name, user.email)
- Git repository status
- Configuration file validity (JSON/TOML)
- Network connectivity (Gitmoji API)
- File system permissions
- System information

This helps identify and resolve common issues with your pummit setup.`,
	Run: runDoctorCommand,
}

// runDoctorCommand doctor コマンドの実行ロジック
func runDoctorCommand(cmd *cobra.Command, args []string) {
	log := logger.New()

	// ヘッダー表示
	fmt.Println(color.New(color.FgCyan, color.Bold).Sprint("🩺 Pummit Doctor - Environment Diagnostics"))
	fmt.Println()

	// システム情報の表示
	systemInfo := doctor.GetSystemInfo()
	fmt.Print(systemInfo.FormatSystemInfo())
	fmt.Println()

	// 診断チェックの実行
	fmt.Println(color.New(color.FgYellow, color.Bold).Sprint("Running diagnostic checks..."))
	fmt.Println()

	results := doctor.RunAllChecks()

	// 結果の表示
	errorCount := 0
	warningCount := 0

	for _, result := range results {
		status := formatStatus(result.Status)
		fmt.Printf("%-25s %s %s\n", result.Name+":", status, result.Message)

		// 提案がある場合は表示
		if len(result.Suggestions) > 0 {
			for _, suggestion := range result.Suggestions {
				fmt.Printf("  💡 %s\n", suggestion)
			}
		}

		// エラーと警告の数をカウント
		switch result.Status {
		case "ERROR":
			errorCount++
		case "WARNING":
			warningCount++
		}

		fmt.Println()
	}

	// 総括の表示
	fmt.Println(color.New(color.FgCyan, color.Bold).Sprint("Diagnostic Summary:"))

	if errorCount == 0 && warningCount == 0 {
		fmt.Println(color.New(color.FgGreen).Sprint("✅ All checks passed! Your pummit environment is healthy."))
	} else {
		if errorCount > 0 {
			fmt.Printf("%s Found %d error(s) that need immediate attention.\n",
				color.New(color.FgRed).Sprint("❌"), errorCount)
		}
		if warningCount > 0 {
			fmt.Printf("%s Found %d warning(s) that may affect functionality.\n",
				color.New(color.FgYellow).Sprint("⚠️ "), warningCount)
		}

		if errorCount > 0 {
			log.Error("Please resolve the errors above before using pummit.")
			os.Exit(1)
		}
	}
}

// formatStatus ステータスに応じて色付きの表示文字列を返す
func formatStatus(status string) string {
	switch status {
	case "OK":
		return color.New(color.FgGreen).Sprint("✅ OK")
	case "WARNING":
		return color.New(color.FgYellow).Sprint("⚠️  WARN")
	case "ERROR":
		return color.New(color.FgRed).Sprint("❌ ERROR")
	default:
		return color.New(color.FgWhite).Sprint("❓ UNKNOWN")
	}
}
