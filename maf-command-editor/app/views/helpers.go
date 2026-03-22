package views

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/loottables"
	"tools2/app/internal/domain/spawntables"
	"tools2/app/internal/domain/treasures"
	"tools2/app/internal/webui"
)

func FieldError(errs map[string]string, key string) string {
	if errs == nil {
		return ""
	}
	return errs[key]
}

func HasValue(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func EnchantmentCategories(options []webui.ItemEnchantmentOption) []string {
	categories := make([]string, 0, len(options))
	seen := map[string]bool{}
	for _, option := range options {
		category := strings.TrimSpace(option.Category)
		if category == "" || seen[category] {
			continue
		}
		seen[category] = true
		categories = append(categories, category)
	}
	return categories
}

func EnchantmentsByCategory(options []webui.ItemEnchantmentOption, category string) []webui.ItemEnchantmentOption {
	filtered := make([]webui.ItemEnchantmentOption, 0, len(options))
	for _, option := range options {
		if option.Category == category {
			filtered = append(filtered, option)
		}
	}
	return filtered
}

func NoticeClass(notice *webui.Notice) string {
	if notice == nil || notice.Kind == "" {
		return "notice-info"
	}
	switch notice.Kind {
	case "success":
		return "notice-success"
	case "error":
		return "notice-error"
	default:
		return "notice-info"
	}
}

func FloatText(value *float64) string {
	if value == nil {
		return "-"
	}
	return trimFloat(*value)
}

func trimFloat(value float64) string {
	text := fmt.Sprintf("%.4f", value)
	text = strings.TrimRight(text, "0")
	text = strings.TrimRight(text, ".")
	if text == "" {
		return "0"
	}
	return text
}

func IntText(value int) string {
	return strconv.Itoa(value)
}

func BoolText(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func DisplayText(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func BracketText(value string) string {
	return "【" + strings.TrimSpace(value) + "】"
}

func ItemListTitle(itemID string, name string) string {
	title := BracketText(itemID)
	name = strings.TrimSpace(name)
	if name == "" {
		return title
	}
	return title + " " + name
}

func EnemyListTitle(mobType string, name string) string {
	title := BracketText(mobType)
	name = strings.TrimSpace(name)
	if name == "" {
		return title
	}
	return title + name
}

func JoinCompactWithComma(values ...string) string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			out = append(out, value)
		}
	}
	if len(out) == 0 {
		return "-"
	}
	return strings.Join(out, ", ")
}

func JoinCompactWithSpace(values ...string) string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			out = append(out, value)
		}
	}
	if len(out) == 0 {
		return "-"
	}
	return strings.Join(out, " ")
}

func SubmitLabel(editing bool) string {
	if editing {
		return "更新"
	}
	return "新規追加"
}

func FormTitle(label string, editing bool) string {
	if editing {
		return "Edit " + label
	}
	return "New " + label
}

func FormAction(basePath string, editing bool) string {
	if editing {
		return basePath + "/edit"
	}
	return basePath + "/new"
}

func BackLink(listPath string, returnTo string) string {
	if strings.TrimSpace(returnTo) != "" {
		return returnTo
	}
	return listPath
}

func WithReturnTo(path string, returnTo string) string {
	if strings.TrimSpace(returnTo) == "" {
		return path
	}
	ref, err := url.Parse(path)
	if err != nil {
		return path
	}
	query := ref.Query()
	query.Set("returnTo", returnTo)
	ref.RawQuery = query.Encode()
	return ref.String()
}

func JoinLines(values []string) string {
	if len(values) == 0 {
		return "-"
	}
	return strings.Join(values, ", ")
}

func SearchBlob(values ...string) string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.ToLower(strings.TrimSpace(value))
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return strings.Join(out, " ")
}

func JoinTokens(values ...string) string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return strings.Join(out, " ")
}

func NameFilterToken(value string) string {
	if strings.TrimSpace(value) == "" {
		return "unnamed"
	}
	return "named"
}

