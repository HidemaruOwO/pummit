package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	infraconfig "github.com/HidemaruOwO/pummit/internal/infra/config"
	"github.com/HidemaruOwO/pummit/internal/usecase"
	"github.com/spf13/cobra"
)

func newGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get [key]",
		Short: "Get a configuration value by key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			service := usecase.NewConfigService(infraconfig.NewStore(""))
			value, err := service.Get(args[0])
			if err != nil {
				return &exitError{code: configErrorCode, err: err}
			}
			return printValue(cmd, args[0], value)
		},
	}
}

func printValue(cmd *cobra.Command, key string, value any) error {
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		_, err := fmt.Fprintln(cmd.OutOrStdout())
		return err
	}

	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			_, err := fmt.Fprintln(cmd.OutOrStdout())
			return err
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
		last := key
		if idx := strings.LastIndex(key, "."); idx >= 0 && idx < len(key)-1 {
			last = key[idx+1:]
		}
		return toml.NewEncoder(cmd.OutOrStdout()).Encode(map[string]any{last: value})
	default:
		_, err := fmt.Fprintln(cmd.OutOrStdout(), value)
		return err
	}
}
