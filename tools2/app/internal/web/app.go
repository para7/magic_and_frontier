package web

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"

	"tools2/app/internal/application"
	"tools2/app/internal/config"
	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/treasures"
	"tools2/app/internal/webui"
	"tools2/app/views"
)

type Dependencies = application.Dependencies

type App struct {
	cfg  config.Config
	deps Dependencies
}

func RegisterRoutes(mux *http.ServeMux, cfg config.Config, deps Dependencies) {
	if deps.Now == nil {
		deps.Now = time.Now
	}
	app := App{cfg: cfg, deps: deps}

	mux.HandleFunc("GET /items", app.itemsPage)
	mux.HandleFunc("GET /items/new", app.itemsNewPage)
	mux.HandleFunc("POST /items/new", app.itemsSubmit)
	mux.HandleFunc("GET /items/edit", app.itemsEditPage)
	mux.HandleFunc("POST /items/edit", app.itemsEditSubmit)
	mux.HandleFunc("POST /items/{id}/delete", app.itemsDelete)

	mux.HandleFunc("GET /grimoire", app.grimoirePage)
	mux.HandleFunc("GET /grimoire/new", app.grimoireNewPage)
	mux.HandleFunc("POST /grimoire/new", app.grimoireSubmit)
	mux.HandleFunc("GET /grimoire/edit", app.grimoireEditPage)
	mux.HandleFunc("POST /grimoire/edit", app.grimoireEditSubmit)
	mux.HandleFunc("POST /grimoire/{id}/delete", app.grimoireDelete)

	mux.HandleFunc("GET /skills", app.skillsPage)
	mux.HandleFunc("GET /skills/new", app.skillsNewPage)
	mux.HandleFunc("POST /skills/new", app.skillsSubmit)
	mux.HandleFunc("GET /skills/edit", app.skillsEditPage)
	mux.HandleFunc("POST /skills/edit", app.skillsEditSubmit)
	mux.HandleFunc("POST /skills/{id}/delete", app.skillsDelete)

	mux.HandleFunc("GET /enemy-skills", app.enemySkillsPage)
	mux.HandleFunc("GET /enemy-skills/new", app.enemySkillsNewPage)
	mux.HandleFunc("POST /enemy-skills/new", app.enemySkillsSubmit)
	mux.HandleFunc("GET /enemy-skills/edit", app.enemySkillsEditPage)
	mux.HandleFunc("POST /enemy-skills/edit", app.enemySkillsEditSubmit)
	mux.HandleFunc("POST /enemy-skills/{id}/delete", app.enemySkillsDelete)

	mux.HandleFunc("GET /treasures", app.treasuresPage)
	mux.HandleFunc("GET /treasures/new", app.treasuresNewPage)
	mux.HandleFunc("POST /treasures/new", app.treasuresSubmit)
	mux.HandleFunc("GET /treasures/edit", app.treasuresEditPage)
	mux.HandleFunc("POST /treasures/edit", app.treasuresEditSubmit)
	mux.HandleFunc("POST /treasures/{id}/delete", app.treasuresDelete)

	mux.HandleFunc("GET /enemies", app.enemiesPage)
	mux.HandleFunc("GET /enemies/new", app.enemiesNewPage)
	mux.HandleFunc("POST /enemies/new", app.enemiesSubmit)
	mux.HandleFunc("GET /enemies/edit", app.enemiesEditPage)
	mux.HandleFunc("POST /enemies/edit", app.enemiesEditSubmit)
	mux.HandleFunc("POST /enemies/{id}/delete", app.enemiesDelete)

	mux.HandleFunc("POST /save", app.saveExport)
}

func (a App) itemsPage(w http.ResponseWriter, r *http.Request) {
	notice := consumeFlashNotice(w, r)
	state, err := a.deps.ItemRepo.LoadItemState()
	if err != nil {
		a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error())})
		return
	}
	a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Entries: state.Items, Notice: notice})
}

func (a App) itemsNewPage(w http.ResponseWriter, r *http.Request) {
	a.renderItemForm(w, r, webui.ItemsPageData{Meta: itemMeta(), Form: defaultItemForm()})
}

func (a App) itemsEditPage(w http.ResponseWriter, r *http.Request) {
	state, err := a.deps.ItemRepo.LoadItemState()
	if err != nil {
		a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if entry, ok := findEntry(state.Items, id, func(entry items.ItemEntry) string { return entry.ID }); ok {
		a.renderItemForm(w, r, webui.ItemsPageData{Meta: itemMeta(), Form: itemEntryToForm(entry)})
		return
	}
	a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Entries: state.Items, Notice: errorNotice("Item not found.")})
}

func (a App) itemsSubmit(w http.ResponseWriter, r *http.Request) {
	a.itemsSave(w, r, false)
}

func (a App) itemsEditSubmit(w http.ResponseWriter, r *http.Request) {
	a.itemsSave(w, r, true)
}

func (a App) itemsSave(w http.ResponseWriter, r *http.Request, editing bool) {
	_ = r.ParseForm()
	form, input, parseErrs := parseItemForm(r)
	form.IsEditing = editing
	state, err := a.deps.ItemRepo.LoadItemState()
	if err != nil {
		a.renderItemForm(w, r, webui.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	if editing {
		if _, ok := findEntry(state.Items, form.ID, func(entry items.ItemEntry) string { return entry.ID }); !ok {
			a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Entries: state.Items, Notice: errorNotice("Item not found.")})
			return
		}
	}
	result := items.ValidateSave(input, a.deps.Now())
	errors := mergeFieldErrors(parseErrs, result.FieldErrors)
	if len(errors) > 0 {
		form.FieldErrors = errors
		form.FormError = formErrorText(result.FormError)
		a.renderItemForm(w, r, webui.ItemsPageData{Meta: itemMeta(), Form: form})
		return
	}
	nextState, mode := items.Upsert(state, *result.Entry)
	if err := a.deps.ItemRepo.SaveItemState(nextState); err != nil {
		a.renderItemForm(w, r, webui.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	notice := successNotice(noticeText("Item", mode))
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, "/items", notice) {
		return
	}
	a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Entries: nextState.Items, Notice: notice})
}

func (a App) itemsDelete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	state, err := a.deps.ItemRepo.LoadItemState()
	if err != nil {
		a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error())})
		return
	}
	nextState, ok := items.Delete(state, id)
	if !ok {
		a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Entries: state.Items, Notice: errorNotice("Item not found.")})
		return
	}
	if err := a.deps.ItemRepo.SaveItemState(nextState); err != nil {
		a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Entries: state.Items, Notice: errorNotice(err.Error())})
		return
	}
	notice := successNotice("Item deleted.")
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, "/items", notice) {
		return
	}
	a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Entries: nextState.Items, Notice: notice})
}

func (a App) grimoirePage(w http.ResponseWriter, r *http.Request) {
	notice := consumeFlashNotice(w, r)
	state, err := a.deps.GrimoireRepo.LoadGrimoireState()
	if err != nil {
		a.renderGrimoire(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Notice: errorNotice(err.Error())})
		return
	}
	a.renderGrimoire(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Entries: state.Entries, Notice: notice})
}

func (a App) grimoireNewPage(w http.ResponseWriter, r *http.Request) {
	state, err := a.deps.GrimoireRepo.LoadGrimoireState()
	if err != nil {
		a.renderGrimoireForm(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Notice: errorNotice(err.Error()), Form: defaultGrimoireForm(nil)})
		return
	}
	a.renderGrimoireForm(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Form: defaultGrimoireForm(state.Entries)})
}

func (a App) grimoireEditPage(w http.ResponseWriter, r *http.Request) {
	state, err := a.deps.GrimoireRepo.LoadGrimoireState()
	if err != nil {
		a.renderGrimoire(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if entry, ok := findEntry(state.Entries, id, func(entry grimoire.GrimoireEntry) string { return entry.ID }); ok {
		a.renderGrimoireForm(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Form: grimoireEntryToForm(entry)})
		return
	}
	a.renderGrimoire(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Entries: state.Entries, Notice: errorNotice("Grimoire entry not found.")})
}

