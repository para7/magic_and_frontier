package api

import (
	"errors"
	"net/http"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/entity/treasures"
	dmaster "tools2/app/internal/domain/master"
)

func (a apiRouter) registerTreasureRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/treasures", func(w http.ResponseWriter, r *http.Request) {
		state, err := a.treasureState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	})
	mux.HandleFunc("POST /api/treasures", func(w http.ResponseWriter, r *http.Request) {
		var input treasures.SaveInput
		if !decodeJSON(w, r, &input) {
			return
		}
		master, err := a.masterOrErr()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		state, err := a.treasureState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Entries, strings.TrimSpace(input.ID), func(entry treasures.TreasureEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[treasures.TreasureEntry](w)
			return
		}
		result := master.Treasures().Validate(input, master)
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		if err := master.Treasures().Create(*result.Entry, master); err != nil {
			if errors.Is(err, dmaster.ErrDuplicateID) {
				writeDuplicateIDValidationError[treasures.TreasureEntry](w)
				return
			}
			writeCodedError(w, http.StatusBadRequest, "VALIDATION_FAILED", err.Error())
			return
		}
		mode := common.SaveModeCreated
		result.Mode = mode
		if err := master.Treasures().Save(); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("DELETE /api/treasures/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		if id == "" {
			writeFormError(w, http.StatusBadRequest, "Missing treasure id.")
			return
		}
		master, err := a.masterOrErr()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if err := master.Treasures().Delete(id, master); err != nil {
			writeJSON(w, http.StatusNotFound, common.DeleteNotFound("Treasure"))
			return
		}
		if err := master.Treasures().Save(); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, common.DeleteSuccess(id))
	})
}
