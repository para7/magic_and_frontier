package export_convert

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	bowModel "maf_command_editor/app/domain/model/bow"
	grimoireModel "maf_command_editor/app/domain/model/grimoire"
	itemModel "maf_command_editor/app/domain/model/item"
	passiveModel "maf_command_editor/app/domain/model/passive"
)

func ItemToGiveCommand(
	entry itemModel.Item,
	grimoiresByID map[string]grimoireModel.Grimoire,
	passivesByID map[string]passiveModel.Passive,
	bowsByID map[string]bowModel.BowPassive,
) (string, error) {
	normalizedComponents, errMsg := itemModel.NormalizeComponents(entry.Minecraft.Components)
	if errMsg != "" {
		return "", fmt.Errorf("item(%s): %s", entry.ID, errMsg)
	}

	components := make([]string, 0, len(normalizedComponents)+2)
	for _, entry := range normalizedComponents {
		components = append(components, fmt.Sprintf("%s=%s", entry.Key, entry.Value))
	}

	spellMeta, err := resolveItemSpellMeta(entry, grimoiresByID, passivesByID, bowsByID)
	if err != nil {
		return "", err
	}
	if spellMeta.hasUseSpell && componentValue(entry, "minecraft:consumable") == "" {
		components = append(components, "minecraft:consumable="+bookConsumableSNBT)
	}

	customData, err := itemCustomData(entry, grimoiresByID, passivesByID, bowsByID)
	if err != nil {
		return "", err
	}
	components = append(components, "minecraft:custom_data="+customData)

	itemID := strings.TrimSpace(entry.Minecraft.ItemID)
	if len(components) == 0 {
		return fmt.Sprintf("give @p %s 1", itemID), nil
	}
	return fmt.Sprintf("give @p %s[%s] 1", itemID, strings.Join(components, ",")), nil
}

func itemCustomData(
	entry itemModel.Item,
	grimoiresByID map[string]grimoireModel.Grimoire,
	passivesByID map[string]passiveModel.Passive,
	bowsByID map[string]bowModel.BowPassive,
) (string, error) {
	itemSNBT, _ := itemModel.BuildItemComponents(entry)
	parts := []string{
		fmt.Sprintf("item_id:%s", JsonString(entry.Minecraft.ItemID)),
		fmt.Sprintf("source_id:%s", JsonString(entry.ID)),
		fmt.Sprintf("nbt_snapshot:%s", JsonString(itemSNBT)),
	}

	spellMeta, err := resolveItemSpellMeta(entry, grimoiresByID, passivesByID, bowsByID)
	if err != nil {
		return "", err
	}
	if spellMeta.grimoireID != "" {
		parts = append(parts, fmt.Sprintf("grimoire_id:%s", JsonString(spellMeta.grimoireID)))
	}
	if len(spellMeta.customFragments) > 0 {
		parts = append(parts, spellMeta.customFragments...)
	}
	if spellMeta.spellFragment != "" {
		parts = append(parts, spellMeta.spellFragment)
	}
	return "{maf:{" + strings.Join(parts, ",") + "}}", nil
}

func itemComponentsForLoot(
	entry itemModel.Item,
	grimoiresByID map[string]grimoireModel.Grimoire,
	passivesByID map[string]passiveModel.Passive,
	bowsByID map[string]bowModel.BowPassive,
) (map[string]any, error) {
	components := map[string]any{}

	if name, ok := decodeTextComponentSNBT(componentValue(entry, "minecraft:custom_name")); ok {
		components["minecraft:custom_name"] = name
	}
	if lore, ok := decodeTextComponentListSNBT(componentValue(entry, "minecraft:lore")); ok && len(lore) > 0 {
		components["minecraft:lore"] = lore
	}
	if componentValue(entry, "minecraft:unbreakable") != "" {
		components["minecraft:unbreakable"] = map[string]any{}
	}
	spellMeta, err := resolveItemSpellMeta(entry, grimoiresByID, passivesByID, bowsByID)
	if err != nil {
		return nil, err
	}
	if spellMeta.hasUseSpell && componentValue(entry, "minecraft:consumable") == "" {
		components["minecraft:consumable"] = map[string]any{
			"consume_seconds":       99999.0,
			"animation":             "bow",
			"has_consume_particles": false,
		}
	}

	return components, nil
}

func itemEnchantmentsForLoot(entry itemModel.Item) map[string]any {
	ench := componentValue(entry, "minecraft:enchantments")
	if ench == "" {
		return nil
	}
	var enchMap map[string]any
	if err := json.Unmarshal([]byte(ench), &enchMap); err != nil {
		return nil
	}
	return enchMap
}

func componentValue(entry itemModel.Item, key string) string {
	if entry.Minecraft.Components == nil {
		return ""
	}
	return strings.TrimSpace(entry.Minecraft.Components[key])
}