func (a App) grimoireSubmit(w http.ResponseWriter, r *http.Request) {
	a.grimoireSave(w, r, false)
}

func (a App) grimoireEditSubmit(w http.ResponseWriter, r *http.Request) {
	a.grimoireSave(w, r, true)
}

func (a App) grimoireSave(w http.ResponseWriter, r *http.Request, editing bool) {
	_ = r.ParseForm()
	state, err := a.deps.GrimoireRepo.LoadGrimoireState()
	if err != nil {
		form := defaultGrimoireForm(nil)
		form.IsEditing = editing
		a.renderGrimoireForm(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	form, input, parseErrs := parseGrimoireForm(r)
	form.IsEditing = editing
	if editing {
		if _, ok := findEntry(state.Entries, form.ID, func(entry grimoire.GrimoireEntry) string { return entry.ID }); !ok {
			a.renderGrimoire(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Entries: state.Entries, Notice: errorNotice("Grimoire entry not found.")})
			return
		}
	}
	result := grimoire.ValidateSave(input, a.deps.Now())
	errors := mergeFieldErrors(parseErrs, mapFieldErrors(result.FieldErrors, mapGrimoireField))
	if len(errors) > 0 {
		form.FieldErrors = errors
		form.FormError = formErrorText(result.FormError)
		a.renderGrimoireForm(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Form: form})
		return
	}
	nextState, mode := grimoire.Upsert(state, *result.Entry)
	if err := a.deps.GrimoireRepo.SaveGrimoireState(nextState); err != nil {
		a.renderGrimoireForm(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	notice := successNotice(noticeText("Grimoire entry", mode))
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, "/grimoire", notice) {
		return
	}
	a.renderGrimoire(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) grimoireDelete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	state, err := a.deps.GrimoireRepo.LoadGrimoireState()
	if err != nil {
		a.renderGrimoire(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Notice: errorNotice(err.Error())})
		return
	}
	nextState, ok := grimoire.Delete(state, id)
	if !ok {
		a.renderGrimoire(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Entries: state.Entries, Notice: errorNotice("Grimoire entry not found.")})
		return
	}
	if err := a.deps.GrimoireRepo.SaveGrimoireState(nextState); err != nil {
		a.renderGrimoire(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Entries: state.Entries, Notice: errorNotice(err.Error())})
		return
	}
	notice := successNotice("Grimoire entry deleted.")
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, "/grimoire", notice) {
		return
	}
	a.renderGrimoire(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) skillsPage(w http.ResponseWriter, r *http.Request) {
	notice := consumeFlashNotice(w, r)
	state, err := a.deps.SkillRepo.LoadState()
	if err != nil {
		a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Notice: errorNotice(err.Error())})
		return
	}
	a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Entries: state.Entries, Notice: notice})
}

func (a App) skillsNewPage(w http.ResponseWriter, r *http.Request) {
	itemState, err := a.deps.ItemRepo.LoadItemState()
	if err != nil {
		a.renderSkillForm(w, r, webui.SkillsPageData{Meta: skillsMeta(), Notice: errorNotice(err.Error()), Form: defaultSkillForm(nil)})
		return
	}
	a.renderSkillForm(w, r, webui.SkillsPageData{Meta: skillsMeta(), Form: defaultSkillForm(itemState.Items)})
}

func (a App) skillsEditPage(w http.ResponseWriter, r *http.Request) {
	itemState, err := a.deps.ItemRepo.LoadItemState()
	if err != nil {
		a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Notice: errorNotice(err.Error())})
		return
	}
	state, err := a.deps.SkillRepo.LoadState()
	if err != nil {
		a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if entry, ok := findEntry(state.Entries, id, func(entry skills.SkillEntry) string { return entry.ID }); ok {
		form := skillEntryToForm(entry, itemOptions(itemState.Items))
		a.renderSkillForm(w, r, webui.SkillsPageData{Meta: skillsMeta(), Form: form})
		return
	}
	a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Entries: state.Entries, Notice: errorNotice("Skill not found.")})
}

func (a App) skillsSubmit(w http.ResponseWriter, r *http.Request) {
	a.skillsSave(w, r, false)
}

func (a App) skillsEditSubmit(w http.ResponseWriter, r *http.Request) {
	a.skillsSave(w, r, true)
}

func (a App) skillsSave(w http.ResponseWriter, r *http.Request, editing bool) {
	_ = r.ParseForm()
	itemState, err := a.deps.ItemRepo.LoadItemState()
	if err != nil {
		form := defaultSkillForm(nil)
		form.IsEditing = editing
		a.renderSkillForm(w, r, webui.SkillsPageData{Meta: skillsMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	state, err := a.deps.SkillRepo.LoadState()
	if err != nil {
		form := defaultSkillForm(itemState.Items)
		form.IsEditing = editing
		a.renderSkillForm(w, r, webui.SkillsPageData{Meta: skillsMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	form, input, parseErrs := parseSkillForm(r, itemState.Items)
	form.IsEditing = editing
	if editing {
		if _, ok := findEntry(state.Entries, form.ID, func(entry skills.SkillEntry) string { return entry.ID }); !ok {
			a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Entries: state.Entries, Notice: errorNotice("Skill not found.")})
			return
		}
	}
	result := skills.ValidateSave(input, itemIDSet(itemState), a.deps.Now())
	errors := mergeFieldErrors(parseErrs, result.FieldErrors)
	if len(errors) > 0 {
		form.FieldErrors = errors
		form.FormError = formErrorText(result.FormError)
		a.renderSkillForm(w, r, webui.SkillsPageData{Meta: skillsMeta(), Form: form})
		return
	}
	nextState, mode := common.UpsertEntries(state, *result.Entry, func(entry skills.SkillEntry) string { return entry.ID })
	if err := a.deps.SkillRepo.SaveState(nextState); err != nil {
		a.renderSkillForm(w, r, webui.SkillsPageData{Meta: skillsMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	notice := successNotice(noticeText("Skill", mode))
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, "/skills", notice) {
		return
	}
	a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) skillsDelete(w http.ResponseWriter, r *http.Request) {
	state, err := a.deps.SkillRepo.LoadState()
	if err != nil {
		a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.PathValue("id"))
	nextState, ok := common.DeleteEntries(state, id, func(entry skills.SkillEntry) string { return entry.ID })
	if !ok {
		a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Entries: state.Entries, Notice: errorNotice("Skill not found.")})
		return
	}
	if err := a.deps.SkillRepo.SaveState(nextState); err != nil {
		a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Entries: state.Entries, Notice: errorNotice(err.Error())})
		return
	}
	notice := successNotice("Skill deleted.")
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, "/skills", notice) {
		return
	}
	a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) enemySkillsPage(w http.ResponseWriter, r *http.Request) {
	notice := consumeFlashNotice(w, r)
	state, err := a.deps.EnemySkillRepo.LoadState()
	if err != nil {
		a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Notice: errorNotice(err.Error())})
		return
	}
	a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Entries: state.Entries, Notice: notice})
}

func (a App) enemySkillsNewPage(w http.ResponseWriter, r *http.Request) {
	a.renderEnemySkillForm(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Form: defaultEnemySkillForm()})
}

func (a App) enemySkillsEditPage(w http.ResponseWriter, r *http.Request) {
	state, err := a.deps.EnemySkillRepo.LoadState()
	if err != nil {
		a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if entry, ok := findEntry(state.Entries, id, func(entry enemyskills.EnemySkillEntry) string { return entry.ID }); ok {
		a.renderEnemySkillForm(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Form: enemySkillEntryToForm(entry)})
		return
	}
	a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Entries: state.Entries, Notice: errorNotice("Enemy skill not found.")})
}

func (a App) enemySkillsSubmit(w http.ResponseWriter, r *http.Request) {
	a.enemySkillsSave(w, r, false)
}

