package loottables

import (
	"fmt"
	"time"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/treasures"
)

func ValidateSave(input SaveInput, itemIDs, grimoireIDs map[string]struct{}, now time.Time) common.SaveResult[LootTableEntry] {
	errs := common.ViolationsToFieldErrors(common.ValidateStruct(input), common.DefaultValidationMessage)
	id := common.RequirePrefixedSequenceID(errs, "id", input.ID, "loottable_")
	tablePath := common.NormalizeText(input.TablePath)
	if !common.IsSafeNamespacedResourcePath(tablePath) {
		errs.Add("tablePath", "Must be a namespaced loot table path.")
	}
	pools := make([]treasures.DropRef, 0, len(input.LootPools))
	for i, p := range input.LootPools {
		kind := common.NormalizeText(p.Kind)
		refID := common.NormalizeText(p.RefID)
		if refID != "" {
			switch kind {
			case "item":
				if !common.IsPrefixedSequenceID(refID, "items_") {
					errs.Add(fmt.Sprintf("lootPools.%d.refId", i), "Invalid ID format.")
				} else if _, ok := itemIDs[refID]; !ok {
					errs.Add(fmt.Sprintf("lootPools.%d.refId", i), "Referenced entry does not exist.")
				}
			case "grimoire":
				if !common.IsPrefixedSequenceID(refID, "grimoire_") {
					errs.Add(fmt.Sprintf("lootPools.%d.refId", i), "Invalid ID format.")
				} else if _, ok := grimoireIDs[refID]; !ok {
					errs.Add(fmt.Sprintf("lootPools.%d.refId", i), "Referenced entry does not exist.")
				}
			case "minecraft_item":
				if !common.IsNamespacedResourceID(refID) {
					errs.Add(fmt.Sprintf("lootPools.%d.refId", i), "Must be a minecraft item id.")
				}
			}
		}
		cmin := p.CountMin
		cmax := p.CountMax
		if cmin != nil && cmax != nil && *cmin > *cmax {
			errs.Add(fmt.Sprintf("lootPools.%d.countMin", i), "Must be <= countMax.")
		}
		if _, invalid := errs[fmt.Sprintf("lootPools.%d.kind", i)]; invalid {
			continue
		}
		if _, invalid := errs[fmt.Sprintf("lootPools.%d.refId", i)]; invalid {
			continue
		}
		if _, invalid := errs[fmt.Sprintf("lootPools.%d.weight", i)]; invalid {
			continue
		}
		pools = append(pools, treasures.DropRef{Kind: kind, RefID: refID, Weight: p.Weight, CountMin: cmin, CountMax: cmax})
	}
	if errs.Any() {
		return common.SaveValidationError[LootTableEntry](errs, "Validation failed. Fix the highlighted fields.")
	}
	entry := LootTableEntry{
		ID:        id,
		TablePath: tablePath,
		LootPools: pools,
		UpdatedAt: now.UTC().Format(time.RFC3339),
	}
	return common.SaveSuccess(entry, common.SaveModeCreated)
}
