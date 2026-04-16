package grimoire

import (
	"fmt"

	cv "maf_command_editor/app/domain/custom_validator"
	model "maf_command_editor/app/domain/model"
	"maf_command_editor/app/files"
)

// MafEntity の実装
type GrimoireEntity struct {
	store files.JsonStore[Grimoire]
	data  []Grimoire
}

func NewGrimoireEntity(path string) *GrimoireEntity {
	store := files.NewJsonStore[Grimoire](path)
	return &GrimoireEntity{store: store}
}

func (s *GrimoireEntity) ValidateJSON(newEntity Grimoire, mas model.DBMaster) (Grimoire, []model.ValidationError) {
	var errs []model.ValidationError
	errs = append(errs, s.ValidateStruct(newEntity)...)
	errs = append(errs, s.ValidateRelation(newEntity, mas)...)
	if len(errs) > 0 {
		return Grimoire{}, errs
	}
	return newEntity, nil
}

func (s *GrimoireEntity) ValidateStruct(newEntity Grimoire) []model.ValidationError {
	err := cv.Validate.Struct(newEntity)
	if err == nil {
		return nil
	}
	var errs []model.ValidationError
	for _, fe := range err.(cv.ValidationErrors) {
		errs = append(errs, cv.NewValidationError("grimoire", newEntity.ID, fe))
	}
	return errs
}

func (s *GrimoireEntity) ValidateRelation(newEntity Grimoire, _ model.DBMaster) []model.ValidationError {
	// ID の重複チェックを行う

	return nil
}

func (s *GrimoireEntity) Load() error {
	data, err := s.store.Load()
	if err != nil {
		return err
	}
	s.data = data
	fmt.Printf("[grimoire.Load] Loaded %d records\n", len(data))
	return nil
}

func (s *GrimoireEntity) ValidateAll(mas model.DBMaster) [][]model.ValidationError {
	// いまの data の中身すべてに対して validate を実行する
	var result [][]model.ValidationError
	seenIDs := map[string]bool{}
	for _, g := range s.data {
		if _, errs := s.ValidateJSON(g, mas); len(errs) > 0 {
			result = append(result, errs)
		}
		if seenIDs[g.ID] {
			result = append(result, []model.ValidationError{{
				Entity: "grimoire",
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
		fmt.Printf("[grimoire.ValidateAll] Found errors in %d record(s)\n", len(result))
	} else {
		fmt.Printf("[grimoire.ValidateAll] No errors found\n")
	}
	return result
}

func (s *GrimoireEntity) Find(id string) (Grimoire, bool) {
	// data の中から id と一致するものを返す
	for _, g := range s.data {
		if g.ID == id {
			return g, true
		}
	}
	return Grimoire{}, false
}

func (s *GrimoireEntity) GetAll() []Grimoire {
	return s.data
}
