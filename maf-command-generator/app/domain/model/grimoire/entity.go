package grimoire

type Store struct {
	Grimoires []Grimoire
}

func (s *Store) Validate() (Grimoire, error) {
	// TODO
	return Grimoire{}, nil
}

func Load() {

}

func (s *Store) ValidateAll() []error {
	// TODO
	return nil
}
