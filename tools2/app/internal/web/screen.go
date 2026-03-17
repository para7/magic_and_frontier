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
	"tools2/app/internal/webui"
	"tools2/app/views"
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
		state, err := a.deps.SkillRepo.LoadState()
		if err != nil {
			a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Notice: errorNotice(err.Error()), Form: defaultSkillForm()})
			return
		}
		a.renderSkills(w, r, webui.SkillsPageData{
			Meta:    skillsMeta(),
			Notice:  notice,
			Entries: state.Entries,
			Form:    defaultSkillForm(),
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
			a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error()), Form: defaultItemForm(nil)})
			return
		}
		a.renderItems(w, r, webui.ItemsPageData{
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
