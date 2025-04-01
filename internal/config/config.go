package config

import (
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	UseRawEmoji    bool       `json:"writeEmoji"`
	UseAlias       bool       `json:"useAlias"`
	UseFilesLength bool       `json:"useLimitPathesLength"`
	FilesLength    int        `json:"limitPathesLength"`
	Aliases        [][]string `json:"alias"`
}

// TODO variables/consts.go に移動させる
//
//go:embed config.json
var DEFAULT_CONFIG []byte

var (
	DefaultConfig = Config{}
	CurrentConfig = DefaultConfig
	ConfigPath    string
)

func Init() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// load default config
	json.Unmarshal(DEFAULT_CONFIG, &DefaultConfig)

	configDir := filepath.Join(home, ".config", "pummit")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	ConfigPath = filepath.Join(configDir, "config.json")

	return Load()
}

// コンフィグを読み込む
func Load() error {
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
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
