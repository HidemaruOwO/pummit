package usecase

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	domainconfig "github.com/HidemaruOwO/pummit/internal/domain/config"
	infraconfig "github.com/HidemaruOwO/pummit/internal/infra/config"
)

type ConfigService struct {
	store *infraconfig.Store
}

func NewConfigService(store *infraconfig.Store) *ConfigService {
	return &ConfigService{store: store}
}

func (s *ConfigService) Load() (domainconfig.Config, error) {
	return s.store.Load()
}

func (s *ConfigService) Validate() (domainconfig.Config, error) {
	return s.Load()
}

func (s *ConfigService) Save(cfg domainconfig.Config) error {
	if err := domainconfig.Validate(cfg); err != nil {
		return err
	}

	return s.store.Save(cfg)
}

func (s *ConfigService) List() (domainconfig.Config, error) {
	return s.Load()
}

func (s *ConfigService) Get(key string) (any, error) {
	cfg, err := s.Load()
	if err != nil {
		return nil, err
	}

	value, err := resolveValue(reflect.ValueOf(cfg), strings.Split(key, "."))
	if err != nil {
		return nil, err
	}

	return normalizeValue(value).Interface(), nil
}

func (s *ConfigService) Set(key, raw string) (domainconfig.Config, error) {
	cfg, err := s.Load()
	if err != nil {
		return domainconfig.Config{}, err
	}

	if err := setValue(reflect.ValueOf(&cfg), strings.Split(key, "."), raw); err != nil {
		return domainconfig.Config{}, err
	}

	if err := domainconfig.Validate(cfg); err != nil {
		return domainconfig.Config{}, err
	}

	if err := s.store.Save(cfg); err != nil {
		return domainconfig.Config{}, err
	}

	return cfg, nil
}

func (s *ConfigService) Reset() (domainconfig.Config, error) {
	cfg := domainconfig.Default()
	if err := s.store.Save(cfg); err != nil {
		return domainconfig.Config{}, err
	}

	return cfg, nil
}

func resolveValue(v reflect.Value, path []string) (reflect.Value, error) {
	v = normalizeValue(v)
	if len(path) == 0 {
		return v, nil
	}

	switch v.Kind() {
	case reflect.Struct:
		field, ok := fieldByTOMLTag(v, path[0])
		if !ok {
			return reflect.Value{}, fmt.Errorf("unknown key: %s", strings.Join(path, "."))
		}
		return resolveValue(field, path[1:])
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return reflect.Value{}, fmt.Errorf("invalid map key type for %s", strings.Join(path, "."))
		}
		key := reflect.ValueOf(path[0])
		value := v.MapIndex(key)
		if !value.IsValid() {
			return reflect.Value{}, fmt.Errorf("key not found: %s", strings.Join(path, "."))
		}
		return resolveValue(value, path[1:])
	case reflect.Slice, reflect.Array:
		index, err := strconv.Atoi(path[0])
		if err != nil {
			return reflect.Value{}, fmt.Errorf("invalid index %q", path[0])
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

func setValue(v reflect.Value, path []string, raw string) error {
	v = normalizeValue(v)
	if len(path) == 0 {
		return fmt.Errorf("key is required")
	}

	switch v.Kind() {
	case reflect.Struct:
		field, ok := fieldByTOMLTag(v, path[0])
		if !ok {
			return fmt.Errorf("unknown key: %s", strings.Join(path, "."))
		}
		if len(path) == 1 {
			return assignValue(field, raw)
		}
		return setValue(field, path[1:], raw)
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return fmt.Errorf("invalid map key type for %s", strings.Join(path, "."))
		}
		key := reflect.ValueOf(path[0])
		value := v.MapIndex(key)
		if !value.IsValid() {
			return fmt.Errorf("key not found: %s", strings.Join(path, "."))
		}
		copyValue := reflect.New(normalizeValue(value).Type()).Elem()
		copyValue.Set(normalizeValue(value))
		if len(path) == 1 {
			if err := assignValue(copyValue, raw); err != nil {
				return err
			}
		} else {
			if err := setValue(copyValue, path[1:], raw); err != nil {
				return err
			}
		}
		v.SetMapIndex(key, copyValue)
		return nil
	case reflect.Slice, reflect.Array:
		index, err := strconv.Atoi(path[0])
		if err != nil {
			return fmt.Errorf("invalid index %q", path[0])
		}
		if index < 0 || index >= v.Len() {
			return fmt.Errorf("index out of range: %d", index)
		}
		if len(path) == 1 {
			return assignValue(v.Index(index), raw)
		}
		return setValue(v.Index(index), path[1:], raw)
	default:
		if len(path) == 1 {
			return assignValue(v, raw)
		}
		return fmt.Errorf("cannot descend into %s", strings.Join(path, "."))
	}
}

func fieldByTOMLTag(v reflect.Value, name string) (reflect.Value, bool) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := strings.Split(field.Tag.Get("toml"), ",")[0]
		if tag == name {
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

func assignValue(v reflect.Value, raw string) error {
	v = normalizeValue(v)
	if !v.CanSet() {
		return fmt.Errorf("cannot set value")
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(raw)
		return nil
	case reflect.Bool:
		value, err := strconv.ParseBool(raw)
		if err != nil {
			return fmt.Errorf("invalid bool: %s", raw)
		}
		v.SetBool(value)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid int: %s", raw)
		}
		v.SetInt(value)
		return nil
	case reflect.Slice:
		if v.Type().Elem().Kind() != reflect.String {
			return fmt.Errorf("unsupported slice type: %s", v.Type())
		}
		parts := strings.Split(raw, ",")
		values := make([]string, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				values = append(values, part)
			}
		}
		v.Set(reflect.ValueOf(values))
		return nil
	default:
		return fmt.Errorf("unsupported type: %s", v.Kind())
	}
}
