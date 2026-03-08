package enemies

import (
	"fmt"
	"time"

	"tools2/app/internal/domain/common"
)

func ValidateSave(input SaveInput, enemySkillIDs, itemIDs, grimoireIDs map[string]struct{}, now time.Time) common.SaveResult[EnemyEntry] {
	errs := common.ViolationsToFieldErrors(common.ValidateStruct(input), common.DefaultValidationMessage)
	attack := input.Attack
	defense := input.Defense
	moveSpeed := input.MoveSpeed

	normalizedEnemySkillIDs := make([]string, 0, len(input.EnemySkillIDs))
	seen := map[string]struct{}{}
	for i, sid := range input.EnemySkillIDs {
		idv := common.NormalizeText(sid)
		if idv == "" {
			continue
		}
		if _, ok := enemySkillIDs[idv]; !ok {
			errs.Add(fmt.Sprintf("enemySkillIds.%d", i), "Referenced enemy skill does not exist.")
			continue
		}
		if _, exists := seen[idv]; !exists {
			seen[idv] = struct{}{}
			normalizedEnemySkillIDs = append(normalizedEnemySkillIDs, idv)
		}
	}

	if input.SpawnRule.Distance.Min > input.SpawnRule.Distance.Max {
		errs.Add("spawnRule.distance.min", "Must be <= distance.max.")
	}

	var axis *AxisBounds
	if input.SpawnRule.AxisBounds != nil {
		a := input.SpawnRule.AxisBounds
		axis = &AxisBounds{
			XMin: a.XMin,
			XMax: a.XMax,
			YMin: a.YMin,
			YMax: a.YMax,
			ZMin: a.ZMin,
			ZMax: a.ZMax,
		}
		if axis.XMin != nil && axis.XMax != nil && *axis.XMin > *axis.XMax {
			errs.Add("spawnRule.axisBounds.xMin", "Must be <= xMax.")
		}
		if axis.YMin != nil && axis.YMax != nil && *axis.YMin > *axis.YMax {
			errs.Add("spawnRule.axisBounds.yMin", "Must be <= yMax.")
		}
		if axis.ZMin != nil && axis.ZMax != nil && *axis.ZMin > *axis.ZMax {
			errs.Add("spawnRule.axisBounds.zMin", "Must be <= zMax.")
		}
	}

	dropTable := make([]DropRef, 0, len(input.DropTable))
	for i, d := range input.DropTable {
		kind := common.NormalizeText(d.Kind)
		refID := common.NormalizeText(d.RefID)
		if refID != "" {
			if kind == "item" {
				if _, ok := itemIDs[refID]; !ok {
					errs.Add(fmt.Sprintf("dropTable.%d.refId", i), "Referenced entry does not exist.")
				}
			} else if _, ok := grimoireIDs[refID]; !ok {
				errs.Add(fmt.Sprintf("dropTable.%d.refId", i), "Referenced entry does not exist.")
			}
		}
		countMin := d.CountMin
		countMax := d.CountMax
		if countMin != nil && countMax != nil && *countMin > *countMax {
			errs.Add(fmt.Sprintf("dropTable.%d.countMin", i), "Must be <= countMax.")
		}
		if _, invalid := errs[fmt.Sprintf("dropTable.%d.kind", i)]; invalid {
			continue
		}
		if _, invalid := errs[fmt.Sprintf("dropTable.%d.refId", i)]; invalid {
			continue
		}
		if _, invalid := errs[fmt.Sprintf("dropTable.%d.weight", i)]; invalid {
			continue
		}
		dropTable = append(dropTable, DropRef{Kind: kind, RefID: refID, Weight: d.Weight, CountMin: countMin, CountMax: countMax})
	}

	if errs.Any() {
		return common.SaveValidationError[EnemyEntry](errs, "Validation failed. Fix the highlighted fields.")
	}

	entry := EnemyEntry{
		ID:            common.NormalizeText(input.ID),
		Name:          common.NormalizeText(input.Name),
		HP:            input.HP,
		Attack:        attack,
		Defense:       defense,
		MoveSpeed:     moveSpeed,
		DropTableID:   common.NormalizeText(input.DropTableID),
		EnemySkillIDs: normalizedEnemySkillIDs,
		SpawnRule: SpawnRule{
			Origin:     input.SpawnRule.Origin,
			Distance:   input.SpawnRule.Distance,
			AxisBounds: axis,
		},
		DropTable: dropTable,
		UpdatedAt: now.UTC().Format(time.RFC3339),
	}
	return common.SaveSuccess(entry, common.SaveModeCreated)
}
