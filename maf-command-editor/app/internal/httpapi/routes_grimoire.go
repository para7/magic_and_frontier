package httpapi

import (
	"net/http"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/grimoire"
)

func (a apiRouter) registerGrimoireRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/grimoire", func(w http.ResponseWriter, r *http.Request) {
		state, err := a.deps.GrimoireRepo.LoadGrimoireState()
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
		state, err := a.deps.GrimoireRepo.LoadGrimoireState()
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
		result := grimoire.ValidateSave(input, a.deps.Now())
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		if conflictID := duplicateCastID(state.Entries, result.Entry.ID, result.Entry.CastID); conflictID != "" {
			writeJSON(w, http.StatusBadRequest, common.SaveValidationError[grimoire.GrimoireEntry](common.FieldErrors{"castid": "Cast ID is already used by " + conflictID + "."}, "Validation failed. Fix the highlighted fields."))
			return
		}
		nextState, mode := grimoire.Upsert(state, *result.Entry)
		result.Mode = mode
		if err := a.deps.GrimoireRepo.SaveGrimoireState(nextState); err != nil {
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
		state, err := a.deps.GrimoireRepo.LoadGrimoireState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		nextState, ok := grimoire.Delete(state, id)
		if !ok {
			writeJSON(w, http.StatusNotFound, common.DeleteNotFound("Grimoire"))
			return
		}
		if err := a.deps.GrimoireRepo.SaveGrimoireState(nextState); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, common.DeleteSuccess(id))
	})
}
