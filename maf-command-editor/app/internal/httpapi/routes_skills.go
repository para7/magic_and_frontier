package httpapi

import (
	"net/http"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/skills"
)

func (a apiRouter) registerSkillRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/skills", func(w http.ResponseWriter, r *http.Request) {
		state, err := a.deps.SkillRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	})
	mux.HandleFunc("POST /api/skills", func(w http.ResponseWriter, r *http.Request) {
		var input skills.SaveInput
		if !decodeJSON(w, r, &input) {
			return
		}
		state, err := a.deps.SkillRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Entries, strings.TrimSpace(input.ID), func(entry skills.SkillEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[skills.SkillEntry](w)
			return
		}
		result := skills.ValidateSave(input, a.deps.Now())
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		nextState, mode := common.UpsertEntries(state, *result.Entry, func(entry skills.SkillEntry) string { return entry.ID })
		result.Mode = mode
		if err := a.deps.SkillRepo.SaveState(nextState); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("DELETE /api/skills/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		itemState, err := a.deps.ItemRepo.LoadItemState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		for _, entry := range itemState.Items {
			if entry.SkillID == id {
				writeCodedError(w, http.StatusBadRequest, "REFERENCE_ERROR", "Skill is referenced by item "+entry.ID+".")
				return
			}
		}
		deleteEntry(w, r, a.deps.SkillRepo, "skill", "Skill", func(entry skills.SkillEntry) string { return entry.ID })
	})
}
