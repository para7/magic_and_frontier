package grimoire

import (
	"time"

	"tools2/app/internal/domain/common"
)

func ValidateSave(input SaveInput, now time.Time) common.SaveResult[GrimoireEntry] {
	errs := common.ViolationsToFieldErrors(common.ValidateStruct(input), common.DefaultValidationMessage)
	id := common.RequireNonEmptyID(errs, "id", input.ID)
	if errs.Any() {
		return common.SaveValidationError[GrimoireEntry](errs, "Validation failed. Fix the highlighted fields.")
	}
	entry := GrimoireEntry{
		ID:          id,
		CastID:      input.CastID,
		CastTime:    input.CastTime,
		MPCost:      input.MPCost,
		Script:      common.NormalizeText(input.Script),
		Title:       common.NormalizeText(input.Title),
		Description: common.OptionalText(input.Description),
		UpdatedAt:   now.UTC().Format(time.RFC3339),
	}
	return common.SaveSuccess(entry, common.SaveModeCreated)
}

func Upsert(state GrimoireState, entry GrimoireEntry) (GrimoireState, common.SaveMode) {
	for i := range state.Entries {
		if state.Entries[i].ID == entry.ID {
			next := append([]GrimoireEntry{}, state.Entries...)
			next[i] = entry
			return GrimoireState{Entries: next}, common.SaveModeUpdated
		}
	}
	next := append(append([]GrimoireEntry{}, state.Entries...), entry)
	return GrimoireState{Entries: next}, common.SaveModeCreated
}

func Delete(state GrimoireState, id string) (GrimoireState, bool) {
	next := make([]GrimoireEntry, 0, len(state.Entries))
	found := false
	for _, it := range state.Entries {
		if it.ID == id {
			found = true
			continue
		}
		next = append(next, it)
	}
	return GrimoireState{Entries: next}, found
}
