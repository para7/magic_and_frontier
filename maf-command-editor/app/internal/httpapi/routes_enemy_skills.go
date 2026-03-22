package httpapi

import (
	"net/http"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemyskills"
)

func (a apiRouter) registerEnemySkillRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/enemy-skills", func(w http.ResponseWriter, r *http.Request) {
		state, err := a.deps.EnemySkillRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	})
	mux.HandleFunc("POST /api/enemy-skills", func(w http.ResponseWriter, r *http.Request) {
		var input enemyskills.SaveInput
		if !decodeJSON(w, r, &input) {
			return
		}
		state, err := a.deps.EnemySkillRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Entries, strings.TrimSpace(input.ID), func(entry enemyskills.EnemySkillEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[enemyskills.EnemySkillEntry](w)
			return
		}
		result := enemyskills.ValidateSave(input, a.deps.Now())
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		nextState, mode := common.UpsertEntries(state, *result.Entry, func(entry enemyskills.EnemySkillEntry) string { return entry.ID })
		result.Mode = mode
		if err := a.deps.EnemySkillRepo.SaveState(nextState); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("DELETE /api/enemy-skills/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		if id == "" {
			writeFormError(w, http.StatusBadRequest, "Missing enemy skill id.")
			return
		}
		enemyState, err := a.deps.EnemyRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		for _, enemy := range enemyState.Entries {
			for _, skillID := range enemy.EnemySkillIDs {
				if skillID == id {
					writeCodedError(w, http.StatusBadRequest, "REFERENCE_ERROR", "Enemy skill is referenced by enemy "+enemy.ID+".")
					return
				}
			}
		}
		state, err := a.deps.EnemySkillRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		nextState, ok := common.DeleteEntries(state, id, func(entry enemyskills.EnemySkillEntry) string { return entry.ID })
		if !ok {
			writeJSON(w, http.StatusNotFound, common.DeleteNotFound("Enemy skill"))
			return
		}
		if err := a.deps.EnemySkillRepo.SaveState(nextState); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, common.DeleteSuccess(id))
	})
}
