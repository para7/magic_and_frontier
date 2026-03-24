package web

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/grimoire"
	dmaster "tools2/app/internal/domain/master"
	"tools2/app/internal/web/ui"
	"tools2/app/internal/web/views"
)

func (a App) grimoirePage(w http.ResponseWriter, r *http.Request) {
	notice := consumeFlashNotice(w, r)
	state, err := a.loadGrimoireStateFromMaster()
	if err != nil {
		a.renderGrimoire(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Notice: errorNotice(err.Error())})
		return
	}
	a.renderGrimoire(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Entries: state.Entries, Notice: notice, Form: defaultGrimoireForm(nil)})
}

func (a App) grimoireNewPage(w http.ResponseWriter, r *http.Request) {
	form := defaultGrimoireForm(nil)
	form.ReturnTo = queryReturnTo(r, grimoireMeta().CurrentPath)
	a.renderGrimoireForm(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Form: form})
}

func (a App) grimoireEditPage(w http.ResponseWriter, r *http.Request) {
	returnTo := queryReturnTo(r, grimoireMeta().CurrentPath)
	state, err := a.loadGrimoireStateFromMaster()
	if err != nil {
		a.renderGrimoire(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if entry, ok := findEntry(state.Entries, id, func(entry grimoire.GrimoireEntry) string { return entry.ID }); ok {
		form := grimoireEntryToForm(entry)
		form.ReturnTo = returnTo
		a.renderGrimoireForm(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Form: form})
		return
	}
	a.renderGrimoire(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Entries: state.Entries, Notice: errorNotice("Grimoire entry not found.")})
}

func (a App) grimoireSubmit(w http.ResponseWriter, r *http.Request) {
	a.grimoireSave(w, r, false)
}

func (a App) grimoireEditSubmit(w http.ResponseWriter, r *http.Request) {
	a.grimoireSave(w, r, true)
}

func (a App) grimoireSave(w http.ResponseWriter, r *http.Request, editing bool) {
	_ = r.ParseForm()
	returnTo := submittedReturnTo(r, grimoireMeta().CurrentPath)
	master, err := a.masterOrErr()
	if err != nil {
		form := defaultGrimoireForm(nil)
		form.IsEditing = editing
		form.ReturnTo = returnTo
		a.renderGrimoireForm(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	state, err := a.loadGrimoireStateFromMaster()
	if err != nil {
		form := defaultGrimoireForm(nil)
		form.IsEditing = editing
		form.ReturnTo = returnTo
		a.renderGrimoireForm(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	form, input, parseErrs := parseGrimoireForm(r)
	form.IsEditing = editing
	form.ReturnTo = returnTo
	if editing {
		existing, ok := findEntry(state.Entries, form.ID, func(entry grimoire.GrimoireEntry) string { return entry.ID })
		if !ok {
			a.renderGrimoire(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Entries: state.Entries, Notice: errorNotice("Grimoire entry not found.")})
			return
		}
		input.ID = existing.ID
		input.CastID = existing.CastID
		form.ID = existing.ID
		form.CastID = strconv.Itoa(existing.CastID)
	} else {
		castID, parseErr := strconv.Atoi(form.CastID)
		if parseErr != nil {
			parseErrs["castid"] = "Must be a number."
		} else {
			input.CastID = castID
		}
		if _, ok := findEntry(state.Entries, form.ID, func(entry grimoire.GrimoireEntry) string { return entry.ID }); ok {
			parseErrs["id"] = "この ID は既に使用されています。"
		}
	}
	result := master.Grimoires().Validate(input, master)
	fieldErrs := mergeFieldErrors(parseErrs, mapFieldErrors(result.FieldErrors, mapGrimoireField))
	if len(fieldErrs) > 0 {
		form.FieldErrors = fieldErrs
		form.FormError = formErrorText(result.FormError)
		a.renderGrimoireForm(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Form: form})
		return
	}
	mode := common.SaveModeCreated
	if editing {
		mode = common.SaveModeUpdated
		if err := master.Grimoires().Update(*result.Entry, master); err != nil {
			form.FormError = formErrorText(err.Error())
			a.renderGrimoireForm(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Form: form})
			return
		}
	} else {
		if err := master.Grimoires().Create(*result.Entry, master); err != nil {
			if errors.Is(err, dmaster.ErrDuplicateID) {
				form.FieldErrors = mergeFieldErrors(form.FieldErrors, map[string]string{"id": "この ID は既に使用されています。"})
			} else {
				form.FormError = formErrorText(err.Error())
			}
			a.renderGrimoireForm(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Form: form})
			return
		}
	}
	if err := master.Grimoires().Save(); err != nil {
		a.renderGrimoireForm(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	nextState, err := a.loadGrimoireStateFromMaster()
	if err != nil {
		a.renderGrimoireForm(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	notice := successNotice(noticeText("Grimoire entry", mode))
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	a.renderGrimoire(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) grimoireDelete(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	returnTo := submittedReturnTo(r, grimoireMeta().CurrentPath)
	id := strings.TrimSpace(r.PathValue("id"))
	master, err := a.masterOrErr()
	if err != nil {
		a.renderGrimoire(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Notice: errorNotice(err.Error())})
		return
	}
	state, err := a.loadGrimoireStateFromMaster()
	if err != nil {
		a.renderGrimoire(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Notice: errorNotice(err.Error())})
		return
	}
	if err := master.Grimoires().Delete(id, master); err != nil {
		a.renderGrimoire(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Entries: state.Entries, Notice: errorNotice("Grimoire entry not found.")})
		return
	}
	if err := master.Grimoires().Save(); err != nil {
		a.renderGrimoire(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Entries: state.Entries, Notice: errorNotice(err.Error())})
		return
	}
	nextState, err := a.loadGrimoireStateFromMaster()
	if err != nil {
		a.renderGrimoire(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Entries: state.Entries, Notice: errorNotice(err.Error())})
		return
	}
	notice := successNotice("Grimoire entry deleted.")
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	a.renderGrimoire(w, r, ui.GrimoirePageData{Meta: grimoireMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) renderGrimoire(w http.ResponseWriter, r *http.Request, data ui.GrimoirePageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.GrimoireShell(data))
		return
	}
	a.renderComponent(w, views.GrimoirePage(data))
}

