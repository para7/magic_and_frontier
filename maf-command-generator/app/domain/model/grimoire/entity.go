package grimoire

import (
	"errors"
	"fmt"

	model "maf_command_editor/app/domain/model"
	"maf_command_editor/app/files"

	"maf_command_editor/app/domain/custom_validator"
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

func (s *GrimoireEntity) ValidateJSON(data Grimoire, mas model.DBMaster) (Grimoire, error) {
	// validateStruct, validateRelationを順に呼び出し、データを受け入れ可能かを検証する
	if errs := s.ValidateStruct(data); len(errs) > 0 {
		return Grimoire{}, errs[0]
	}
	if errs := s.ValidateRelation(data, mas); len(errs) > 0 {
		return Grimoire{}, errs[0]
	}
	return data, nil
}

func (s *GrimoireEntity) ValidateStruct(newData Grimoire) []error {
	err := custom_validator.Validate.Struct(newData)
	if err == nil {
		return nil
	}
	var errs []error
	for _, fe := range err.(custom_validator.ValidationErrors) {
		errs = append(errs, fmt.Errorf("%s: failed on '%s' rule", fe.Field(), fe.Tag()))
	}
	return errs
}

func (s *GrimoireEntity) ValidateRelation(newData Grimoire, mas model.DBMaster) []error {
	// ほかのテーブルとのリレーションに関する内容を検証する
	// grimoire は特に検証すべき内容はない
	return nil
}

func (s *GrimoireEntity) Create(data Grimoire, mas model.DBMaster) error {
	// Validateを実行し、問題なければ data に追加する
	// Grimoire が参照する先のリレーションはないのでIDだけチェック
	if _, err := s.ValidateJSON(data, mas); err != nil {
		return err
	}
	for _, g := range s.data {
		if g.ID == data.ID {
			return errors.New("grimoire id already exists: " + data.ID)
		}
	}
	s.data = append(s.data, data)
	return nil
}

func (s *GrimoireEntity) Update(data Grimoire, mas model.DBMaster) error {
	// Validateを実行し、問題なければdataを更新する
	if _, err := s.ValidateJSON(data, mas); err != nil {
		return err
	}
	for i, g := range s.data {
		if g.ID == data.ID {
			s.data[i] = data
			return nil
		}
	}
	return errors.New("grimoire not found: " + data.ID)
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

func (s *GrimoireEntity) ValidateAll(mas model.DBMaster) []error {
	// いまの data の中身すべてに対して validate を実行する
	var errs []error
	for _, g := range s.data {
		if _, err := s.ValidateJSON(g, mas); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		fmt.Printf("[grimoire.ValidateAll] Found %d errors\n", len(errs))
	} else {
		fmt.Printf("[grimoire.ValidateAll] No errors found\n")
	}
	return errs
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
