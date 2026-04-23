package export_convert

import (
	"encoding/json"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// SNBTString encodes s as a double-quoted SNBT string value.
// Double-quoted SNBT requires escaping backslash and double quotes.
func SNBTString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return `"` + s + `"`
}

// MapToSNBT converts a map into SNBT compound form.
func MapToSNBT(data map[string]any) string {
	if len(data) == 0 {
		return "{}"
	}

	keys := make([]string, 0, len(data))
	for key, value := range data {
		if value == nil {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		value, ok := valueToSNBT(data[key])
		if !ok {
			continue
		}
		parts = append(parts, snbtKey(key)+":"+value)
	}
	return "{" + strings.Join(parts, ",") + "}"
}

func listToSNBT(values []any) string {
	parts := make([]string, 0, len(values))
	for _, value := range values {
		formatted, ok := valueToSNBT(value)
		if !ok {
			continue
		}
		parts = append(parts, formatted)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func valueToSNBT(value any) (string, bool) {
	switch v := value.(type) {
	case nil:
		return "", false
	case json.Number:
		return v.String(), true
	case bool:
		if v {
			return "1b", true
		}
		return "0b", true
	case string:
		return SNBTString(v), true
	case map[string]any:
		return MapToSNBT(v), true
	case []any:
		return listToSNBT(v), true
	case int:
		return strconv.Itoa(v), true
	case int8:
		return strconv.FormatInt(int64(v), 10), true
	case int16:
		return strconv.FormatInt(int64(v), 10), true
	case int32:
		return strconv.FormatInt(int64(v), 10), true
	case int64:
		return strconv.FormatInt(v, 10), true
	case uint:
		return strconv.FormatUint(uint64(v), 10), true
	case uint8:
		return strconv.FormatUint(uint64(v), 10), true
	case uint16:
		return strconv.FormatUint(uint64(v), 10), true
	case uint32:
		return strconv.FormatUint(uint64(v), 10), true
	case uint64:
		return strconv.FormatUint(v, 10), true
	case float32:
		return formatFloat(float64(v)), true
	case float64:
		return formatFloat(v), true
	}

	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return "", false
	}
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return "", false
		}
		return valueToSNBT(rv.Elem().Interface())
	}
	if rv.Kind() == reflect.Map && rv.Type().Key().Kind() == reflect.String {
		mapValue := make(map[string]any, rv.Len())
		iter := rv.MapRange()
		for iter.Next() {
			mapValue[iter.Key().String()] = iter.Value().Interface()
		}
		return MapToSNBT(mapValue), true
	}
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		values := make([]any, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			values = append(values, rv.Index(i).Interface())
		}
		return listToSNBT(values), true
	}
	return "", false
}

func snbtKey(key string) string {
	if key == "" {
		return SNBTString(key)
	}
	for _, r := range key {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '_' || r == '-' || r == '.' || r == '+':
		default:
			return SNBTString(key)
		}
	}
	return key
}
