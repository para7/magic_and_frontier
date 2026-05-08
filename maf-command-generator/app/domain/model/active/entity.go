package active

import (
	"fmt"

	cv "maf_command_editor/app/domain/custom_validator"
	model "maf_command_editor/app/domain/model"
	"maf_command_editor/app/files"
)

// MafEntity の実装
type ActiveEntity struct {
	store files.JsonStore[Active]
	data  []Active
}

func NewActiveEntity(path string) *ActiveEntity {
	store := files.NewJsonStore[Active](path)
	return &ActiveEntity{store: store}
}

func (s *ActiveEntity) ValidateJSON(newEntity Active, mas model.DBMaster) (Active, []model.ValidationError) {
	var errs []model.ValidationError
	errs = append(errs, s.ValidateStruct(newEntity)...)
	errs = append(errs, s.ValidateRelation(newEntity, mas)...)
	if len(errs) > 0 {
		return Active{}, errs
	}
	return newEntity, nil
}

func (s *ActiveEntity) ValidateStruct(newEntity Active) []model.ValidationError {
	err := cv.Validate.Struct(newEntity)
	if err == nil {
		return nil
	}
	var errs []model.ValidationError
	for _, fe := range err.(cv.ValidationErrors) {
		errs = append(errs, cv.NewValidationError("active", newEntity.ID, fe))
	}
	return errs
}

func (s *ActiveEntity) ValidateRelation(newEntity Active, _ model.DBMaster) []model.ValidationError {
	// ID の重複チェックを行う

	return nil
}

func (s *ActiveEntity) Load() error {
	data, err := s.store.Load()
	if err != nil {
		return err
	}
	s.data = data
	fmt.Printf("[active.Load] Loaded %d records\n", len(data))
	return nil
}

func (s *ActiveEntity) ValidateAll(mas model.DBMaster) [][]model.ValidationError {
	// いまの data の中身すべてに対して validate を実行する
	var result [][]model.ValidationError
	seenIDs := map[string]bool{}
	for _, g := range s.data {
		if _, errs := s.ValidateJSON(g, mas); len(errs) > 0 {
			result = append(result, errs)
		}
		if seenIDs[g.ID] {
			result = append(result, []model.ValidationError{{
				Entity: "active",
				ID:     g.ID,
				Field:  "id",
				Tag:    "unique",
				Param:  "ID重複を検出",
			}})
			continue
		}
		seenIDs[g.ID] = true
	}

	if len(result) > 0 {
		fmt.Printf("[active.ValidateAll] Found errors in %d record(s)\n", len(result))
	} else {
		fmt.Printf("[active.ValidateAll] No errors found\n")
	}
	return result
}

func (s *ActiveEntity) Find(id string) (Active, bool) {
	// data の中から id と一致するものを返す
	for _, g := range s.data {
		if g.ID == id {
			return g, true
		}
	}
	return Active{}, false
}

func (s *ActiveEntity) GetAll() []Active {
	return s.data
}