func (a App) enemySkillsEditSubmit(w http.ResponseWriter, r *http.Request) {
	a.enemySkillsSave(w, r, true)
}

func (a App) enemySkillsSave(w http.ResponseWriter, r *http.Request, editing bool) {
	_ = r.ParseForm()
	state, err := a.deps.EnemySkillRepo.LoadState()
	if err != nil {
		form := defaultEnemySkillForm()
		form.IsEditing = editing
		a.renderEnemySkillForm(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	form, input, parseErrs := parseEnemySkillForm(r)
	form.IsEditing = editing
	if editing {
		if _, ok := findEntry(state.Entries, form.ID, func(entry enemyskills.EnemySkillEntry) string { return entry.ID }); !ok {
			a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Entries: state.Entries, Notice: errorNotice("Enemy skill not found.")})
			return
		}
	}
	result := enemyskills.ValidateSave(input, a.deps.Now())
	errors := mergeFieldErrors(parseErrs, result.FieldErrors)
	if len(errors) > 0 {
		form.FieldErrors = errors
		form.FormError = formErrorText(result.FormError)
		a.renderEnemySkillForm(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Form: form})
		return
	}
	nextState, mode := common.UpsertEntries(state, *result.Entry, func(entry enemyskills.EnemySkillEntry) string { return entry.ID })
	if err := a.deps.EnemySkillRepo.SaveState(nextState); err != nil {
		a.renderEnemySkillForm(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	notice := successNotice(noticeText("Enemy skill", mode))
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, "/enemy-skills", notice) {
		return
	}
	a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) enemySkillsDelete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	state, err := a.deps.EnemySkillRepo.LoadState()
	if err != nil {
		a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Notice: errorNotice(err.Error())})
		return
	}
	enemyState, err := a.deps.EnemyRepo.LoadState()
	if err != nil {
		a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Entries: state.Entries, Notice: errorNotice(err.Error())})
		return
	}
	for _, enemy := range enemyState.Entries {
		for _, refID := range enemy.EnemySkillIDs {
			if refID == id {
				a.renderEnemySkills(w, r, webui.EnemySkillsPageData{
					Meta:    enemySkillsMeta(),
					Entries: state.Entries,
					Notice:  errorNotice("Enemy skill is referenced by enemy " + enemy.ID + "."),
				})
				return
			}
		}
	}
	nextState, ok := common.DeleteEntries(state, id, func(entry enemyskills.EnemySkillEntry) string { return entry.ID })
	if !ok {
		a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Entries: state.Entries, Notice: errorNotice("Enemy skill not found.")})
		return
	}
	if err := a.deps.EnemySkillRepo.SaveState(nextState); err != nil {
		a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Entries: state.Entries, Notice: errorNotice(err.Error())})
		return
	}
	notice := successNotice("Enemy skill deleted.")
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, "/enemy-skills", notice) {
		return
	}
	a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) treasuresPage(w http.ResponseWriter, r *http.Request) {
	notice := consumeFlashNotice(w, r)
	state, err := a.deps.TreasureRepo.LoadState()
	if err != nil {
		a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error())})
		return
	}
	a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Entries: state.Entries, Notice: notice})
}

func (a App) treasuresNewPage(w http.ResponseWriter, r *http.Request) {
	itemState, err := a.deps.ItemRepo.LoadItemState()
	if err != nil {
		a.renderTreasureForm(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error()), Form: defaultTreasureForm()})
		return
	}
	grimoireState, err := a.deps.GrimoireRepo.LoadGrimoireState()
	if err != nil {
		a.renderTreasureForm(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error()), ItemOptions: itemOptions(itemState.Items), Form: defaultTreasureForm()})
		return
	}
	a.renderTreasureForm(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), ItemOptions: itemOptions(itemState.Items), GrimoireOptions: grimoireOptions(grimoireState.Entries), Form: defaultTreasureForm()})
}

func (a App) treasuresEditPage(w http.ResponseWriter, r *http.Request) {
	itemState, err := a.deps.ItemRepo.LoadItemState()
	if err != nil {
		a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error())})
		return
	}
	grimoireState, err := a.deps.GrimoireRepo.LoadGrimoireState()
	if err != nil {
		a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error())})
		return
	}
	state, err := a.deps.TreasureRepo.LoadState()
	if err != nil {
		a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if entry, ok := findEntry(state.Entries, id, func(entry treasures.TreasureEntry) string { return entry.ID }); ok {
		a.renderTreasureForm(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), ItemOptions: itemOptions(itemState.Items), GrimoireOptions: grimoireOptions(grimoireState.Entries), Form: treasureEntryToForm(entry)})
		return
	}
	a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Entries: state.Entries, Notice: errorNotice("Treasure not found.")})
}

func (a App) treasuresSubmit(w http.ResponseWriter, r *http.Request) {
	a.treasuresSave(w, r, false)
}

func (a App) treasuresEditSubmit(w http.ResponseWriter, r *http.Request) {
	a.treasuresSave(w, r, true)
}

func (a App) treasuresSave(w http.ResponseWriter, r *http.Request, editing bool) {
	_ = r.ParseForm()
	itemState, err := a.deps.ItemRepo.LoadItemState()
	if err != nil {
		form := defaultTreasureForm()
		form.IsEditing = editing
		a.renderTreasureForm(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	grimoireState, err := a.deps.GrimoireRepo.LoadGrimoireState()
	if err != nil {
		form := defaultTreasureForm()
		form.IsEditing = editing
		a.renderTreasureForm(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error()), ItemOptions: itemOptions(itemState.Items), Form: form})
		return
	}
	state, err := a.deps.TreasureRepo.LoadState()
	if err != nil {
		form := defaultTreasureForm()
		form.IsEditing = editing
		a.renderTreasureForm(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error()), ItemOptions: itemOptions(itemState.Items), GrimoireOptions: grimoireOptions(grimoireState.Entries), Form: form})
		return
	}
	form, input, parseErrs := parseTreasureForm(r)
	form.IsEditing = editing
	if editing {
		if _, ok := findEntry(state.Entries, form.ID, func(entry treasures.TreasureEntry) string { return entry.ID }); !ok {
			a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Entries: state.Entries, Notice: errorNotice("Treasure not found.")})
			return
		}
	}
	result := treasures.ValidateSave(input, itemIDSet(itemState), grimoireIDSet(grimoireState), a.deps.Now())
	errors := mergeFieldErrors(parseErrs, mapFieldErrors(result.FieldErrors, mapTreasureField))
	if len(errors) > 0 {
		form.FieldErrors = errors
		form.FormError = formErrorText(result.FormError)
		a.renderTreasureForm(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), ItemOptions: itemOptions(itemState.Items), GrimoireOptions: grimoireOptions(grimoireState.Entries), Form: form})
		return
	}
	nextState, mode := common.UpsertEntries(state, *result.Entry, func(entry treasures.TreasureEntry) string { return entry.ID })
	if err := a.deps.TreasureRepo.SaveState(nextState); err != nil {
		a.renderTreasureForm(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error()), ItemOptions: itemOptions(itemState.Items), GrimoireOptions: grimoireOptions(grimoireState.Entries), Form: form})
		return
	}
	notice := successNotice(noticeText("Treasure", mode))
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, "/treasures", notice) {
		return
	}
	a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) treasuresDelete(w http.ResponseWriter, r *http.Request) {
	state, err := a.deps.TreasureRepo.LoadState()
	if err != nil {
		a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.PathValue("id"))
	nextState, ok := common.DeleteEntries(state, id, func(entry treasures.TreasureEntry) string { return entry.ID })
	if !ok {
		a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Entries: state.Entries, Notice: errorNotice("Treasure not found.")})
		return
	}
	if err := a.deps.TreasureRepo.SaveState(nextState); err != nil {
		a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Entries: state.Entries, Notice: errorNotice(err.Error())})
		return
	}
	notice := successNotice("Treasure deleted.")
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, "/treasures", notice) {
		return
	}
	a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) enemiesPage(w http.ResponseWriter, r *http.Request) {
	notice := consumeFlashNotice(w, r)
	state, err := a.deps.EnemyRepo.LoadState()
	if err != nil {
		a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Entries: state.Entries, Notice: notice})
}