func SkillFilterTokens(skillID string) string {
	if strings.TrimSpace(skillID) == "" {
		return JoinTokens("all", "without-skill")
	}
	return JoinTokens("all", "with-skill")
}

func NameFilterTokens(value string) string {
	return JoinTokens("all", NameFilterToken(value))
}

func TreasureSearchBlob(entry treasures.TreasureEntry) string {
	values := []string{entry.ID, entry.TablePath, entry.UpdatedAt}
	for _, pool := range entry.LootPools {
		values = append(values, pool.Kind, pool.RefID, trimFloat(pool.Weight))
		if pool.CountMin != nil {
			values = append(values, trimFloat(*pool.CountMin))
		}
		if pool.CountMax != nil {
			values = append(values, trimFloat(*pool.CountMax))
		}
	}
	return SearchBlob(values...)
}

func LootTableSearchBlob(entry loottables.LootTableEntry) string {
	values := []string{entry.ID, entry.UpdatedAt}
	for _, pool := range entry.LootPools {
		values = append(values, pool.Kind, pool.RefID, trimFloat(pool.Weight))
		if pool.CountMin != nil {
			values = append(values, trimFloat(*pool.CountMin))
		}
		if pool.CountMax != nil {
			values = append(values, trimFloat(*pool.CountMax))
		}
	}
	return SearchBlob(values...)
}

func EnemySearchBlob(entry enemies.EnemyEntry) string {
	values := []string{
		entry.ID,
		entry.MobType,
		entry.Name,
		entry.Memo,
		entry.DropMode,
		entry.UpdatedAt,
		trimFloat(entry.HP),
		JoinLines(entry.EnemySkillIDs),
	}
	if entry.Attack != nil {
		values = append(values, trimFloat(*entry.Attack))
	}
	if entry.Defense != nil {
		values = append(values, trimFloat(*entry.Defense))
	}
	if entry.MoveSpeed != nil {
		values = append(values, trimFloat(*entry.MoveSpeed))
	}
	collectSlot := func(slot *enemies.EquipmentSlot) {
		if slot == nil {
			return
		}
		values = append(values, slot.Kind, slot.RefID, strconv.Itoa(slot.Count))
		if slot.DropChance != nil {
			values = append(values, trimFloat(*slot.DropChance))
		}
	}
	collectSlot(entry.Equipment.Mainhand)
	collectSlot(entry.Equipment.Offhand)
	collectSlot(entry.Equipment.Head)
	collectSlot(entry.Equipment.Chest)
	collectSlot(entry.Equipment.Legs)
	collectSlot(entry.Equipment.Feet)
	for _, drop := range entry.Drops {
		values = append(values, drop.Kind, drop.RefID, trimFloat(drop.Weight))
		if drop.CountMin != nil {
			values = append(values, trimFloat(*drop.CountMin))
		}
		if drop.CountMax != nil {
			values = append(values, trimFloat(*drop.CountMax))
		}
	}
	return SearchBlob(values...)
}

func SpawnTableSearchBlob(entry spawntables.SpawnTableEntry) string {
	values := []string{
		entry.ID,
		entry.SourceMobType,
		entry.Dimension,
		strconv.Itoa(entry.MinX),
		strconv.Itoa(entry.MaxX),
		strconv.Itoa(entry.MinY),
		strconv.Itoa(entry.MaxY),
		strconv.Itoa(entry.MinZ),
		strconv.Itoa(entry.MaxZ),
		strconv.Itoa(entry.BaseMobWeight),
		entry.UpdatedAt,
	}
	for _, replacement := range entry.Replacements {
		values = append(values, replacement.EnemyID, strconv.Itoa(replacement.Weight))
	}
	return SearchBlob(values...)
}

func PageSizeOptions() []webui.SelectOption {
	return []webui.SelectOption{
		{Value: "25", Label: "25 / page"},
		{Value: "50", Label: "50 / page"},
		{Value: "100", Label: "100 / page"},
	}
}

