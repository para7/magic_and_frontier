package web

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/entity/loottables"
	dmaster "tools2/app/internal/domain/master"
	"tools2/app/internal/web/views"
)

func (a App) lootTablesPage(w http.ResponseWriter, r *http.Request) {
	notice := consumeFlashNotice(w, r)
	state, err := a.loadLootTableStateFromMaster()
	if err != nil {
		a.renderLootTables(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	a.renderLootTables(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Entries: state.Entries, Notice: notice})
}

func (a App) lootTablesNewPage(w http.ResponseWriter, r *http.Request) {
	returnTo := queryReturnTo(r, lootTablesMeta().CurrentPath)
	itemState, err := a.loadItemStateFromMaster()
	if err != nil {
		form := defaultLootTableForm()
		form.ReturnTo = returnTo
		a.renderLootTableForm(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	grimoireState, err := a.loadGrimoireStateFromMaster()
	if err != nil {
		form := defaultLootTableForm()
		form.ReturnTo = returnTo
		a.renderLootTableForm(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Notice: errorNotice(err.Error()), ItemOptions: itemOptions(itemState.Entries), Form: form})
		return
	}
	form := defaultLootTableForm()
	form.ReturnTo = returnTo
	a.renderLootTableForm(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), ItemOptions: itemOptions(itemState.Entries), GrimoireOptions: grimoireOptions(grimoireState.Entries), Form: form})
}

func (a App) lootTablesEditPage(w http.ResponseWriter, r *http.Request) {
	returnTo := queryReturnTo(r, lootTablesMeta().CurrentPath)
	itemState, err := a.loadItemStateFromMaster()
	if err != nil {
		a.renderLootTables(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	grimoireState, err := a.loadGrimoireStateFromMaster()
	if err != nil {
		a.renderLootTables(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	state, err := a.loadLootTableStateFromMaster()
	if err != nil {
		a.renderLootTables(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if entry, ok := findEntry(state.Entries, id, func(entry loottables.LootTableEntry) string { return entry.ID }); ok {
		form := lootTableEntryToForm(entry)
		form.ReturnTo = returnTo
		a.renderLootTableForm(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), ItemOptions: itemOptions(itemState.Entries), GrimoireOptions: grimoireOptions(grimoireState.Entries), Form: form})
		return
	}
	a.renderLootTables(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Entries: state.Entries, Notice: errorNotice("Loottable not found.")})
}

func (a App) lootTablesSubmit(w http.ResponseWriter, r *http.Request) {
	a.lootTablesSave(w, r, false)
}

func (a App) lootTablesEditSubmit(w http.ResponseWriter, r *http.Request) {
	a.lootTablesSave(w, r, true)
}

func (a App) lootTablesSave(w http.ResponseWriter, r *http.Request, editing bool) {
	_ = r.ParseForm()
	returnTo := submittedReturnTo(r, lootTablesMeta().CurrentPath)
	master, err := a.masterOrErr()
	if err != nil {
		form := defaultLootTableForm()
		form.IsEditing = editing
		form.ReturnTo = returnTo
		a.renderLootTableForm(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	itemState, err := a.loadItemStateFromMaster()
	if err != nil {
		form := defaultLootTableForm()
		form.IsEditing = editing
		form.ReturnTo = returnTo
		a.renderLootTableForm(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	grimoireState, err := a.loadGrimoireStateFromMaster()
	if err != nil {
		form := defaultLootTableForm()
		form.IsEditing = editing
		form.ReturnTo = returnTo
		a.renderLootTableForm(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Notice: errorNotice(err.Error()), ItemOptions: itemOptions(itemState.Entries), Form: form})
		return
	}
	state, err := a.loadLootTableStateFromMaster()
	if err != nil {
		form := defaultLootTableForm()
		form.IsEditing = editing
		form.ReturnTo = returnTo
		a.renderLootTableForm(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Notice: errorNotice(err.Error()), ItemOptions: itemOptions(itemState.Entries), GrimoireOptions: grimoireOptions(grimoireState.Entries), Form: form})
		return
	}
	form, input, parseErrs := parseLootTableForm(r)
	form.IsEditing = editing
	form.ReturnTo = returnTo
	if editing {
		if _, ok := findEntry(state.Entries, form.ID, func(entry loottables.LootTableEntry) string { return entry.ID }); !ok {
			a.renderLootTables(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Entries: state.Entries, Notice: errorNotice("Loottable not found.")})
			return
		}
	} else if _, ok := findEntry(state.Entries, form.ID, func(entry loottables.LootTableEntry) string { return entry.ID }); ok {
		parseErrs["id"] = "この ID は既に使用されています。"
	}
	result := master.LootTables().Validate(input, master)
	fieldErrs := mergeFieldErrors(parseErrs, mapFieldErrors(result.FieldErrors, mapLootTableField))
	if len(fieldErrs) > 0 {
		form.FieldErrors = fieldErrs
		form.FormError = formErrorText(result.FormError)
		a.renderLootTableForm(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), ItemOptions: itemOptions(itemState.Entries), GrimoireOptions: grimoireOptions(grimoireState.Entries), Form: form})
		return
	}
	mode := common.SaveModeCreated
	if editing {
		mode = common.SaveModeUpdated
		if err := master.LootTables().Update(*result.Entry, master); err != nil {
			form.FormError = formErrorText(err.Error())
			a.renderLootTableForm(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), ItemOptions: itemOptions(itemState.Entries), GrimoireOptions: grimoireOptions(grimoireState.Entries), Form: form})
			return
		}
	} else {
		if err := master.LootTables().Create(*result.Entry, master); err != nil {
			if errors.Is(err, dmaster.ErrDuplicateID) {
				form.FieldErrors = mergeFieldErrors(form.FieldErrors, map[string]string{"id": "この ID は既に使用されています。"})
			} else {
				form.FormError = formErrorText(err.Error())
			}
			a.renderLootTableForm(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), ItemOptions: itemOptions(itemState.Entries), GrimoireOptions: grimoireOptions(grimoireState.Entries), Form: form})
			return
		}
	}
	if err := master.LootTables().Save(); err != nil {
		a.renderLootTableForm(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Notice: errorNotice(err.Error()), ItemOptions: itemOptions(itemState.Entries), GrimoireOptions: grimoireOptions(grimoireState.Entries), Form: form})
		return
	}
	nextState, err := a.loadLootTableStateFromMaster()
	if err != nil {
		a.renderLootTableForm(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Notice: errorNotice(err.Error()), ItemOptions: itemOptions(itemState.Entries), GrimoireOptions: grimoireOptions(grimoireState.Entries), Form: form})
		return
	}
	notice := successNotice(noticeText("Loottable", mode))
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	a.renderLootTables(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) lootTablesDelete(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	returnTo := submittedReturnTo(r, lootTablesMeta().CurrentPath)
	master, err := a.masterOrErr()
	if err != nil {
		a.renderLootTables(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	state, err := a.loadLootTableStateFromMaster()
	if err != nil {
		a.renderLootTables(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.PathValue("id"))
	if err := master.LootTables().Delete(id, master); err != nil {
		a.renderLootTables(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Entries: state.Entries, Notice: errorNotice("Loottable not found.")})
		return
	}
	if err := master.LootTables().Save(); err != nil {
		a.renderLootTables(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Entries: state.Entries, Notice: errorNotice(err.Error())})
		return
	}
	nextState, err := a.loadLootTableStateFromMaster()
	if err != nil {
		a.renderLootTables(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Entries: state.Entries, Notice: errorNotice(err.Error())})
		return
	}
	notice := successNotice("Loottable deleted.")
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	a.renderLootTables(w, r, views.LootTablesPageData{Meta: lootTablesMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) renderLootTables(w http.ResponseWriter, r *http.Request, data views.LootTablesPageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.LootTablesShell(data))
		return
	}
	a.renderComponent(w, views.LootTablesPage(data))
}

