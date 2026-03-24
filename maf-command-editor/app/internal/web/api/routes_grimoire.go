package api

import (
	"errors"
	"net/http"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/grimoire"
	dmaster "tools2/app/internal/domain/master"
)

func (a apiRouter) registerGrimoireRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/grimoire", func(w http.ResponseWriter, r *http.Request) {
		state, err := a.grimoireState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	})
	mux.HandleFunc("POST /api/grimoire", func(w http.ResponseWriter, r *http.Request) {
		var input grimoire.SaveInput
		if !decodeJSON(w, r, &input) {
			return
		}
		master, err := a.masterOrErr()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		state, err := a.grimoireState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Entries, strings.TrimSpace(input.ID), func(entry grimoire.GrimoireEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[grimoire.GrimoireEntry](w)
			return
		}
		castID, allocErr := a.service.AllocateCastID()
		if allocErr != nil {
			writeInternalError(w, allocErr)
			return
		}
		input.CastID = castID
		result := master.Grimoires().Validate(input, master)
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		if err := master.Grimoires().Create(*result.Entry, master); err != nil {
			if errors.Is(err, dmaster.ErrDuplicateID) {
				writeDuplicateIDValidationError[grimoire.GrimoireEntry](w)
				return
			}
			writeCodedError(w, http.StatusBadRequest, "VALIDATION_FAILED", err.Error())
			return
		}
		mode := common.SaveModeCreated
		result.Mode = mode
		if err := master.Grimoires().Save(); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("DELETE /api/grimoire/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		if id == "" {
			writeFormError(w, http.StatusBadRequest, "Missing grimoire id.")
			return
		}
		master, err := a.masterOrErr()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if err := master.Grimoires().Delete(id, master); err != nil {
			writeJSON(w, http.StatusNotFound, common.DeleteNotFound("Grimoire"))
			return
		}
		if err := master.Grimoires().Save(); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, common.DeleteSuccess(id))
	})
}
