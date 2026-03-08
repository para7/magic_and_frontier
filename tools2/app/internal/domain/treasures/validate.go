package treasures

import (
	"fmt"
	"time"

	"tools2/app/internal/domain/common"
)

func ValidateSave(input SaveInput, itemIDs, grimoireIDs map[string]struct{}, now time.Time) common.SaveResult[TreasureEntry] {
	errs := common.ViolationsToFieldErrors(common.ValidateStruct(input), common.DefaultValidationMessage)
	pools := make([]DropRef, 0, len(input.LootPools))
	for i, p := range input.LootPools {
		kind := common.NormalizeText(p.Kind)
		refID := common.NormalizeText(p.RefID)
		if refID != "" {
			if kind == "item" {
				if _, ok := itemIDs[refID]; !ok {
					errs.Add(fmt.Sprintf("lootPools.%d.refId", i), "Referenced entry does not exist.")
				}
			} else if _, ok := grimoireIDs[refID]; !ok {
				errs.Add(fmt.Sprintf("lootPools.%d.refId", i), "Referenced entry does not exist.")
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
		pools = append(pools, DropRef{Kind: kind, RefID: refID, Weight: p.Weight, CountMin: cmin, CountMax: cmax})
	}
	if errs.Any() {
		return common.SaveValidationError[TreasureEntry](errs, "Validation failed. Fix the highlighted fields.")
	}
	entry := TreasureEntry{
		ID:        common.NormalizeText(input.ID),
		Name:      common.NormalizeText(input.Name),
		LootPools: pools,
		UpdatedAt: now.UTC().Format(time.RFC3339),
	}
	return common.SaveSuccess(entry, common.SaveModeCreated)
}
