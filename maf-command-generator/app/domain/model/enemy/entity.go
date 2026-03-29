package enemy

import (
	"errors"
	"fmt"
	"strings"

	cv "maf_command_editor/app/domain/custom_validator"
	model "maf_command_editor/app/domain/model"
	"maf_command_editor/app/files"
)

type EnemyEntity struct {
	store files.JsonStore[Enemy]
	data  []Enemy
}

func NewEnemyEntity(path string) *EnemyEntity {
	return &EnemyEntity{store: files.NewJsonStore[Enemy](path)}
}

func (s *EnemyEntity) ValidateJSON(newEntity Enemy, mas model.DBMaster) (Enemy, []model.ValidationError) {
	var errs []model.ValidationError
	errs = append(errs, s.ValidateStruct(newEntity)...)
	errs = append(errs, s.ValidateRelation(newEntity, mas)...)
	if len(errs) > 0 {
		return Enemy{}, errs
	}
	return newEntity, nil
}

func (s *EnemyEntity) ValidateStruct(newEntity Enemy) []model.ValidationError {
	err := cv.Validate.Struct(newEntity)
	if err == nil {
		return nil
	}
	var errs []model.ValidationError
	for _, fe := range err.(cv.ValidationErrors) {
		errs = append(errs, cv.NewValidationError("enemy", newEntity.ID, fe))
	}
	return errs
}

func (s *EnemyEntity) ValidateRelation(newEntity Enemy, mas model.DBMaster) []model.ValidationError {
	var errs []model.ValidationError

	// EnemySkillIDs の参照チェック
	seen := map[string]bool{}
	for i, skillID := range newEntity.EnemySkillIDs {
		id := strings.TrimSpace(skillID)
		if id == "" {
			continue
		}
		if !mas.HasEnemySkill(id) {
			errs = append(errs, model.ValidationError{
				Entity: "enemy", ID: newEntity.ID,
				Field: fmt.Sprintf("enemySkillIds[%d]", i),
				Tag:   "relation", Param: "enemyskill not found",
			})
		} else if seen[id] {
			errs = append(errs, model.ValidationError{
				Entity: "enemy", ID: newEntity.ID,
				Field: fmt.Sprintf("enemySkillIds[%d]", i),
				Tag:   "relation", Param: "duplicate enemyskill id",
			})
		}
		seen[id] = true
	}

	// Drops の参照チェック
	errs = append(errs, model.ValidateDropRefs("enemy", newEntity.ID, "drops", newEntity.Drops, mas)...)

	// Equipment スロットの参照チェック
	errs = append(errs, model.ValidateEquipmentSlots("enemy", newEntity.ID, newEntity.Equipment, mas)...)

	return errs
}

func (s *EnemyEntity) Create(newEntity Enemy, mas model.DBMaster) error {
	validated, errs := s.ValidateJSON(newEntity, mas)
	if len(errs) > 0 {
		return fmt.Errorf("%s.%s: %s", errs[0].Entity, errs[0].Field, errs[0].Tag)
	}
	for _, e := range s.data {
		if e.ID == validated.ID {
			return errors.New("enemy id already exists: " + validated.ID)
		}
	}
	s.data = append(s.data, validated)
	return nil
}

func (s *EnemyEntity) Update(newEntity Enemy, mas model.DBMaster) error {
	validated, errs := s.ValidateJSON(newEntity, mas)
	if len(errs) > 0 {
		return fmt.Errorf("%s.%s: %s", errs[0].Entity, errs[0].Field, errs[0].Tag)
	}
	for i, e := range s.data {
		if e.ID == validated.ID {
			s.data[i] = validated
			return nil
		}
	}
	return errors.New("enemy not found: " + validated.ID)
}

func (s *EnemyEntity) Delete(id string, mas model.DBMaster) error {
	for i, e := range s.data {
		if e.ID == id {
			s.data = append(s.data[:i], s.data[i+1:]...)
			return nil
		}
	}
	return errors.New("enemy not found: " + id)
}

func (s *EnemyEntity) Save() error {
	return s.store.Save(s.data)
}

func (s *EnemyEntity) Load() error {
	data, err := s.store.Load()
	if err != nil {
		return err
	}
	s.data = data
	fmt.Printf("[enemy.Load] Loaded %d records\n", len(data))
	return nil
}

func (s *EnemyEntity) ValidateAll(mas model.DBMaster) [][]model.ValidationError {
	var result [][]model.ValidationError
	seenIDs := map[string]bool{}
	for _, e := range s.data {
		if _, errs := s.ValidateJSON(e, mas); len(errs) > 0 {
			result = append(result, errs)
		}
		if seenIDs[e.ID] {
			result = append(result, []model.ValidationError{{
				Entity: "enemy",
				ID:     e.ID,
				Field:  "id",
				Tag:    "unique",
				Param:  "ID重複を検出",
			}})
			continue
		}
		seenIDs[e.ID] = true
	}
	if len(result) > 0 {
		fmt.Printf("[enemy.ValidateAll] Found errors in %d record(s)\n", len(result))
	} else {
		fmt.Printf("[enemy.ValidateAll] No errors found\n")
	}
	return result
}

func (s *EnemyEntity) Find(id string) (Enemy, bool) {
	for _, e := range s.data {
		if e.ID == id {
			return e, true
		}
	}
	return Enemy{}, false
}

func (s *EnemyEntity) GetAll() []Enemy {
	return s.data
}
