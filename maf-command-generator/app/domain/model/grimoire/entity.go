package grimoire

import (
	"maf_command_editor/app/domain/master"
	"maf_command_editor/app/files"
)

type GrimoireEntity struct {
	files.JsonStore[Grimoire]
}

// // func NewStore(path string) *Store {
// // 	return &Store{JsonStore: files.NewJsonStore[Grimoire](path)}
// // }

// func (s *Store) ValidateJSON(g Grimoire, mas master.DBMaster) (Grimoire, error) {
// 	// TODO
// 	return Grimoire{}, nil
// }

// func Create()

// func (s *Store) ValidateAll(mas master.DBMaster) []error {

// 	// TODO
// 	return nil
// }

func (s *GrimoireEntity) ValidateJSON(data Grimoire, mas master.DBMaster) (Grimoire, error) {
	// TODO

	return Grimoire{}, nil
}

func (s *GrimoireEntity) Create(data Grimoire, mas master.DBMaster) error
func (s *GrimoireEntity) Update(data Grimoire, mas master.DBMaster) error
func (s *GrimoireEntity) Delete(id string, mas master.DBMaster) error
func (s *GrimoireEntity) Save() error
func (s *GrimoireEntity) Load() error

func (s *GrimoireEntity) ValidateAll(mas master.DBMaster) []error
func (s *GrimoireEntity) Find() Grimoire
