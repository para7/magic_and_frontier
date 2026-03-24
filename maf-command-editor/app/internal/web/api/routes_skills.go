package api

import (
	"errors"
	"net/http"
	"strings"

	"tools2/app/internal/domain/common"
	dmaster "tools2/app/internal/domain/master"
	"tools2/app/internal/domain/skills"
)

func (a apiRouter) registerSkillRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/skills", func(w http.ResponseWriter, r *http.Request) {
		state, err := a.skillState()
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
		master, err := a.masterOrErr()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		state, err := a.skillState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Entries, strings.TrimSpace(input.ID), func(entry skills.SkillEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[skills.SkillEntry](w)
			return
		}
		result := master.Skills().Validate(input, master)
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		if err := master.Skills().Create(*result.Entry, master); err != nil {
			if errors.Is(err, dmaster.ErrDuplicateID) {
				writeDuplicateIDValidationError[skills.SkillEntry](w)
				return
			}
			writeCodedError(w, http.StatusBadRequest, "VALIDATION_FAILED", err.Error())
			return
		}
		mode := common.SaveModeCreated
		result.Mode = mode
		if err := master.Skills().Save(); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("DELETE /api/skills/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		itemState, err := a.itemState()
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
		master, err := a.masterOrErr()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if err := master.Skills().Delete(id, master); err != nil {
			writeJSON(w, http.StatusNotFound, common.DeleteNotFound("Skill"))
			return
		}
		if err := master.Skills().Save(); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, common.DeleteSuccess(id))
	})
}