func (a App) enemiesNewPage(w http.ResponseWriter, r *http.Request) {
	enemySkillState, err := a.deps.EnemySkillRepo.LoadState()
	if err != nil {
		a.renderEnemyForm(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: defaultEnemyForm(nil)})
		return
	}
	a.renderEnemyForm(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Form: defaultEnemyForm(enemySkillState.Entries)})
}

func (a App) enemiesEditPage(w http.ResponseWriter, r *http.Request) {
	state, err := a.deps.EnemyRepo.LoadState()
	if err != nil {
		a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	enemySkillState, err := a.deps.EnemySkillRepo.LoadState()
	if err != nil {
		a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if entry, ok := findEntry(state.Entries, id, func(entry enemies.EnemyEntry) string { return entry.ID }); ok {
		form := enemyEntryToForm(entry, enemySkillOptions(enemySkillState.Entries))
		a.renderEnemyForm(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Form: form})
		return
	}
	a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Entries: state.Entries, Notice: errorNotice("Enemy not found.")})
}

func (a App) enemiesSubmit(w http.ResponseWriter, r *http.Request) {
	a.enemiesSave(w, r, false)
}

func (a App) enemiesEditSubmit(w http.ResponseWriter, r *http.Request) {
	a.enemiesSave(w, r, true)
}

func (a App) enemiesSave(w http.ResponseWriter, r *http.Request, editing bool) {
	_ = r.ParseForm()
	enemyState, err := a.deps.EnemyRepo.LoadState()
	if err != nil {
		form := defaultEnemyForm(nil)
		form.IsEditing = editing
		a.renderEnemyForm(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	enemySkillState, err := a.deps.EnemySkillRepo.LoadState()
	if err != nil {
		form := defaultEnemyForm(nil)
		form.IsEditing = editing
		a.renderEnemyForm(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	itemState, err := a.deps.ItemRepo.LoadItemState()
	if err != nil {
		form := defaultEnemyForm(enemySkillState.Entries)
		form.IsEditing = editing
		a.renderEnemyForm(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	grimoireState, err := a.deps.GrimoireRepo.LoadGrimoireState()
	if err != nil {
		form := defaultEnemyForm(enemySkillState.Entries)
		form.IsEditing = editing
		a.renderEnemyForm(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	form, input, parseErrs := parseEnemyForm(r, enemySkillState.Entries)
	form.IsEditing = editing
	if editing {
		if _, ok := findEntry(enemyState.Entries, form.ID, func(entry enemies.EnemyEntry) string { return entry.ID }); !ok {
			a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Entries: enemyState.Entries, Notice: errorNotice("Enemy not found.")})
			return
		}
	}
	result := enemies.ValidateSave(input, enemySkillIDSet(enemySkillState), itemIDSet(itemState), grimoireIDSet(grimoireState), a.deps.Now())
	errors := mergeFieldErrors(parseErrs, mapFieldErrors(result.FieldErrors, mapEnemyField))
	if len(errors) > 0 {
		form.FieldErrors = errors
		form.FormError = formErrorText(result.FormError)
		a.renderEnemyForm(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Form: form})
		return
	}
	nextState, mode := common.UpsertEntries(enemyState, *result.Entry, func(entry enemies.EnemyEntry) string { return entry.ID })
	if err := a.deps.EnemyRepo.SaveState(nextState); err != nil {
		a.renderEnemyForm(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	notice := successNotice(noticeText("Enemy", mode))
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, "/enemies", notice) {
		return
	}
	a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) enemiesDelete(w http.ResponseWriter, r *http.Request) {
	state, err := a.deps.EnemyRepo.LoadState()
	if err != nil {
		a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.PathValue("id"))
	nextState, ok := common.DeleteEntries(state, id, func(entry enemies.EnemyEntry) string { return entry.ID })
	if !ok {
		a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Entries: state.Entries, Notice: errorNotice("Enemy not found.")})
		return
	}
	if err := a.deps.EnemyRepo.SaveState(nextState); err != nil {
		a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Entries: state.Entries, Notice: errorNotice(err.Error())})
		return
	}
	notice := successNotice("Enemy deleted.")
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, "/enemies", notice) {
		return
	}
	a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) saveExport(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	currentPath := normalizeScreenPath(r.Form.Get("currentPath"))
	result := application.NewService(a.cfg, a.deps).ExportDatapack()
	if !result.OK {
		a.renderSaveResponse(w, r, currentPath, errorNotice(result.Message))
		return
	}
	message := fmt.Sprintf("Exported %d files to %s.", result.Generated.TotalFiles, result.OutputRoot)
	setToast(w, message)
	a.renderSaveResponse(w, r, currentPath, successNotice(message))
}

func (a App) renderSaveResponse(w http.ResponseWriter, r *http.Request, currentPath string, notice *webui.Notice) {
	if isHX(r) {
		a.renderComponent(w, views.NoticeBox(notice))
		return
	}
	a.renderScreen(w, r, currentPath, notice)
}

func (a App) renderScreen(w http.ResponseWriter, r *http.Request, currentPath string, notice *webui.Notice) {
	switch normalizeScreenPath(currentPath) {
	case "/grimoire":
		state, err := a.deps.GrimoireRepo.LoadGrimoireState()
		if err != nil {
			a.renderGrimoire(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Notice: errorNotice(err.Error()), Form: defaultGrimoireForm(nil)})
			return
		}
		a.renderGrimoire(w, r, webui.GrimoirePageData{
			Meta:    grimoireMeta(),
			Notice:  notice,
			Entries: state.Entries,
			Form:    defaultGrimoireForm(state.Entries),
		})
	case "/skills":
		itemState, err := a.deps.ItemRepo.LoadItemState()
		if err != nil {
			a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Notice: errorNotice(err.Error()), Form: defaultSkillForm(nil)})
			return
		}
		state, err := a.deps.SkillRepo.LoadState()
		if err != nil {
			a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Notice: errorNotice(err.Error()), Form: defaultSkillForm(itemState.Items)})
			return
		}
		a.renderSkills(w, r, webui.SkillsPageData{
			Meta:    skillsMeta(),
			Notice:  notice,
			Entries: state.Entries,
			Form:    defaultSkillForm(itemState.Items),
		})
	case "/enemy-skills":
		state, err := a.deps.EnemySkillRepo.LoadState()
		if err != nil {
			a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Notice: errorNotice(err.Error()), Form: defaultEnemySkillForm()})
			return
		}
		a.renderEnemySkills(w, r, webui.EnemySkillsPageData{
			Meta:    enemySkillsMeta(),
			Notice:  notice,
			Entries: state.Entries,
			Form:    defaultEnemySkillForm(),
		})
	case "/treasures":
		itemState, err := a.deps.ItemRepo.LoadItemState()
		if err != nil {
			a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error()), Form: defaultTreasureForm()})
			return
		}
		grimoireState, err := a.deps.GrimoireRepo.LoadGrimoireState()
		if err != nil {
			a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error()), Form: defaultTreasureForm()})
			return
		}
		state, err := a.deps.TreasureRepo.LoadState()
		if err != nil {
			a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error()), Form: defaultTreasureForm()})
			return
		}
		a.renderTreasures(w, r, webui.TreasuresPageData{
			Meta:            treasuresMeta(),
			Notice:          notice,
			Entries:         state.Entries,
			ItemOptions:     itemOptions(itemState.Items),
			GrimoireOptions: grimoireOptions(grimoireState.Entries),
			Form:            defaultTreasureForm(),
		})
	case "/enemies":
		state, err := a.deps.EnemyRepo.LoadState()
		if err != nil {
			a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: defaultEnemyForm(nil)})
			return
		}
		enemySkillState, err := a.deps.EnemySkillRepo.LoadState()
		if err != nil {
			a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: defaultEnemyForm(nil)})
			return
		}
		a.renderEnemies(w, r, webui.EnemiesPageData{
			Meta:    enemiesMeta(),
			Notice:  notice,
			Entries: state.Entries,
			Form:    defaultEnemyForm(enemySkillState.Entries),
		})
	case "/items":
		fallthrough
	default:
		state, err := a.deps.ItemRepo.LoadItemState()
		if err != nil {
			a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error()), Form: defaultItemForm()})
			return
		}
		a.renderItems(w, r, webui.ItemsPageData{
			Meta:    itemMeta(),
			Notice:  notice,
			Entries: state.Items,
			Form:    defaultItemForm(),
		})
	}
}

