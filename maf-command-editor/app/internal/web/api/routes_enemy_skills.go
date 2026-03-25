package api

import (
	"errors"
	"net/http"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/entity/enemyskills"
	dmaster "tools2/app/internal/domain/master"
)

func (a apiRouter) registerEnemySkillRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/enemy-skills", func(w http.ResponseWriter, r *http.Request) {
		state, err := a.enemySkillState()
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
		master, err := a.masterOrErr()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		state, err := a.enemySkillState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Entries, strings.TrimSpace(input.ID), func(entry enemyskills.EnemySkillEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[enemyskills.EnemySkillEntry](w)
			return
		}
		result := master.EnemySkills().Validate(input, master)
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		if err := master.EnemySkills().Create(*result.Entry, master); err != nil {
			if errors.Is(err, dmaster.ErrDuplicateID) {
				writeDuplicateIDValidationError[enemyskills.EnemySkillEntry](w)
				return
			}
			writeCodedError(w, http.StatusBadRequest, "VALIDATION_FAILED", err.Error())
			return
		}
		mode := common.SaveModeCreated
		result.Mode = mode
		if err := master.EnemySkills().Save(); err != nil {
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
		enemyState, err := a.enemyState()
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
		master, err := a.masterOrErr()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if err := master.EnemySkills().Delete(id, master); err != nil {
			writeJSON(w, http.StatusNotFound, common.DeleteNotFound("Enemy skill"))
			return
		}
		if err := master.EnemySkills().Save(); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, common.DeleteSuccess(id))
	})
}
