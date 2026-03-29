package model

import (
	"fmt"
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
				if !mas.HasItem(refID) {
					errs = append(errs, ValidationError{
						Entity: entity, ID: id,
						Field: fmt.Sprintf("%s[%d].refId", prefix, i),
						Tag:   "relation", Param: "item not found",
					})
				}
			case "grimoire":
				if !mas.HasGrimoire(refID) {
					errs = append(errs, ValidationError{
						Entity: entity, ID: id,
						Field: fmt.Sprintf("%s[%d].refId", prefix, i),
						Tag:   "relation", Param: "grimoire not found",
					})
				}
			case "minecraft_item":
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
