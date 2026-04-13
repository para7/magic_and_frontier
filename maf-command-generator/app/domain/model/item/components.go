package item

import (
	"fmt"
	"sort"
	"strings"
)

type NormalizedComponent struct {
	Key   string
	Value string
}

func BuildItemComponents(entry Item) (string, string) {
	itemParts := []string{
		fmt.Sprintf("id:%q", strings.TrimSpace(entry.Minecraft.ItemID)),
		"count:1",
	}
	componentParts, errMsg := buildComponentParts(entry.Minecraft.Components)
	if errMsg != "" {
		return "", errMsg
	}
	if len(componentParts) > 0 {
		itemParts = append(itemParts, fmt.Sprintf("components:{%s}", strings.Join(componentParts, ",")))
	}
	return fmt.Sprintf("{%s}", strings.Join(itemParts, ",")), ""
}

func buildComponentParts(components map[string]string) ([]string, string) {
	entries, errMsg := NormalizeComponents(components)
	if errMsg != "" {
		return nil, errMsg
	}

	parts := make([]string, 0, len(entries))
	for _, entry := range entries {
		parts = append(parts, fmt.Sprintf("%q:%s", entry.Key, entry.Value))
	}
	return parts, ""
}

func NormalizeComponents(components map[string]string) ([]NormalizedComponent, string) {
	if len(components) == 0 {
		return nil, ""
	}

	keys := make([]string, 0, len(components))
	normalizedValues := make(map[string]string, len(components))
	for key, value := range components {
		normalizedKey := strings.TrimSpace(key)
		normalizedValue := strings.TrimSpace(value)
		if normalizedKey == "" {
			return nil, "component key is empty"
		}
		if !strings.Contains(normalizedKey, ":") {
			return nil, fmt.Sprintf("component key must be namespaced: %q", normalizedKey)
		}
		if normalizedValue == "" {
			return nil, fmt.Sprintf("component value is empty: %q", normalizedKey)
		}
		keys = append(keys, normalizedKey)
		normalizedValues[normalizedKey] = normalizedValue
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
