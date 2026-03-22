package httpapi

import (
	"net/http"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
)

func (a apiRouter) registerEnemyRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/enemies", func(w http.ResponseWriter, r *http.Request) {
		state, err := a.deps.EnemyRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	})
	mux.HandleFunc("POST /api/enemies", func(w http.ResponseWriter, r *http.Request) {
		var input enemies.SaveInput
		if !decodeJSON(w, r, &input) {
			return
		}
		state, err := a.deps.EnemyRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Entries, strings.TrimSpace(input.ID), func(entry enemies.EnemyEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[enemies.EnemyEntry](w)
			return
		}
		enemySkillState, err := a.deps.EnemySkillRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		itemState, err := a.deps.ItemRepo.LoadItemState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		grimoireState, err := a.deps.GrimoireRepo.LoadGrimoireState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		result := enemies.ValidateSave(input, entryIDs(enemySkillState.Entries, func(entry enemyskills.EnemySkillEntry) string { return entry.ID }), itemIDs(itemState), grimoireIDs(grimoireState), a.deps.Now())
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		nextState, mode := common.UpsertEntries(state, *result.Entry, func(entry enemies.EnemyEntry) string { return entry.ID })
		result.Mode = mode
		if err := a.deps.EnemyRepo.SaveState(nextState); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("DELETE /api/enemies/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteEntry(w, r, a.deps.EnemyRepo, "enemy", "Enemy", func(entry enemies.EnemyEntry) string { return entry.ID })
	})
}
