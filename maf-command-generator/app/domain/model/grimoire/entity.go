package grimoire

import (
	"encoding/json"
	"maf_command_editor/app/files"
	"os"
)

type Store struct {
	Path      string
	Grimoires []Grimoire
}

func NewStore(path string) *Store {
	return &Store{Path: path}
}

func (s *Store) Load() error {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		return err
	}
	var f grimoireJsonFile
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	s.Grimoires = f.Entries
	return nil
}

func (s *Store) Save() error {
	entries := s.Grimoires
	if entries == nil {
		entries = []Grimoire{}
	}
	return files.WriteJson(s.Path, grimoireJsonFile{Entries: entries})
}

func (s *Store) Validate() (Grimoire, error) {
	// TODO
	return Grimoire{}, nil
}

func (s *Store) ValidateAll() []error {
	// TODO
	return nil
}
