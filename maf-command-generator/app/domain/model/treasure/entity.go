package treasure

import (
	"errors"
	"fmt"
	"strings"

	cv "maf_command_editor/app/domain/custom_validator"
	model "maf_command_editor/app/domain/model"
	"maf_command_editor/app/files"
)

type TreasureEntity struct {
	store files.JsonStore[Treasure]
	data  []Treasure
}

func NewTreasureEntity(path string) *TreasureEntity {
	return &TreasureEntity{store: files.NewJsonStore[Treasure](path)}
}

func (s *TreasureEntity) ValidateJSON(newEntity Treasure, mas model.DBMaster) (Treasure, []model.ValidationError) {
	var errs []model.ValidationError
	errs = append(errs, s.ValidateStruct(newEntity)...)
	errs = append(errs, s.ValidateRelation(newEntity, mas)...)
	if len(errs) > 0 {
		return Treasure{}, errs
	}
	return newEntity, nil
}

func (s *TreasureEntity) ValidateStruct(newEntity Treasure) []model.ValidationError {
	err := cv.Validate.Struct(newEntity)
	if err == nil {
		return nil
	}
	var errs []model.ValidationError
	for _, fe := range err.(cv.ValidationErrors) {
		errs = append(errs, cv.NewValidationError("treasure", newEntity.ID, fe))
	}
	return errs
}

func (s *TreasureEntity) ValidateRelation(newEntity Treasure, mas model.DBMaster) []model.ValidationError {
	var errs []model.ValidationError

	tablePath := strings.TrimSpace(newEntity.TablePath)
	if tablePath != "" {
		if !model.IsSafeNamespacedResourcePath(tablePath) {
			errs = append(errs, model.ValidationError{
				Entity: "treasure", ID: newEntity.ID,
				Field: "tablePath",
				Tag:   "format", Param: "invalid namespaced loot table path",
			})
		} else if !strings.HasPrefix(tablePath, "minecraft:chests/") {
			errs = append(errs, model.ValidationError{
				Entity: "treasure", ID: newEntity.ID,
				Field: "tablePath",
				Tag:   "format", Param: "treasure tablePath must be under minecraft:chests/",
			})
		} else if !mas.HasMinecraftLootTable(tablePath) {
			errs = append(errs, model.ValidationError{
				Entity: "treasure", ID: newEntity.ID,
				Field: "tablePath",
				Tag:   "relation", Param: "minecraft loot table not found",
			})
		}
	}

	errs = append(errs, model.ValidateDropRefs("treasure", newEntity.ID, "lootPools", newEntity.LootPools, mas)...)

	return errs
}

func (s *TreasureEntity) Create(newEntity Treasure, mas model.DBMaster) error {
	validated, errs := s.ValidateJSON(newEntity, mas)
	if len(errs) > 0 {
		return fmt.Errorf("%s.%s: %s", errs[0].Entity, errs[0].Field, errs[0].Tag)
	}
	for _, t := range s.data {
		if t.ID == validated.ID {
			return errors.New("treasure id already exists: " + validated.ID)
		}
	}
	s.data = append(s.data, validated)
	return nil
}

func (s *TreasureEntity) Update(newEntity Treasure, mas model.DBMaster) error {
	validated, errs := s.ValidateJSON(newEntity, mas)
	if len(errs) > 0 {
		return fmt.Errorf("%s.%s: %s", errs[0].Entity, errs[0].Field, errs[0].Tag)
	}
	for i, t := range s.data {
		if t.ID == validated.ID {
			s.data[i] = validated
			return nil
		}
	}
	return errors.New("treasure not found: " + validated.ID)
}

func (s *TreasureEntity) Delete(id string, mas model.DBMaster) error {
	for i, t := range s.data {
		if t.ID == id {
			s.data = append(s.data[:i], s.data[i+1:]...)
			return nil
		}
	}
	return errors.New("treasure not found: " + id)
}

func (s *TreasureEntity) Save() error {
	return s.store.Save(s.data)
}

func (s *TreasureEntity) Load() error {
	data, err := s.store.Load()
	if err != nil {
		return err
	}
	s.data = data
	return nil
}

func (s *TreasureEntity) ValidateAll(mas model.DBMaster) [][]model.ValidationError {
	var result [][]model.ValidationError
	seenIDs := map[string]bool{}
	for _, t := range s.data {
		if _, errs := s.ValidateJSON(t, mas); len(errs) > 0 {
			result = append(result, errs)
		}
		if seenIDs[t.ID] {
			result = append(result, []model.ValidationError{{
				Entity: "treasure",
				ID:     t.ID,
				Field:  "id",
				Tag:    "unique",
				Param:  "ID重複を検出",
			}})
			continue
		}
		seenIDs[t.ID] = true
	}
	return result
}

func (s *TreasureEntity) Find(id string) (Treasure, bool) {
	for _, t := range s.data {
		if t.ID == id {
			return t, true
		}
	}
	return Treasure{}, false
}

func (s *TreasureEntity) GetAll() []Treasure {
	return s.data
}
