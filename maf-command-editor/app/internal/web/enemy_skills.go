package web

import (
	"errors"
	"net/http"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemyskills"
	dmaster "tools2/app/internal/domain/master"
	"tools2/app/internal/webui"
	"tools2/app/views"
)

func (a App) enemySkillsPage(w http.ResponseWriter, r *http.Request) {
	notice := consumeFlashNotice(w, r)
	state, err := a.loadEnemySkillStateFromMaster()
	if err != nil {
		a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Notice: errorNotice(err.Error())})
		return
	}
	a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Entries: state.Entries, Notice: notice, Form: defaultEnemySkillForm()})
}

func (a App) enemySkillsNewPage(w http.ResponseWriter, r *http.Request) {
	form := defaultEnemySkillForm()
	form.ReturnTo = queryReturnTo(r, enemySkillsMeta().CurrentPath)
	a.renderEnemySkillForm(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Form: form})
}

func (a App) enemySkillsEditPage(w http.ResponseWriter, r *http.Request) {
	returnTo := queryReturnTo(r, enemySkillsMeta().CurrentPath)
	state, err := a.loadEnemySkillStateFromMaster()
	if err != nil {
		a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if entry, ok := findEntry(state.Entries, id, func(entry enemyskills.EnemySkillEntry) string { return entry.ID }); ok {
		form := enemySkillEntryToForm(entry)
		form.ReturnTo = returnTo
		a.renderEnemySkillForm(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Form: form})
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
	returnTo := submittedReturnTo(r, enemySkillsMeta().CurrentPath)
	master, err := a.masterOrErr()
	if err != nil {
		form := defaultEnemySkillForm()
		form.IsEditing = editing
		form.ReturnTo = returnTo
		a.renderEnemySkillForm(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	state, err := a.loadEnemySkillStateFromMaster()
	if err != nil {
		form := defaultEnemySkillForm()
		form.IsEditing = editing
		form.ReturnTo = returnTo
		a.renderEnemySkillForm(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	form, input, parseErrs := parseEnemySkillForm(r)
	form.IsEditing = editing
	form.ReturnTo = returnTo
	if editing {
		if _, ok := findEntry(state.Entries, form.ID, func(entry enemyskills.EnemySkillEntry) string { return entry.ID }); !ok {
			a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Entries: state.Entries, Notice: errorNotice("Enemy skill not found.")})
			return
		}
	} else if _, ok := findEntry(state.Entries, form.ID, func(entry enemyskills.EnemySkillEntry) string { return entry.ID }); ok {
		parseErrs["id"] = "この ID は既に使用されています。"
	}
	result := master.EnemySkills().Validate(input, master)
	fieldErrs := mergeFieldErrors(parseErrs, result.FieldErrors)
	if len(fieldErrs) > 0 {
		form.FieldErrors = fieldErrs
		form.FormError = formErrorText(result.FormError)
		a.renderEnemySkillForm(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Form: form})
		return
	}
	mode := common.SaveModeCreated
	if editing {
		mode = common.SaveModeUpdated
		if err := master.EnemySkills().Update(*result.Entry, master); err != nil {
			form.FormError = formErrorText(err.Error())
			a.renderEnemySkillForm(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Form: form})
			return
		}
	} else {
		if err := master.EnemySkills().Create(*result.Entry, master); err != nil {
			if errors.Is(err, dmaster.ErrDuplicateID) {
				form.FieldErrors = mergeFieldErrors(form.FieldErrors, map[string]string{"id": "この ID は既に使用されています。"})
			} else {
				form.FormError = formErrorText(err.Error())
			}
			a.renderEnemySkillForm(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Form: form})
			return
		}
	}
	if err := master.EnemySkills().Save(); err != nil {
		a.renderEnemySkillForm(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	nextState, err := a.loadEnemySkillStateFromMaster()
	if err != nil {
		a.renderEnemySkillForm(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	notice := successNotice(noticeText("Enemy skill", mode))
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) enemySkillsDelete(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	returnTo := submittedReturnTo(r, enemySkillsMeta().CurrentPath)
	id := strings.TrimSpace(r.PathValue("id"))
	master, err := a.masterOrErr()
	if err != nil {
		a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Notice: errorNotice(err.Error())})
		return
	}
	state, err := a.loadEnemySkillStateFromMaster()
	if err != nil {
		a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Notice: errorNotice(err.Error())})
		return
	}
	enemyState, err := a.loadEnemyStateFromMaster()
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
	if err := master.EnemySkills().Delete(id, master); err != nil {
		a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Entries: state.Entries, Notice: errorNotice(err.Error())})
		return
	}
	if err := master.EnemySkills().Save(); err != nil {
		a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Entries: state.Entries, Notice: errorNotice(err.Error())})
		return
	}
	nextState, err := a.loadEnemySkillStateFromMaster()
	if err != nil {
		a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Entries: state.Entries, Notice: errorNotice(err.Error())})
		return
	}
	notice := successNotice("Enemy skill deleted.")
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	a.renderEnemySkills(w, r, webui.EnemySkillsPageData{Meta: enemySkillsMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) renderEnemySkills(w http.ResponseWriter, r *http.Request, data webui.EnemySkillsPageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.EnemySkillsShell(data))
		return
	}
	a.renderComponent(w, views.EnemySkillsPage(data))
}

func (a App) renderEnemySkillForm(w http.ResponseWriter, r *http.Request, data webui.EnemySkillsPageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.EnemySkillFormShell(data))
		return
	}
	a.renderComponent(w, views.EnemySkillFormPage(data))
}

func enemySkillsMeta() webui.PageMeta {
	return webui.PageMeta{Title: "Enemy Skills", CurrentPath: "/enemy-skills"}
}

func defaultEnemySkillForm() webui.EnemySkillFormData {
	return webui.EnemySkillFormData{
		FieldErrors: map[string]string{},
	}
}

func enemySkillEntryToForm(entry enemyskills.EnemySkillEntry) webui.EnemySkillFormData {
	return webui.EnemySkillFormData{
		ID:          entry.ID,
		Name:        entry.Name,
		Description: entry.Description,
		Script:      entry.Script,
		FieldErrors: map[string]string{},
		IsEditing:   true,
	}
}

func parseEnemySkillForm(r *http.Request) (webui.EnemySkillFormData, enemyskills.SaveInput, map[string]string) {
	form := defaultEnemySkillForm()
	form.ID = strings.TrimSpace(r.Form.Get("id"))
	form.Name = r.Form.Get("name")
	form.Description = r.Form.Get("description")
	form.Script = r.Form.Get("script")
	errs := map[string]string{}
	input := enemyskills.SaveInput{
		ID:          form.ID,
		Name:        form.Name,
		Description: form.Description,
		Script:      form.Script,
	}
	return form, input, errs
}

func enemySkillOptions(entries []enemyskills.EnemySkillEntry) []webui.ReferenceOption {
	options := make([]webui.ReferenceOption, 0, len(entries))
	for _, entry := range entries {
		options = append(options, webui.ReferenceOption{ID: entry.ID, Label: entry.Name})
	}
	return options
}

func enemySkillIDSet(state common.EntryState[enemyskills.EnemySkillEntry]) map[string]struct{} {
	return toIDSet(state.Entries, func(entry enemyskills.EnemySkillEntry) string { return entry.ID })
}
