package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	domainconfig "github.com/HidemaruOwO/pummit/internal/domain/config"
)

func DecodeTOML(data []byte) (domainconfig.Config, error) {
	cfg := domainconfig.Default()
	md, err := toml.Decode(string(data), &cfg)
	if err != nil {
		return domainconfig.Config{}, err
	}

	undecoded := md.Undecoded()
	if len(undecoded) > 0 {
		keys := make([]string, 0, len(undecoded))
		for _, key := range undecoded {
			keys = append(keys, key.String())
		}
		return domainconfig.Config{}, fmt.Errorf("unknown keys: %s", strings.Join(keys, ", "))
	}

	if err := domainconfig.Validate(cfg); err != nil {
		return domainconfig.Config{}, err
	}

	return cfg, nil
}

func EncodeTOML(cfg domainconfig.Config) ([]byte, error) {
	if err := domainconfig.Validate(cfg); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(cfg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
