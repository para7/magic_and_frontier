package httpapi

import (
	"net/http"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/loottables"
)

func (a apiRouter) registerLootTableRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/loottables", func(w http.ResponseWriter, r *http.Request) {
		state, err := a.deps.LootTableRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	})
	mux.HandleFunc("POST /api/loottables", func(w http.ResponseWriter, r *http.Request) {
		var input loottables.SaveInput
		if !decodeJSON(w, r, &input) {
			return
		}
		state, err := a.deps.LootTableRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		itemState, err := a.deps.ItemRepo.LoadItemState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		grimoireState, err := a.deps.GrimoireRepo.LoadGrimoireState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Entries, strings.TrimSpace(input.ID), func(entry loottables.LootTableEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[loottables.LootTableEntry](w)
			return
		}
		result := loottables.ValidateSave(input, itemIDs(itemState), grimoireIDs(grimoireState), a.deps.Now())
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		nextState, mode := common.UpsertEntries(state, *result.Entry, func(entry loottables.LootTableEntry) string { return entry.ID })
		result.Mode = mode
		if err := a.deps.LootTableRepo.SaveState(nextState); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("DELETE /api/loottables/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteEntry(w, r, a.deps.LootTableRepo, "loottable", "Loottable", func(entry loottables.LootTableEntry) string { return entry.ID })
	})
}
