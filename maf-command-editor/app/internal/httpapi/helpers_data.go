package httpapi

import (
	"net/http"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/treasures"
	"tools2/app/internal/mcsource"
	"tools2/app/internal/store"
)

func deleteEntry[T any](w http.ResponseWriter, r *http.Request, repo store.EntryStateRepository[T], missingLabel, notFoundLabel string, idOf func(T) string) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFormError(w, http.StatusBadRequest, "Missing "+missingLabel+" id.")
		return
	}
	state, err := repo.LoadState()
	if err != nil {
		writeInternalError(w, err)
		return
	}
	nextState, ok := common.DeleteEntries(state, id, idOf)
	if !ok {
		writeJSON(w, http.StatusNotFound, common.DeleteNotFound(notFoundLabel))
		return
	}
	if err := repo.SaveState(nextState); err != nil {
		writeInternalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, common.DeleteSuccess(id))
}

func duplicateCastID(entries []grimoire.GrimoireEntry, entryID string, castID int) string {
	for _, entry := range entries {
		if entry.ID != entryID && entry.CastID == castID {
			return entry.ID
		}
	}
	return ""
}

func duplicateTreasureTablePath(entries []treasures.TreasureEntry, entryID, tablePath string) string {
	for _, entry := range entries {
		if entry.ID != entryID && entry.TablePath == tablePath {
			return entry.ID
		}
	}
	return ""
}

func treasureSourcePaths(root string) (map[string]struct{}, error) {
	sources, err := mcsource.ListLootTables(root)
	if err != nil {
		return nil, err
	}
	out := make(map[string]struct{}, len(sources))
	for _, source := range sources {
		out[source.TablePath] = struct{}{}
	}
	return out, nil
}
