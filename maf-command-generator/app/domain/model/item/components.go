package item

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// BuildItemComponents はアイテムデータから Minecraft の item component 文字列を生成する。
// エラーの場合は空文字列とエラーメッセージを返す。
func BuildItemComponents(entry Item) (string, string) {
	componentParts := make([]string, 0, 16)
	formKeys := make(map[string]bool)

	addComponent := func(key, value string, aliases ...string) {
		componentParts = append(componentParts, fmt.Sprintf("%q:%s", key, value))
		registerFormKey(formKeys, key)
		for _, alias := range aliases {
			registerFormKey(formKeys, alias)
		}
	}

	if name := strings.TrimSpace(entry.CustomName); name != "" {
		addComponent("minecraft:custom_name", toTextComponentSNBT(name), "display")
	}
	loreLines := make([]string, 0)
	for _, line := range strings.Split(entry.Lore, "\n") {
		line = strings.TrimSpace(strings.TrimSuffix(line, "\r"))
		if line != "" {
			loreLines = append(loreLines, line)
		}
	}
	if len(loreLines) > 0 {
		lore := make([]string, 0, len(loreLines))
		for _, line := range loreLines {
			lore = append(lore, toTextComponentSNBT(line))
		}
		addComponent("minecraft:lore", fmt.Sprintf("[%s]", strings.Join(lore, ",")), "display")
	}

	if ench := strings.TrimSpace(entry.Enchantments); ench != "" {
		lines := strings.Split(strings.ReplaceAll(ench, "\r\n", "\n"), "\n")
		parts := make([]string, 0, len(lines))
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) < 2 {
				return "", fmt.Sprintf("エンチャントの書式が正しくありません: %q (例: \"minecraft:sharpness 5\")", line)
			}
			level, err := strconv.Atoi(fields[1])
			if err != nil || level < 1 || level > 255 {
				return "", fmt.Sprintf("エンチャントの書式が正しくありません: %q (例: \"minecraft:sharpness 5\")", line)
			}
			parts = append(parts, fmt.Sprintf("%q:%d", fields[0], level))
		}
		if len(parts) > 0 {
			addComponent("minecraft:enchantments", fmt.Sprintf("{%s}", strings.Join(parts, ",")), "Enchantments")
		}
	}

	if entry.Unbreakable {
		addComponent("minecraft:unbreakable", "{}", "Unbreakable")
	}
	if v := normalizeCustomModelDataComponent(entry.CustomModelData); v != "" {
		addComponent("minecraft:custom_model_data", v, "CustomModelData")
	}
	if v := strings.TrimSpace(entry.RepairCost); v != "" {
		addComponent("minecraft:repair_cost", v, "RepairCost")
	}

	if v := normalizeHideFlagsComponent(entry.HideFlags); v != "" {
		addComponent("minecraft:tooltip_display", v, "HideFlags")
	}

	potionParts := make([]string, 0, 3)
	if v := strings.Trim(strings.TrimSpace(entry.PotionID), `"`); v != "" {
		potionParts = append(potionParts, fmt.Sprintf("potion:%q", v))
	}
	if v := strings.TrimSpace(entry.CustomPotionColor); v != "" {
		potionParts = append(potionParts, "custom_color:"+v)
	}
	if v := normalizeListValue(entry.CustomPotionEffects); v != "" {
		potionParts = append(potionParts, "custom_effects:"+v)
	}
	if len(potionParts) > 0 {
		addComponent("minecraft:potion_contents", fmt.Sprintf("{%s}", strings.Join(potionParts, ",")),
			"Potion", "CustomPotionColor", "CustomPotionEffects")
	}
	if v := normalizeListValue(entry.AttributeModifiers); v != "" {
		addComponent("minecraft:attribute_modifiers", v, "AttributeModifiers")
	}

	if v := normalizeCustomNbtFragment(entry.CustomNBT); v != "" {
		for _, e := range splitSNBTEntries(v) {
			if !formKeys[canonicalComponentKey(e.Key)] {
				componentParts = append(componentParts, e.Raw)
			}
		}
	}

	itemParts := []string{
		fmt.Sprintf("id:%q", strings.TrimSpace(entry.ItemID)),
		"count:1",
	}
	if len(componentParts) > 0 {
		itemParts = append(itemParts, fmt.Sprintf("components:{%s}", strings.Join(componentParts, ",")))
	}
	return fmt.Sprintf("{%s}", strings.Join(itemParts, ",")), ""
}

func toTextComponentSNBT(value string) string {
	raw, _ := json.Marshal(map[string]string{"text": value})
	return "'" + strings.ReplaceAll(string(raw), "'", "\\'") + "'"
}

func registerFormKey(formKeys map[string]bool, key string) {
	formKeys[canonicalComponentKey(key)] = true
}

func canonicalComponentKey(key string) string {
	k := strings.ToLower(strings.Trim(trimSNBTKeyQuotes(strings.TrimSpace(key)), " "))
	switch k {
	case "display":
		return "display"
	case "enchantments":
		return "minecraft:enchantments"
	case "storedenchantments":
		return "minecraft:stored_enchantments"
	case "unbreakable":
		return "minecraft:unbreakable"
	case "custommodeldata":
		return "minecraft:custom_model_data"
	case "repaircost":
		return "minecraft:repair_cost"
	case "hideflags":
		return "minecraft:tooltip_display"
	case "potion", "custompotioncolor", "custompotioneffects":
		return "minecraft:potion_contents"
	case "attributemodifiers":
		return "minecraft:attribute_modifiers"
	}
	return k
}

