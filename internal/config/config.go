package config

import (
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"

	"github.com/HidemaruOwO/pummit/internal/variable"
)

type Config struct {
	UseRawEmoji    bool       `json:"writeEmoji"`
	UseAlias       bool       `json:"useAlias"`
	UseFilesLength bool       `json:"useLimitPathesLength"`
	FilesLength    int        `json:"limitPathesLength"`
	Aliases        [][]string `json:"alias"`
}

var (
	DefaultConfig = Config{}
	CurrentConfig = DefaultConfig
	ConfigPath    string
)

// クロスプラットフォーム対応の設定ディレクトリ取得
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var configDir string
	switch runtime.GOOS {
	case "windows":
		// Windows: %APPDATA%\pummit
		appData := os.Getenv("APPDATA")
		if appData != "" {
			configDir = filepath.Join(appData, "pummit")
		} else {
			// フォールバック: ユーザーホームディレクトリ
			configDir = filepath.Join(home, "pummit")
		}
	default:
		// Unix系 (Linux, macOS): ~/.config/pummit
		configDir = filepath.Join(home, ".config", "pummit")
	}

	return configDir, nil
}

func Init() error {
	// デフォルトの値を読み込む
	json.Unmarshal([]byte(variable.DEFAULT_CONFIG), &DefaultConfig)

	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	ConfigPath = filepath.Join(configDir, "config.json")

	// 自動マイグレーション対応の設定読み込み
	return AutoMigrate()
}

// コンフィグを読み込む
func Load() error {
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		// CurrentConfigにDefaultConfigを代入
		CurrentConfig = DefaultConfig
		return Save() // ファイルが存在しない場合は新規作成
	}

	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &CurrentConfig)
}

// コンフィグを保存
func Save() error {
	data, err := json.MarshalIndent(CurrentConfig, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ConfigPath, data, 0644)
}
