package bow

import (
	"errors"
	"fmt"

	cv "maf_command_editor/app/domain/custom_validator"
	model "maf_command_editor/app/domain/model"
	"maf_command_editor/app/files"
)

type BowEntity struct {
	store files.JsonStore[BowPassive]
	data  []BowPassive
}

func NewBowEntity(path string) *BowEntity {
	return &BowEntity{store: files.NewJsonStore[BowPassive](path)}
}

func (s *BowEntity) ValidateJSON(newEntity BowPassive, mas model.DBMaster) (BowPassive, []model.ValidationError) {
	var errs []model.ValidationError
	errs = append(errs, s.ValidateStruct(newEntity)...)
	errs = append(errs, s.ValidateRelation(newEntity, mas)...)
	if len(errs) > 0 {
		return BowPassive{}, errs
	}
	return newEntity, nil
}

func (s *BowEntity) ValidateStruct(newEntity BowPassive) []model.ValidationError {
	err := cv.Validate.Struct(newEntity)
	if err == nil {
		return nil
	}
	var errs []model.ValidationError
	for _, fe := range err.(cv.ValidationErrors) {
		errs = append(errs, cv.NewValidationError("bow", newEntity.ID, fe))
	}
	return errs
}

func (s *BowEntity) ValidateRelation(newEntity BowPassive, mas model.DBMaster) []model.ValidationError {
	if mas.HasPassive("bow_" + newEntity.ID) {
		return []model.ValidationError{{
			Entity: "bow",
			ID:     newEntity.ID,
			Field:  "id",
			Tag:    "relation",
			Param:  "conflicts with passive id bow_" + newEntity.ID + " in generated/passive/effect",
		}}
	}
	return nil
}

func (s *BowEntity) Create(newEntity BowPassive, mas model.DBMaster) error {
	if _, errs := s.ValidateJSON(newEntity, mas); len(errs) > 0 {
		return fmt.Errorf("%s.%s: %s", errs[0].Entity, errs[0].Field, errs[0].Tag)
	}
	for _, p := range s.data {
		if p.ID == newEntity.ID {
			return errors.New("bow id already exists: " + newEntity.ID)
		}
	}
	s.data = append(s.data, newEntity)
	return nil
}

func (s *BowEntity) Update(newEntity BowPassive, mas model.DBMaster) error {
	if _, errs := s.ValidateJSON(newEntity, mas); len(errs) > 0 {
		return fmt.Errorf("%s.%s: %s", errs[0].Entity, errs[0].Field, errs[0].Tag)
	}
	for i, p := range s.data {
		if p.ID == newEntity.ID {
			s.data[i] = newEntity
			return nil
		}
	}
	return errors.New("bow not found: " + newEntity.ID)
}

func (s *BowEntity) Delete(id string, mas model.DBMaster) error {
	for i, p := range s.data {
		if p.ID == id {
			s.data = append(s.data[:i], s.data[i+1:]...)
			return nil
		}
	}
	return errors.New("bow not found: " + id)
}

func (s *BowEntity) Save() error {
	return s.store.Save(s.data)
}

func (s *BowEntity) Load() error {
	data, err := s.store.Load()
	if err != nil {
		return err
	}
	s.data = data
	fmt.Printf("[bow.Load] Loaded %d records\n", len(data))
	return nil
}

func (s *BowEntity) ValidateAll(mas model.DBMaster) [][]model.ValidationError {
	var result [][]model.ValidationError
	seenIDs := map[string]bool{}
	for _, p := range s.data {
		if _, errs := s.ValidateJSON(p, mas); len(errs) > 0 {
			result = append(result, errs)
		}
		if seenIDs[p.ID] {
			result = append(result, []model.ValidationError{{
				Entity: "bow",
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
		fmt.Printf("[bow.ValidateAll] Found errors in %d record(s)\n", len(result))
	} else {
		fmt.Printf("[bow.ValidateAll] No errors found\n")
	}
	return result
}

func (s *BowEntity) Find(id string) (BowPassive, bool) {
	for _, p := range s.data {
		if p.ID == id {
			return p, true
		}
	}
	return BowPassive{}, false
}

func (s *BowEntity) GetAll() []BowPassive {
	return s.data
}
