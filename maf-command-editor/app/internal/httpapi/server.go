package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"tools2/app/internal/application"
	"tools2/app/internal/config"
	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/loottables"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/spawntables"
	"tools2/app/internal/domain/treasures"
	"tools2/app/internal/mcsource"
	"tools2/app/internal/store"
	"tools2/app/internal/web"
)

type Dependencies = application.Dependencies

type apiErrorResponse struct {
	OK        bool   `json:"ok"`
	Code      string `json:"code,omitempty"`
	FormError string `json:"formError,omitempty"`
	Details   string `json:"details,omitempty"`
}

func NewHandler(cfg config.Config, deps Dependencies) http.Handler {
	defaults := DefaultDependencies(cfg)
	if deps.ItemRepo == nil {
		deps.ItemRepo = defaults.ItemRepo
	}
	if deps.GrimoireRepo == nil {
		deps.GrimoireRepo = defaults.GrimoireRepo
	}
	if deps.SkillRepo == nil {
		deps.SkillRepo = defaults.SkillRepo
	}
	if deps.EnemySkillRepo == nil {
		deps.EnemySkillRepo = defaults.EnemySkillRepo
	}
	if deps.EnemyRepo == nil {
		deps.EnemyRepo = defaults.EnemyRepo
	}
	if deps.SpawnTableRepo == nil {
		deps.SpawnTableRepo = defaults.SpawnTableRepo
	}
	if deps.TreasureRepo == nil {
		deps.TreasureRepo = defaults.TreasureRepo
	}
	if deps.LootTableRepo == nil {
		deps.LootTableRepo = defaults.LootTableRepo
	}
	if deps.CounterRepo == nil {
		deps.CounterRepo = defaults.CounterRepo
	}
	if deps.ExportSettingsPath == "" {
		deps.ExportSettingsPath = defaults.ExportSettingsPath
	}
	if deps.Now == nil {
		deps.Now = time.Now
	}
	appService := application.NewService(cfg, deps)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/items", http.StatusFound)
	})
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	})
	web.RegisterRoutes(mux, cfg, web.Dependencies{
		ItemRepo:       deps.ItemRepo,
		GrimoireRepo:   deps.GrimoireRepo,
		SkillRepo:      deps.SkillRepo,
		EnemySkillRepo: deps.EnemySkillRepo,
		EnemyRepo:      deps.EnemyRepo,
		SpawnTableRepo: deps.SpawnTableRepo,
		TreasureRepo:   deps.TreasureRepo,
		LootTableRepo:  deps.LootTableRepo,
		CounterRepo:    deps.CounterRepo,
		Now:            deps.Now,
	})

	mux.HandleFunc("GET /api/items", func(w http.ResponseWriter, r *http.Request) {
		state, err := deps.ItemRepo.LoadItemState()
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
		state, err := deps.ItemRepo.LoadItemState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		skillState, err := deps.SkillRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Items, strings.TrimSpace(input.ID), func(entry items.ItemEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[items.ItemEntry](w)
			return
		}
		result := items.ValidateSave(input, entryIDs(skillState.Entries, func(entry skills.SkillEntry) string { return entry.ID }), deps.Now())
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		nextState, mode := items.Upsert(state, *result.Entry)
		result.Mode = mode
		if err := deps.ItemRepo.SaveItemState(nextState); err != nil {
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
		state, err := deps.ItemRepo.LoadItemState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		nextState, ok := items.Delete(state, id)
		if !ok {
			writeJSON(w, http.StatusNotFound, common.DeleteNotFound("Item"))
			return
		}
		if err := deps.ItemRepo.SaveItemState(nextState); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, common.DeleteSuccess(id))
	})

	mux.HandleFunc("GET /api/grimoire", func(w http.ResponseWriter, r *http.Request) {
		state, err := deps.GrimoireRepo.LoadGrimoireState()
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
		state, err := deps.GrimoireRepo.LoadGrimoireState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Entries, strings.TrimSpace(input.ID), func(entry grimoire.GrimoireEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[grimoire.GrimoireEntry](w)
			return
		}
		castID, allocErr := appService.AllocateCastID()
		if allocErr != nil {
			writeInternalError(w, allocErr)
			return
		}
		input.CastID = castID
		result := grimoire.ValidateSave(input, deps.Now())
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
		if err := deps.GrimoireRepo.SaveGrimoireState(nextState); err != nil {
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
		state, err := deps.GrimoireRepo.LoadGrimoireState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		nextState, ok := grimoire.Delete(state, id)
		if !ok {
			writeJSON(w, http.StatusNotFound, common.DeleteNotFound("Grimoire"))
			return
		}
		if err := deps.GrimoireRepo.SaveGrimoireState(nextState); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, common.DeleteSuccess(id))
	})

	mux.HandleFunc("GET /api/skills", func(w http.ResponseWriter, r *http.Request) {
		state, err := deps.SkillRepo.LoadState()
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
		state, err := deps.SkillRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Entries, strings.TrimSpace(input.ID), func(entry skills.SkillEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[skills.SkillEntry](w)
			return
		}
		result := skills.ValidateSave(input, deps.Now())
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		nextState, mode := common.UpsertEntries(state, *result.Entry, func(entry skills.SkillEntry) string { return entry.ID })
		result.Mode = mode
		if err := deps.SkillRepo.SaveState(nextState); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("DELETE /api/skills/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		itemState, err := deps.ItemRepo.LoadItemState()
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
		deleteEntry(w, r, deps.SkillRepo, "skill", "Skill", func(entry skills.SkillEntry) string { return entry.ID })
	})

	mux.HandleFunc("GET /api/enemy-skills", func(w http.ResponseWriter, r *http.Request) {
		state, err := deps.EnemySkillRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	})
	mux.HandleFunc("POST /api/enemy-skills", func(w http.ResponseWriter, r *http.Request) {
		var input enemyskills.SaveInput
		if !decodeJSON(w, r, &input) {
			return
		}
		state, err := deps.EnemySkillRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Entries, strings.TrimSpace(input.ID), func(entry enemyskills.EnemySkillEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[enemyskills.EnemySkillEntry](w)
			return
		}
		result := enemyskills.ValidateSave(input, deps.Now())
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		nextState, mode := common.UpsertEntries(state, *result.Entry, func(entry enemyskills.EnemySkillEntry) string { return entry.ID })
		result.Mode = mode
		if err := deps.EnemySkillRepo.SaveState(nextState); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("DELETE /api/enemy-skills/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		if id == "" {
			writeFormError(w, http.StatusBadRequest, "Missing enemy skill id.")
			return
		}
		enemyState, err := deps.EnemyRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		for _, enemy := range enemyState.Entries {
			for _, skillID := range enemy.EnemySkillIDs {
				if skillID == id {
					writeCodedError(w, http.StatusBadRequest, "REFERENCE_ERROR", "Enemy skill is referenced by enemy "+enemy.ID+".")
					return
				}
			}
		}
		state, err := deps.EnemySkillRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		nextState, ok := common.DeleteEntries(state, id, func(entry enemyskills.EnemySkillEntry) string { return entry.ID })
		if !ok {
			writeJSON(w, http.StatusNotFound, common.DeleteNotFound("Enemy skill"))
			return
		}
		if err := deps.EnemySkillRepo.SaveState(nextState); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, common.DeleteSuccess(id))
	})

	mux.HandleFunc("GET /api/enemies", func(w http.ResponseWriter, r *http.Request) {
		state, err := deps.EnemyRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	})
	mux.HandleFunc("POST /api/enemies", func(w http.ResponseWriter, r *http.Request) {
		var input enemies.SaveInput
		if !decodeJSON(w, r, &input) {
			return
		}
		state, err := deps.EnemyRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Entries, strings.TrimSpace(input.ID), func(entry enemies.EnemyEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[enemies.EnemyEntry](w)
			return
		}
		enemySkillState, err := deps.EnemySkillRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		itemState, err := deps.ItemRepo.LoadItemState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		grimoireState, err := deps.GrimoireRepo.LoadGrimoireState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		result := enemies.ValidateSave(input, entryIDs(enemySkillState.Entries, func(entry enemyskills.EnemySkillEntry) string { return entry.ID }), itemIDs(itemState), grimoireIDs(grimoireState), deps.Now())
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		nextState, mode := common.UpsertEntries(state, *result.Entry, func(entry enemies.EnemyEntry) string { return entry.ID })
		result.Mode = mode
		if err := deps.EnemyRepo.SaveState(nextState); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("DELETE /api/enemies/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteEntry(w, r, deps.EnemyRepo, "enemy", "Enemy", func(entry enemies.EnemyEntry) string { return entry.ID })
	})

	mux.HandleFunc("GET /api/spawn-tables", func(w http.ResponseWriter, r *http.Request) {
		state, err := deps.SpawnTableRepo.LoadState()
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
		state, err := deps.SpawnTableRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Entries, strings.TrimSpace(input.ID), func(entry spawntables.SpawnTableEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[spawntables.SpawnTableEntry](w)
			return
		}
		enemyState, err := deps.EnemyRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		result := spawntables.ValidateSave(input, entryIDs(enemyState.Entries, func(entry enemies.EnemyEntry) string { return entry.ID }), deps.Now())
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		if conflictID, ok := spawntables.FirstOverlap(state.Entries, *result.Entry); ok {
			writeJSON(w, http.StatusBadRequest, common.SaveValidationError[spawntables.SpawnTableEntry](
				common.FieldErrors{"range": "Range overlaps with " + conflictID + "."},
				"Validation failed. Fix the highlighted fields.",
			))
			return
		}
		nextState, mode := common.UpsertEntries(state, *result.Entry, func(entry spawntables.SpawnTableEntry) string { return entry.ID })
		result.Mode = mode
		if err := deps.SpawnTableRepo.SaveState(nextState); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("DELETE /api/spawn-tables/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteEntry(w, r, deps.SpawnTableRepo, "spawn table", "Spawn table", func(entry spawntables.SpawnTableEntry) string { return entry.ID })
	})

	mux.HandleFunc("GET /api/treasures", func(w http.ResponseWriter, r *http.Request) {
		state, err := deps.TreasureRepo.LoadState()
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
		itemState, err := deps.ItemRepo.LoadItemState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		grimoireState, err := deps.GrimoireRepo.LoadGrimoireState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		state, err := deps.TreasureRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Entries, strings.TrimSpace(input.ID), func(entry treasures.TreasureEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[treasures.TreasureEntry](w)
			return
		}
		validTablePaths, err := treasureSourcePaths(cfg.MinecraftLootTableRoot)
		if err != nil {
			writeInternalError(w, err)
			return
		}
		result := treasures.ValidateSave(input, itemIDs(itemState), grimoireIDs(grimoireState), validTablePaths, deps.Now())
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
		if err := deps.TreasureRepo.SaveState(nextState); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("DELETE /api/treasures/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteEntry(w, r, deps.TreasureRepo, "treasure", "Treasure", func(entry treasures.TreasureEntry) string { return entry.ID })
	})

	mux.HandleFunc("GET /api/loottables", func(w http.ResponseWriter, r *http.Request) {
		state, err := deps.LootTableRepo.LoadState()
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
		state, err := deps.LootTableRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		itemState, err := deps.ItemRepo.LoadItemState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		grimoireState, err := deps.GrimoireRepo.LoadGrimoireState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		if _, ok := findEntry(state.Entries, strings.TrimSpace(input.ID), func(entry loottables.LootTableEntry) string { return entry.ID }); ok {
			writeDuplicateIDValidationError[loottables.LootTableEntry](w)
			return
		}
		result := loottables.ValidateSave(input, itemIDs(itemState), grimoireIDs(grimoireState), deps.Now())
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		nextState, mode := common.UpsertEntries(state, *result.Entry, func(entry loottables.LootTableEntry) string { return entry.ID })
		result.Mode = mode
		if err := deps.LootTableRepo.SaveState(nextState); err != nil {
			writeInternalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("DELETE /api/loottables/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteEntry(w, r, deps.LootTableRepo, "loottable", "Loottable", func(entry loottables.LootTableEntry) string { return entry.ID })
	})

	mux.HandleFunc("POST /api/save", func(w http.ResponseWriter, r *http.Request) {
		result := appService.ExportDatapack()
		status := http.StatusOK
		if !result.OK {
			status = http.StatusBadRequest
		}
		writeJSON(w, status, result)
	})

	return mux
}

func DefaultDependencies(cfg config.Config) Dependencies {
	return application.DefaultDependencies(cfg)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dest any) bool {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		writeFormError(w, http.StatusBadRequest, "Invalid JSON request body.")
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeFormError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, apiErrorResponse{
		OK:        false,
		FormError: message,
	})
}

func writeCodedError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, apiErrorResponse{
		OK:        false,
		Code:      code,
		FormError: message,
	})
}

func writeInternalError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusInternalServerError, apiErrorResponse{
		OK:        false,
		FormError: http.StatusText(http.StatusInternalServerError),
		Details:   err.Error(),
	})
}

