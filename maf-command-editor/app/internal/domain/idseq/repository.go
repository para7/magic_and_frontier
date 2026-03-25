package idseq

import (
	"errors"
	"os"

	"maf-command-editor/app/internal/domain/common"
)

type Repository struct {
	Path string
}

func (r Repository) Load() (CounterState, error) {
	var state CounterState
	err := common.ReadJSON(r.Path, &state)
	if err == nil {
		return state, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return EmptyCounterState(), nil
	}
	return CounterState{}, err
}

func (r Repository) Save(state CounterState) error {
	return common.WriteJSON(r.Path, state)
}