func trimSNBTKeyQuotes(key string) string {
	if len(key) >= 2 {
		if (key[0] == '"' && key[len(key)-1] == '"') || (key[0] == '\'' && key[len(key)-1] == '\'') {
			return key[1 : len(key)-1]
		}
	}
	return key
}

func normalizeCustomModelDataComponent(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
		return trimmed
	}
	if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
		return "{floats:" + trimmed + "}"
	}

	parts := strings.Split(trimmed, ",")
	values := make([]string, 0, len(parts))
	for _, p := range parts {
		token := strings.TrimSpace(p)
		if token == "" {
			continue
		}
		if isPlainNumberToken(token) {
			token += "f"
		}
		values = append(values, token)
	}
	if len(values) == 0 {
		return ""
	}
	return fmt.Sprintf("{floats:[%s]}", strings.Join(values, ","))
}

func isPlainNumberToken(token string) bool {
	if _, err := strconv.ParseFloat(token, 64); err != nil {
		return false
	}
	last := token[len(token)-1]
	return (last >= '0' && last <= '9') || last == '.'
}

func normalizeHideFlagsComponent(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}

	mask, err := strconv.ParseInt(trimmed, 0, 64)
	if err != nil {
		return trimmed
	}
	if mask == 0 {
		return ""
	}

	var hidden []string
	seen := make(map[string]bool)
	addHidden := func(id string) {
		if seen[id] {
			return
		}
		seen[id] = true
		hidden = append(hidden, id)
	}

	if mask&1 != 0 {
		addHidden("minecraft:enchantments")
		addHidden("minecraft:stored_enchantments")
	}
	if mask&2 != 0 {
		addHidden("minecraft:attribute_modifiers")
	}
	if mask&4 != 0 {
		addHidden("minecraft:unbreakable")
	}
	if mask&8 != 0 {
		addHidden("minecraft:can_break")
	}
	if mask&16 != 0 {
		addHidden("minecraft:can_place_on")
	}
	if mask&32 != 0 {
		addHidden("minecraft:potion_contents")
		addHidden("minecraft:written_book_content")
		addHidden("minecraft:fireworks")
		addHidden("minecraft:map_id")
	}
	if mask&64 != 0 {
		addHidden("minecraft:dyed_color")
	}
	if mask&128 != 0 {
		addHidden("minecraft:trim")
	}

	if len(hidden) == 0 {
		return ""
	}
	quoted := make([]string, 0, len(hidden))
	for _, id := range hidden {
		quoted = append(quoted, fmt.Sprintf("%q", id))
	}
	return fmt.Sprintf("{hidden_components:[%s]}", strings.Join(quoted, ","))
}

type snbtEntry struct {
	Key string
	Raw string
}

func splitSNBTEntries(fragment string) []snbtEntry {
	fragment = strings.TrimSpace(fragment)
	if fragment == "" {
		return nil
	}
	var entries []snbtEntry
	depth := 0
	inDouble := false
	inSingle := false
	escaped := false
	start := 0

	for i := 0; i < len(fragment); i++ {
		ch := fragment[i]
		if escaped {
			escaped = false
			continue
		}
		if ch == '\\' && inDouble {
			escaped = true
			continue
		}
		if ch == '"' && !inSingle {
			inDouble = !inDouble
			continue
		}
		if ch == '\'' && !inDouble {
			inSingle = !inSingle
			continue
		}
		if inDouble || inSingle {
			continue
		}
		switch ch {
		case '{', '[':
			depth++
		case '}', ']':
			depth--
		case ',':
			if depth == 0 {
				if entry, ok := parseSnbtEntry(fragment[start:i]); ok {
					entries = append(entries, entry)
				}
				start = i + 1
			}
		}
	}
	if entry, ok := parseSnbtEntry(fragment[start:]); ok {
		entries = append(entries, entry)
	}
	return entries
}

func parseSnbtEntry(s string) (snbtEntry, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return snbtEntry{}, false
	}
	colon := findTopLevelColon(s)
	if colon < 0 {
		return snbtEntry{Raw: s, Key: s}, true
	}
	key := strings.TrimSpace(s[:colon])
	return snbtEntry{Key: key, Raw: s}, true
}

func findTopLevelColon(s string) int {
	depth := 0
	inDouble := false
	inSingle := false
	escaped := false

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if escaped {
			escaped = false
			continue
		}
		if ch == '\\' && inDouble {
			escaped = true
			continue
		}
		if ch == '"' && !inSingle {
			inDouble = !inDouble
			continue
		}
		if ch == '\'' && !inDouble {
			inSingle = !inSingle
			continue
		}
		if inDouble || inSingle {
			continue
		}
		switch ch {
		case '{', '[':
			depth++
		case '}', ']':
			depth--
		case ':':
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func normalizeCustomNbtFragment(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
		return strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(trimmed, "{"), "}"))
	}
	return trimmed
}

func normalizeListValue(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
		return trimmed
	}
	return "[" + trimmed + "]"
}
