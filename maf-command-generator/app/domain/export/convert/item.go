package export_convert

import (
	"fmt"
	"sort"
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
	componentValues, err := itemComponentsForGive(entry)
	if err != nil {
		return "", err
	}

	spellMeta, err := resolveItemSpellMeta(entry, grimoiresByID, passivesByID, bowsByID)
	if err != nil {
		return "", err
	}
	if spellMeta.hasUseSpell {
		componentValues["minecraft:consumable"] = bookConsumableSNBT
	}

	customData, err := itemCustomData(entry, grimoiresByID, passivesByID, bowsByID)
	if err != nil {
		return "", err
	}
	componentValues["minecraft:custom_data"] = customData

	components := sortedItemGiveComponents(componentValues)
	itemID := strings.TrimSpace(entry.ItemID)
	if len(components) == 0 {
		return fmt.Sprintf("give @p %s 1", itemID), nil
	}
	return fmt.Sprintf("give @p %s[%s] 1", itemID, strings.Join(components, ",")), nil
}

func itemComponentsForGive(entry itemModel.Item) (map[string]string, error) {
	normalizedComponents, errMsg := itemModel.NormalizeComponents(componentData(entry))
	if errMsg != "" {
		return nil, fmt.Errorf("item(%s): %s", entry.ID, errMsg)
	}

	values := make(map[string]string, len(normalizedComponents))
	for _, component := range normalizedComponents {
		value, ok := valueToSNBT(component.Value)
		if !ok {
			continue
		}
		values[component.Key] = value
	}
	return values, nil
}

func sortedItemGiveComponents(values map[string]string) []string {
	if len(values) == 0 {
		return nil
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	components := make([]string, 0, len(keys))
	for _, key := range keys {
		components = append(components, fmt.Sprintf("%s=%s", key, values[key]))
	}
	return components
}

func itemCustomData(
	entry itemModel.Item,
	grimoiresByID map[string]grimoireModel.Grimoire,
	passivesByID map[string]passiveModel.Passive,
	bowsByID map[string]bowModel.BowPassive,
) (string, error) {
	itemSNBT, _ := itemModel.BuildItemComponents(entry)
	parts := []string{
		fmt.Sprintf("item_id:%s", SNBTString(entry.ItemID)),
		fmt.Sprintf("source_id:%s", SNBTString(entry.ID)),
		fmt.Sprintf("nbt_snapshot:%s", itemSNBT),
	}

	spellMeta, err := resolveItemSpellMeta(entry, grimoiresByID, passivesByID, bowsByID)
	if err != nil {
		return "", err
	}
	if spellMeta.grimoireID != "" {
		parts = append(parts, fmt.Sprintf("grimoire_id:%s", SNBTString(spellMeta.grimoireID)))
	}
	if len(spellMeta.customFragments) > 0 {
		parts = append(parts, spellMeta.customFragments...)
	}
	if entry.Maf.MaxMP != nil {
		parts = append(parts, fmt.Sprintf("maxmp:%d", *entry.Maf.MaxMP))
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
	components := deepCopyMapFromAny(componentData(entry))
	delete(components, "minecraft:enchantments")
	delete(components, "minecraft:custom_data")

	spellMeta, err := resolveItemSpellMeta(entry, grimoiresByID, passivesByID, bowsByID)
	if err != nil {
		return nil, err
	}
	if spellMeta.hasUseSpell {
		components["minecraft:consumable"] = map[string]any{
			"consume_seconds":       99999.0,
			"animation":             "bow",
			"has_consume_particles": false,
		}
	}

	return components, nil
}

func itemEnchantmentsForLoot(entry itemModel.Item) map[string]any {
	rawComponents := deepCopyMapFromAny(componentData(entry))
	rawEnchantments, ok := rawComponents["minecraft:enchantments"]
	if !ok || rawEnchantments == nil {
		return nil
	}
	enchMap := deepCopyMapFromAny(rawEnchantments)
	if len(enchMap) == 0 {
		return nil
	}
	return enchMap
}

func componentData(entry itemModel.Item) any {
	if entry.Minecraft == nil {
		return nil
	}
	return entry.Minecraft["components"]
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
			fmt.Sprintf("bowId:%s", SNBTString(bowID)),
			fmt.Sprintf("passiveId:%s", SNBTString("bow_"+bowID)),
			fmt.Sprintf("passiveCondition:%s", SNBTString("always")),
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
		fmt.Sprintf("passiveId:%s", SNBTString(passive.ID)),
		fmt.Sprintf("passiveSlot:%d", slot),
		fmt.Sprintf("passiveCondition:%s", SNBTString(strings.TrimSpace(passive.Condition))),
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
