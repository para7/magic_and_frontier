package model

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

var (
	resourceIDPattern = regexp.MustCompile(`^[a-z0-9_.-]+:[a-z0-9_./-]+$`)
)

// IsNamespacedResourceID は "namespace:path" 形式かどうかを検証する。
func IsNamespacedResourceID(value string) bool {
	return resourceIDPattern.MatchString(strings.TrimSpace(value))
}

// IsSafeNamespacedResourcePath は namespaced リソースパスとして安全かどうかを検証する。
func IsSafeNamespacedResourcePath(value string) bool {
	value = strings.TrimSpace(value)
	if !resourceIDPattern.MatchString(value) {
		return false
	}
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return false
	}
	return hasSafeResourcePathSegments(parts[1])
}

func hasSafeResourcePathSegments(value string) bool {
	path := NormalizeResourcePath(value)
	if path == "" {
		return false
	}
	for _, segment := range strings.Split(path, "/") {
		if segment == "" || segment == "." || segment == ".." {
			return false
		}
	}
	return true
}

// NormalizeResourcePath はリソースパスを正規化する（バックスラッシュをスラッシュに変換し、前後のスラッシュを除去）。
func NormalizeResourcePath(value string) string {
	value = strings.ReplaceAll(value, "\\", "/")
	value = strings.TrimSpace(value)
	return strings.Trim(value, "/")
}

// ValidateDropRefs は DropRef スライスのリレーションバリデーションを行う。
// entity: エンティティ名（エラーメッセージ用）, id: エンティティID, prefix: フィールドプレフィックス（例: "drops"）
func ValidateDropRefs(entity, id, prefix string, drops []DropRef, mas DBMaster) []ValidationError {
	var errs []ValidationError
	for i, d := range drops {
		kind := strings.TrimSpace(d.Kind)
		refID := strings.TrimSpace(d.RefID)
		if refID != "" {
			switch kind {
			case "item":
				if d.Slot != nil {
					errs = append(errs, ValidationError{
						Entity: entity, ID: id,
						Field: fmt.Sprintf("%s[%d].slot", prefix, i),
						Tag:   "relation", Param: "slot is only supported when kind=passive",
					})
				}
				if !mas.HasItem(refID) {
					errs = append(errs, ValidationError{
						Entity: entity, ID: id,
						Field: fmt.Sprintf("%s[%d].refId", prefix, i),
						Tag:   "relation", Param: "item not found",
					})
				}
			case "grimoire":
				if d.Slot != nil {
					errs = append(errs, ValidationError{
						Entity: entity, ID: id,
						Field: fmt.Sprintf("%s[%d].slot", prefix, i),
						Tag:   "relation", Param: "slot is only supported when kind=passive",
					})
				}
				if !mas.HasGrimoire(refID) {
					errs = append(errs, ValidationError{
						Entity: entity, ID: id,
						Field: fmt.Sprintf("%s[%d].refId", prefix, i),
						Tag:   "relation", Param: "grimoire not found",
					})
				}
			case "passive":
				if !mas.HasPassive(refID) {
					errs = append(errs, ValidationError{
						Entity: entity, ID: id,
						Field: fmt.Sprintf("%s[%d].refId", prefix, i),
						Tag:   "relation", Param: "passive not found",
					})
				}
				if d.Slot == nil {
					errs = append(errs, ValidationError{
						Entity: entity, ID: id,
						Field: fmt.Sprintf("%s[%d].slot", prefix, i),
						Tag:   "relation", Param: "slot is required when kind=passive",
					})
				}
			case "minecraft_item":
				if d.Slot != nil {
					errs = append(errs, ValidationError{
						Entity: entity, ID: id,
						Field: fmt.Sprintf("%s[%d].slot", prefix, i),
						Tag:   "relation", Param: "slot is only supported when kind=passive",
					})
				}
				if !IsNamespacedResourceID(refID) {
					errs = append(errs, ValidationError{
						Entity: entity, ID: id,
						Field: fmt.Sprintf("%s[%d].refId", prefix, i),
						Tag:   "relation", Param: "invalid minecraft item id",
					})
				}
			}
		}
		if d.CountMin != nil && d.CountMax != nil && *d.CountMin > *d.CountMax {
			errs = append(errs, ValidationError{
				Entity: entity, ID: id,
				Field: fmt.Sprintf("%s[%d].countMin", prefix, i),
				Tag:   "lte", Param: "countMax",
			})
		}
	}
	return errs
}

