package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// マイグレーション関連のエラー
type MigrationError struct {
	Op      string
	Path    string
	Message string
	Err     error
}

func (e *MigrationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf(
			"migration %s failed for %s: %s (%v)",
			e.Op, e.Path, e.Message, e.Err,
		)
	}
	return fmt.Sprintf("migration %s failed for %s: %s", e.Op, e.Path, e.Message)
}

// マイグレーション結果
type MigrationResult struct {
	Success       bool
	JSONPath      string
	TOMLPath      string
	BackupPath    string
	Message       string
	ConvertedFrom string
}

// 設定ファイルマイグレーション (JSON → TOML)
func MigrateConfig(force bool, dryRun bool) (*MigrationResult, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, &MigrationError{
			Op:      "get_config_dir",
			Message: "failed to get config directory",
			Err:     err,
		}
	}

	jsonPath := filepath.Join(configDir, "config.json")
	tomlPath := filepath.Join(configDir, "config.toml")

	result := &MigrationResult{
		JSONPath: jsonPath,
		TOMLPath: tomlPath,
	}

	// TOML設定ファイルが既に存在する場合
	if _, err := os.Stat(tomlPath); err == nil {
		if !force {
			result.Message = "TOML config already exists. Use --force to overwrite."
			return result, nil
		}
		result.Message = "Existing TOML config will be overwritten (--force used)."
	}

	// JSON設定ファイルの存在確認
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		// JSONファイルが存在しない場合、デフォルトTOMLを作成
		if dryRun {
			result.Message = "[DRY-RUN] Would create default TOML config " +
				"(no JSON config found)."
			return result, nil
		}

		CurrentTOMLConfig = GetDefaultTOMLConfig()
		if err := SaveTOMLConfig(); err != nil {
			return result, &MigrationError{
				Op:      "save_default_toml",
				Path:    tomlPath,
				Message: "failed to save default TOML config",
				Err:     err,
			}
		}

		result.Success = true
		result.Message = "Created default TOML config (no JSON config found)."
		result.ConvertedFrom = "default"
		return result, nil
	}

	// JSON設定ファイルを読み込み
	var jsonConfig Config
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return result, &MigrationError{
			Op:      "read_json",
			Path:    jsonPath,
			Message: "failed to read JSON config",
			Err:     err,
		}
	}

	if err := json.Unmarshal(data, &jsonConfig); err != nil {
		return result, &MigrationError{
			Op:      "parse_json",
			Path:    jsonPath,
			Message: "failed to parse JSON config",
			Err:     err,
		}
	}

	// JSON → TOML変換
	tomlConfig := ConvertJSONToTOML(jsonConfig)

	if dryRun {
		result.Message = "[DRY-RUN] Would convert JSON config to TOML format."
		result.ConvertedFrom = "json"
		return result, nil
	}

	// TOML設定ファイルを保存
	CurrentTOMLConfig = tomlConfig
	if err := SaveTOMLConfig(); err != nil {
		return result, &MigrationError{
			Op:      "save_toml",
			Path:    tomlPath,
			Message: "failed to save TOML config",
			Err:     err,
		}
	}

	// back up JSON config
	backupPath, err := createBackup(jsonPath)
	if err != nil {
		// バックアップ失敗は警告として扱い、処理は継続
		result.Message = fmt.Sprintf(
			"Migration completed, but backup failed: %v",
			err,
		)
	} else {
		result.BackupPath = backupPath
		result.Message = fmt.Sprintf(
			"Migration completed. JSON config backed up to: %s",
			backupPath,
		)
	}

	result.Success = true
	result.ConvertedFrom = "json"
	return result, nil
}

// 設定ファイル状態をチェック
func CheckConfigStatus() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	jsonPath := filepath.Join(configDir, "config.json")
	tomlPath := filepath.Join(configDir, "config.toml")

	jsonExists := fileExists(jsonPath)
	tomlExists := fileExists(tomlPath)

	switch {
	case tomlExists && jsonExists:
		return "both", nil
	case tomlExists:
		return "toml_only", nil
	case jsonExists:
		return "json_only", nil
	default:
		return "none", nil
	}
}

// createBackup renames the original JSON to a .bak file so it can be recovered
func createBackup(originalPath string) (string, error) {
	backupPath := originalPath + ".bak"
	if err := os.Rename(originalPath, backupPath); err != nil {
		return "", err
	}
	return backupPath, nil
}

// rollback configuration from backup
func RollbackConfig(backupPath string) error {
	if !fileExists(backupPath) {
		return &MigrationError{
			Op:      "rollback",
			Path:    backupPath,
			Message: "backup file does not exist",
		}
	}

	configDir, err := GetConfigDir()
	if err != nil {
		return &MigrationError{
			Op:      "rollback_get_dir",
			Message: "failed to get config directory",
			Err:     err,
		}
	}

	jsonPath := filepath.Join(configDir, "config.json")
	tomlPath := filepath.Join(configDir, "config.toml")

	// バックアップからJSONを復元
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return &MigrationError{
			Op:      "rollback_read_backup",
			Path:    backupPath,
			Message: "failed to read backup file",
			Err:     err,
		}
	}

	if err := os.WriteFile(jsonPath, data, 0644); err != nil {
		return &MigrationError{
			Op:      "rollback_restore_json",
			Path:    jsonPath,
			Message: "failed to restore JSON config",
			Err:     err,
		}
	}

	// TOML設定ファイルを削除
	if fileExists(tomlPath) {
		if err := os.Remove(tomlPath); err != nil {
			return &MigrationError{
				Op:      "rollback_remove_toml",
				Path:    tomlPath,
				Message: "failed to remove TOML config",
				Err:     err,
			}
		}
	}

	return nil
}

// ファイル存在確認のヘルパー関数
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// CreateDefaultTOMLConfig writes the default TOML configuration to disk.
func CreateDefaultTOMLConfig() error {
	CurrentTOMLConfig = GetDefaultTOMLConfig()
	return SaveTOMLConfig()
}

// 自動マイグレーション（設定読み込み時に実行）
func AutoMigrate() error {
	status, err := CheckConfigStatus()
	if err != nil {
		return err
	}

	// TOMLファイルが存在する場合、そちらを優先
	if status == "toml_only" || status == "both" {
		return LoadTOMLConfig()
	}

	// JSONのみ存在する場合、自動変換
	if status == "json_only" {
		result, err := MigrateConfig(false, false)
		if err != nil {
			// マイグレーション失敗時は従来のJSON設定を使用
			return Load()
		}

		if result.Success {
			fmt.Printf("Configuration automatically migrated from JSON to TOML.\n")
			if result.BackupPath != "" {
				fmt.Printf("JSON config backed up to: %s\n", result.BackupPath)
			}
		}
		return nil
	}

	// 設定ファイルが存在しない場合はデフォルトTOML設定を作成
	return CreateDefaultTOMLConfig()
}
