package config

import (
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"

	"github.com/HidemaruOwO/pummit/internal/variable"
)

// Config is legacy JSON config structure
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

// GetConfigDir gets the configuration directory for the current platform.
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var configDir string
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData != "" {
			configDir = filepath.Join(appData, "pummit")
		} else {
			configDir = filepath.Join(home, ".pummit") // フォールバック時にはホームディレクトリを使用するが、ホームディレクトリでは隠しファイルではないと邪魔なので.pummitを使用
		}
	default:
		configDir = filepath.Join(home, ".config", "pummit")
	}

	return configDir, nil
}

// Init initializes the configuration.
func Init() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	ConfigPath = filepath.Join(configDir, "config.json")
	TOMLConfigPath = filepath.Join(configDir, "config.toml")

	// 1. TOML設定ファイルが存在すれば、それをロードして完了
	if _, err := os.Stat(TOMLConfigPath); err == nil {
		return LoadTOMLConfig()
	}

	// 2. TOMLがなくJSONが存在すれば、マイグレーションを実行
	if _, err := os.Stat(ConfigPath); err == nil {
		// デフォルトのJSON設定をロード（マイグレーションに必要）
		if err := json.Unmarshal([]byte(variable.DEFAULT_CONFIG), &DefaultConfig); err != nil {
			return err
		}
		if err := AutoMigrate(); err != nil {
			return err
		}
		// マイグレーション成功後、作成されたTOMLをロード
		return LoadTOMLConfig()
	}

	// 3. 両方存在しない場合（初回起動）、デフォルトのTOML設定を作成・ロード
	return LoadTOMLConfig()
}

// Load loads the configuration from JSON.
func Load() error {
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		// デフォルトのJSON設定をロード
		if err := json.Unmarshal([]byte(variable.DEFAULT_CONFIG), &DefaultConfig); err != nil {
			return err
		}
		CurrentConfig = DefaultConfig
		return Save()
	}

	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &CurrentConfig)
}

// Save saves the configuration to JSON.
func Save() error {
	data, err := json.MarshalIndent(CurrentConfig, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigPath, data, 0644)
}
