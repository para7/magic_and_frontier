package export_convert

import (
	"fmt"
	"sort"
	"strings"

	activeModel "maf_command_editor/app/domain/model/active"
	bowModel "maf_command_editor/app/domain/model/bow"
	itemModel "maf_command_editor/app/domain/model/item"
	passiveModel "maf_command_editor/app/domain/model/passive"
)

func ItemToGiveCommand(
	entry itemModel.Item,
	activesByID map[string]activeModel.Active,
	passivesByID map[string]passiveModel.Passive,
	bowsByID map[string]bowModel.BowPassive,
) (string, error) {
	spellMeta, err := resolveItemSpellMeta(entry, activesByID, passivesByID, bowsByID)
	if err != nil {
		return "", err
	}

	componentValues, err := itemComponentsForGive(entry, spellMeta)
	if err != nil {
		return "", err
	}
	if spellMeta.hasUseSpell {
		componentValues["minecraft:consumable"] = bookConsumableSNBT
	}

	customData, err := itemCustomData(entry, activesByID, passivesByID, bowsByID)
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

func itemComponentsForGive(entry itemModel.Item, spellMeta itemSpellMeta) (map[string]string, error) {
	normalizedComponents, errMsg := itemModel.NormalizeComponents(componentData(entry))
	if errMsg != "" {
		return nil, fmt.Errorf("item(%s): %s", entry.ID, errMsg)
	}

	values := make(map[string]string, len(normalizedComponents))
	rawValues := make(map[string]any, len(normalizedComponents))
	for _, component := range normalizedComponents {
		rawValues[component.Key] = component.Value
		value, ok := valueToSNBT(component.Value)
		if !ok {
			continue
		}
		values[component.Key] = value
	}
	if lore, ok := mergedItemSkillLore(rawValues["minecraft:lore"], spellMeta); ok {
		if value, valueOK := valueToSNBT(lore); valueOK {
			values["minecraft:lore"] = value
		}
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
	activesByID map[string]activeModel.Active,
	passivesByID map[string]passiveModel.Passive,
	bowsByID map[string]bowModel.BowPassive,
) (string, error) {
	itemSNBT, _ := itemModel.BuildItemComponents(entry)
	parts := []string{
		fmt.Sprintf("item_id:%s", SNBTString(entry.ItemID)),
		fmt.Sprintf("source_id:%s", SNBTString(entry.ID)),
		fmt.Sprintf("nbt_snapshot:%s", itemSNBT),
	}

	spellMeta, err := resolveItemSpellMeta(entry, activesByID, passivesByID, bowsByID)
	if err != nil {
		return "", err
	}
	if spellMeta.activeID != "" {
		parts = append(parts, fmt.Sprintf("active_id:%s", SNBTString(spellMeta.activeID)))
	}
	if len(spellMeta.customFragments) > 0 {
		parts = append(parts, spellMeta.customFragments...)
	}
	if entry.Maf.MaxMP != nil {
		parts = append(parts, fmt.Sprintf("maxmp:%d", *entry.Maf.MaxMP))
	}
	return "{maf:{" + strings.Join(parts, ",") + "}}", nil
}

func itemComponentsForLoot(
	entry itemModel.Item,
	activesByID map[string]activeModel.Active,
	passivesByID map[string]passiveModel.Passive,
	bowsByID map[string]bowModel.BowPassive,
) (map[string]any, error) {
	components := deepCopyMapFromAny(componentData(entry))
	delete(components, "minecraft:enchantments")
	delete(components, "minecraft:custom_data")

	spellMeta, err := resolveItemSpellMeta(entry, activesByID, passivesByID, bowsByID)
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
	if lore, ok := mergedItemSkillLore(components["minecraft:lore"], spellMeta); ok {
		components["minecraft:lore"] = lore
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
	hasActive       bool
	activeID        string
	active          activeModel.Active
	hasPassive      bool
	passive         passiveModel.Passive
	hasBow          bool
	bow             bowModel.BowPassive
	customFragments []string
}

func resolveItemSpellMeta(
	entry itemModel.Item,
	activesByID map[string]activeModel.Active,
	passivesByID map[string]passiveModel.Passive,
	bowsByID map[string]bowModel.BowPassive,
) (itemSpellMeta, error) {
	meta := itemSpellMeta{}

	activeID := strings.TrimSpace(entry.Maf.ActiveID)
	if activeID != "" {
		active, ok := activesByID[activeID]
		if !ok {
			return itemSpellMeta{}, fmt.Errorf("item(%s): referenced active not found (%s)", entry.ID, activeID)
		}
		meta.hasUseSpell = true
		meta.hasActive = true
		meta.activeID = active.ID
		meta.active = active
	}

	bowID := strings.TrimSpace(entry.Maf.BowID)
	if bowID != "" {
		if activeID != "" {
			return itemSpellMeta{}, fmt.Errorf("item(%s): bowId cannot be combined with activeId", entry.ID)
		}
		if strings.TrimSpace(entry.Maf.PassiveID) != "" {
			return itemSpellMeta{}, fmt.Errorf("item(%s): bowId cannot be combined with passiveId", entry.ID)
		}
		bow, ok := bowsByID[bowID]
		if !ok {
			return itemSpellMeta{}, fmt.Errorf("item(%s): referenced bow not found (%s)", entry.ID, bowID)
		}
		meta.hasBow = true
		meta.bow = bow
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
	meta.hasPassive = true
	meta.passive = passive
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

func mergedItemSkillLore(existing any, spellMeta itemSpellMeta) ([]any, bool) {
	skillLore := itemSkillLoreComponents(spellMeta)
	if len(skillLore) == 0 {
		return nil, false
	}

	lore := anyList(existing)
	lore = append(lore, skillLore...)
	return lore, true
}

func itemSkillLoreComponents(spellMeta itemSpellMeta) []any {
	var lore []any
	if spellMeta.hasActive {
		lore = append(lore, activeItemSkillLoreComponents(spellMeta.active)...)
	}
	if spellMeta.hasPassive {
		lore = append(lore, passiveItemSkillLoreComponents(spellMeta.passive)...)
	}
	if spellMeta.hasBow {
		lore = append(lore, bowItemSkillLoreComponents(spellMeta.bow)...)
	}
	return lore
}

func activeItemSkillLoreComponents(entry activeModel.Active) []any {
	return []any{
		itemSkillLoreBlankLine(),
		itemSkillLoreComponent("アクティブスキル", "light_purple", false),
		itemSkillLoreComponent(strings.TrimSpace(entry.Title), "white", true),
		itemSkillLoreComponent(strings.TrimSpace(entry.Description), "aqua", false),
		itemSkillLoreComponent(fmt.Sprintf("cast   : %d", entry.CastTime), "white", false),
		itemSkillLoreComponent(fmt.Sprintf("MP     : %d", entry.MPCost), "white", false),
		itemSkillLoreComponent(fmt.Sprintf("recast : %d", entry.CoolTime), "white", false),
		itemSkillLoreComponent(fmt.Sprintf("target : %s", strings.TrimSpace(entry.Target)), "white", false),
		itemSkillLoreComponent(fmt.Sprintf("range  : %d", entry.Range), "white", false),
	}
}

func passiveItemSkillLoreComponents(entry passiveModel.Passive) []any {
	return []any{
		itemSkillLoreBlankLine(),
		itemSkillLoreComponent("パッシブスキル", "gold", false),
		itemSkillLoreComponent(passiveDisplayName(entry), "white", true),
		itemSkillLoreComponent(passiveLoreLine(entry), "aqua", false),
		itemSkillLoreComponent(fmt.Sprintf("MP     : %d", PassiveMPCost), "white", false),
		itemSkillLoreComponent(fmt.Sprintf("target : %s", strings.TrimSpace(entry.Target)), "white", false),
		itemSkillLoreComponent(fmt.Sprintf("range  : %d", entry.Range), "white", false),
	}
}

func bowItemSkillLoreComponents(entry bowModel.BowPassive) []any {
	return []any{
		itemSkillLoreBlankLine(),
		itemSkillLoreComponent("パッシブスキル", "gold", false),
		itemSkillLoreComponent(bowDisplayName(entry), "white", true),
		itemSkillLoreComponent(strings.TrimSpace(entry.Lore), "aqua", false),
		itemSkillLoreComponent(fmt.Sprintf("MP     : %d", entry.MPCost), "white", false),
		itemSkillLoreComponent(fmt.Sprintf("target : %s", strings.TrimSpace(entry.Target)), "white", false),
		itemSkillLoreComponent(fmt.Sprintf("range  : %d", entry.Range), "white", false),
	}
}

func passiveDisplayName(entry passiveModel.Passive) string {
	name := strings.TrimSpace(entry.Name)
	if name == "" {
		return entry.ID
	}
	return name
}

func bowDisplayName(entry bowModel.BowPassive) string {
	name := strings.TrimSpace(entry.Name)
	if name == "" {
		return entry.ID
	}
	return name
}

func itemSkillLoreBlankLine() map[string]any {
	return itemSkillLoreComponent("", "white", false)
}

func itemSkillLoreComponent(text, color string, bold bool) map[string]any {
	component := map[string]any{
		"font":   "minecraft:uniform",
		"text":   text,
		"color":  color,
		"italic": false,
	}
	if bold {
		component["bold"] = true
	}
	return component
}
