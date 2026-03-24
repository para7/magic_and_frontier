package web

import (
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"tools2/app/internal/web/ui"
)

type itemEnchantmentDef struct {
	Category string
	Key      string
	MaxLevel int
}

var itemEnchantmentCatalog = []itemEnchantmentDef{
	{Category: "Armor", Key: "blast_protection", MaxLevel: 4},
	{Category: "Armor", Key: "feather_falling", MaxLevel: 4},
	{Category: "Armor", Key: "fire_protection", MaxLevel: 4},
	{Category: "Armor", Key: "projectile_protection", MaxLevel: 4},
	{Category: "Armor", Key: "protection", MaxLevel: 4},
	{Category: "Armor", Key: "thorns", MaxLevel: 3},

	{Category: "Movement/Environment", Key: "aqua_affinity", MaxLevel: 1},
	{Category: "Movement/Environment", Key: "depth_strider", MaxLevel: 3},
	{Category: "Movement/Environment", Key: "frost_walker", MaxLevel: 2},
	{Category: "Movement/Environment", Key: "respiration", MaxLevel: 3},
	{Category: "Movement/Environment", Key: "soul_speed", MaxLevel: 3},
	{Category: "Movement/Environment", Key: "swift_sneak", MaxLevel: 3},

	{Category: "Weapon", Key: "bane_of_arthropods", MaxLevel: 5},
	{Category: "Weapon", Key: "breach", MaxLevel: 4},
	{Category: "Weapon", Key: "density", MaxLevel: 5},
	{Category: "Weapon", Key: "fire_aspect", MaxLevel: 2},
	{Category: "Weapon", Key: "knockback", MaxLevel: 2},
	{Category: "Weapon", Key: "looting", MaxLevel: 3},
	{Category: "Weapon", Key: "sharpness", MaxLevel: 5},
	{Category: "Weapon", Key: "smite", MaxLevel: 5},
	{Category: "Weapon", Key: "sweeping_edge", MaxLevel: 3},
	{Category: "Weapon", Key: "wind_burst", MaxLevel: 3},
	{Category: "Weapon", Key: "lunge", MaxLevel: 3},

	{Category: "Bow", Key: "flame", MaxLevel: 1},
	{Category: "Bow", Key: "power", MaxLevel: 5},
	{Category: "Bow", Key: "punch", MaxLevel: 2},

	{Category: "Crossbow", Key: "quick_charge", MaxLevel: 3},

	{Category: "Bow/Crossbow", Key: "multishot", MaxLevel: 1},
	{Category: "Bow/Crossbow", Key: "piercing", MaxLevel: 4},
	{Category: "Bow/Crossbow", Key: "infinity", MaxLevel: 1},

	{Category: "Trident", Key: "channeling", MaxLevel: 1},
	{Category: "Trident", Key: "impaling", MaxLevel: 5},
	{Category: "Trident", Key: "loyalty", MaxLevel: 3},
	{Category: "Trident", Key: "riptide", MaxLevel: 3},

	{Category: "Tools", Key: "efficiency", MaxLevel: 5},
	{Category: "Tools", Key: "fortune", MaxLevel: 3},
	{Category: "Tools", Key: "silk_touch", MaxLevel: 1},

	{Category: "Durability", Key: "mending", MaxLevel: 1},
	{Category: "Durability", Key: "unbreaking", MaxLevel: 3},

	{Category: "Fishing", Key: "luck_of_the_sea", MaxLevel: 3},
	{Category: "Fishing", Key: "lure", MaxLevel: 3},

	{Category: "Cursed", Key: "binding_curse", MaxLevel: 1},
	{Category: "Cursed", Key: "vanishing_curse", MaxLevel: 1},
}

func itemFormEnchantmentsFromText(text string) ([]ui.ItemEnchantmentOption, int) {
	levels := map[string]string{}
	selected := map[string]bool{}
	for _, line := range strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) < 2 {
			continue
		}
		id := fields[0]
		selected[id] = true
		levels[id] = fields[1]
	}
	return buildItemEnchantmentOptions(selected, levels)
}

func itemFormEnchantmentsFromRequest(r *http.Request) (string, []ui.ItemEnchantmentOption, int) {
	selected := map[string]bool{}
	for _, id := range r.Form["enchantmentIds"] {
		id = strings.TrimSpace(id)
		if id != "" {
			selected[id] = true
		}
	}

	levels := map[string]string{}
	lines := make([]string, 0, len(itemEnchantmentCatalog))
	for _, def := range itemEnchantmentCatalog {
		id := itemEnchantmentID(def.Key)
		level := strings.TrimSpace(r.Form.Get(itemEnchantmentLevelFieldName(def.Key)))
		if level == "" {
			level = "0"
		}
		levels[id] = level
		if selected[id] {
			lines = append(lines, fmt.Sprintf("%s %s", id, level))
		}
	}

	options, selectedCount := buildItemEnchantmentOptions(selected, levels)
	return strings.Join(lines, "\n"), options, selectedCount
}

func buildItemEnchantmentOptions(selected map[string]bool, levels map[string]string) ([]ui.ItemEnchantmentOption, int) {
	options := make([]ui.ItemEnchantmentOption, 0, len(itemEnchantmentCatalog))
	selectedCount := 0

	for _, def := range itemEnchantmentCatalog {
		id := itemEnchantmentID(def.Key)
		level := "0"
		if raw := strings.TrimSpace(levels[id]); raw != "" {
			level = raw
		}
		checked := selected[id]
		if checked {
			selectedCount++
		}
		options = append(options, ui.ItemEnchantmentOption{
			ID:             id,
			Category:       def.Category,
			Key:            def.Key,
			Label:          itemEnchantmentLabel(def.Key),
			MaxLevel:       def.MaxLevel,
			Checked:        checked,
			Level:          level,
			LevelFieldName: itemEnchantmentLevelFieldName(def.Key),
		})
	}

	return options, selectedCount
}

func itemEnchantmentID(key string) string {
	return "minecraft:" + key
}

func itemEnchantmentLevelFieldName(key string) string {
	return "enchantmentLevel." + key
}

func itemEnchantmentLabel(key string) string {
	parts := strings.Split(key, "_")
	for i, part := range parts {
		if i > 0 && (part == "of" || part == "the") {
			continue
		}
		parts[i] = capitalizeWord(part)
	}
	return strings.Join(parts, " ")
}

func capitalizeWord(value string) string {
	if value == "" {
		return value
	}
	runes := []rune(value)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
