package loottable

import (
	"errors"
	"fmt"

	cv "maf_command_editor/app/domain/custom_validator"
	model "maf_command_editor/app/domain/model"
	"maf_command_editor/app/files"
)

type LootTableEntity struct {
	store files.JsonStore[LootTable]
	data  []LootTable
}

func NewLootTableEntity(path string) *LootTableEntity {
	return &LootTableEntity{store: files.NewJsonStore[LootTable](path)}
}

func (s *LootTableEntity) ValidateJSON(newEntity LootTable, mas model.DBMaster) (LootTable, []model.ValidationError) {
	var errs []model.ValidationError
	errs = append(errs, s.ValidateStruct(newEntity)...)
	errs = append(errs, s.ValidateRelation(newEntity, mas)...)
	if len(errs) > 0 {
		return LootTable{}, errs
	}
	return newEntity, nil
}

func (s *LootTableEntity) ValidateStruct(newEntity LootTable) []model.ValidationError {
	err := cv.Validate.Struct(newEntity)
	if err == nil {
		return nil
	}
	var errs []model.ValidationError
	for _, fe := range err.(cv.ValidationErrors) {
		errs = append(errs, cv.NewValidationError("loottable", newEntity.ID, fe))
	}
	return errs
}

func (s *LootTableEntity) ValidateRelation(newEntity LootTable, mas model.DBMaster) []model.ValidationError {
	return model.ValidateDropRefs("loottable", newEntity.ID, "lootPools", newEntity.LootPools, mas)
}

func (s *LootTableEntity) Create(newEntity LootTable, mas model.DBMaster) error {
	validated, errs := s.ValidateJSON(newEntity, mas)
	if len(errs) > 0 {
		return fmt.Errorf("%s.%s: %s", errs[0].Entity, errs[0].Field, errs[0].Tag)
	}
	for _, lt := range s.data {
		if lt.ID == validated.ID {
			return errors.New("loottable id already exists: " + validated.ID)
		}
	}
	s.data = append(s.data, validated)
	return nil
}

func (s *LootTableEntity) Update(newEntity LootTable, mas model.DBMaster) error {
	validated, errs := s.ValidateJSON(newEntity, mas)
	if len(errs) > 0 {
		return fmt.Errorf("%s.%s: %s", errs[0].Entity, errs[0].Field, errs[0].Tag)
	}
	for i, lt := range s.data {
		if lt.ID == validated.ID {
			s.data[i] = validated
			return nil
		}
	}
	return errors.New("loottable not found: " + validated.ID)
}

func (s *LootTableEntity) Delete(id string, mas model.DBMaster) error {
	for i, lt := range s.data {
		if lt.ID == id {
			s.data = append(s.data[:i], s.data[i+1:]...)
			return nil
		}
	}
	return errors.New("loottable not found: " + id)
}

func (s *LootTableEntity) Save() error {
	return s.store.Save(s.data)
}

func (s *LootTableEntity) Load() error {
	data, err := s.store.Load()
	if err != nil {
		return err
	}
	s.data = data
	fmt.Printf("[loottable.Load] Loaded %d records\n", len(data))
	return nil
}

func (s *LootTableEntity) ValidateAll(mas model.DBMaster) [][]model.ValidationError {
	var result [][]model.ValidationError
	seenIDs := map[string]bool{}
	for _, lt := range s.data {
		if _, errs := s.ValidateJSON(lt, mas); len(errs) > 0 {
			result = append(result, errs)
		}
		if seenIDs[lt.ID] {
			result = append(result, []model.ValidationError{{
				Entity: "loottable",
				ID:     lt.ID,
				Field:  "id",
				Tag:    "unique",
				Param:  "ID重複を検出",
			}})
			continue
		}
		seenIDs[lt.ID] = true
	}
	if len(result) > 0 {
		fmt.Printf("[loottable.ValidateAll] Found errors in %d record(s)\n", len(result))
	} else {
		fmt.Printf("[loottable.ValidateAll] No errors found\n")
	}
	return result
}

func (s *LootTableEntity) Find(id string) (LootTable, bool) {
	for _, lt := range s.data {
		if lt.ID == id {
			return lt, true
		}
	}
	return LootTable{}, false
}

func (s *LootTableEntity) GetAll() []LootTable {
	return s.data
}
