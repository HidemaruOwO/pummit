package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	cfg "github.com/HidemaruOwO/pummit/legacy/config"
	"github.com/HidemaruOwO/pummit/legacy/logger"
	"github.com/spf13/cobra"
)

var SetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Update a configuration value and save it",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := strings.Split(args[0], ".")
		value := args[1]

		if err := setValue(reflect.ValueOf(&cfg.CurrentTOMLConfig), path, value); err != nil {
			return err
		}

		if err := ValidateRequiredFields(cfg.CurrentTOMLConfig); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		if err := cfg.SaveTOMLConfig(); err != nil {
			return err
		}

		logger.New().Infof("Updated %s", args[0])
		return nil
	},
}

func setValue(v reflect.Value, path []string, raw string) error {
	v = normalizeValue(v)
	if len(path) == 0 {
		return fmt.Errorf("key is required")
	}

	switch v.Kind() {
	case reflect.Struct:
		field, ok := fieldByTomlTag(v, path[0])
		if !ok {
			return fmt.Errorf("unknown key: %s", strings.Join(path, "."))
		}
		if len(path) == 1 {
			return assignPrimitive(field, raw)
		}
		return setValue(field, path[1:], raw)

	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return fmt.Errorf("invalid map key type for %s", strings.Join(path, "."))
		}
		key := reflect.ValueOf(path[0])
		val := v.MapIndex(key)
		if !val.IsValid() {
			return fmt.Errorf("key not found: %s", strings.Join(path, "."))
		}
		valCopy := reflect.New(normalizeValue(val).Type()).Elem()
		valCopy.Set(normalizeValue(val))
		if len(path) == 1 {
			if err := assignPrimitive(valCopy, raw); err != nil {
				return err
			}
		} else {
			if err := setValue(valCopy, path[1:], raw); err != nil {
				return err
			}
		}
		v.SetMapIndex(key, valCopy)
		return nil

	case reflect.Slice, reflect.Array:
		index, err := strconv.Atoi(path[0])
		if err != nil {
			return fmt.Errorf("invalid index '%s'", path[0])
		}
		if index < 0 || index >= v.Len() {
			return fmt.Errorf("index out of range: %d", index)
		}
		elem := v.Index(index)
		elemCopy := reflect.New(normalizeValue(elem).Type()).Elem()
		elemCopy.Set(normalizeValue(elem))
		if len(path) == 1 {
			if err := assignPrimitive(elemCopy, raw); err != nil {
				return err
			}
		} else {
			if err := setValue(elemCopy, path[1:], raw); err != nil {
				return err
			}
		}
		elem.Set(elemCopy)
		return nil
	default:
		if len(path) == 1 {
			return assignPrimitive(v, raw)
		}
		return fmt.Errorf("cannot descend into %s", strings.Join(path, "."))
	}
}

func assignPrimitive(v reflect.Value, raw string) error {
	v = normalizeValue(v)
	if !v.CanSet() {
		return fmt.Errorf("cannot set value")
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(raw)
		return nil
	case reflect.Bool:
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			return fmt.Errorf("invalid bool: %s", raw)
		}
		v.SetBool(parsed)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid int: %s", raw)
		}
		v.SetInt(parsed)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		parsed, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid uint: %s", raw)
		}
		v.SetUint(parsed)
		return nil
	case reflect.Float32, reflect.Float64:
		parsed, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return fmt.Errorf("invalid float: %s", raw)
		}
		v.SetFloat(parsed)
		return nil
	case reflect.Slice, reflect.Array:
		if v.Type().Elem().Kind() != reflect.String {
			return fmt.Errorf("unsupported slice type: %s", v.Type())
		}
		parts := []string{}
		for _, p := range strings.Split(raw, ",") {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				parts = append(parts, trimmed)
			}
		}
		v.Set(reflect.ValueOf(parts))
		return nil
	default:
		return fmt.Errorf("unsupported type: %s", v.Kind())
	}
}