func (a App) renderItems(w http.ResponseWriter, r *http.Request, data webui.ItemsPageData) {
	if isHX(r) {
		a.renderComponent(w, views.ItemsShell(data))
		return
	}
	a.renderComponent(w, views.ItemsPage(data))
}

func (a App) renderItemForm(w http.ResponseWriter, r *http.Request, data webui.ItemsPageData) {
	if isHX(r) {
		a.renderComponent(w, views.ItemFormShell(data))
		return
	}
	a.renderComponent(w, views.ItemFormPage(data))
}

func (a App) renderGrimoire(w http.ResponseWriter, r *http.Request, data webui.GrimoirePageData) {
	if isHX(r) {
		a.renderComponent(w, views.GrimoireShell(data))
		return
	}
	a.renderComponent(w, views.GrimoirePage(data))
}

func (a App) renderGrimoireForm(w http.ResponseWriter, r *http.Request, data webui.GrimoirePageData) {
	if isHX(r) {
		a.renderComponent(w, views.GrimoireFormShell(data))
		return
	}
	a.renderComponent(w, views.GrimoireFormPage(data))
}

func (a App) renderSkills(w http.ResponseWriter, r *http.Request, data webui.SkillsPageData) {
	if isHX(r) {
		a.renderComponent(w, views.SkillsShell(data))
		return
	}
	a.renderComponent(w, views.SkillsPage(data))
}

func (a App) renderSkillForm(w http.ResponseWriter, r *http.Request, data webui.SkillsPageData) {
	if isHX(r) {
		a.renderComponent(w, views.SkillFormShell(data))
		return
	}
	a.renderComponent(w, views.SkillFormPage(data))
}

func (a App) renderEnemySkills(w http.ResponseWriter, r *http.Request, data webui.EnemySkillsPageData) {
	if isHX(r) {
		a.renderComponent(w, views.EnemySkillsShell(data))
		return
	}
	a.renderComponent(w, views.EnemySkillsPage(data))
}

func (a App) renderEnemySkillForm(w http.ResponseWriter, r *http.Request, data webui.EnemySkillsPageData) {
	if isHX(r) {
		a.renderComponent(w, views.EnemySkillFormShell(data))
		return
	}
	a.renderComponent(w, views.EnemySkillFormPage(data))
}

func (a App) renderTreasures(w http.ResponseWriter, r *http.Request, data webui.TreasuresPageData) {
	if isHX(r) {
		a.renderComponent(w, views.TreasuresShell(data))
		return
	}
	a.renderComponent(w, views.TreasuresPage(data))
}

func (a App) renderTreasureForm(w http.ResponseWriter, r *http.Request, data webui.TreasuresPageData) {
	if isHX(r) {
		a.renderComponent(w, views.TreasureFormShell(data))
		return
	}
	a.renderComponent(w, views.TreasureFormPage(data))
}

func (a App) renderEnemies(w http.ResponseWriter, r *http.Request, data webui.EnemiesPageData) {
	if isHX(r) {
		a.renderComponent(w, views.EnemiesShell(data))
		return
	}
	a.renderComponent(w, views.EnemiesPage(data))
}

func (a App) renderEnemyForm(w http.ResponseWriter, r *http.Request, data webui.EnemiesPageData) {
	if isHX(r) {
		a.renderComponent(w, views.EnemyFormShell(data))
		return
	}
	a.renderComponent(w, views.EnemyFormPage(data))
}