func ItemsFilterOptions() []webui.SelectOption {
	return []webui.SelectOption{
		{Value: "all", Label: "All items"},
		{Value: "with-skill", Label: "With skill"},
		{Value: "without-skill", Label: "Without skill"},
	}
}

func ItemsSortOptions() []webui.SelectOption {
	return []webui.SelectOption{
		{Value: "updated_desc", Label: "Updated desc"},
		{Value: "updated_asc", Label: "Updated asc"},
		{Value: "id_asc", Label: "ID asc"},
		{Value: "id_desc", Label: "ID desc"},
		{Value: "item_id_asc", Label: "Item ID asc"},
		{Value: "item_id_desc", Label: "Item ID desc"},
	}
}

func GrimoireFilterOptions() []webui.SelectOption {
	return []webui.SelectOption{{Value: "all", Label: "All entries"}}
}

func GrimoireSortOptions() []webui.SelectOption {
	return []webui.SelectOption{
		{Value: "updated_desc", Label: "Updated desc"},
		{Value: "updated_asc", Label: "Updated asc"},
		{Value: "id_asc", Label: "ID asc"},
		{Value: "id_desc", Label: "ID desc"},
		{Value: "title_asc", Label: "Title asc"},
		{Value: "title_desc", Label: "Title desc"},
		{Value: "cast_id_asc", Label: "Cast ID asc"},
		{Value: "cast_id_desc", Label: "Cast ID desc"},
	}
}

func NamedFilterOptions(label string) []webui.SelectOption {
	return []webui.SelectOption{
		{Value: "all", Label: "All " + label},
		{Value: "named", Label: "Named"},
		{Value: "unnamed", Label: "Unnamed"},
	}
}

func NameSortOptions() []webui.SelectOption {
	return []webui.SelectOption{
		{Value: "updated_desc", Label: "Updated desc"},
		{Value: "updated_asc", Label: "Updated asc"},
		{Value: "id_asc", Label: "ID asc"},
		{Value: "id_desc", Label: "ID desc"},
		{Value: "name_asc", Label: "Name asc"},
		{Value: "name_desc", Label: "Name desc"},
	}
}

func TreasureFilterOptions() []webui.SelectOption {
	return []webui.SelectOption{{Value: "all", Label: "All treasures"}}
}

func TreasureSortOptions() []webui.SelectOption {
	return []webui.SelectOption{
		{Value: "updated_desc", Label: "Updated desc"},
		{Value: "updated_asc", Label: "Updated asc"},
		{Value: "id_asc", Label: "ID asc"},
		{Value: "id_desc", Label: "ID desc"},
	}
}

func LootTableFilterOptions() []webui.SelectOption {
	return []webui.SelectOption{{Value: "all", Label: "All loottables"}}
}

func LootTableSortOptions() []webui.SelectOption {
	return []webui.SelectOption{
		{Value: "updated_desc", Label: "Updated desc"},
		{Value: "updated_asc", Label: "Updated asc"},
		{Value: "id_asc", Label: "ID asc"},
		{Value: "id_desc", Label: "ID desc"},
		{Value: "table_path_asc", Label: "Table path asc"},
		{Value: "table_path_desc", Label: "Table path desc"},
	}
}

func EnemyFilterOptions() []webui.SelectOption {
	return []webui.SelectOption{
		{Value: "all", Label: "All enemies"},
		{Value: "append", Label: "Drop append"},
		{Value: "replace", Label: "Drop replace"},
	}
}

func EnemySortOptions() []webui.SelectOption {
	return []webui.SelectOption{
		{Value: "updated_desc", Label: "Updated desc"},
		{Value: "updated_asc", Label: "Updated asc"},
		{Value: "id_asc", Label: "ID asc"},
		{Value: "id_desc", Label: "ID desc"},
		{Value: "mob_type_asc", Label: "Mob type asc"},
		{Value: "mob_type_desc", Label: "Mob type desc"},
		{Value: "hp_desc", Label: "HP desc"},
		{Value: "hp_asc", Label: "HP asc"},
	}
}
