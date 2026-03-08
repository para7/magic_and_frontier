package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"tools2/app/internal/config"
	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/treasures"
	"tools2/app/internal/export"
	"tools2/app/internal/store"
	"tools2/app/internal/web"
)

type Dependencies struct {
	ItemRepo       store.ItemStateRepository
	GrimoireRepo   store.GrimoireStateRepository
	SkillRepo      store.EntryStateRepository[skills.SkillEntry]
	EnemySkillRepo store.EntryStateRepository[enemyskills.EnemySkillEntry]
	EnemyRepo      store.EntryStateRepository[enemies.EnemyEntry]
	TreasureRepo   store.EntryStateRepository[treasures.TreasureEntry]
	Now            func() time.Time
}

func NewHandler(cfg config.Config, deps Dependencies) http.Handler {
	if deps.Now == nil {
		deps.Now = time.Now
	}

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
		TreasureRepo:   deps.TreasureRepo,
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
		result := items.ValidateSave(input, deps.Now())
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		state, err := deps.ItemRepo.LoadItemState()
		if err != nil {
			writeInternalError(w, err)
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
			writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "formError": "Missing item id."})
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
		result := grimoire.ValidateSave(input, deps.Now())
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		state, err := deps.GrimoireRepo.LoadGrimoireState()
		if err != nil {
			writeInternalError(w, err)
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
			writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "formError": "Missing grimoire id."})
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
		itemState, err := deps.ItemRepo.LoadItemState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		result := skills.ValidateSave(input, itemIDs(itemState), deps.Now())
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		state, err := deps.SkillRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
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
		result := enemyskills.ValidateSave(input, deps.Now())
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		state, err := deps.EnemySkillRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
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
			writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "formError": "Missing enemy skill id."})
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
					writeJSON(w, http.StatusBadRequest, map[string]any{
						"ok":        false,
						"code":      "REFERENCE_ERROR",
						"formError": "Enemy skill is referenced by enemy " + enemy.ID + ".",
					})
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
		state, err := deps.EnemyRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
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
		result := treasures.ValidateSave(input, itemIDs(itemState), grimoireIDs(grimoireState), deps.Now())
		if !result.OK {
			writeJSON(w, http.StatusBadRequest, result)
			return
		}
		state, err := deps.TreasureRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
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

	mux.HandleFunc("POST /api/save", func(w http.ResponseWriter, r *http.Request) {
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
		skillState, err := deps.SkillRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		enemySkillState, err := deps.EnemySkillRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		enemyState, err := deps.EnemyRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		treasureState, err := deps.TreasureRepo.LoadState()
		if err != nil {
			writeInternalError(w, err)
			return
		}
		result := export.ExportDatapack(export.ExportParams{
			ItemState:          itemState,
			GrimoireState:      grimoireState,
			Skills:             skillState.Entries,
			EnemySkills:        enemySkillState.Entries,
			Enemies:            enemyState.Entries,
			Treasures:          treasureState.Entries,
			ExportSettingsPath: cfg.ExportSettingsPath,
		})
		status := http.StatusOK
		if !result.OK {
			status = http.StatusBadRequest
		}
		writeJSON(w, status, result)
	})

	return mux
}

func DefaultDependencies(cfg config.Config) Dependencies {
	return Dependencies{
		ItemRepo:       store.NewItemStateRepository(cfg.ItemStatePath),
		GrimoireRepo:   store.NewGrimoireStateRepository(cfg.GrimoireStatePath),
		SkillRepo:      store.NewEntryStateRepository[skills.SkillEntry](cfg.SkillStatePath),
		EnemySkillRepo: store.NewEntryStateRepository[enemyskills.EnemySkillEntry](cfg.EnemySkillStatePath),
		EnemyRepo:      store.NewEntryStateRepository[enemies.EnemyEntry](cfg.EnemyStatePath),
		TreasureRepo:   store.NewEntryStateRepository[treasures.TreasureEntry](cfg.TreasureStatePath),
		Now:            time.Now,
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dest any) bool {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "formError": "Invalid JSON request body."})
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeInternalError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusInternalServerError, map[string]any{"ok": false, "formError": http.StatusText(http.StatusInternalServerError), "details": err.Error()})
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

func deleteEntry[T any](w http.ResponseWriter, r *http.Request, repo store.EntryStateRepository[T], missingLabel, notFoundLabel string, idOf func(T) string) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "formError": "Missing " + missingLabel + " id."})
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