func (a App) renderComponent(w http.ResponseWriter, component templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(context.Background(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func itemMeta() webui.PageMeta {
	return webui.PageMeta{Title: "Items", CurrentPath: "/items", Description: "アイテム出力を作成・管理します。複雑な NBT 項目はこの移行段階ではテキストのまま扱います。"}
}

func grimoireMeta() webui.PageMeta {
	return webui.PageMeta{Title: "Grimoire", CurrentPath: "/grimoire", Description: "呪文エントリを管理します。Variants は `cast,cost` を1行ずつ入力します。"}
}

func skillsMeta() webui.PageMeta {
	return webui.PageMeta{Title: "Skills", CurrentPath: "/skills", Description: "スキルの script と参照する item entry を管理します。"}
}

func enemySkillsMeta() webui.PageMeta {
	return webui.PageMeta{Title: "Enemy Skills", CurrentPath: "/enemy-skills", Description: "再利用可能な enemy-skill script を管理します。trigger と cooldown は任意です。"}
}

func treasuresMeta() webui.PageMeta {
	return webui.PageMeta{Title: "Treasures", CurrentPath: "/treasures", Description: "treasure loot pools を管理します。各行は `kind,refId,weight,countMin,countMax` 形式です。"}
}

func enemiesMeta() webui.PageMeta {
	return webui.PageMeta{Title: "Enemies", CurrentPath: "/enemies", Description: "enemy stats、spawn rules、参照する enemy skills を管理します。"}
}

func defaultItemForm() webui.ItemFormData {
	return webui.ItemFormData{
		ID:                  newUUID(),
		ItemID:              "minecraft:stone",
		Count:               "1",
		FieldErrors:         map[string]string{},
		CustomName:          "",
		Lore:                "",
		Enchantments:        "",
		CustomModelData:     "",
		RepairCost:          "",
		HideFlags:           "",
		PotionID:            "",
		CustomPotionColor:   "",
		CustomPotionEffects: "",
		AttributeModifiers:  "",
		CustomNBT:           "",
	}
}

func itemEntryToForm(entry items.ItemEntry) webui.ItemFormData {
	return webui.ItemFormData{
		ID:                  entry.ID,
		ItemID:              entry.ItemID,
		Count:               strconv.Itoa(entry.Count),
		CustomName:          entry.CustomName,
		Lore:                entry.Lore,
		Enchantments:        entry.Enchantments,
		Unbreakable:         entry.Unbreakable,
		CustomModelData:     entry.CustomModelData,
		RepairCost:          entry.RepairCost,
		HideFlags:           entry.HideFlags,
		PotionID:            entry.PotionID,
		CustomPotionColor:   entry.CustomPotionColor,
		CustomPotionEffects: entry.CustomPotionEffects,
		AttributeModifiers:  entry.AttributeModifiers,
		CustomNBT:           entry.CustomNBT,
		FieldErrors:         map[string]string{},
		IsEditing:           true,
	}
}

func defaultGrimoireForm(entries []grimoire.GrimoireEntry) webui.GrimoireFormData {
	nextCastID := 10
	for _, entry := range entries {
		if entry.CastID >= nextCastID {
			nextCastID = entry.CastID + 10
		}
	}
	return webui.GrimoireFormData{
		ID:           newUUID(),
		CastID:       strconv.Itoa(nextCastID),
		VariantsText: "0,0",
		FieldErrors:  map[string]string{},
	}
}

func grimoireEntryToForm(entry grimoire.GrimoireEntry) webui.GrimoireFormData {
	lines := make([]string, 0, len(entry.Variants))
	for _, variant := range entry.Variants {
		lines = append(lines, fmt.Sprintf("%d,%d", variant.Cast, variant.Cost))
	}
	return webui.GrimoireFormData{
		ID:           entry.ID,
		CastID:       strconv.Itoa(entry.CastID),
		Script:       entry.Script,
		Title:        entry.Title,
		Description:  entry.Description,
		VariantsText: strings.Join(lines, "\n"),
		FieldErrors:  map[string]string{},
		IsEditing:    true,
	}
}

func defaultSkillForm(items []items.ItemEntry) webui.SkillFormData {
	form := webui.SkillFormData{
		ID:          newUUID(),
		ItemOptions: itemOptions(items),
		FieldErrors: map[string]string{},
	}
	if len(items) > 0 {
		form.ItemID = items[0].ID
	}
	return form
}

func skillEntryToForm(entry skills.SkillEntry, options []webui.ReferenceOption) webui.SkillFormData {
	return webui.SkillFormData{
		ID:          entry.ID,
		Name:        entry.Name,
		Script:      entry.Script,
		ItemID:      entry.ItemID,
		ItemOptions: options,
		FieldErrors: map[string]string{},
		IsEditing:   true,
	}
}

func defaultEnemySkillForm() webui.EnemySkillFormData {
	return webui.EnemySkillFormData{
		ID:          newUUID(),
		FieldErrors: map[string]string{},
	}
}

func enemySkillEntryToForm(entry enemyskills.EnemySkillEntry) webui.EnemySkillFormData {
	form := webui.EnemySkillFormData{
		ID:          entry.ID,
		Name:        entry.Name,
		Script:      entry.Script,
		FieldErrors: map[string]string{},
		IsEditing:   true,
	}
	if entry.Cooldown != nil {
		form.Cooldown = strconv.FormatFloat(*entry.Cooldown, 'f', -1, 64)
	}
	if entry.Trigger != nil {
		form.Trigger = string(*entry.Trigger)
	}
	return form
}

func defaultTreasureForm() webui.TreasureFormData {
	return webui.TreasureFormData{
		ID:            newUUID(),
		LootPoolsText: "item,,1,1,1",
		FieldErrors:   map[string]string{},
	}
}

func treasureEntryToForm(entry treasures.TreasureEntry) webui.TreasureFormData {
	lines := make([]string, 0, len(entry.LootPools))
	for _, pool := range entry.LootPools {
		countMin := ""
		countMax := ""
		if pool.CountMin != nil {
			countMin = strconv.FormatFloat(*pool.CountMin, 'f', -1, 64)
		}
		if pool.CountMax != nil {
			countMax = strconv.FormatFloat(*pool.CountMax, 'f', -1, 64)
		}
		lines = append(lines, strings.Join([]string{
			pool.Kind,
			pool.RefID,
			strconv.FormatFloat(pool.Weight, 'f', -1, 64),
			countMin,
			countMax,
		}, ","))
	}
	return webui.TreasureFormData{
		ID:            entry.ID,
		Name:          entry.Name,
		LootPoolsText: strings.Join(lines, "\n"),
		FieldErrors:   map[string]string{},
		IsEditing:     true,
	}
}

func defaultEnemyForm(entries []enemyskills.EnemySkillEntry) webui.EnemyFormData {
	form := webui.EnemyFormData{
		ID:                newUUID(),
		HP:                "20",
		OriginX:           "0",
		OriginY:           "0",
		OriginZ:           "0",
		DistanceMin:       "0",
		DistanceMax:       "16",
		FieldErrors:       map[string]string{},
		EnemySkillOptions: enemySkillOptions(entries),
	}
	form.DropTableID = form.ID
	return form
}

func enemyEntryToForm(entry enemies.EnemyEntry, options []webui.ReferenceOption) webui.EnemyFormData {
	form := webui.EnemyFormData{
		ID:                entry.ID,
		Name:              entry.Name,
		HP:                strconv.FormatFloat(entry.HP, 'f', -1, 64),
		DropTableID:       entry.DropTableID,
		EnemySkillIDs:     append([]string{}, entry.EnemySkillIDs...),
		EnemySkillOptions: options,
		OriginX:           strconv.FormatFloat(entry.SpawnRule.Origin.X, 'f', -1, 64),
		OriginY:           strconv.FormatFloat(entry.SpawnRule.Origin.Y, 'f', -1, 64),
		OriginZ:           strconv.FormatFloat(entry.SpawnRule.Origin.Z, 'f', -1, 64),
		DistanceMin:       strconv.FormatFloat(entry.SpawnRule.Distance.Min, 'f', -1, 64),
		DistanceMax:       strconv.FormatFloat(entry.SpawnRule.Distance.Max, 'f', -1, 64),
		FieldErrors:       map[string]string{},
		IsEditing:         true,
	}
	if entry.Attack != nil {
		form.Attack = strconv.FormatFloat(*entry.Attack, 'f', -1, 64)
	}
	if entry.Defense != nil {
		form.Defense = strconv.FormatFloat(*entry.Defense, 'f', -1, 64)
	}
	if entry.MoveSpeed != nil {
		form.MoveSpeed = strconv.FormatFloat(*entry.MoveSpeed, 'f', -1, 64)
	}
	if axis := entry.SpawnRule.AxisBounds; axis != nil {
		if axis.XMin != nil {
			form.XMin = strconv.FormatFloat(*axis.XMin, 'f', -1, 64)
		}
		if axis.XMax != nil {
			form.XMax = strconv.FormatFloat(*axis.XMax, 'f', -1, 64)
		}
		if axis.YMin != nil {
			form.YMin = strconv.FormatFloat(*axis.YMin, 'f', -1, 64)
		}
		if axis.YMax != nil {
			form.YMax = strconv.FormatFloat(*axis.YMax, 'f', -1, 64)
		}
		if axis.ZMin != nil {
			form.ZMin = strconv.FormatFloat(*axis.ZMin, 'f', -1, 64)
		}
		if axis.ZMax != nil {
			form.ZMax = strconv.FormatFloat(*axis.ZMax, 'f', -1, 64)
		}
	}
	return form
}

func parseItemForm(r *http.Request) (webui.ItemFormData, items.SaveInput, map[string]string) {
	form := defaultItemForm()
	form.ID = valueOrDefault(r.Form.Get("id"), newUUID())
	form.ItemID = strings.TrimSpace(r.Form.Get("itemId"))
	form.Count = strings.TrimSpace(r.Form.Get("count"))
	form.CustomName = r.Form.Get("customName")
	form.Lore = r.Form.Get("lore")
	form.Enchantments = r.Form.Get("enchantments")
	form.Unbreakable = r.Form.Get("unbreakable") != ""
	form.CustomModelData = r.Form.Get("customModelData")
	form.RepairCost = r.Form.Get("repairCost")
	form.HideFlags = r.Form.Get("hideFlags")
	form.PotionID = r.Form.Get("potionId")
	form.CustomPotionColor = r.Form.Get("customPotionColor")
	form.CustomPotionEffects = r.Form.Get("customPotionEffects")
	form.AttributeModifiers = r.Form.Get("attributeModifiers")
	form.CustomNBT = r.Form.Get("customNbt")
	errs := map[string]string{}
	input := items.SaveInput{
		ID:                  form.ID,
		ItemID:              form.ItemID,
		Count:               parseRequiredInt(errs, "count", form.Count),
		CustomName:          form.CustomName,
		Lore:                form.Lore,
		Enchantments:        form.Enchantments,
		Unbreakable:         form.Unbreakable,
		CustomModelData:     form.CustomModelData,
		RepairCost:          form.RepairCost,
		HideFlags:           form.HideFlags,
		PotionID:            form.PotionID,
		CustomPotionColor:   form.CustomPotionColor,
		CustomPotionEffects: form.CustomPotionEffects,
		AttributeModifiers:  form.AttributeModifiers,
		CustomNBT:           form.CustomNBT,
	}
	return form, input, errs
}

func parseGrimoireForm(r *http.Request) (webui.GrimoireFormData, grimoire.SaveInput, map[string]string) {
	form := defaultGrimoireForm(nil)
	form.ID = valueOrDefault(r.Form.Get("id"), newUUID())
	form.CastID = strings.TrimSpace(r.Form.Get("castid"))
	form.Script = r.Form.Get("script")
	form.Title = r.Form.Get("title")
	form.Description = r.Form.Get("description")
	form.VariantsText = r.Form.Get("variantsText")
	errs := map[string]string{}
	input := grimoire.SaveInput{
		ID:          form.ID,
		CastID:      parseRequiredInt(errs, "castid", form.CastID),
		Script:      form.Script,
		Title:       form.Title,
		Description: form.Description,
		Variants:    parseVariants(errs, form.VariantsText),
	}
	return form, input, errs
}

func parseSkillForm(r *http.Request, itemEntries []items.ItemEntry) (webui.SkillFormData, skills.SaveInput, map[string]string) {
	form := defaultSkillForm(itemEntries)
	form.ID = valueOrDefault(r.Form.Get("id"), newUUID())
	form.Name = r.Form.Get("name")
	form.Script = r.Form.Get("script")
	form.ItemID = strings.TrimSpace(r.Form.Get("itemId"))
	errs := map[string]string{}
	input := skills.SaveInput{
		ID:     form.ID,
		Name:   form.Name,
		Script: form.Script,
		ItemID: form.ItemID,
	}
	return form, input, errs
}

func parseEnemySkillForm(r *http.Request) (webui.EnemySkillFormData, enemyskills.SaveInput, map[string]string) {
	form := defaultEnemySkillForm()
	form.ID = valueOrDefault(r.Form.Get("id"), newUUID())
	form.Name = r.Form.Get("name")
	form.Script = r.Form.Get("script")
	form.Cooldown = strings.TrimSpace(r.Form.Get("cooldown"))
	form.Trigger = strings.TrimSpace(r.Form.Get("trigger"))
	errs := map[string]string{}
	input := enemyskills.SaveInput{
		ID:       form.ID,
		Name:     form.Name,
		Script:   form.Script,
		Cooldown: parseOptionalFloat(errs, "cooldown", form.Cooldown),
		Trigger:  form.Trigger,
	}
	return form, input, errs
}

func parseTreasureForm(r *http.Request) (webui.TreasureFormData, treasures.SaveInput, map[string]string) {
	form := defaultTreasureForm()
	form.ID = valueOrDefault(r.Form.Get("id"), newUUID())
	form.Name = r.Form.Get("name")
	form.LootPoolsText = r.Form.Get("lootPoolsText")
	errs := map[string]string{}
	input := treasures.SaveInput{
		ID:        form.ID,
		Name:      form.Name,
		LootPools: parseTreasurePools(errs, form.LootPoolsText),
	}
	return form, input, errs
}

func parseEnemyForm(r *http.Request, enemySkillEntries []enemyskills.EnemySkillEntry) (webui.EnemyFormData, enemies.SaveInput, map[string]string) {
	form := defaultEnemyForm(enemySkillEntries)
	form.ID = valueOrDefault(r.Form.Get("id"), newUUID())
	form.Name = r.Form.Get("name")
	form.HP = strings.TrimSpace(r.Form.Get("hp"))
	form.Attack = strings.TrimSpace(r.Form.Get("attack"))
	form.Defense = strings.TrimSpace(r.Form.Get("defense"))
	form.MoveSpeed = strings.TrimSpace(r.Form.Get("moveSpeed"))
	form.DropTableID = strings.TrimSpace(r.Form.Get("dropTableId"))
	form.EnemySkillIDs = append([]string{}, r.Form["enemySkillIds"]...)
	form.OriginX = strings.TrimSpace(r.Form.Get("originX"))
	form.OriginY = strings.TrimSpace(r.Form.Get("originY"))
	form.OriginZ = strings.TrimSpace(r.Form.Get("originZ"))
	form.DistanceMin = strings.TrimSpace(r.Form.Get("distanceMin"))
	form.DistanceMax = strings.TrimSpace(r.Form.Get("distanceMax"))
	form.XMin = strings.TrimSpace(r.Form.Get("xMin"))
	form.XMax = strings.TrimSpace(r.Form.Get("xMax"))
	form.YMin = strings.TrimSpace(r.Form.Get("yMin"))
	form.YMax = strings.TrimSpace(r.Form.Get("yMax"))
	form.ZMin = strings.TrimSpace(r.Form.Get("zMin"))
	form.ZMax = strings.TrimSpace(r.Form.Get("zMax"))
	errs := map[string]string{}

	axisBounds := &enemies.AxisBounds{
		XMin: parseOptionalFloat(errs, "xMin", form.XMin),
		XMax: parseOptionalFloat(errs, "xMax", form.XMax),
		YMin: parseOptionalFloat(errs, "yMin", form.YMin),
		YMax: parseOptionalFloat(errs, "yMax", form.YMax),
		ZMin: parseOptionalFloat(errs, "zMin", form.ZMin),
		ZMax: parseOptionalFloat(errs, "zMax", form.ZMax),
	}
	if axisBounds.XMin == nil && axisBounds.XMax == nil && axisBounds.YMin == nil && axisBounds.YMax == nil && axisBounds.ZMin == nil && axisBounds.ZMax == nil {
		axisBounds = nil
	}

	input := enemies.SaveInput{
		ID:            form.ID,
		Name:          form.Name,
		HP:            parseRequiredFloat(errs, "hp", form.HP),
		Attack:        parseOptionalFloat(errs, "attack", form.Attack),
		Defense:       parseOptionalFloat(errs, "defense", form.Defense),
		MoveSpeed:     parseOptionalFloat(errs, "moveSpeed", form.MoveSpeed),
		DropTableID:   form.DropTableID,
		EnemySkillIDs: append([]string{}, form.EnemySkillIDs...),
		SpawnRule: enemies.SpawnRule{
			Origin: enemies.Vec3{
				X: parseRequiredFloat(errs, "originX", form.OriginX),
				Y: parseRequiredFloat(errs, "originY", form.OriginY),
				Z: parseRequiredFloat(errs, "originZ", form.OriginZ),
			},
			Distance: enemies.Distance{
				Min: parseRequiredFloat(errs, "distanceMin", form.DistanceMin),
				Max: parseRequiredFloat(errs, "distanceMax", form.DistanceMax),
			},
			AxisBounds: axisBounds,
		},
	}
	return form, input, errs
}

func parseVariants(errs map[string]string, value string) []grimoire.Variant {
	lines := compactLines(value)
	out := make([]grimoire.Variant, 0, len(lines))
	for _, line := range lines {
		parts := splitCSV(line, 2)
		if len(parts) != 2 {
			errs["variantsText"] = "Each variant line must be `cast,cost`."
			return nil
		}
		cast, okCast := parseIntText(parts[0])
		cost, okCost := parseIntText(parts[1])
		if !okCast || !okCost {
			errs["variantsText"] = "Each variant line must contain numeric cast and cost."
			return nil
		}
		out = append(out, grimoire.Variant{Cast: cast, Cost: cost})
	}
	return out
}

func parseTreasurePools(errs map[string]string, value string) []treasures.DropRef {
	lines := compactLines(value)
	out := make([]treasures.DropRef, 0, len(lines))
	for _, line := range lines {
		parts := splitCSV(line, 5)
		if len(parts) < 3 {
			errs["lootPoolsText"] = "Each loot line must be `kind,refId,weight,countMin,countMax`."
			return nil
		}
		weight, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			errs["lootPoolsText"] = "Weight must be numeric."
			return nil
		}
		out = append(out, treasures.DropRef{
			Kind:     parts[0],
			RefID:    parts[1],
			Weight:   weight,
			CountMin: parseOptionalFloat(errs, "lootPoolsText", valueOrIndex(parts, 3)),
			CountMax: parseOptionalFloat(errs, "lootPoolsText", valueOrIndex(parts, 4)),
		})
		if errs["lootPoolsText"] != "" {
			errs["lootPoolsText"] = "Count values must be numeric when provided."
			return nil
		}
	}
	return out
}

func itemOptions(entries []items.ItemEntry) []webui.ReferenceOption {
	options := make([]webui.ReferenceOption, 0, len(entries))
	for _, entry := range entries {
		options = append(options, webui.ReferenceOption{ID: entry.ID, Label: entry.ItemID})
	}
	return options
}

func grimoireOptions(entries []grimoire.GrimoireEntry) []webui.ReferenceOption {
	options := make([]webui.ReferenceOption, 0, len(entries))
	for _, entry := range entries {
		options = append(options, webui.ReferenceOption{ID: entry.ID, Label: entry.Title})
	}
	return options
}

func enemySkillOptions(entries []enemyskills.EnemySkillEntry) []webui.ReferenceOption {
	options := make([]webui.ReferenceOption, 0, len(entries))
	for _, entry := range entries {
		options = append(options, webui.ReferenceOption{ID: entry.ID, Label: entry.Name})
	}
	return options
}

func itemIDSet(state items.ItemState) map[string]struct{} {
	return toIDSet(state.Items, func(entry items.ItemEntry) string { return entry.ID })
}

func grimoireIDSet(state grimoire.GrimoireState) map[string]struct{} {
	return toIDSet(state.Entries, func(entry grimoire.GrimoireEntry) string { return entry.ID })
}

func enemySkillIDSet(state common.EntryState[enemyskills.EnemySkillEntry]) map[string]struct{} {
	return toIDSet(state.Entries, func(entry enemyskills.EnemySkillEntry) string { return entry.ID })
}

func toIDSet[T any](entries []T, idOf func(T) string) map[string]struct{} {
	out := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		id := strings.TrimSpace(idOf(entry))
		if id != "" {
			out[id] = struct{}{}
		}
	}
	return out
}

func compactLines(value string) []string {
	raw := strings.Split(value, "\n")
	out := make([]string, 0, len(raw))
	for _, line := range raw {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func splitCSV(value string, max int) []string {
	parts := strings.Split(value, ",")
	if len(parts) > max {
		return nil
	}
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		out = append(out, strings.TrimSpace(part))
	}
	return out
}

func parseRequiredInt(errs map[string]string, key, value string) int {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return 0
	}
	parsed, err := strconv.Atoi(trimmed)
	if err != nil {
		errs[key] = "Must be a number."
		return 0
	}
	return parsed
}

func parseRequiredFloat(errs map[string]string, key, value string) float64 {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return 0
	}
	parsed, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		errs[key] = "Must be a number."
		return 0
	}
	return parsed
}

