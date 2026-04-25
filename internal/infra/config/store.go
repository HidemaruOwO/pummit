package config

import (
	"os"
	"path/filepath"

	domainconfig "github.com/HidemaruOwO/pummit/internal/domain/config"
)

type Store struct {
	path string
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

func (s *Store) ConfigDir() (string, error) {
	if s.path != "" {
		return filepath.Dir(s.path), nil
	}

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(userConfigDir, "pummit"), nil
}

func (s *Store) ConfigPath() (string, error) {
	if s.path != "" {
		return s.path, nil
	}

	dir, err := s.ConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "config.toml"), nil
}

func (s *Store) JSONPath() (string, error) {
	dir, err := s.ConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "config.json"), nil
}

func (s *Store) Load() (domainconfig.Config, error) {
	configPath, err := s.ConfigPath()
	if err != nil {
		return domainconfig.Config{}, err
	}

	if exists(configPath) {
		return s.loadTOML(configPath)
	}

	jsonPath, err := s.JSONPath()
	if err != nil {
		return domainconfig.Config{}, err
	}

	if exists(jsonPath) {
		cfg, err := s.loadJSON(jsonPath)
		if err != nil {
			return domainconfig.Config{}, err
		}
		if err := s.Save(cfg); err != nil {
			return domainconfig.Config{}, err
		}
		return cfg, nil
	}

	cfg := domainconfig.Default()
	if err := s.Save(cfg); err != nil {
		return domainconfig.Config{}, err
	}

	return cfg, nil
}

func (s *Store) Save(cfg domainconfig.Config) error {
	data, err := EncodeTOML(cfg)
	if err != nil {
		return err
	}

	configPath, err := s.ConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0o644)
}

func (s *Store) loadTOML(path string) (domainconfig.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return domainconfig.Config{}, err
	}

	return DecodeTOML(data)
}

func (s *Store) loadJSON(path string) (domainconfig.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return domainconfig.Config{}, err
	}

	return MigrateLegacyJSON(data)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
