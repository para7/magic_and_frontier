package export_convert

import (
	"encoding/json"
	"fmt"
	"strings"
)

const bookConsumableSNBT = `{consume_seconds:99999,animation:"bow",has_consume_particles:false}`

type spellBookModel struct {
	itemName   string
	lore       []string
	customData string
}

func (m spellBookModel) ToGiveItem() string {
	return fmt.Sprintf(
		`minecraft:book[minecraft:item_name=%s,minecraft:lore=%s,minecraft:consumable=%s,minecraft:custom_data=%s]`,
		textComponentSNBT(m.itemName),
		loreComponentsSNBT(m.lore),
		bookConsumableSNBT,
		m.customData,
	)
}

func (m spellBookModel) LootComponents() map[string]any {
	return map[string]any{
		"minecraft:item_name": textComponentObject(m.itemName),
		"minecraft:lore":      loreComponentObjects(m.lore),
		"minecraft:consumable": map[string]any{
			"consume_seconds":       99999.0,
			"animation":             "bow",
			"has_consume_particles": false,
		},
	}
}

func textComponentSNBT(text string) string {
	return fmt.Sprintf(`{text:%s}`, JsonString(text))
}

func textComponentObject(text string) map[string]any {
	return map[string]any{"text": text}
}

func loreComponentObjects(lines []string) []any {
	lore := make([]any, 0, len(lines))
	for _, line := range lines {
		lore = append(lore, textComponentObject(line))
	}
	return lore
}

func loreComponentsSNBT(lines []string) string {
	if len(lines) == 0 {
		return "[]"
	}
	parts := make([]string, 0, len(lines))
	for _, line := range lines {
		parts = append(parts, textComponentSNBT(line))
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func jsonComponentValue(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(data)
}
