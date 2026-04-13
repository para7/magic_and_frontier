package passive

import (
	"errors"
	"fmt"

	cv "maf_command_editor/app/domain/custom_validator"
	model "maf_command_editor/app/domain/model"
	"maf_command_editor/app/files"
)

type PassiveEntity struct {
	store files.JsonStore[Passive]
	data  []Passive
}

func NewPassiveEntity(path string) *PassiveEntity {
	return &PassiveEntity{store: files.NewJsonStore[Passive](path)}
}

func (s *PassiveEntity) ValidateJSON(newEntity Passive, mas model.DBMaster) (Passive, []model.ValidationError) {
	var errs []model.ValidationError
	errs = append(errs, s.ValidateStruct(newEntity)...)
	errs = append(errs, s.ValidateRelation(newEntity, mas)...)
	if len(errs) > 0 {
		return Passive{}, errs
	}
	return newEntity, nil
}

func (s *PassiveEntity) ValidateStruct(newEntity Passive) []model.ValidationError {
	err := cv.Validate.Struct(newEntity)
	if err == nil {
		return nil
	}
	var errs []model.ValidationError
	for _, fe := range err.(cv.ValidationErrors) {
		errs = append(errs, cv.NewValidationError("passive", newEntity.ID, fe))
	}
	return errs
}

func (s *PassiveEntity) ValidateRelation(newEntity Passive, mas model.DBMaster) []model.ValidationError {
	return nil
}

func (s *PassiveEntity) Create(newEntity Passive, mas model.DBMaster) error {
	if _, errs := s.ValidateJSON(newEntity, mas); len(errs) > 0 {
		return fmt.Errorf("%s.%s: %s", errs[0].Entity, errs[0].Field, errs[0].Tag)
	}
	for _, p := range s.data {
		if p.ID == newEntity.ID {
			return errors.New("passive id already exists: " + newEntity.ID)
		}
	}
	s.data = append(s.data, newEntity)
	return nil
}

func (s *PassiveEntity) Update(newEntity Passive, mas model.DBMaster) error {
	if _, errs := s.ValidateJSON(newEntity, mas); len(errs) > 0 {
		return fmt.Errorf("%s.%s: %s", errs[0].Entity, errs[0].Field, errs[0].Tag)
	}
	for i, p := range s.data {
		if p.ID == newEntity.ID {
			s.data[i] = newEntity
			return nil
		}
	}
	return errors.New("passive not found: " + newEntity.ID)
}

func (s *PassiveEntity) Delete(id string, mas model.DBMaster) error {
	for i, p := range s.data {
		if p.ID == id {
			s.data = append(s.data[:i], s.data[i+1:]...)
			return nil
		}
	}
	return errors.New("passive not found: " + id)
}

func (s *PassiveEntity) Save() error {
	return s.store.Save(s.data)
}

func (s *PassiveEntity) Load() error {
	data, err := s.store.Load()
	if err != nil {
		return err
	}
	s.data = data
	fmt.Printf("[passive.Load] Loaded %d records\n", len(data))
	return nil
}

func (s *PassiveEntity) ValidateAll(mas model.DBMaster) [][]model.ValidationError {
	var result [][]model.ValidationError
	seenIDs := map[string]bool{}
	for _, p := range s.data {
		if _, errs := s.ValidateJSON(p, mas); len(errs) > 0 {
			result = append(result, errs)
		}
		if seenIDs[p.ID] {
			result = append(result, []model.ValidationError{{
				Entity: "passive",
				ID:     p.ID,
				Field:  "id",
				Tag:    "unique",
				Param:  "ID重複を検出",
			}})
			continue
		}
		seenIDs[p.ID] = true
	}
	if len(result) > 0 {
		fmt.Printf("[passive.ValidateAll] Found errors in %d record(s)\n", len(result))
	} else {
		fmt.Printf("[passive.ValidateAll] No errors found\n")
	}
	return result
}

func (s *PassiveEntity) Find(id string) (Passive, bool) {
	for _, p := range s.data {
		if p.ID == id {
			return p, true
		}
	}
	return Passive{}, false
}

func (s *PassiveEntity) GetAll() []Passive {
	return s.data
}
