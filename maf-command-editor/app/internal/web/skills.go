package web

import (
	"net/http"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/webui"
	"tools2/app/views"
)

func (a App) skillsPage(w http.ResponseWriter, r *http.Request) {
	notice := consumeFlashNotice(w, r)
	state, err := a.deps.SkillRepo.LoadState()
	if err != nil {
		a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Notice: errorNotice(err.Error())})
		return
	}
	a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Entries: state.Entries, Notice: notice, Form: defaultSkillForm()})
}

func (a App) skillsNewPage(w http.ResponseWriter, r *http.Request) {
	form := defaultSkillForm()
	form.ReturnTo = queryReturnTo(r, skillsMeta().CurrentPath)
	a.renderSkillForm(w, r, webui.SkillsPageData{Meta: skillsMeta(), Form: form})
}

func (a App) skillsEditPage(w http.ResponseWriter, r *http.Request) {
	returnTo := queryReturnTo(r, skillsMeta().CurrentPath)
	state, err := a.deps.SkillRepo.LoadState()
	if err != nil {
		a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if entry, ok := findEntry(state.Entries, id, func(entry skills.SkillEntry) string { return entry.ID }); ok {
		form := skillEntryToForm(entry)
		form.ReturnTo = returnTo
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
	returnTo := submittedReturnTo(r, skillsMeta().CurrentPath)
	state, err := a.deps.SkillRepo.LoadState()
	if err != nil {
		form := defaultSkillForm()
		form.IsEditing = editing
		form.ReturnTo = returnTo
		a.renderSkillForm(w, r, webui.SkillsPageData{Meta: skillsMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	form, input, parseErrs := parseSkillForm(r)
	form.IsEditing = editing
	form.ReturnTo = returnTo
	if editing {
		if _, ok := findEntry(state.Entries, form.ID, func(entry skills.SkillEntry) string { return entry.ID }); !ok {
			a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Entries: state.Entries, Notice: errorNotice("Skill not found.")})
			return
		}
	} else if _, ok := findEntry(state.Entries, form.ID, func(entry skills.SkillEntry) string { return entry.ID }); ok {
		parseErrs["id"] = "この ID は既に使用されています。"
	}
	result := skills.ValidateSave(input, a.deps.Now())
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
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) skillsDelete(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	returnTo := submittedReturnTo(r, skillsMeta().CurrentPath)
	itemState, err := a.deps.ItemRepo.LoadItemState()
	if err != nil {
		a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.PathValue("id"))
	for _, entry := range itemState.Items {
		if entry.SkillID == id {
			state, _ := a.deps.SkillRepo.LoadState()
			a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Entries: state.Entries, Notice: errorNotice("Skill is referenced by item " + entry.ID + ".")})
			return
		}
	}
	state, err := a.deps.SkillRepo.LoadState()
	if err != nil {
		a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Notice: errorNotice(err.Error())})
		return
	}
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
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	a.renderSkills(w, r, webui.SkillsPageData{Meta: skillsMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) renderSkills(w http.ResponseWriter, r *http.Request, data webui.SkillsPageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.SkillsShell(data))
		return
	}
	a.renderComponent(w, views.SkillsPage(data))
}

func (a App) renderSkillForm(w http.ResponseWriter, r *http.Request, data webui.SkillsPageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.SkillFormShell(data))
		return
	}
	a.renderComponent(w, views.SkillFormPage(data))
}

func skillsMeta() webui.PageMeta {
	return webui.PageMeta{Title: "Skills", CurrentPath: "/skills"}
}

func defaultSkillForm() webui.SkillFormData {
	return webui.SkillFormData{
		SkillType:   "sword",
		FieldErrors: map[string]string{},
	}
}

func skillEntryToForm(entry skills.SkillEntry) webui.SkillFormData {
	return webui.SkillFormData{
		ID:          entry.ID,
		Name:        entry.Name,
		SkillType:   entry.SkillType,
		Description: entry.Description,
		Script:      entry.Script,
		FieldErrors: map[string]string{},
		IsEditing:   true,
	}
}

func parseSkillForm(r *http.Request) (webui.SkillFormData, skills.SaveInput, map[string]string) {
	form := defaultSkillForm()
	form.ID = strings.TrimSpace(r.Form.Get("id"))
	form.Name = r.Form.Get("name")
	form.SkillType = strings.TrimSpace(r.Form.Get("skilltype"))
	form.Description = r.Form.Get("description")
	form.Script = r.Form.Get("script")
	errs := map[string]string{}
	input := skills.SaveInput{
		ID:          form.ID,
		Name:        form.Name,
		SkillType:   form.SkillType,
		Description: form.Description,
		Script:      form.Script,
	}
	return form, input, errs
}
