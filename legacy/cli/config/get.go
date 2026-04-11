package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	cfg "github.com/HidemaruOwO/pummit/legacy/config"
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value by key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyPath := strings.Split(args[0], ".")
		val, err := resolveValue(reflect.ValueOf(cfg.CurrentTOMLConfig), keyPath)
		if err != nil {
			return err
		}
		return printValue(args[0], val)
	},
}

func resolveValue(v reflect.Value, path []string) (reflect.Value, error) {
	v = normalizeValue(v)
	if len(path) == 0 {
		return v, nil
	}

	switch v.Kind() {
	case reflect.Struct:
		field, ok := fieldByTomlTag(v, path[0])
		if !ok {
			return reflect.Value{}, fmt.Errorf("unknown key: %s", strings.Join(path, "."))
		}
		return resolveValue(field, path[1:])
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return reflect.Value{}, fmt.Errorf("invalid map key type for %s", strings.Join(path, "."))
		}
		key := reflect.ValueOf(path[0])
		val := v.MapIndex(key)
		if !val.IsValid() {
			return reflect.Value{}, fmt.Errorf("key not found: %s", strings.Join(path, "."))
		}
		return resolveValue(val, path[1:])
	case reflect.Slice, reflect.Array:
		index, err := strconv.Atoi(path[0])
		if err != nil {
			return reflect.Value{}, fmt.Errorf("invalid index '%s'", path[0])
		}
		if index < 0 || index >= v.Len() {
			return reflect.Value{}, fmt.Errorf("index out of range: %d", index)
		}
		return resolveValue(v.Index(index), path[1:])
	default:
		if len(path) > 0 {
			return reflect.Value{}, fmt.Errorf("cannot descend into %s", strings.Join(path, "."))
		}
		return v, nil
	}
}

func fieldByTomlTag(v reflect.Value, name string) (reflect.Value, bool) {
	typeOf := v.Type()
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		tag := field.Tag.Get("toml")
		if tag == "" {
			continue
		}
		tagName := strings.Split(tag, ",")[0]
		if tagName == name {
			return v.Field(i), true
		}
	}
	return reflect.Value{}, false
}

func normalizeValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return reflect.Zero(v.Type().Elem())
		}
		v = v.Elem()
	}
	return v
}

func printValue(key string, v reflect.Value) error {
	v = normalizeValue(v)
	if isComplexKind(v.Kind()) {
		last := key
		if idx := strings.LastIndex(key, "."); idx != -1 && idx < len(key)-1 {
			last = key[idx+1:]
		}
		if last == "" {
			last = "value"
		}
		wrapper := map[string]any{last: v.Interface()}
		return toml.NewEncoder(os.Stdout).Encode(wrapper)
	}

	fmt.Fprintln(os.Stdout, formatPrimitive(v))
	return nil
}

func isComplexKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
		return true
	default:
		return false
	}
}

func formatPrimitive(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}
