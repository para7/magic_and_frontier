package httpapi

import (
	"net/http"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/treasures"
)

func (a apiRouter) registerTreasureRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/treasures", func(w http.ResponseWriter, r *http.Request) {
		state, err := a.deps.TreasureRepo.LoadState()
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
		state, err := a.deps.TreasureRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Entries, strings.TrimSpace(input.ID), func(entry treasures.TreasureEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[treasures.TreasureEntry](w)
			return
		}
		validTablePaths, err := treasureSourcePaths(a.cfg.MinecraftLootTableRoot)
		if err != nil {
			writeInternalError(w, err)
			return
		}
		result := treasures.ValidateSave(input, itemIDs(itemState), grimoireIDs(grimoireState), validTablePaths, a.deps.Now())
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		if conflictID := duplicateTreasureTablePath(state.Entries, result.Entry.ID, result.Entry.TablePath); conflictID != "" {
			writeJSON(w, http.StatusBadRequest, common.SaveValidationError[treasures.TreasureEntry](common.FieldErrors{"tablePath": "Loot table path is already used by " + conflictID + "."}, "Validation failed. Fix the highlighted fields."))
			return
		}
		nextState, mode := common.UpsertEntries(state, *result.Entry, func(entry treasures.TreasureEntry) string { return entry.ID })
		result.Mode = mode
		if err := a.deps.TreasureRepo.SaveState(nextState); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("DELETE /api/treasures/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteEntry(w, r, a.deps.TreasureRepo, "treasure", "Treasure", func(entry treasures.TreasureEntry) string { return entry.ID })
	})
}
