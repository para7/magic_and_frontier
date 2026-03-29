package enemyskill

import (
	"errors"
	"fmt"

	cv "maf_command_editor/app/domain/custom_validator"
	model "maf_command_editor/app/domain/model"
	"maf_command_editor/app/files"
)

type EnemySkillEntity struct {
	store files.JsonStore[EnemySkill]
	data  []EnemySkill
}

func NewEnemySkillEntity(path string) *EnemySkillEntity {
	return &EnemySkillEntity{store: files.NewJsonStore[EnemySkill](path)}
}

func (s *EnemySkillEntity) ValidateJSON(newEntity EnemySkill, mas model.DBMaster) (EnemySkill, []model.ValidationError) {
	var errs []model.ValidationError
	errs = append(errs, s.ValidateStruct(newEntity)...)
	errs = append(errs, s.ValidateRelation(newEntity, mas)...)
	if len(errs) > 0 {
		return EnemySkill{}, errs
	}
	return newEntity, nil
}

func (s *EnemySkillEntity) ValidateStruct(newEntity EnemySkill) []model.ValidationError {
	err := cv.Validate.Struct(newEntity)
	if err == nil {
		return nil
	}
	var errs []model.ValidationError
	for _, fe := range err.(cv.ValidationErrors) {
		errs = append(errs, cv.NewValidationError("enemyskill", newEntity.ID, fe))
	}
	return errs
}

func (s *EnemySkillEntity) ValidateRelation(_ EnemySkill, _ model.DBMaster) []model.ValidationError {
	return nil
}

func (s *EnemySkillEntity) Create(newEntity EnemySkill, mas model.DBMaster) error {
	if _, errs := s.ValidateJSON(newEntity, mas); len(errs) > 0 {
		return fmt.Errorf("%s.%s: %s", errs[0].Entity, errs[0].Field, errs[0].Tag)
	}
	for _, e := range s.data {
		if e.ID == newEntity.ID {
			return errors.New("enemyskill id already exists: " + newEntity.ID)
		}
	}
	s.data = append(s.data, newEntity)
	return nil
}

func (s *EnemySkillEntity) Update(newEntity EnemySkill, mas model.DBMaster) error {
	if _, errs := s.ValidateJSON(newEntity, mas); len(errs) > 0 {
		return fmt.Errorf("%s.%s: %s", errs[0].Entity, errs[0].Field, errs[0].Tag)
	}
	for i, e := range s.data {
		if e.ID == newEntity.ID {
			s.data[i] = newEntity
			return nil
		}
	}
	return errors.New("enemyskill not found: " + newEntity.ID)
}

func (s *EnemySkillEntity) Delete(id string, mas model.DBMaster) error {
	for i, e := range s.data {
		if e.ID == id {
			s.data = append(s.data[:i], s.data[i+1:]...)
			return nil
		}
	}
	return errors.New("enemyskill not found: " + id)
}

func (s *EnemySkillEntity) Save() error {
	return s.store.Save(s.data)
}

func (s *EnemySkillEntity) Load() error {
	data, err := s.store.Load()
	if err != nil {
		return err
	}
	s.data = data
	return nil
}

func (s *EnemySkillEntity) ValidateAll(mas model.DBMaster) [][]model.ValidationError {
	var result [][]model.ValidationError
	seenIDs := map[string]bool{}
	for _, e := range s.data {
		if _, errs := s.ValidateJSON(e, mas); len(errs) > 0 {
			result = append(result, errs)
		}
		if seenIDs[e.ID] {
			result = append(result, []model.ValidationError{{
				Entity: "enemyskill",
				ID:     e.ID,
				Field:  "id",
				Tag:    "unique",
				Param:  "ID重複を検出",
			}})
			continue
		}
		seenIDs[e.ID] = true
	}
	return result
}

func (s *EnemySkillEntity) Find(id string) (EnemySkill, bool) {
	for _, e := range s.data {
		if e.ID == id {
			return e, true
		}
	}
	return EnemySkill{}, false
}

func (s *EnemySkillEntity) GetAll() []EnemySkill {
	return s.data
}
