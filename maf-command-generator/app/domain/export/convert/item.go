package export_convert

import (
	"fmt"
	"strconv"
	"strings"

	itemModel "maf_command_editor/app/domain/model/item"
)

func itemCustomData(entry itemModel.Item) string {
	parts := []string{
		fmt.Sprintf("item_id:%s", JsonString(entry.ItemID)),
		fmt.Sprintf("source_id:%s", JsonString(entry.ID)),
		fmt.Sprintf("nbt_snapshot:%s", JsonString(entry.NBT)),
	}
	if entry.SkillID != "" {
		parts = append(parts, "maf_skill:1b", fmt.Sprintf("maf_skill_id:%s", JsonString(entry.SkillID)))
	}
	return "{maf:{" + strings.Join(parts, ",") + "}}"
}

func itemComponentsForLoot(entry itemModel.Item) map[string]any {
	components := map[string]any{}

	if name := strings.TrimSpace(entry.CustomName); name != "" {
		components["minecraft:custom_name"] = map[string]any{"text": name}
	}

	var loreLines []string
	for _, line := range strings.Split(entry.Lore, "\n") {
		line = strings.TrimSpace(strings.TrimSuffix(line, "\r"))
		if line != "" {
			loreLines = append(loreLines, line)
		}
	}
	if len(loreLines) > 0 {
		lore := make([]any, 0, len(loreLines))
		for _, line := range loreLines {
			lore = append(lore, map[string]any{"text": line})
		}
		components["minecraft:lore"] = lore
	}

	if entry.Unbreakable {
		components["minecraft:unbreakable"] = map[string]any{}
	}

	return components
}

func itemEnchantmentsForLoot(entry itemModel.Item) map[string]any {
	ench := strings.TrimSpace(entry.Enchantments)
	if ench == "" {
		return nil
	}
	enchMap := map[string]any{}
	for _, line := range strings.Split(strings.ReplaceAll(ench, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			if level, err := strconv.Atoi(fields[1]); err == nil {
				enchMap[fields[0]] = level
			}
		}
	}
	return enchMap
}