func writeDuplicateIDValidationError[T any](w http.ResponseWriter) {
	writeJSON(w, http.StatusBadRequest, common.SaveValidationError[T](
		common.FieldErrors{"id": "この ID は既に使用されています。"},
		"Validation failed. Fix the highlighted fields.",
	))
}

func itemIDs(state items.ItemState) map[string]struct{} {
	return entryIDs(state.Items, func(entry items.ItemEntry) string { return entry.ID })
}

func grimoireIDs(state grimoire.GrimoireState) map[string]struct{} {
	return entryIDs(state.Entries, func(entry grimoire.GrimoireEntry) string { return entry.ID })
}

func entryIDs[T any](entries []T, idOf func(T) string) map[string]struct{} {
	out := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		id := strings.TrimSpace(idOf(entry))
		if id != "" {
			out[id] = struct{}{}
		}
	}
	return out
}

func findEntry[T any](entries []T, id string, idOf func(T) string) (T, bool) {
	var zero T
	id = strings.TrimSpace(id)
	if id == "" {
		return zero, false
	}
	for _, entry := range entries {
		if strings.TrimSpace(idOf(entry)) == id {
			return entry, true
		}
	}
	return zero, false
}

func deleteEntry[T any](w http.ResponseWriter, r *http.Request, repo store.EntryStateRepository[T], missingLabel, notFoundLabel string, idOf func(T) string) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFormError(w, http.StatusBadRequest, "Missing "+missingLabel+" id.")
		return
	}
	state, err := repo.LoadState()
	if err != nil {
		writeInternalError(w, err)
		return
	}
	nextState, ok := common.DeleteEntries(state, id, idOf)
	if !ok {
		writeJSON(w, http.StatusNotFound, common.DeleteNotFound(notFoundLabel))
		return
	}
	if err := repo.SaveState(nextState); err != nil {
		writeInternalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, common.DeleteSuccess(id))
}

func duplicateCastID(entries []grimoire.GrimoireEntry, entryID string, castID int) string {
	for _, entry := range entries {
		if entry.ID != entryID && entry.CastID == castID {
			return entry.ID
		}
	}
	return ""
}

func duplicateTreasureTablePath(entries []treasures.TreasureEntry, entryID, tablePath string) string {
	for _, entry := range entries {
		if entry.ID != entryID && entry.TablePath == tablePath {
			return entry.ID
		}
	}
	return ""
}

func treasureSourcePaths(root string) (map[string]struct{}, error) {
	sources, err := mcsource.ListLootTables(root)
	if err != nil {
		return nil, err
	}
	out := make(map[string]struct{}, len(sources))
	for _, source := range sources {
		out[source.TablePath] = struct{}{}
	}
	return out, nil
}