func parseOptionalFloat(errs map[string]string, key, value string) *float64 {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	parsed, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		errs[key] = "Must be a number."
		return nil
	}
	return &parsed
}

func parseIntText(value string) (int, bool) {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	return parsed, err == nil
}

func mergeFieldErrors(primary, secondary map[string]string) map[string]string {
	if len(primary) == 0 && len(secondary) == 0 {
		return map[string]string{}
	}
	out := map[string]string{}
	for key, value := range secondary {
		out[key] = value
	}
	for key, value := range primary {
		out[key] = value
	}
	return out
}

func mapFieldErrors(errs common.FieldErrors, mapField func(string) string) map[string]string {
	out := map[string]string{}
	for key, value := range errs {
		mapped := mapField(key)
		if mapped == "" {
			mapped = key
		}
		if _, exists := out[mapped]; !exists {
			out[mapped] = value
		}
	}
	return out
}

func mapGrimoireField(key string) string {
	if strings.HasPrefix(key, "variants.") {
		return "variantsText"
	}
	return key
}

func mapTreasureField(key string) string {
	if strings.HasPrefix(key, "lootPools.") {
		return "lootPoolsText"
	}
	return key
}

func mapEnemyField(key string) string {
	switch key {
	case "spawnRule.origin.x":
		return "originX"
	case "spawnRule.origin.y":
		return "originY"
	case "spawnRule.origin.z":
		return "originZ"
	case "spawnRule.distance.min":
		return "distanceMin"
	case "spawnRule.distance.max":
		return "distanceMax"
	case "spawnRule.axisBounds.xMin":
		return "xMin"
	case "spawnRule.axisBounds.xMax":
		return "xMax"
	case "spawnRule.axisBounds.yMin":
		return "yMin"
	case "spawnRule.axisBounds.yMax":
		return "yMax"
	case "spawnRule.axisBounds.zMin":
		return "zMin"
	case "spawnRule.axisBounds.zMax":
		return "zMax"
	}
	if strings.HasPrefix(key, "enemySkillIds.") {
		return "enemySkillIds"
	}
	return key
}