func decodeTextComponentSNBT(raw string) (map[string]any, bool) {
	text, ok := decodeQuotedSNBTString(raw)
	if !ok {
		return nil, false
	}
	var component map[string]any
	if err := json.Unmarshal([]byte(text), &component); err != nil {
		return nil, false
	}
	return component, true
}

func decodeTextComponentListSNBT(raw string) ([]any, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	if !strings.HasPrefix(raw, "[") || !strings.HasSuffix(raw, "]") {
		return nil, false
	}

	items := splitTopLevelSNBTList(raw[1 : len(raw)-1])
	out := make([]any, 0, len(items))
	for _, item := range items {
		text, ok := decodeQuotedSNBTString(item)
		if !ok {
			return nil, false
		}
		var component map[string]any
		if err := json.Unmarshal([]byte(text), &component); err != nil {
			return nil, false
		}
		out = append(out, component)
	}
	return out, true
}

func decodeQuotedSNBTString(raw string) (string, bool) {
	raw = strings.TrimSpace(raw)
	if len(raw) < 2 {
		return "", false
	}

	switch raw[0] {
	case '\'':
		if raw[len(raw)-1] != '\'' {
			return "", false
		}
		return strings.ReplaceAll(raw[1:len(raw)-1], `\'`, `'`), true
	case '"':
		value, err := strconv.Unquote(raw)
		if err != nil {
			return "", false
		}
		return value, true
	default:
		return "", false
	}
}

func splitTopLevelSNBTList(fragment string) []string {
	fragment = strings.TrimSpace(fragment)
	if fragment == "" {
		return nil
	}

	var items []string
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
				item := strings.TrimSpace(fragment[start:i])
				if item != "" {
					items = append(items, item)
				}
				start = i + 1
			}
		}
	}

	last := strings.TrimSpace(fragment[start:])
	if last != "" {
		items = append(items, last)
	}
	return items
}

type itemSpellMeta struct {
	hasUseSpell     bool
	grimoireID      string
	spellFragment   string
	customFragments []string
}

func resolveItemSpellMeta(
	entry itemModel.Item,
	grimoiresByID map[string]grimoireModel.Grimoire,
	passivesByID map[string]passiveModel.Passive,
	bowsByID map[string]bowModel.BowPassive,
) (itemSpellMeta, error) {
	meta := itemSpellMeta{}

	grimoireID := strings.TrimSpace(entry.Maf.GrimoireID)
	if grimoireID != "" {
		grimoire, ok := grimoiresByID[grimoireID]
		if !ok {
			return itemSpellMeta{}, fmt.Errorf("item(%s): referenced grimoire not found (%s)", entry.ID, grimoireID)
		}
		meta.hasUseSpell = true
		meta.grimoireID = grimoire.ID
		meta.spellFragment = grimoireSpellFragment(grimoire)
	}

	bowID := strings.TrimSpace(entry.Maf.BowID)
	if bowID != "" {
		if grimoireID != "" {
			return itemSpellMeta{}, fmt.Errorf("item(%s): bowId cannot be combined with grimoireId", entry.ID)
		}
		if _, ok := bowsByID[bowID]; !ok {
			return itemSpellMeta{}, fmt.Errorf("item(%s): referenced bow not found (%s)", entry.ID, bowID)
		}
		meta.customFragments = append(meta.customFragments,
			fmt.Sprintf("bowId:%s", JsonString(bowID)),
			fmt.Sprintf("passiveId:%s", JsonString("bow_"+bowID)),
		)
		return meta, nil
	}

	passiveID := strings.TrimSpace(entry.Maf.PassiveID)
	if passiveID == "" {
		return meta, nil
	}
	passive, ok := passivesByID[passiveID]
	if !ok {
		return itemSpellMeta{}, fmt.Errorf("item(%s): referenced passive not found (%s)", entry.ID, passiveID)
	}
	slot, err := resolvePassiveSlot(entry, passive)
	if err != nil {
		return itemSpellMeta{}, err
	}
	meta.customFragments = []string{
		"hasPassive:1b",
		fmt.Sprintf("passiveId:%s", JsonString(passive.ID)),
		fmt.Sprintf("passiveSlot:%d", slot),
		fmt.Sprintf("passiveCondition:%s", JsonString(strings.TrimSpace(passive.Condition))),
	}
	return meta, nil
}

func resolvePassiveSlot(entry itemModel.Item, passive passiveModel.Passive) (int, error) {
	if entry.Maf.PassiveSlot != 0 {
		for _, slot := range passive.Slots {
			if slot == entry.Maf.PassiveSlot {
				return slot, nil
			}
		}
		return 0, fmt.Errorf("item(%s): passive(%s) does not support slot %d", entry.ID, passive.ID, entry.Maf.PassiveSlot)
	}
	if len(passive.Slots) == 0 {
		return 0, fmt.Errorf("item(%s): passive(%s) has no available slots", entry.ID, passive.ID)
	}
	return passive.Slots[0], nil
}
