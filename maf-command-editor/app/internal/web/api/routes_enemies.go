package api

import (
	"errors"
	"net/http"
	"strings"

	"maf-command-editor/app/internal/domain/common"
	"maf-command-editor/app/internal/domain/entity/enemies"
	dmaster "maf-command-editor/app/internal/domain/master"
)

func (a apiRouter) registerEnemyRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/enemies", func(w http.ResponseWriter, r *http.Request) {
		state, err := a.enemyState()
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
		master, err := a.masterOrErr()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		state, err := a.enemyState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Entries, strings.TrimSpace(input.ID), func(entry enemies.EnemyEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[enemies.EnemyEntry](w)
			return
		}
		result := master.Enemies().Validate(input, master)
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		if err := master.Enemies().Create(*result.Entry, master); err != nil {
			if errors.Is(err, dmaster.ErrDuplicateID) {
				writeDuplicateIDValidationError[enemies.EnemyEntry](w)
				return
			}
			writeCodedError(w, http.StatusBadRequest, "VALIDATION_FAILED", err.Error())
			return
		}
		mode := common.SaveModeCreated
		result.Mode = mode
		if err := master.Enemies().Save(); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("DELETE /api/enemies/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		if id == "" {
			writeFormError(w, http.StatusBadRequest, "Missing enemy id.")
			return
		}
		master, err := a.masterOrErr()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if err := master.Enemies().Delete(id, master); err != nil {
			writeJSON(w, http.StatusNotFound, common.DeleteNotFound("Enemy"))
			return
		}
		if err := master.Enemies().Save(); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, common.DeleteSuccess(id))
	})
}
