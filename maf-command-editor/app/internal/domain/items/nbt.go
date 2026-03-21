package items

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func buildItemNBT(input SaveInput) (string, string) {
	tagParts := make([]string, 0, 16)

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
		tagParts = append(tagParts, fmt.Sprintf("display:{%s}", strings.Join(displayParts, ",")))
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
			tagParts = append(tagParts, fmt.Sprintf("Enchantments:[%s]", strings.Join(parts, ",")))
		}
	}

	if input.Unbreakable {
		tagParts = append(tagParts, "Unbreakable:1b")
	}
	if v := strings.TrimSpace(input.CustomModelData); v != "" {
		tagParts = append(tagParts, "CustomModelData:"+v)
	}
	if v := strings.TrimSpace(input.RepairCost); v != "" {
		tagParts = append(tagParts, "RepairCost:"+v)
	}
	if v := strings.TrimSpace(input.HideFlags); v != "" {
		tagParts = append(tagParts, "HideFlags:"+v)
	}
	if v := strings.Trim(strings.TrimSpace(input.PotionID), `"`); v != "" {
		tagParts = append(tagParts, fmt.Sprintf("Potion:%q", v))
	}
	if v := strings.TrimSpace(input.CustomPotionColor); v != "" {
		tagParts = append(tagParts, "CustomPotionColor:"+v)
	}
	if v := normalizeListValue(input.CustomPotionEffects); v != "" {
		tagParts = append(tagParts, "CustomPotionEffects:"+v)
	}
	if v := normalizeListValue(input.AttributeModifiers); v != "" {
		tagParts = append(tagParts, "AttributeModifiers:"+v)
	}
	if v := normalizeCustomNbtFragment(input.CustomNBT); v != "" {
		tagParts = append(tagParts, v)
	}

	itemParts := []string{
		fmt.Sprintf("id:%q", strings.TrimSpace(input.ItemID)),
		fmt.Sprintf("Count:%db", input.Count),
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