func (a App) renderLootTableForm(w http.ResponseWriter, r *http.Request, data views.LootTablesPageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.LootTableFormShell(data))
		return
	}
	a.renderComponent(w, views.LootTableFormPage(data))
}

func lootTablesMeta() views.PageMeta {
	return views.PageMeta{Title: "Loottables", CurrentPath: "/loottables"}
}

func defaultLootTableForm() views.LootTableFormData {
	return views.LootTableFormData{
		ID:            "",
		Memo:          "",
		LootPoolsText: "item,,1,1,1",
		FieldErrors:   map[string]string{},
	}
}

func lootTableEntryToForm(entry loottables.LootTableEntry) views.LootTableFormData {
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
	return views.LootTableFormData{
		ID:            entry.ID,
		Memo:          entry.Memo,
		LootPoolsText: strings.Join(lines, "\n"),
		FieldErrors:   map[string]string{},
		IsEditing:     true,
	}
}

func parseLootTableForm(r *http.Request) (views.LootTableFormData, loottables.SaveInput, map[string]string) {
	form := defaultLootTableForm()
	form.ID = strings.TrimSpace(r.Form.Get("id"))
	form.Memo = r.Form.Get("memo")
	form.LootPoolsText = r.Form.Get("lootPoolsText")
	errs := map[string]string{}
	input := loottables.SaveInput{
		ID:        form.ID,
		Memo:      form.Memo,
		LootPools: parseTreasurePools(errs, form.LootPoolsText),
	}
	return form, input, errs
}