// ValidateMafLootPools は maf:* エントリを含む loot pools の共通バリデーションを行う。
func ValidateMafLootPools(entity, id, prefix string, pools []any, mas DBMaster) []ValidationError {
	var errs []ValidationError
	for i, rawPool := range pools {
		pool, ok := rawPool.(map[string]any)
		if !ok {
			errs = append(errs, ValidationError{
				Entity: entity, ID: id,
				Field: fmt.Sprintf("%s[%d]", prefix, i),
				Tag:   "format", Param: "pool must be an object",
			})
			continue
		}

		rawEntries, exists := pool["entries"]
		if !exists || rawEntries == nil {
			errs = append(errs, ValidationError{
				Entity: entity, ID: id,
				Field: fmt.Sprintf("%s[%d].entries", prefix, i),
				Tag:   "format", Param: "entries must be an array",
			})
			continue
		}
		entries, ok := rawEntries.([]any)
		if !ok {
			errs = append(errs, ValidationError{
				Entity: entity, ID: id,
				Field: fmt.Sprintf("%s[%d].entries", prefix, i),
				Tag:   "format", Param: "entries must be an array",
			})
			continue
		}

		for j, rawEntry := range entries {
			entryField := fmt.Sprintf("%s[%d].entries[%d]", prefix, i, j)
			entry, ok := rawEntry.(map[string]any)
			if !ok {
				errs = append(errs, ValidationError{
					Entity: entity, ID: id,
					Field: entryField,
					Tag:   "format", Param: "entry must be an object",
				})
				continue
			}

			rawType, _ := entry["type"].(string)
			entryType := strings.TrimSpace(rawType)
			if entryType == "" {
				errs = append(errs, ValidationError{
					Entity: entity, ID: id,
					Field: entryField + ".type",
					Tag:   "format", Param: "type is required",
				})
				continue
			}

			if !strings.HasPrefix(entryType, "maf:") {
				continue
			}

			rawName, _ := entry["name"].(string)
			refID := strings.TrimSpace(rawName)
			if refID == "" {
				errs = append(errs, ValidationError{
					Entity: entity, ID: id,
					Field: entryField + ".name",
					Tag:   "format", Param: "name is required",
				})
				continue
			}

			switch entryType {
			case "maf:item":
				if hasEntrySlot(entry) {
					errs = append(errs, ValidationError{
						Entity: entity, ID: id,
						Field: entryField + ".slot",
						Tag:   "relation", Param: "slot is only supported when type=maf:passive",
					})
				}
				if !mas.HasItem(refID) {
					errs = append(errs, ValidationError{
						Entity: entity, ID: id,
						Field: entryField + ".name",
						Tag:   "relation", Param: "item not found",
					})
				}
			case "maf:grimoire":
				if hasEntrySlot(entry) {
					errs = append(errs, ValidationError{
						Entity: entity, ID: id,
						Field: entryField + ".slot",
						Tag:   "relation", Param: "slot is only supported when type=maf:passive",
					})
				}
				if !mas.HasGrimoire(refID) {
					errs = append(errs, ValidationError{
						Entity: entity, ID: id,
						Field: entryField + ".name",
						Tag:   "relation", Param: "grimoire not found",
					})
				}
			case "maf:passive":
				if _, err := ParseLootEntrySlot(entry); err != nil {
					errs = append(errs, ValidationError{
						Entity: entity, ID: id,
						Field: entryField + ".slot",
						Tag:   "format", Param: err.Error(),
					})
				}
				if !mas.HasPassive(refID) {
					errs = append(errs, ValidationError{
						Entity: entity, ID: id,
						Field: entryField + ".name",
						Tag:   "relation", Param: "passive not found",
					})
				}
			default:
				errs = append(errs, ValidationError{
					Entity: entity, ID: id,
					Field: entryField + ".type",
					Tag:   "format", Param: "unsupported maf entry type",
				})
				continue
			}

			if _, _, err := ParseLootEntryCount(entry); err != nil {
				errs = append(errs, ValidationError{
					Entity: entity, ID: id,
					Field: entryField + ".count",
					Tag:   "format", Param: err.Error(),
				})
			}
		}
	}
	return errs
}

