package web

import (
	"net/http"
	"strconv"
	"strings"

	"tools2/app/internal/application"
	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/treasures"
	"tools2/app/internal/webui"
	"tools2/app/views"
)

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
	} else if strings.TrimSpace(input.ID) == "" {
		id, allocErr := application.NewService(a.cfg, a.deps).AllocateID("treasure")
		if allocErr != nil {
			a.renderTreasureForm(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(allocErr.Error()), ItemOptions: itemOptions(itemState.Items), GrimoireOptions: grimoireOptions(grimoireState.Entries), Form: form})
			return
		}
		input.ID = id
		form.ID = id
	}
	result := treasures.ValidateSave(input, itemIDSet(itemState), grimoireIDSet(grimoireState), a.deps.Now())
	errors := mergeFieldErrors(parseErrs, mapFieldErrors(result.FieldErrors, mapTreasureField))
	if conflictID := duplicateCustomTablePath(state.Entries, input.ID, input.Mode, input.TablePath); conflictID != "" {
		errors["tablePath"] = "Custom loot table path is already used by " + conflictID + "."
	}
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

func treasuresMeta() webui.PageMeta {
	return webui.PageMeta{Title: "Treasures", CurrentPath: "/treasures", Description: "treasure loot pools と保存先 table path を管理します。"}
}

func defaultTreasureForm() webui.TreasureFormData {
	return webui.TreasureFormData{
		ID:            "",
		Mode:          "custom",
		TablePath:     "maf:treasure/example",
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
		Mode:          entry.Mode,
		TablePath:     entry.TablePath,
		LootPoolsText: strings.Join(lines, "\n"),
		FieldErrors:   map[string]string{},
		IsEditing:     true,
	}
}

func parseTreasureForm(r *http.Request) (webui.TreasureFormData, treasures.SaveInput, map[string]string) {
	form := defaultTreasureForm()
	form.ID = strings.TrimSpace(r.Form.Get("id"))
	form.Mode = strings.TrimSpace(r.Form.Get("mode"))
	form.TablePath = strings.TrimSpace(r.Form.Get("tablePath"))
	form.LootPoolsText = r.Form.Get("lootPoolsText")
	errs := map[string]string{}
	input := treasures.SaveInput{
		ID:        form.ID,
		Mode:      form.Mode,
		TablePath: form.TablePath,
		LootPools: parseTreasurePools(errs, form.LootPoolsText),
	}
	return form, input, errs
}
