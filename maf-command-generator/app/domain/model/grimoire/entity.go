package grimoire

import (
	"errors"
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

func (s *GrimoireEntity) Create(newEntity Grimoire, mas model.DBMaster) error {
	// Validateを実行し、問題なければ data に追加する
	if _, errs := s.ValidateJSON(newEntity, mas); len(errs) > 0 {
		return fmt.Errorf("%s.%s: %s", errs[0].Entity, errs[0].Field, errs[0].Tag)
	}
	for _, g := range s.data {
		if g.ID == newEntity.ID {
			return errors.New("grimoire id already exists: " + newEntity.ID)
		}
	}
	s.data = append(s.data, newEntity)
	return nil
}

func (s *GrimoireEntity) Update(newEntity Grimoire, mas model.DBMaster) error {
	// Validateを実行し、問題なければdataを更新する
	if _, errs := s.ValidateJSON(newEntity, mas); len(errs) > 0 {
		return fmt.Errorf("%s.%s: %s", errs[0].Entity, errs[0].Field, errs[0].Tag)
	}
	for i, g := range s.data {
		if g.ID == newEntity.ID {
			s.data[i] = newEntity
			return nil
		}
	}
	return errors.New("grimoire not found: " + newEntity.ID)
}

func (s *GrimoireEntity) Delete(id string, mas model.DBMaster) error {
	// ほかのデータから参照されていなければ配列から削除する
	// Grimoireは各種ルートテーブル系から参照されている可能性あり。
	// ほかのルートテーブルはまだ実装してないのでTODOとしておく
	for i, g := range s.data {
		if g.ID == id {
			s.data = append(s.data[:i], s.data[i+1:]...)
			return nil
		}
	}
	return errors.New("grimoire not found: " + id)
}

func (s *GrimoireEntity) Save() error {
	return s.store.Save(s.data)
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
