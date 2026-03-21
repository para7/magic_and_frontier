package web

import (
	"net/http"
	"strconv"
	"strings"

	"tools2/app/internal/application"
	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/spawntables"
	"tools2/app/internal/webui"
	"tools2/app/views"
)

func (a App) spawnTablesPage(w http.ResponseWriter, r *http.Request) {
	notice := consumeFlashNotice(w, r)
	state, err := a.deps.SpawnTableRepo.LoadState()
	if err != nil {
		a.renderSpawnTables(w, r, webui.SpawnTablesPageData{Meta: spawnTablesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	a.renderSpawnTables(w, r, webui.SpawnTablesPageData{Meta: spawnTablesMeta(), Entries: state.Entries, Notice: notice})
}

func (a App) spawnTablesNewPage(w http.ResponseWriter, r *http.Request) {
	form := defaultSpawnTableForm()
	form.ReturnTo = queryReturnTo(r, spawnTablesMeta().CurrentPath)
	a.renderSpawnTableForm(w, r, webui.SpawnTablesPageData{Meta: spawnTablesMeta(), Form: form})
}

func (a App) spawnTablesEditPage(w http.ResponseWriter, r *http.Request) {
	returnTo := queryReturnTo(r, spawnTablesMeta().CurrentPath)
	state, err := a.deps.SpawnTableRepo.LoadState()
	if err != nil {
		a.renderSpawnTables(w, r, webui.SpawnTablesPageData{Meta: spawnTablesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if entry, ok := findEntry(state.Entries, id, func(entry spawntables.SpawnTableEntry) string { return entry.ID }); ok {
		form := spawnTableEntryToForm(entry)
		form.ReturnTo = returnTo
		a.renderSpawnTableForm(w, r, webui.SpawnTablesPageData{Meta: spawnTablesMeta(), Form: form})
		return
	}
	a.renderSpawnTables(w, r, webui.SpawnTablesPageData{Meta: spawnTablesMeta(), Entries: state.Entries, Notice: errorNotice("Spawn table not found.")})
}

func (a App) spawnTablesSubmit(w http.ResponseWriter, r *http.Request) {
	a.spawnTablesSave(w, r, false)
}

func (a App) spawnTablesEditSubmit(w http.ResponseWriter, r *http.Request) {
	a.spawnTablesSave(w, r, true)
}

func (a App) spawnTablesSave(w http.ResponseWriter, r *http.Request, editing bool) {
	_ = r.ParseForm()
	returnTo := submittedReturnTo(r, spawnTablesMeta().CurrentPath)
	state, err := a.deps.SpawnTableRepo.LoadState()
	if err != nil {
		form := defaultSpawnTableForm()
		form.IsEditing = editing
		form.ReturnTo = returnTo
		a.renderSpawnTableForm(w, r, webui.SpawnTablesPageData{Meta: spawnTablesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	enemyState, err := a.deps.EnemyRepo.LoadState()
	if err != nil {
		form := defaultSpawnTableForm()
		form.IsEditing = editing
		form.ReturnTo = returnTo
		a.renderSpawnTableForm(w, r, webui.SpawnTablesPageData{Meta: spawnTablesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}

	form, input, parseErrs := parseSpawnTableForm(r)
	form.IsEditing = editing
	form.ReturnTo = returnTo
	if editing {
		if _, ok := findEntry(state.Entries, form.ID, func(entry spawntables.SpawnTableEntry) string { return entry.ID }); !ok {
			a.renderSpawnTables(w, r, webui.SpawnTablesPageData{Meta: spawnTablesMeta(), Entries: state.Entries, Notice: errorNotice("Spawn table not found.")})
			return
		}
	} else if strings.TrimSpace(input.ID) == "" {
		id, allocErr := application.NewService(a.cfg, a.deps).AllocateID("spawntable")
		if allocErr != nil {
			a.renderSpawnTableForm(w, r, webui.SpawnTablesPageData{Meta: spawnTablesMeta(), Notice: errorNotice(allocErr.Error()), Form: form})
			return
		}
		input.ID = id
		form.ID = id
	}

	result := spawntables.ValidateSave(input, toIDSet(enemyState.Entries, func(entry enemies.EnemyEntry) string { return entry.ID }), a.deps.Now())
	errors := mergeFieldErrors(parseErrs, mapFieldErrors(result.FieldErrors, mapSpawnTableField))
	if len(errors) > 0 {
		form.FieldErrors = errors
		form.FormError = formErrorText(result.FormError)
		a.renderSpawnTableForm(w, r, webui.SpawnTablesPageData{Meta: spawnTablesMeta(), Form: form})
		return
	}
	if conflictID, ok := spawntables.FirstOverlap(state.Entries, *result.Entry); ok {
		form.FieldErrors = map[string]string{"replacementsText": "Range overlaps with " + conflictID + "."}
		form.FormError = formErrorText("Validation failed. Fix the highlighted fields.")
		a.renderSpawnTableForm(w, r, webui.SpawnTablesPageData{Meta: spawnTablesMeta(), Form: form})
		return
	}

	nextState, mode := common.UpsertEntries(state, *result.Entry, func(entry spawntables.SpawnTableEntry) string { return entry.ID })
	if err := a.deps.SpawnTableRepo.SaveState(nextState); err != nil {
		a.renderSpawnTableForm(w, r, webui.SpawnTablesPageData{Meta: spawnTablesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	notice := successNotice(noticeText("Spawn table", mode))
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	a.renderSpawnTables(w, r, webui.SpawnTablesPageData{Meta: spawnTablesMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) spawnTablesDelete(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	returnTo := submittedReturnTo(r, spawnTablesMeta().CurrentPath)
	state, err := a.deps.SpawnTableRepo.LoadState()
	if err != nil {
		a.renderSpawnTables(w, r, webui.SpawnTablesPageData{Meta: spawnTablesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.PathValue("id"))
	nextState, ok := common.DeleteEntries(state, id, func(entry spawntables.SpawnTableEntry) string { return entry.ID })
	if !ok {
		a.renderSpawnTables(w, r, webui.SpawnTablesPageData{Meta: spawnTablesMeta(), Entries: state.Entries, Notice: errorNotice("Spawn table not found.")})
		return
	}
	if err := a.deps.SpawnTableRepo.SaveState(nextState); err != nil {
		a.renderSpawnTables(w, r, webui.SpawnTablesPageData{Meta: spawnTablesMeta(), Entries: state.Entries, Notice: errorNotice(err.Error())})
		return
	}
	notice := successNotice("Spawn table deleted.")
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	a.renderSpawnTables(w, r, webui.SpawnTablesPageData{Meta: spawnTablesMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) renderSpawnTables(w http.ResponseWriter, r *http.Request, data webui.SpawnTablesPageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.SpawnTablesShell(data))
		return
	}
	a.renderComponent(w, views.SpawnTablesPage(data))
}

func (a App) renderSpawnTableForm(w http.ResponseWriter, r *http.Request, data webui.SpawnTablesPageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.SpawnTableFormShell(data))
		return
	}
	a.renderComponent(w, views.SpawnTableFormPage(data))
}

func spawnTablesMeta() webui.PageMeta {
	return webui.PageMeta{Title: "Spawn Tables", CurrentPath: "/spawn-tables", Description: "source mob + dimension + xyz range ごとの敵置換テーブルを管理します。"}
}

func defaultSpawnTableForm() webui.SpawnTableFormData {
	return webui.SpawnTableFormData{
		Dimension:        "minecraft:overworld",
		DimensionOptions: spawnTableDimensionOptions(),
		BaseMobWeight:    "8000",
		FieldErrors:      map[string]string{},
	}
}

func spawnTableEntryToForm(entry spawntables.SpawnTableEntry) webui.SpawnTableFormData {
	return webui.SpawnTableFormData{
		ID:               entry.ID,
		SourceMobType:    entry.SourceMobType,
		Dimension:        entry.Dimension,
		DimensionOptions: spawnTableDimensionOptions(),
		MinX:             strconv.Itoa(entry.MinX),
		MaxX:             strconv.Itoa(entry.MaxX),
		MinY:             strconv.Itoa(entry.MinY),
		MaxY:             strconv.Itoa(entry.MaxY),
		MinZ:             strconv.Itoa(entry.MinZ),
		MaxZ:             strconv.Itoa(entry.MaxZ),
		BaseMobWeight:    strconv.Itoa(entry.BaseMobWeight),
		ReplacementsText: formatSpawnTableReplacements(entry.Replacements),
		FieldErrors:      map[string]string{},
		IsEditing:        true,
	}
}

func parseSpawnTableForm(r *http.Request) (webui.SpawnTableFormData, spawntables.SaveInput, map[string]string) {
	form := defaultSpawnTableForm()
	form.ID = strings.TrimSpace(r.Form.Get("id"))
	form.SourceMobType = strings.TrimSpace(r.Form.Get("sourceMobType"))
	form.Dimension = strings.TrimSpace(r.Form.Get("dimension"))
	form.MinX = strings.TrimSpace(r.Form.Get("minX"))
	form.MaxX = strings.TrimSpace(r.Form.Get("maxX"))
	form.MinY = strings.TrimSpace(r.Form.Get("minY"))
	form.MaxY = strings.TrimSpace(r.Form.Get("maxY"))
	form.MinZ = strings.TrimSpace(r.Form.Get("minZ"))
	form.MaxZ = strings.TrimSpace(r.Form.Get("maxZ"))
	form.BaseMobWeight = strings.TrimSpace(r.Form.Get("baseMobWeight"))
	form.ReplacementsText = r.Form.Get("replacementsText")
	err := map[string]string{}

	input := spawntables.SaveInput{
		ID:            form.ID,
		SourceMobType: form.SourceMobType,
		Dimension:     form.Dimension,
		MinX:          parseRequiredIntField(err, "minX", form.MinX),
		MaxX:          parseRequiredIntField(err, "maxX", form.MaxX),
		MinY:          parseRequiredIntField(err, "minY", form.MinY),
		MaxY:          parseRequiredIntField(err, "maxY", form.MaxY),
		MinZ:          parseRequiredIntField(err, "minZ", form.MinZ),
		MaxZ:          parseRequiredIntField(err, "maxZ", form.MaxZ),
		BaseMobWeight: parseRequiredIntField(err, "baseMobWeight", form.BaseMobWeight),
		Replacements:  parseSpawnTableReplacements(err, form.ReplacementsText),
	}
	return form, input, err
}

func parseRequiredIntField(errs map[string]string, key, value string) int {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		errs[key] = "Required."
		return 0
	}
	parsed, err := strconv.Atoi(trimmed)
	if err != nil {
		errs[key] = "Must be a number."
		return 0
	}
	return parsed
}

func formatSpawnTableReplacements(replacements []spawntables.ReplacementEntry) string {
	lines := make([]string, 0, len(replacements))
	for _, replacement := range replacements {
		lines = append(lines, replacement.EnemyID+","+strconv.Itoa(replacement.Weight))
	}
	return strings.Join(lines, "\n")
}

func spawnTableDimensionOptions() []webui.SelectOption {
	return []webui.SelectOption{
		{Value: "minecraft:overworld", Label: "minecraft:overworld"},
		{Value: "minecraft:the_nether", Label: "minecraft:the_nether"},
		{Value: "minecraft:the_end", Label: "minecraft:the_end"},
	}
}
