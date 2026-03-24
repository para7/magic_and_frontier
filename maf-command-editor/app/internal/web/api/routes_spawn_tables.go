package api

import (
	"errors"
	"net/http"
	"strings"

	"tools2/app/internal/domain/common"
	dmaster "tools2/app/internal/domain/master"
	"tools2/app/internal/domain/spawntables"
)

func (a apiRouter) registerSpawnTableRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/spawn-tables", func(w http.ResponseWriter, r *http.Request) {
		state, err := a.spawnTableState()
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
		master, err := a.masterOrErr()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		state, err := a.spawnTableState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Entries, strings.TrimSpace(input.ID), func(entry spawntables.SpawnTableEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[spawntables.SpawnTableEntry](w)
			return
		}
		result := master.SpawnTables().Validate(input, master)
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		if err := master.SpawnTables().Create(*result.Entry, master); err != nil {
			if errors.Is(err, dmaster.ErrDuplicateID) {
				writeDuplicateIDValidationError[spawntables.SpawnTableEntry](w)
				return
			}
			writeCodedError(w, http.StatusBadRequest, "VALIDATION_FAILED", err.Error())
			return
		}
		mode := common.SaveModeCreated
		result.Mode = mode
		if err := master.SpawnTables().Save(); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("DELETE /api/spawn-tables/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		if id == "" {
			writeFormError(w, http.StatusBadRequest, "Missing spawn table id.")
			return
		}
		master, err := a.masterOrErr()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if err := master.SpawnTables().Delete(id, master); err != nil {
			writeJSON(w, http.StatusNotFound, common.DeleteNotFound("Spawn table"))
			return
		}
		if err := master.SpawnTables().Save(); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, common.DeleteSuccess(id))
	})
}
