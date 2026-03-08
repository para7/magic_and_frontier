package common

func UpsertEntries[T any](state EntryState[T], entry T, idOf func(T) string) (EntryState[T], SaveMode) {
	entryID := idOf(entry)
	for i := range state.Entries {
		if idOf(state.Entries[i]) == entryID {
			next := append([]T{}, state.Entries...)
			next[i] = entry
			return EntryState[T]{Entries: next}, SaveModeUpdated
		}
	}
	next := append(append([]T{}, state.Entries...), entry)
	return EntryState[T]{Entries: next}, SaveModeCreated
}

func DeleteEntries[T any](state EntryState[T], id string, idOf func(T) string) (EntryState[T], bool) {
	next := make([]T, 0, len(state.Entries))
	found := false
	for _, e := range state.Entries {
		if idOf(e) == id {
			found = true
			continue
		}
		next = append(next, e)
	}
	return EntryState[T]{Entries: next}, found
}
