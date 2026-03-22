package web

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"tools2/app/internal/application"
	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/treasures"
	"tools2/app/internal/mcsource"
	"tools2/app/internal/webui"
	"tools2/app/views"
)

func (a App) treasuresPage(w http.ResponseWriter, r *http.Request) {
	notice := consumeFlashNotice(w, r)
	data, err := a.treasuresPageData(notice)
	if err != nil {
		a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error())})
		return
	}
	a.renderTreasures(w, r, data)
}

func (a App) treasuresNewPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, treasuresMeta().CurrentPath, http.StatusSeeOther)
}

func (a App) treasuresEditPage(w http.ResponseWriter, r *http.Request) {
	returnTo := queryReturnTo(r, treasuresMeta().CurrentPath)
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
	state, sources, sourcePaths, err := a.loadTreasureCatalog()
	if err != nil {
		a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error())})
		return
	}
	tablePath := strings.TrimSpace(r.URL.Query().Get("tablePath"))
	if tablePath == "" {
		if id := strings.TrimSpace(r.URL.Query().Get("id")); id != "" {
			if entry, ok := findEntry(state.Entries, id, func(entry treasures.TreasureEntry) string { return entry.ID }); ok {
				tablePath = entry.TablePath
			}
		}
	}
	if strings.TrimSpace(tablePath) == "" {
		data := treasureListPageData(state, sources, noticeWithError("Treasure not found."))
		a.renderTreasures(w, r, data)
		return
	}
	entry, hasOverlay := findTreasureByTablePath(state.Entries, tablePath)
	form := defaultTreasureForm()
	form.ReturnTo = returnTo
	form.TablePath = tablePath
	form.HasSource = hasTreasureSource(sourcePaths, tablePath)
	form.HasOverlay = hasOverlay
	form.IsEditing = hasOverlay
	if hasOverlay {
		form = treasureEntryToForm(entry)
		form.ReturnTo = returnTo
		form.HasSource = hasTreasureSource(sourcePaths, tablePath)
		form.HasOverlay = true
		form.IsEditing = true
	}
	if !form.HasSource && !form.HasOverlay {
		data := treasureListPageData(state, sources, noticeWithError("Treasure not found."))
		a.renderTreasures(w, r, data)
		return
	}
	a.renderTreasureForm(w, r, webui.TreasuresPageData{
		Meta:            treasuresMeta(),
		ItemOptions:     itemOptions(itemState.Items),
		GrimoireOptions: grimoireOptions(grimoireState.Entries),
		Form:            form,
	})
}

func (a App) treasuresSubmit(w http.ResponseWriter, r *http.Request) {
	a.treasuresSave(w, r)
}

func (a App) treasuresEditSubmit(w http.ResponseWriter, r *http.Request) {
	a.treasuresSave(w, r)
}

func (a App) treasuresSave(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	returnTo := submittedReturnTo(r, treasuresMeta().CurrentPath)
	itemState, err := a.deps.ItemRepo.LoadItemState()
	if err != nil {
		form := defaultTreasureForm()
		form.ReturnTo = returnTo
		a.renderTreasureForm(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	grimoireState, err := a.deps.GrimoireRepo.LoadGrimoireState()
	if err != nil {
		form := defaultTreasureForm()
		form.ReturnTo = returnTo
		a.renderTreasureForm(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error()), ItemOptions: itemOptions(itemState.Items), Form: form})
		return
	}
	state, _, sourcePaths, err := a.loadTreasureCatalog()
	if err != nil {
		form := defaultTreasureForm()
		form.ReturnTo = returnTo
		a.renderTreasureForm(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error()), ItemOptions: itemOptions(itemState.Items), GrimoireOptions: grimoireOptions(grimoireState.Entries), Form: form})
		return
	}
	form, input, parseErrs := parseTreasureForm(r)
	form.ReturnTo = returnTo
	form.HasSource = hasTreasureSource(sourcePaths, input.TablePath)
	if existing, ok := findTreasureByTablePath(state.Entries, input.TablePath); ok {
		form.HasOverlay = true
		form.IsEditing = true
		if strings.TrimSpace(input.ID) == "" {
			input.ID = existing.ID
			form.ID = existing.ID
		}
	}
	if strings.TrimSpace(input.ID) != "" {
		if _, ok := findEntry(state.Entries, input.ID, func(entry treasures.TreasureEntry) string { return entry.ID }); !ok {
			data := treasureListPageData(state, nil, noticeWithError("Treasure not found."))
			a.renderTreasures(w, r, data)
			return
		}
	} else {
		id, allocErr := application.NewService(a.cfg, a.deps).AllocateID("treasure")
		if allocErr != nil {
			a.renderTreasureForm(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(allocErr.Error()), ItemOptions: itemOptions(itemState.Items), GrimoireOptions: grimoireOptions(grimoireState.Entries), Form: form})
			return
		}
		input.ID = id
		form.ID = id
	}
	result := treasures.ValidateSave(input, itemIDSet(itemState), grimoireIDSet(grimoireState), sourcePaths, a.deps.Now())
	errors := mergeFieldErrors(parseErrs, mapFieldErrors(result.FieldErrors, mapTreasureField))
	if conflictID := duplicateTreasureTablePath(state.Entries, input.ID, strings.TrimSpace(input.TablePath)); conflictID != "" {
		errors["tablePath"] = "Loot table path is already used by " + conflictID + "."
	}
	if len(errors) > 0 {
		form.FieldErrors = errors
		form.FormError = formErrorText(result.FormError)
		form.HasOverlay = strings.TrimSpace(form.ID) != ""
		form.IsEditing = form.HasOverlay
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
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	data, err := a.treasuresPageData(notice)
	if err != nil {
		a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error())})
		return
	}
	a.renderTreasures(w, r, data)
}

