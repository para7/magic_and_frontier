package web

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"

	"tools2/app/internal/application"
	"tools2/app/internal/domain/common"
	"tools2/app/internal/web/ui"
	"tools2/app/internal/web/views"
)

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

func (a App) renderSaveResponse(w http.ResponseWriter, r *http.Request, currentPath string, notice *ui.Notice) {
	if isHX(r) {
		a.renderComponent(w, views.NoticeBox(notice))
		return
	}
	a.renderScreen(w, r, currentPath, notice)
}

func (a App) renderScreen(w http.ResponseWriter, r *http.Request, currentPath string, notice *ui.Notice) {
	switch normalizeScreenPath(currentPath) {
	case "/grimoire":
		state, err := a.loadGrimoireStateFromMaster()
		if err != nil {
			a.renderGrimoire(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Notice: errorNotice(err.Error()), Form: defaultGrimoireForm(nil)})
			return
		}
		a.renderGrimoire(w, r, ui.GrimoirePageData{
			Meta:    grimoireMeta(),
			Notice:  notice,
			Entries: state.Entries,
			Form:    defaultGrimoireForm(state.Entries),
		})
	case "/skills":
		state, err := a.loadSkillStateFromMaster()
		if err != nil {
			a.renderSkills(w, r, ui.SkillsPageData{Meta: skillsMeta(), Notice: errorNotice(err.Error()), Form: defaultSkillForm()})
			return
		}
		a.renderSkills(w, r, ui.SkillsPageData{
			Meta:    skillsMeta(),
			Notice:  notice,
			Entries: state.Entries,
			Form:    defaultSkillForm(),
		})
	case "/enemy-skills":
		state, err := a.loadEnemySkillStateFromMaster()
		if err != nil {
			a.renderEnemySkills(w, r, ui.EnemySkillsPageData{Meta: enemySkillsMeta(), Notice: errorNotice(err.Error()), Form: defaultEnemySkillForm()})
			return
		}
		a.renderEnemySkills(w, r, ui.EnemySkillsPageData{
			Meta:    enemySkillsMeta(),
			Notice:  notice,
			Entries: state.Entries,
			Form:    defaultEnemySkillForm(),
		})
	case "/treasures":
		data, err := a.treasuresPageData(notice)
		if err != nil {
			a.renderTreasures(w, r, ui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error()), Form: defaultTreasureForm()})
			return
		}
		a.renderTreasures(w, r, data)
	case "/loottables":
		state, err := a.loadLootTableStateFromMaster()
		if err != nil {
			a.renderLootTables(w, r, ui.LootTablesPageData{Meta: lootTablesMeta(), Notice: errorNotice(err.Error()), Form: defaultLootTableForm()})
			return
		}
		a.renderLootTables(w, r, ui.LootTablesPageData{
			Meta:    lootTablesMeta(),
			Notice:  notice,
			Entries: state.Entries,
			Form:    defaultLootTableForm(),
		})
	case "/enemies":
		state, err := a.loadEnemyStateFromMaster()
		if err != nil {
			a.renderEnemies(w, r, ui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: defaultEnemyForm(nil)})
			return
		}
		enemySkillState, err := a.loadEnemySkillStateFromMaster()
		if err != nil {
			a.renderEnemies(w, r, ui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: defaultEnemyForm(nil)})
			return
		}
		a.renderEnemies(w, r, ui.EnemiesPageData{
			Meta:    enemiesMeta(),
			Notice:  notice,
			Entries: state.Entries,
			Form:    defaultEnemyForm(enemySkillState.Entries),
		})
	case "/spawn-tables":
		state, err := a.loadSpawnTableStateFromMaster()
		if err != nil {
			a.renderSpawnTables(w, r, ui.SpawnTablesPageData{Meta: spawnTablesMeta(), Notice: errorNotice(err.Error()), Form: defaultSpawnTableForm()})
			return
		}
		a.renderSpawnTables(w, r, ui.SpawnTablesPageData{
			Meta:    spawnTablesMeta(),
			Notice:  notice,
			Entries: state.Entries,
			Form:    defaultSpawnTableForm(),
		})
	case "/items":
		fallthrough
	default:
		state, err := a.loadItemStateFromMaster()
		if err != nil {
			a.renderItems(w, r, ui.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error()), Form: defaultItemForm(nil)})
			return
		}
		a.renderItems(w, r, ui.ItemsPageData{
			Meta:    itemMeta(),
			Notice:  notice,
			Entries: state.Items,
			Form:    defaultItemForm(nil),
		})
	}
}

func (a App) renderComponent(w http.ResponseWriter, component templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(context.Background(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func noticeText(label string, mode common.SaveMode) string {
	if mode == common.SaveModeUpdated {
		return label + " updated."
	}
	return label + " created."
}

func successNotice(text string) *ui.Notice {
	return &ui.Notice{Kind: "success", Text: text}
}

func errorNotice(text string) *ui.Notice {
	return &ui.Notice{Kind: "error", Text: text}
}

func formErrorText(value string) string {
	if strings.TrimSpace(value) == "" {
		return "Validation failed. Fix the highlighted fields."
	}
	return value
}

const flashNoticeCookieName = "tools2-flash-notice"

func redirectWithNotice(w http.ResponseWriter, r *http.Request, path string, notice *ui.Notice) bool {
	if isHX(r) {
		return false
	}
	setFlashNotice(w, notice)
	http.Redirect(w, r, path, http.StatusSeeOther)
	return true
}

func setFlashNotice(w http.ResponseWriter, notice *ui.Notice) {
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

func consumeFlashNotice(w http.ResponseWriter, r *http.Request) *ui.Notice {
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
	return &ui.Notice{Kind: parts[0], Text: parts[1]}
}

func normalizeScreenPath(value string) string {
	switch strings.TrimSpace(value) {
	case "/items", "/grimoire", "/skills", "/enemy-skills", "/treasures", "/loottables", "/enemies", "/spawn-tables":
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
