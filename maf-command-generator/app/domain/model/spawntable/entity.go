package spawntable

import (
	"fmt"
	"strings"

	cv "maf_command_editor/app/domain/custom_validator"
	model "maf_command_editor/app/domain/model"
	"maf_command_editor/app/files"
)

type SpawnTableEntity struct {
	store files.JsonStore[SpawnTable]
	data  []SpawnTable
}

func NewSpawnTableEntity(path string) *SpawnTableEntity {
	return &SpawnTableEntity{store: files.NewJsonStore[SpawnTable](path)}
}

func (s *SpawnTableEntity) ValidateJSON(newEntity SpawnTable, mas model.DBMaster) (SpawnTable, []model.ValidationError) {
	var errs []model.ValidationError
	errs = append(errs, s.ValidateStruct(newEntity)...)
	errs = append(errs, s.ValidateRelation(newEntity, mas)...)
	if len(errs) > 0 {
		return SpawnTable{}, errs
	}
	return newEntity, nil
}

func (s *SpawnTableEntity) ValidateStruct(newEntity SpawnTable) []model.ValidationError {
	err := cv.Validate.Struct(newEntity)
	if err == nil {
		return nil
	}
	var errs []model.ValidationError
	for _, fe := range err.(cv.ValidationErrors) {
		errs = append(errs, cv.NewValidationError("spawntable", newEntity.ID, fe))
	}
	return errs
}

func (s *SpawnTableEntity) ValidateRelation(newEntity SpawnTable, mas model.DBMaster) []model.ValidationError {
	var errs []model.ValidationError

	// SourceMobType は namespaced resource ID
	if src := strings.TrimSpace(newEntity.SourceMobType); src != "" && !model.IsNamespacedResourceID(src) {
		errs = append(errs, model.ValidationError{
			Entity: "spawntable", ID: newEntity.ID,
			Field: "sourceMobType",
			Tag:   "format", Param: "invalid minecraft entity id",
		})
	}

	// 座標の Min <= Max チェック
	if newEntity.MinX > newEntity.MaxX {
		errs = append(errs, model.ValidationError{Entity: "spawntable", ID: newEntity.ID, Field: "minX", Tag: "lte", Param: "maxX"})
	}
	if newEntity.MinY > newEntity.MaxY {
		errs = append(errs, model.ValidationError{Entity: "spawntable", ID: newEntity.ID, Field: "minY", Tag: "lte", Param: "maxY"})
	}
	if newEntity.MinZ > newEntity.MaxZ {
		errs = append(errs, model.ValidationError{Entity: "spawntable", ID: newEntity.ID, Field: "minZ", Tag: "lte", Param: "maxZ"})
	}

	// Replacements の参照チェック
	seen := map[string]bool{}
	totalWeight := newEntity.BaseMobWeight
	for i, r := range newEntity.Replacements {
		enemyID := strings.TrimSpace(r.EnemyID)
		if !mas.HasEnemy(enemyID) {
			errs = append(errs, model.ValidationError{
				Entity: "spawntable", ID: newEntity.ID,
				Field: fmt.Sprintf("replacements[%d].enemyId", i),
				Tag:   "relation", Param: "enemy not found",
			})
			continue
		}
		if seen[enemyID] {
			errs = append(errs, model.ValidationError{
				Entity: "spawntable", ID: newEntity.ID,
				Field: fmt.Sprintf("replacements[%d].enemyId", i),
				Tag:   "relation", Param: "duplicate enemy id",
			})
			continue
		}
		if r.Weight <= 0 {
			errs = append(errs, model.ValidationError{
				Entity: "spawntable", ID: newEntity.ID,
				Field: fmt.Sprintf("replacements[%d].weight", i),
				Tag:   "gt", Param: "0",
			})
			continue
		}
		seen[enemyID] = true
		totalWeight += r.Weight
	}
	if totalWeight <= 0 {
		errs = append(errs, model.ValidationError{
			Entity: "spawntable", ID: newEntity.ID,
			Field: "baseMobWeight",
			Tag:   "relation", Param: "total weight must be > 0",
		})
	}

	return errs
}

func (s *SpawnTableEntity) Load() error {
	data, err := s.store.Load()
	if err != nil {
		return err
	}
	s.data = data
	fmt.Printf("[spawntable.Load] Loaded %d records\n", len(data))
	return nil
}

func (s *SpawnTableEntity) ValidateAll(mas model.DBMaster) [][]model.ValidationError {
	var result [][]model.ValidationError
	seenIDs := map[string]bool{}
	for _, st := range s.data {
		if _, errs := s.ValidateJSON(st, mas); len(errs) > 0 {
			result = append(result, errs)
		}
		if seenIDs[st.ID] {
			result = append(result, []model.ValidationError{{
				Entity: "spawntable",
				ID:     st.ID,
				Field:  "id",
				Tag:    "unique",
				Param:  "ID重複を検出",
			}})
			continue
		}
		seenIDs[st.ID] = true
	}
	if len(result) > 0 {
		fmt.Printf("[spawntable.ValidateAll] Found errors in %d record(s)\n", len(result))
	} else {
		fmt.Printf("[spawntable.ValidateAll] No errors found\n")
	}
	return result
}

func (s *SpawnTableEntity) Find(id string) (SpawnTable, bool) {
	for _, st := range s.data {
		if st.ID == id {
			return st, true
		}
	}
	return SpawnTable{}, false
}

func (s *SpawnTableEntity) GetAll() []SpawnTable {
	return s.data
}
