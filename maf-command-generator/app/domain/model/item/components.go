package item

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type NormalizedComponent struct {
	Key   string
	Value any
}

func BuildItemComponents(entry Item) (string, string) {
	itemParts := []string{
		fmt.Sprintf("id:%q", strings.TrimSpace(entry.ItemID)),
		"count:1",
	}
	componentParts, errMsg := buildComponentParts(componentData(entry))
	if errMsg != "" {
		return "", errMsg
	}
	if len(componentParts) > 0 {
		itemParts = append(itemParts, fmt.Sprintf("components:{%s}", strings.Join(componentParts, ",")))
	}
	return fmt.Sprintf("{%s}", strings.Join(itemParts, ",")), ""
}

func componentData(entry Item) any {
	if entry.Minecraft == nil {
		return nil
	}
	return entry.Minecraft["components"]
}

func buildComponentParts(rawComponents any) ([]string, string) {
	entries, errMsg := NormalizeComponents(rawComponents)
	if errMsg != "" {
		return nil, errMsg
	}

	parts := make([]string, 0, len(entries))
	for _, entry := range entries {
		value, ok := valueToSNBT(entry.Value)
		if !ok {
			continue
		}
		parts = append(parts, fmt.Sprintf("%q:%s", entry.Key, value))
	}
	return parts, ""
}

func NormalizeComponents(rawComponents any) ([]NormalizedComponent, string) {
	components, ok := toComponentMap(rawComponents)
	if !ok {
		return nil, "minecraft.components must be an object"
	}
	if len(components) == 0 {
		return nil, ""
	}

	keys := make([]string, 0, len(components))
	normalizedValues := make(map[string]any, len(components))
	for key, value := range components {
		normalizedKey := strings.TrimSpace(key)
		if normalizedKey == "" {
			return nil, "component key is empty"
		}
		if valueString, ok := value.(string); ok {
			value = strings.TrimSpace(valueString)
		}
		keys = append(keys, normalizedKey)
		normalizedValues[normalizedKey] = value
	}

	sort.Strings(keys)
	entries := make([]NormalizedComponent, 0, len(keys))
	for _, key := range keys {
		entries = append(entries, NormalizedComponent{
			Key:   key,
			Value: normalizedValues[key],
		})
	}
	return entries, ""
}

func toComponentMap(raw any) (map[string]any, bool) {
	switch v := raw.(type) {
	case nil:
		return nil, true
	case map[string]any:
		out := make(map[string]any, len(v))
		for key, value := range v {
			out[key] = value
		}
		return out, true
	case map[string]string:
		out := make(map[string]any, len(v))
		for key, value := range v {
			out[key] = value
		}
		return out, true
	}

	rv := reflect.ValueOf(raw)
	if !rv.IsValid() || rv.Kind() != reflect.Map || rv.Type().Key().Kind() != reflect.String {
		return nil, false
	}
	out := make(map[string]any, rv.Len())
	iter := rv.MapRange()
	for iter.Next() {
		out[iter.Key().String()] = iter.Value().Interface()
	}
	return out, true
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
		return snbtString(v), true
	case map[string]any:
		return mapToSNBT(v), true
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
		return strconv.FormatFloat(float64(v), 'f', -1, 64), true
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), true
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
		out := make(map[string]any, rv.Len())
		iter := rv.MapRange()
		for iter.Next() {
			out[iter.Key().String()] = iter.Value().Interface()
		}
		return mapToSNBT(out), true
	}
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		out := make([]any, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			out = append(out, rv.Index(i).Interface())
		}
		return listToSNBT(out), true
	}
	return "", false
}

func mapToSNBT(data map[string]any) string {
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

func snbtString(raw string) string {
	raw = strings.ReplaceAll(raw, `\`, `\\`)
	raw = strings.ReplaceAll(raw, `"`, `\"`)
	return `"` + raw + `"`
}

func snbtKey(key string) string {
	if key == "" {
		return snbtString(key)
	}
	for _, r := range key {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '_' || r == '-' || r == '.' || r == '+':
		default:
			return snbtString(key)
		}
	}
	return key
}