func (a App) treasuresDelete(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	returnTo := submittedReturnTo(r, treasuresMeta().CurrentPath)
	state, err := a.deps.TreasureRepo.LoadState()
	if err != nil {
		a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.PathValue("id"))
	nextState, ok := common.DeleteEntries(state, id, func(entry treasures.TreasureEntry) string { return entry.ID })
	if !ok {
		data, dataErr := a.treasuresPageData(errorNotice("Treasure not found."))
		if dataErr != nil {
			a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(dataErr.Error())})
			return
		}
		a.renderTreasures(w, r, data)
		return
	}
	if err := a.deps.TreasureRepo.SaveState(nextState); err != nil {
		data, dataErr := a.treasuresPageData(errorNotice(err.Error()))
		if dataErr != nil {
			a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(dataErr.Error())})
			return
		}
		a.renderTreasures(w, r, data)
		return
	}
	notice := successNotice("Treasure deleted.")
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	data, dataErr := a.treasuresPageData(notice)
	if dataErr != nil {
		a.renderTreasures(w, r, webui.TreasuresPageData{Meta: treasuresMeta(), Notice: errorNotice(dataErr.Error())})
		return
	}
	a.renderTreasures(w, r, data)
}

func (a App) renderTreasures(w http.ResponseWriter, r *http.Request, data webui.TreasuresPageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.TreasuresShell(data))
		return
	}
	a.renderComponent(w, views.TreasuresPage(data))
}

func (a App) renderTreasureForm(w http.ResponseWriter, r *http.Request, data webui.TreasuresPageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.TreasureFormShell(data))
		return
	}
	a.renderComponent(w, views.TreasureFormPage(data))
}

func (a App) treasuresPageData(notice *webui.Notice) (webui.TreasuresPageData, error) {
	state, sources, _, err := a.loadTreasureCatalog()
	if err != nil {
		return webui.TreasuresPageData{}, err
	}
	return treasureListPageData(state, sources, notice), nil
}

func (a App) loadTreasureCatalog() (common.EntryState[treasures.TreasureEntry], []mcsource.LootTableSource, map[string]struct{}, error) {
	state, err := a.deps.TreasureRepo.LoadState()
	if err != nil {
		return common.EntryState[treasures.TreasureEntry]{}, nil, nil, err
	}
	sources, err := mcsource.ListLootTables(a.cfg.MinecraftLootTableRoot)
	if err != nil {
		return common.EntryState[treasures.TreasureEntry]{}, nil, nil, err
	}
	sourcePaths := map[string]struct{}{}
	for _, source := range sources {
		sourcePaths[source.TablePath] = struct{}{}
	}
	return state, sources, sourcePaths, nil
}

func treasureListPageData(state common.EntryState[treasures.TreasureEntry], sources []mcsource.LootTableSource, notice *webui.Notice) webui.TreasuresPageData {
	return webui.TreasuresPageData{
		Meta:    treasuresMeta(),
		Notice:  notice,
		Entries: buildTreasureListEntries(state.Entries, sources),
	}
}

func buildTreasureListEntries(entries []treasures.TreasureEntry, sources []mcsource.LootTableSource) []webui.TreasureListEntry {
	merged := map[string]webui.TreasureListEntry{}
	for _, source := range sources {
		merged[source.TablePath] = webui.TreasureListEntry{
			TablePath: source.TablePath,
			HasSource: true,
		}
	}
	for _, entry := range entries {
		listEntry := merged[entry.TablePath]
		listEntry.ID = entry.ID
		listEntry.TablePath = entry.TablePath
		listEntry.LootPools = append([]treasures.DropRef{}, entry.LootPools...)
		listEntry.UpdatedAt = entry.UpdatedAt
		listEntry.HasOverlay = true
		merged[entry.TablePath] = listEntry
	}
	out := make([]webui.TreasureListEntry, 0, len(merged))
	for _, entry := range merged {
		out = append(out, entry)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].TablePath < out[j].TablePath
	})
	return out
}

func treasuresMeta() webui.PageMeta {
	return webui.PageMeta{Title: "Treasures", CurrentPath: "/treasures"}
}

func defaultTreasureForm() webui.TreasureFormData {
	return webui.TreasureFormData{
		ID:            "",
		TablePath:     "",
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
		TablePath:     entry.TablePath,
		LootPoolsText: strings.Join(lines, "\n"),
		FieldErrors:   map[string]string{},
		IsEditing:     true,
		HasOverlay:    true,
	}
}

func parseTreasureForm(r *http.Request) (webui.TreasureFormData, treasures.SaveInput, map[string]string) {
	form := defaultTreasureForm()
	form.ID = strings.TrimSpace(r.Form.Get("id"))
	form.TablePath = strings.TrimSpace(r.Form.Get("tablePath"))
	form.LootPoolsText = r.Form.Get("lootPoolsText")
	errs := map[string]string{}
	input := treasures.SaveInput{
		ID:        form.ID,
		TablePath: form.TablePath,
		LootPools: parseTreasurePools(errs, form.LootPoolsText),
	}
	return form, input, errs
}

func findTreasureByTablePath(entries []treasures.TreasureEntry, tablePath string) (treasures.TreasureEntry, bool) {
	return findEntry(entries, tablePath, func(entry treasures.TreasureEntry) string { return entry.TablePath })
}

func hasTreasureSource(sourcePaths map[string]struct{}, tablePath string) bool {
	_, ok := sourcePaths[strings.TrimSpace(tablePath)]
	return ok
}

func noticeWithError(message string) *webui.Notice {
	return errorNotice(message)
}