// ParseLootEntryCount は loot entry の count（数値または {min,max}）を解析する。
func ParseLootEntryCount(entry map[string]any) (*float64, *float64, error) {
	rawCount, exists := entry["count"]
	if !exists || rawCount == nil {
		return nil, nil, nil
	}

	if value, ok := anyNumberToFloat64(rawCount); ok {
		if err := validateCountValue(value, "count"); err != nil {
			return nil, nil, err
		}
		min := value
		max := value
		return &min, &max, nil
	}

	rangeMap, ok := rawCount.(map[string]any)
	if !ok {
		return nil, nil, fmt.Errorf("count must be a number or object")
	}
	rawMin, hasMin := rangeMap["min"]
	rawMax, hasMax := rangeMap["max"]
	if !hasMin || !hasMax {
		return nil, nil, fmt.Errorf("count object requires min and max")
	}
	min, ok := anyNumberToFloat64(rawMin)
	if !ok {
		return nil, nil, fmt.Errorf("count.min must be numeric")
	}
	max, ok := anyNumberToFloat64(rawMax)
	if !ok {
		return nil, nil, fmt.Errorf("count.max must be numeric")
	}
	if err := validateCountValue(min, "count.min"); err != nil {
		return nil, nil, err
	}
	if err := validateCountValue(max, "count.max"); err != nil {
		return nil, nil, err
	}
	if min > max {
		return nil, nil, fmt.Errorf("count.min must be less than or equal to count.max")
	}
	return &min, &max, nil
}

// ParseLootEntrySlot は maf:passive 用 slot を解析する。
func ParseLootEntrySlot(entry map[string]any) (int, error) {
	rawSlot, exists := entry["slot"]
	if !exists || rawSlot == nil {
		return 0, fmt.Errorf("slot is required when type=maf:passive")
	}
	value, ok := anyNumberToFloat64(rawSlot)
	if !ok || value != float64(int(value)) {
		return 0, fmt.Errorf("slot must be an integer")
	}
	slot := int(value)
	if slot < 1 || slot > 3 {
		return 0, fmt.Errorf("slot must be between 1 and 3")
	}
	return slot, nil
}

func hasEntrySlot(entry map[string]any) bool {
	value, exists := entry["slot"]
	return exists && value != nil
}

func validateCountValue(value float64, field string) error {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return fmt.Errorf("%s must be numeric", field)
	}
	if value < 1 || value > 64 {
		return fmt.Errorf("%s must be between 1 and 64", field)
	}
	return nil
}

func anyNumberToFloat64(value any) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	default:
		return 0, false
	}
}

// ValidateEquipmentSlots は Equipment の全スロットのリレーションバリデーションを行う。
func ValidateEquipmentSlots(entity, id string, equip Equipment, mas DBMaster) []ValidationError {
	var errs []ValidationError
	slots := []struct {
		name string
		slot *EquipmentSlot
	}{
		{"equipment.mainhand", equip.Mainhand},
		{"equipment.offhand", equip.Offhand},
		{"equipment.head", equip.Head},
		{"equipment.chest", equip.Chest},
		{"equipment.legs", equip.Legs},
		{"equipment.feet", equip.Feet},
	}
	for _, s := range slots {
		if s.slot == nil {
			continue
		}
		kind := strings.TrimSpace(s.slot.Kind)
		refID := strings.TrimSpace(s.slot.RefID)
		switch kind {
		case "item":
			if !mas.HasItem(refID) {
				errs = append(errs, ValidationError{
					Entity: entity, ID: id,
					Field: s.name + ".refId",
					Tag:   "relation", Param: "item not found",
				})
			}
		case "minecraft_item":
			if !IsNamespacedResourceID(refID) {
				errs = append(errs, ValidationError{
					Entity: entity, ID: id,
					Field: s.name + ".refId",
					Tag:   "relation", Param: "invalid minecraft item id",
				})
			}
		}
	}
	return errs
}