func noticeText(label string, mode common.SaveMode) string {
	if mode == common.SaveModeUpdated {
		return label + " updated."
	}
	return label + " created."
}

func successNotice(text string) *webui.Notice {
	return &webui.Notice{Kind: "success", Text: text}
}

func errorNotice(text string) *webui.Notice {
	return &webui.Notice{Kind: "error", Text: text}
}

func formErrorText(value string) string {
	if strings.TrimSpace(value) == "" {
		return "Validation failed. Fix the highlighted fields."
	}
	return value
}

const flashNoticeCookieName = "tools2-flash-notice"

func redirectWithNotice(w http.ResponseWriter, r *http.Request, path string, notice *webui.Notice) bool {
	if isHX(r) {
		return false
	}
	setFlashNotice(w, notice)
	http.Redirect(w, r, path, http.StatusSeeOther)
	return true
}

func setFlashNotice(w http.ResponseWriter, notice *webui.Notice) {
	if notice == nil || strings.TrimSpace(notice.Text) == "" {
		return
	}
	payload := notice.Kind + "\n" + notice.Text
	http.SetCookie(w, &http.Cookie{
		Name:     flashNoticeCookieName,
		Value:    base64.RawURLEncoding.EncodeToString([]byte(payload)),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func consumeFlashNotice(w http.ResponseWriter, r *http.Request) *webui.Notice {
	cookie, err := r.Cookie(flashNoticeCookieName)
	if err != nil {
		return nil
	}
	http.SetCookie(w, &http.Cookie{
		Name:     flashNoticeCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})
	decoded, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil
	}
	parts := strings.SplitN(string(decoded), "\n", 2)
	if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
		return nil
	}
	return &webui.Notice{Kind: parts[0], Text: parts[1]}
}

func normalizeScreenPath(value string) string {
	switch strings.TrimSpace(value) {
	case "/items", "/grimoire", "/skills", "/enemy-skills", "/treasures", "/enemies":
		return strings.TrimSpace(value)
	default:
		return "/items"
	}
}

func isHX(r *http.Request) bool {
	return strings.EqualFold(r.Header.Get("HX-Request"), "true")
}

func setToast(w http.ResponseWriter, text string) {
	if strings.TrimSpace(text) != "" {
		w.Header().Set("HX-Trigger", text)
	}
}

func valueOrDefault(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func valueOrIndex(parts []string, index int) string {
	if index >= 0 && index < len(parts) {
		return parts[index]
	}
	return ""
}

func newUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uint32(b[0])<<24|uint32(b[1])<<16|uint32(b[2])<<8|uint32(b[3]),
		uint16(b[4])<<8|uint16(b[5]),
		uint16(b[6])<<8|uint16(b[7]),
		uint16(b[8])<<8|uint16(b[9]),
		uint64(b[10])<<40|uint64(b[11])<<32|uint64(b[12])<<24|uint64(b[13])<<16|uint64(b[14])<<8|uint64(b[15]),
	)
}

func findEntry[T any](entries []T, id string, idOf func(T) string) (T, bool) {
	var zero T
	for _, entry := range entries {
		if idOf(entry) == id {
			return entry, true
		}
	}
	return zero, false
}
