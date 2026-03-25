package items

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func buildItemNBT(input SaveInput) (string, string) {
	tagParts := make([]string, 0, 16)
	formKeys := make(map[string]bool)

	addTag := func(key, value string) {
		tagParts = append(tagParts, key+":"+value)
		formKeys[key] = true
	}

	displayParts := make([]string, 0, 2)
	if name := strings.TrimSpace(input.CustomName); name != "" {
		displayParts = append(displayParts, fmt.Sprintf("Name:%s", toNbtText(name)))
	}
	loreLines := make([]string, 0)
	for _, line := range strings.Split(input.Lore, "\n") {
		line = strings.TrimSpace(strings.TrimSuffix(line, "\r"))
		if line != "" {
			loreLines = append(loreLines, line)
		}
	}
	if len(loreLines) > 0 {
		lore := make([]string, 0, len(loreLines))
		for _, line := range loreLines {
			lore = append(lore, toNbtText(line))
		}
		displayParts = append(displayParts, fmt.Sprintf("Lore:[%s]", strings.Join(lore, ",")))
	}
	if len(displayParts) > 0 {
		addTag("display", fmt.Sprintf("{%s}", strings.Join(displayParts, ",")))
	}

	if ench := strings.TrimSpace(input.Enchantments); ench != "" {
		lines := strings.Split(strings.ReplaceAll(ench, "\r\n", "\n"), "\n")
		parts := make([]string, 0, len(lines))
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) < 2 {
				return "", fmt.Sprintf("Invalid enchantment line: %q. Use \"minecraft:sharpness 5\" format.", line)
			}
			level, err := strconv.Atoi(fields[1])
			if err != nil || level < 1 || level > 255 {
				return "", fmt.Sprintf("Invalid enchantment line: %q. Use \"minecraft:sharpness 5\" format.", line)
			}
			parts = append(parts, fmt.Sprintf("{id:%q,lvl:%ds}", fields[0], level))
		}
		if len(parts) > 0 {
			addTag("Enchantments", fmt.Sprintf("[%s]", strings.Join(parts, ",")))
		}
	}

	if input.Unbreakable {
		tagParts = append(tagParts, "Unbreakable:1b")
		formKeys["Unbreakable"] = true
	}
	if v := strings.TrimSpace(input.CustomModelData); v != "" {
		addTag("CustomModelData", v)
	}
	if v := strings.TrimSpace(input.RepairCost); v != "" {
		addTag("RepairCost", v)
	}
	if v := strings.TrimSpace(input.HideFlags); v != "" {
		addTag("HideFlags", v)
	}
	if v := strings.Trim(strings.TrimSpace(input.PotionID), `"`); v != "" {
		addTag("Potion", fmt.Sprintf("%q", v))
	}
	if v := strings.TrimSpace(input.CustomPotionColor); v != "" {
		addTag("CustomPotionColor", v)
	}
	if v := normalizeListValue(input.CustomPotionEffects); v != "" {
		addTag("CustomPotionEffects", v)
	}
	if v := normalizeListValue(input.AttributeModifiers); v != "" {
		addTag("AttributeModifiers", v)
	}
	if v := normalizeCustomNbtFragment(input.CustomNBT); v != "" {
		for _, entry := range splitSNBTEntries(v) {
			if !formKeys[entry.Key] {
				tagParts = append(tagParts, entry.Raw)
			}
		}
	}

	itemParts := []string{
		fmt.Sprintf("id:%q", strings.TrimSpace(input.ItemID)),
		"Count:1b",
	}
	if len(tagParts) > 0 {
		itemParts = append(itemParts, fmt.Sprintf("tag:{%s}", strings.Join(tagParts, ",")))
	}
	return fmt.Sprintf("{%s}", strings.Join(itemParts, ",")), ""
}

func toNbtText(value string) string {
	raw, _ := json.Marshal(map[string]string{"text": value})
	return "'" + strings.ReplaceAll(string(raw), "'", "\\'") + "'"
}

type snbtEntry struct {
	Key string
	Raw string
}

// splitSNBTEntries splits a normalized SNBT fragment (outer braces already removed)
// into top-level key:value entries, respecting nested braces, brackets, and quoted strings.
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
	colon := strings.Index(s, ":")
	if colon < 0 {
		return snbtEntry{Raw: s, Key: s}, true
	}
	key := strings.TrimSpace(s[:colon])
	return snbtEntry{Key: key, Raw: s}, true
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
