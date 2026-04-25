package usecase

import (
	"fmt"
	"os"
	"path/filepath"

	domainconfig "github.com/HidemaruOwO/pummit/internal/domain/config"
	rootconfig "github.com/HidemaruOwO/pummit/internal/infra/config"
)

type MigrationStatus string

const (
	MigrationStatusNone     MigrationStatus = "none"
	MigrationStatusJSONOnly MigrationStatus = "json_only"
	MigrationStatusTOMLOnly MigrationStatus = "toml_only"
	MigrationStatusBoth     MigrationStatus = "both"
)

type MigrationResult struct {
	Status     MigrationStatus
	ConfigPath string
	JSONPath   string
	BackupPath string
	Message    string
}

type MigrateService struct {
	store *rootconfig.Store
}

func NewMigrateService(store *rootconfig.Store) *MigrateService {
	return &MigrateService{store: store}
}

func (s *MigrateService) Status() (MigrationResult, error) {
	configPath, err := s.store.ConfigPath()
	if err != nil {
		return MigrationResult{}, err
	}
	jsonPath, err := s.store.JSONPath()
	if err != nil {
		return MigrationResult{}, err
	}
	backupPath := jsonPath + ".bak"

	status := MigrationStatusNone
	hasTOML := existsFile(configPath)
	hasJSON := existsFile(jsonPath)
	switch {
	case hasTOML && hasJSON:
		status = MigrationStatusBoth
	case hasTOML:
		status = MigrationStatusTOMLOnly
	case hasJSON:
		status = MigrationStatusJSONOnly
	}

	return MigrationResult{Status: status, ConfigPath: configPath, JSONPath: jsonPath, BackupPath: backupPath}, nil
}

func (s *MigrateService) Migrate(force bool) (MigrationResult, error) {
	status, err := s.Status()
	if err != nil {
		return MigrationResult{}, err
	}

	if existsFile(status.ConfigPath) && !force {
		return MigrationResult{}, fmt.Errorf("config.toml already exists; use --force to overwrite")
	}

	if existsFile(status.JSONPath) {
		data, err := os.ReadFile(status.JSONPath)
		if err != nil {
			return MigrationResult{}, err
		}
		cfg, err := rootconfig.MigrateLegacyJSON(data)
		if err != nil {
			return MigrationResult{}, err
		}
		if err := s.store.Save(cfg); err != nil {
			return MigrationResult{}, err
		}
		if existsFile(status.BackupPath) {
			if err := os.Remove(status.BackupPath); err != nil {
				return MigrationResult{}, err
			}
		}
		if err := os.Rename(status.JSONPath, status.BackupPath); err != nil {
			return MigrationResult{}, err
		}
		status.Message = "Migrated config.json to config.toml"
		return status, nil
	}

	if err := s.store.Save(domainconfig.Default()); err != nil {
		return MigrationResult{}, err
	}
	status.Message = "Created default config.toml"
	return status, nil
}

func (s *MigrateService) Rollback(confirm bool) (MigrationResult, error) {
	if !confirm {
		return MigrationResult{}, fmt.Errorf("migrate rollback requires --confirm")
	}

	status, err := s.Status()
	if err != nil {
		return MigrationResult{}, err
	}

	if !existsFile(status.BackupPath) {
		return MigrationResult{}, fmt.Errorf("backup file not found: %s", status.BackupPath)
	}

	if existsFile(status.JSONPath) {
		if err := os.Remove(status.JSONPath); err != nil {
			return MigrationResult{}, err
		}
	}
	if err := os.Rename(status.BackupPath, status.JSONPath); err != nil {
		return MigrationResult{}, err
	}
	if existsFile(status.ConfigPath) {
		if err := os.Remove(status.ConfigPath); err != nil {
			return MigrationResult{}, err
		}
	}
	status.Message = "Rolled back config.toml to config.json"
	return status, nil
}

func existsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func ConfigPathsFromHome(home string) (string, string, string) {
	configDir := filepath.Join(home, ".config", "pummit")
	return filepath.Join(configDir, "config.toml"), filepath.Join(configDir, "config.json"), filepath.Join(configDir, "config.json.bak")
}
