package grimoire

import (
	"errors"
	"strings"

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
	// リレーションとは関係ない、値の範囲、文字数などデータ単品で検証できる内容を検証する
	var errs []error
	if strings.TrimSpace(newData.ID) == "" {
		errs = append(errs, errors.New("id is required"))
	}
	if newData.CastID < 1 {
		errs = append(errs, errors.New("castid must be >= 1"))
	}
	if newData.CastTime < 0 || newData.CastTime > 12000 {
		errs = append(errs, errors.New("castTime must be between 0 and 12000"))
	}
	if newData.MPCost < 0 || newData.MPCost > 1000000 {
		errs = append(errs, errors.New("mpCost must be between 0 and 1000000"))
	}
	if strings.TrimSpace(newData.Script) == "" {
		errs = append(errs, errors.New("script is required"))
	} else if len([]rune(strings.TrimSpace(newData.Script))) > 20000 {
		errs = append(errs, errors.New("script must be <= 20000 characters"))
	}
	if strings.TrimSpace(newData.Title) == "" {
		errs = append(errs, errors.New("title is required"))
	} else if len([]rune(strings.TrimSpace(newData.Title))) > 200 {
		errs = append(errs, errors.New("title must be <= 200 characters"))
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
