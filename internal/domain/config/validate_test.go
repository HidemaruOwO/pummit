package config

import (
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name string
		mut  func(*Config)
		want string
	}{
		{name: "valid default", mut: nil, want: ""},
		{name: "negative filesLength", mut: func(cfg *Config) { cfg.Base.FilesLength = -1 }, want: "base.filesLength"},
		{name: "unknown language", mut: func(cfg *Config) { cfg.Locale.Language = "fr" }, want: "locale.language"},
		{name: "duplicate alias", mut: func(cfg *Config) { cfg.Alias.Entries[1].Shortcuts = append(cfg.Alias.Entries[1].Shortcuts, "s") }, want: "duplicated"},
		{name: "bad branch regex", mut: func(cfg *Config) { cfg.BranchMapping.Rules[0].Pattern = "(" }, want: "branchMapping.rules[0].pattern"},
		{name: "missing default template", mut: func(cfg *Config) { cfg.Templates.DefaultTemplate = "missing" }, want: "defaultTemplate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Default()
			if tt.mut != nil {
				tt.mut(&cfg)
			}

			err := Validate(cfg)
			if tt.want == "" {
				if err != nil {
					t.Fatalf("Validate returned error: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error containing %q", tt.want)
			}

			if got := err.Error(); got == "" || !strings.Contains(got, tt.want) {
				t.Fatalf("error = %q, want substring %q", got, tt.want)
			}
		})
	}
}
