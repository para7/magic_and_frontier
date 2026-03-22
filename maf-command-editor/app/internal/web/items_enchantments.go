package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"tools2/app/internal/webui"
)

type itemEnchantmentDef struct {
	Key      string
	MaxLevel int
}

var itemEnchantmentCatalog = []itemEnchantmentDef{
	{Key: "aqua_affinity", MaxLevel: 1},
	{Key: "bane_of_arthropods", MaxLevel: 5},
	{Key: "binding_curse", MaxLevel: 1},
	{Key: "blast_protection", MaxLevel: 4},
	{Key: "breach", MaxLevel: 4},
	{Key: "channeling", MaxLevel: 1},
	{Key: "density", MaxLevel: 5},
	{Key: "depth_strider", MaxLevel: 3},
	{Key: "efficiency", MaxLevel: 5},
	{Key: "feather_falling", MaxLevel: 4},
	{Key: "fire_aspect", MaxLevel: 2},
	{Key: "fire_protection", MaxLevel: 4},
	{Key: "flame", MaxLevel: 1},
	{Key: "fortune", MaxLevel: 3},
	{Key: "frost_walker", MaxLevel: 2},
	{Key: "impaling", MaxLevel: 5},
	{Key: "infinity", MaxLevel: 1},
	{Key: "knockback", MaxLevel: 2},
	{Key: "looting", MaxLevel: 3},
	{Key: "loyalty", MaxLevel: 3},
	{Key: "luck_of_the_sea", MaxLevel: 3},
	{Key: "lunge", MaxLevel: 3},
	{Key: "lure", MaxLevel: 3},
	{Key: "mending", MaxLevel: 1},
	{Key: "multishot", MaxLevel: 1},
	{Key: "piercing", MaxLevel: 4},
	{Key: "power", MaxLevel: 5},
	{Key: "projectile_protection", MaxLevel: 4},
	{Key: "protection", MaxLevel: 4},
	{Key: "punch", MaxLevel: 2},
	{Key: "quick_charge", MaxLevel: 3},
	{Key: "respiration", MaxLevel: 3},
	{Key: "riptide", MaxLevel: 3},
	{Key: "sharpness", MaxLevel: 5},
	{Key: "silk_touch", MaxLevel: 1},
	{Key: "smite", MaxLevel: 5},
	{Key: "soul_speed", MaxLevel: 3},
	{Key: "sweeping_edge", MaxLevel: 3},
	{Key: "swift_sneak", MaxLevel: 3},
	{Key: "thorns", MaxLevel: 3},
	{Key: "unbreaking", MaxLevel: 3},
	{Key: "vanishing_curse", MaxLevel: 1},
	{Key: "wind_burst", MaxLevel: 3},
}

func itemFormEnchantmentsFromText(text string) ([]webui.ItemEnchantmentOption, int) {
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

func itemFormEnchantmentsFromRequest(r *http.Request) (string, []webui.ItemEnchantmentOption, int) {
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
			level = strconv.Itoa(def.MaxLevel)
		}
		levels[id] = level
		if selected[id] {
			lines = append(lines, fmt.Sprintf("%s %s", id, level))
		}
	}

	options, selectedCount := buildItemEnchantmentOptions(selected, levels)
	return strings.Join(lines, "\n"), options, selectedCount
}

func buildItemEnchantmentOptions(selected map[string]bool, levels map[string]string) ([]webui.ItemEnchantmentOption, int) {
	options := make([]webui.ItemEnchantmentOption, 0, len(itemEnchantmentCatalog))
	selectedCount := 0

	for _, def := range itemEnchantmentCatalog {
		id := itemEnchantmentID(def.Key)
		level := strconv.Itoa(def.MaxLevel)
		if raw := strings.TrimSpace(levels[id]); raw != "" {
			level = raw
		}
		checked := selected[id]
		if checked {
			selectedCount++
		}
		options = append(options, webui.ItemEnchantmentOption{
			ID:             id,
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
