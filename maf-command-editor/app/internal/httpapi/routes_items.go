package httpapi

import (
	"net/http"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/skills"
)

func (a apiRouter) registerItemRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/items", func(w http.ResponseWriter, r *http.Request) {
		state, err := a.deps.ItemRepo.LoadItemState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	})
	mux.HandleFunc("POST /api/items", func(w http.ResponseWriter, r *http.Request) {
		var input items.SaveInput
		if !decodeJSON(w, r, &input) {
			return
		}
		state, err := a.deps.ItemRepo.LoadItemState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		skillState, err := a.deps.SkillRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Items, strings.TrimSpace(input.ID), func(entry items.ItemEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[items.ItemEntry](w)
			return
		}
		result := items.ValidateSave(input, entryIDs(skillState.Entries, func(entry skills.SkillEntry) string { return entry.ID }), a.deps.Now())
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		nextState, mode := items.Upsert(state, *result.Entry)
		result.Mode = mode
		if err := a.deps.ItemRepo.SaveItemState(nextState); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("DELETE /api/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		if id == "" {
			writeFormError(w, http.StatusBadRequest, "Missing item id.")
			return
		}
		state, err := a.deps.ItemRepo.LoadItemState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		nextState, ok := items.Delete(state, id)
		if !ok {
			writeJSON(w, http.StatusNotFound, common.DeleteNotFound("Item"))
			return
		}
		if err := a.deps.ItemRepo.SaveItemState(nextState); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, common.DeleteSuccess(id))
	})
}
