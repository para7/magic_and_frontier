package web

import (
	"net/http"
	"strconv"
	"strings"

	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/webui"
	"tools2/app/views"
)

func (a App) grimoirePage(w http.ResponseWriter, r *http.Request) {
	notice := consumeFlashNotice(w, r)
	state, err := a.deps.GrimoireRepo.LoadGrimoireState()
	if err != nil {
		a.renderGrimoire(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Notice: errorNotice(err.Error())})
		return
	}
	a.renderGrimoire(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Entries: state.Entries, Notice: notice, Form: defaultGrimoireForm(nil)})
}

func (a App) grimoireNewPage(w http.ResponseWriter, r *http.Request) {
	form := defaultGrimoireForm(nil)
	form.ReturnTo = queryReturnTo(r, grimoireMeta().CurrentPath)
	a.renderGrimoireForm(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Form: form})
}

func (a App) grimoireEditPage(w http.ResponseWriter, r *http.Request) {
	returnTo := queryReturnTo(r, grimoireMeta().CurrentPath)
	state, err := a.deps.GrimoireRepo.LoadGrimoireState()
	if err != nil {
		a.renderGrimoire(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if entry, ok := findEntry(state.Entries, id, func(entry grimoire.GrimoireEntry) string { return entry.ID }); ok {
		form := grimoireEntryToForm(entry)
		form.ReturnTo = returnTo
		a.renderGrimoireForm(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Form: form})
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
	returnTo := submittedReturnTo(r, grimoireMeta().CurrentPath)
	state, err := a.deps.GrimoireRepo.LoadGrimoireState()
	if err != nil {
		form := defaultGrimoireForm(nil)
		form.IsEditing = editing
		form.ReturnTo = returnTo
		a.renderGrimoireForm(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	form, input, parseErrs := parseGrimoireForm(r)
	form.IsEditing = editing
	form.ReturnTo = returnTo
	if editing {
		existing, ok := findEntry(state.Entries, form.ID, func(entry grimoire.GrimoireEntry) string { return entry.ID })
		if !ok {
			a.renderGrimoire(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Entries: state.Entries, Notice: errorNotice("Grimoire entry not found.")})
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
	result := grimoire.ValidateSave(input, a.deps.Now())
	errors := mergeFieldErrors(parseErrs, mapFieldErrors(result.FieldErrors, mapGrimoireField))
	if conflictID := duplicateCastID(state.Entries, input.ID, input.CastID); conflictID != "" {
		errors["castid"] = "Cast ID is already used by " + conflictID + "."
	}
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
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	a.renderGrimoire(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) grimoireDelete(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	returnTo := submittedReturnTo(r, grimoireMeta().CurrentPath)
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
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	a.renderGrimoire(w, r, webui.GrimoirePageData{Meta: grimoireMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) renderGrimoire(w http.ResponseWriter, r *http.Request, data webui.GrimoirePageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.GrimoireShell(data))
		return
	}
	a.renderComponent(w, views.GrimoirePage(data))
}

func (a App) renderGrimoireForm(w http.ResponseWriter, r *http.Request, data webui.GrimoirePageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.GrimoireFormShell(data))
		return
	}
	a.renderComponent(w, views.GrimoireFormPage(data))
}

func grimoireMeta() webui.PageMeta {
	return webui.PageMeta{Title: "Grimoire", CurrentPath: "/grimoire"}
}

func defaultGrimoireForm(_ []grimoire.GrimoireEntry) webui.GrimoireFormData {
	return webui.GrimoireFormData{
		ID:          "",
		CastTime:    "0",
		MPCost:      "0",
		FieldErrors: map[string]string{},
	}
}

func grimoireEntryToForm(entry grimoire.GrimoireEntry) webui.GrimoireFormData {
	return webui.GrimoireFormData{
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

func parseGrimoireForm(r *http.Request) (webui.GrimoireFormData, grimoire.SaveInput, map[string]string) {
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

func grimoireOptions(entries []grimoire.GrimoireEntry) []webui.ReferenceOption {
	options := make([]webui.ReferenceOption, 0, len(entries))
	for _, entry := range entries {
		options = append(options, webui.ReferenceOption{ID: entry.ID, Label: entry.Title})
	}
	return options
}

func grimoireIDSet(state grimoire.GrimoireState) map[string]struct{} {
	return toIDSet(state.Entries, func(entry grimoire.GrimoireEntry) string { return entry.ID })
}