func (a App) renderGrimoireForm(w http.ResponseWriter, r *http.Request, data ui.GrimoirePageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.GrimoireFormShell(data))
		return
	}
	a.renderComponent(w, views.GrimoireFormPage(data))
}

func grimoireMeta() ui.PageMeta {
	return ui.PageMeta{Title: "Grimoire", CurrentPath: "/grimoire"}
}

func defaultGrimoireForm(_ []grimoire.GrimoireEntry) ui.GrimoireFormData {
	return ui.GrimoireFormData{
		ID:          "",
		CastTime:    "0",
		MPCost:      "0",
		FieldErrors: map[string]string{},
	}
}

func grimoireEntryToForm(entry grimoire.GrimoireEntry) ui.GrimoireFormData {
	return ui.GrimoireFormData{
		ID:          entry.ID,
		CastID:      strconv.Itoa(entry.CastID),
		CastTime:    strconv.Itoa(entry.CastTime),
		MPCost:      strconv.Itoa(entry.MPCost),
		Script:      entry.Script,
		Title:       entry.Title,
		Description: entry.Description,
		FieldErrors: map[string]string{},
		IsEditing:   true,
	}
}

func parseGrimoireForm(r *http.Request) (ui.GrimoireFormData, grimoire.SaveInput, map[string]string) {
	form := defaultGrimoireForm(nil)
	form.ID = strings.TrimSpace(r.Form.Get("id"))
	form.CastID = strings.TrimSpace(r.Form.Get("castid"))
	form.CastTime = strings.TrimSpace(r.Form.Get("castTime"))
	form.MPCost = strings.TrimSpace(r.Form.Get("mpCost"))
	form.Script = r.Form.Get("script")
	form.Title = r.Form.Get("title")
	form.Description = r.Form.Get("description")
	errs := map[string]string{}
	input := grimoire.SaveInput{
		ID:          form.ID,
		CastID:      0,
		CastTime:    parseRequiredInt(errs, "castTime", form.CastTime),
		MPCost:      parseRequiredInt(errs, "mpCost", form.MPCost),
		Script:      form.Script,
		Title:       form.Title,
		Description: form.Description,
	}
	return form, input, errs
}

func grimoireOptions(entries []grimoire.GrimoireEntry) []ui.ReferenceOption {
	options := make([]ui.ReferenceOption, 0, len(entries))
	for _, entry := range entries {
		options = append(options, ui.ReferenceOption{ID: entry.ID, Label: entry.Title})
	}
	return options
}
