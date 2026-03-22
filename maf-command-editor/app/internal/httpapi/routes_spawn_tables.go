package httpapi

import (
	"net/http"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/spawntables"
)

func (a apiRouter) registerSpawnTableRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/spawn-tables", func(w http.ResponseWriter, r *http.Request) {
		state, err := a.deps.SpawnTableRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	})
	mux.HandleFunc("POST /api/spawn-tables", func(w http.ResponseWriter, r *http.Request) {
		var input spawntables.SaveInput
		if !decodeJSON(w, r, &input) {
			return
		}
		state, err := a.deps.SpawnTableRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Entries, strings.TrimSpace(input.ID), func(entry spawntables.SpawnTableEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[spawntables.SpawnTableEntry](w)
			return
		}
		enemyState, err := a.deps.EnemyRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		result := spawntables.ValidateSave(input, entryIDs(enemyState.Entries, func(entry enemies.EnemyEntry) string { return entry.ID }), a.deps.Now())
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		if conflictID, ok := spawntables.FirstOverlap(state.Entries, *result.Entry); ok {
			writeJSON(w, http.StatusBadRequest, common.SaveValidationError[spawntables.SpawnTableEntry](
				common.FieldErrors{"range": "Range overlaps with " + conflictID + "."},
				"Validation failed. Fix the highlighted fields.",
			))
			return
		}
		nextState, mode := common.UpsertEntries(state, *result.Entry, func(entry spawntables.SpawnTableEntry) string { return entry.ID })
		result.Mode = mode
		if err := a.deps.SpawnTableRepo.SaveState(nextState); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("DELETE /api/spawn-tables/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteEntry(w, r, a.deps.SpawnTableRepo, "spawn table", "Spawn table", func(entry spawntables.SpawnTableEntry) string { return entry.ID })
	})
}
