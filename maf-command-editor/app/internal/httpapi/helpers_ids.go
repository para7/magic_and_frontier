package httpapi

import (
	"strings"

	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
)

func grimoireIDs(state grimoire.GrimoireState) map[string]struct{} {
	return entryIDs(state.Entries, func(entry grimoire.GrimoireEntry) string { return entry.ID })
}

func entryIDs[T any](entries []T, idOf func(T) string) map[string]struct{} {
	out := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		id := strings.TrimSpace(idOf(entry))
		if id != "" {
			out[id] = struct{}{}
		}
	}
	return out
}

func findEntry[T any](entries []T, id string, idOf func(T) string) (T, bool) {
	var zero T
	id = strings.TrimSpace(id)
	if id == "" {
		return zero, false
	}
	for _, entry := range entries {
		if strings.TrimSpace(idOf(entry)) == id {
			return entry, true
		}
	}
	return zero, false
}

func itemIDs(state items.ItemState) map[string]struct{} {
	return entryIDs(state.Items, func(entry items.ItemEntry) string { return entry.ID })
}
