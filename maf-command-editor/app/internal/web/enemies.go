package web

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	dmaster "tools2/app/internal/domain/master"
	"tools2/app/internal/web/views"
)

func (a App) enemiesPage(w http.ResponseWriter, r *http.Request) {
	notice := consumeFlashNotice(w, r)
	state, err := a.loadEnemyStateFromMaster()
	if err != nil {
		a.renderEnemies(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	a.renderEnemies(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Entries: state.Entries, Notice: notice})
}

func (a App) enemiesNewPage(w http.ResponseWriter, r *http.Request) {
	returnTo := queryReturnTo(r, enemiesMeta().CurrentPath)
	enemySkillState, err := a.loadEnemySkillStateFromMaster()
	if err != nil {
		form := defaultEnemyForm(nil)
		form.ReturnTo = returnTo
		a.renderEnemyForm(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	form := defaultEnemyForm(enemySkillState.Entries)
	form.ReturnTo = returnTo
	a.renderEnemyForm(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Form: form})
}

func (a App) enemiesEditPage(w http.ResponseWriter, r *http.Request) {
	returnTo := queryReturnTo(r, enemiesMeta().CurrentPath)
	state, err := a.loadEnemyStateFromMaster()
	if err != nil {
		a.renderEnemies(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	enemySkillState, err := a.loadEnemySkillStateFromMaster()
	if err != nil {
		a.renderEnemies(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if entry, ok := findEntry(state.Entries, id, func(entry enemies.EnemyEntry) string { return entry.ID }); ok {
		form := enemyEntryToForm(entry, enemySkillOptions(enemySkillState.Entries))
		form.ReturnTo = returnTo
		a.renderEnemyForm(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Form: form})
		return
	}
	a.renderEnemies(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Entries: state.Entries, Notice: errorNotice("Enemy not found.")})
}

func (a App) enemiesSubmit(w http.ResponseWriter, r *http.Request) {
	a.enemiesSave(w, r, false)
}

func (a App) enemiesEditSubmit(w http.ResponseWriter, r *http.Request) {
	a.enemiesSave(w, r, true)
}

func (a App) enemiesSave(w http.ResponseWriter, r *http.Request, editing bool) {
	_ = r.ParseForm()
	returnTo := submittedReturnTo(r, enemiesMeta().CurrentPath)
	master, err := a.masterOrErr()
	if err != nil {
		form := defaultEnemyForm(nil)
		form.IsEditing = editing
		form.ReturnTo = returnTo
		a.renderEnemyForm(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	enemyState, err := a.loadEnemyStateFromMaster()
	if err != nil {
		form := defaultEnemyForm(nil)
		form.IsEditing = editing
		form.ReturnTo = returnTo
		a.renderEnemyForm(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	enemySkillState, err := a.loadEnemySkillStateFromMaster()
	if err != nil {
		form := defaultEnemyForm(nil)
		form.IsEditing = editing
		form.ReturnTo = returnTo
		a.renderEnemyForm(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	form, input, parseErrs := parseEnemyForm(r, enemySkillState.Entries)
	form.IsEditing = editing
	form.ReturnTo = returnTo
	if editing {
		if _, ok := findEntry(enemyState.Entries, form.ID, func(entry enemies.EnemyEntry) string { return entry.ID }); !ok {
			a.renderEnemies(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Entries: enemyState.Entries, Notice: errorNotice("Enemy not found.")})
			return
		}
	} else if _, ok := findEntry(enemyState.Entries, form.ID, func(entry enemies.EnemyEntry) string { return entry.ID }); ok {
		parseErrs["id"] = "この ID は既に使用されています。"
	}
	result := master.Enemies().Validate(input, master)
	fieldErrs := mergeFieldErrors(parseErrs, mapFieldErrors(result.FieldErrors, mapEnemyField))
	if len(fieldErrs) > 0 {
		form.FieldErrors = fieldErrs
		form.FormError = formErrorText(result.FormError)
		a.renderEnemyForm(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Form: form})
		return
	}
	mode := common.SaveModeCreated
	if editing {
		mode = common.SaveModeUpdated
		if err := master.Enemies().Update(*result.Entry, master); err != nil {
			form.FormError = formErrorText(err.Error())
			a.renderEnemyForm(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Form: form})
			return
		}
	} else {
		if err := master.Enemies().Create(*result.Entry, master); err != nil {
			if errors.Is(err, dmaster.ErrDuplicateID) {
				form.FieldErrors = mergeFieldErrors(form.FieldErrors, map[string]string{"id": "この ID は既に使用されています。"})
			} else {
				form.FormError = formErrorText(err.Error())
			}
			a.renderEnemyForm(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Form: form})
			return
		}
	}
	if err := master.Enemies().Save(); err != nil {
		a.renderEnemyForm(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	nextState, err := a.loadEnemyStateFromMaster()
	if err != nil {
		a.renderEnemyForm(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	notice := successNotice(noticeText("Enemy", mode))
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	a.renderEnemies(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) enemiesDelete(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	returnTo := submittedReturnTo(r, enemiesMeta().CurrentPath)
	master, err := a.masterOrErr()
	if err != nil {
		a.renderEnemies(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	state, err := a.loadEnemyStateFromMaster()
	if err != nil {
		a.renderEnemies(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.PathValue("id"))
	if err := master.Enemies().Delete(id, master); err != nil {
		a.renderEnemies(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Entries: state.Entries, Notice: errorNotice("Enemy not found.")})
		return
	}
	if err := master.Enemies().Save(); err != nil {
		a.renderEnemies(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Entries: state.Entries, Notice: errorNotice(err.Error())})
		return
	}
	nextState, err := a.loadEnemyStateFromMaster()
	if err != nil {
		a.renderEnemies(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Entries: state.Entries, Notice: errorNotice(err.Error())})
		return
	}
	notice := successNotice("Enemy deleted.")
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	a.renderEnemies(w, r, views.EnemiesPageData{Meta: enemiesMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) renderEnemies(w http.ResponseWriter, r *http.Request, data views.EnemiesPageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.EnemiesShell(data))
		return
	}
	a.renderComponent(w, views.EnemiesPage(data))
}

func (a App) renderEnemyForm(w http.ResponseWriter, r *http.Request, data views.EnemiesPageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.EnemyFormShell(data))
		return
	}
	a.renderComponent(w, views.EnemyFormPage(data))
}

func enemiesMeta() views.PageMeta {
	return views.PageMeta{Title: "Enemies", CurrentPath: "/enemies"}
}

func defaultEnemyForm(entries []enemyskills.EnemySkillEntry) views.EnemyFormData {
	return views.EnemyFormData{
		ID:                "",
		MobType:           "minecraft:zombie",
		HP:                "20",
		DropMode:          "replace",
		FieldErrors:       map[string]string{},
		EnemySkillOptions: enemySkillOptions(entries),
	}
}

func enemyEntryToForm(entry enemies.EnemyEntry, options []views.ReferenceOption) views.EnemyFormData {
	form := views.EnemyFormData{
		ID:                entry.ID,
		MobType:           entry.MobType,
		Name:              entry.Name,
		HP:                strconv.FormatFloat(entry.HP, 'f', -1, 64),
		Memo:              entry.Memo,
		DropMode:          entry.DropMode,
		EnemySkillIDs:     append([]string{}, entry.EnemySkillIDs...),
		EnemySkillOptions: options,
		EquipmentText:     formatEquipmentText(entry.Equipment),
		DropsText:         formatEnemyDropsText(entry.Drops),
		FieldErrors:       map[string]string{},
		IsEditing:         true,
	}
	if entry.Attack != nil {
		form.Attack = strconv.FormatFloat(*entry.Attack, 'f', -1, 64)
	}
	if entry.Defense != nil {
		form.Defense = strconv.FormatFloat(*entry.Defense, 'f', -1, 64)
	}
	if entry.MoveSpeed != nil {
		form.MoveSpeed = strconv.FormatFloat(*entry.MoveSpeed, 'f', -1, 64)
	}
	return form
}

func parseEnemyForm(r *http.Request, enemySkillEntries []enemyskills.EnemySkillEntry) (views.EnemyFormData, enemies.SaveInput, map[string]string) {
	form := defaultEnemyForm(enemySkillEntries)
	form.ID = strings.TrimSpace(r.Form.Get("id"))
	form.MobType = strings.TrimSpace(r.Form.Get("mobType"))
	form.Name = r.Form.Get("name")
	form.HP = strings.TrimSpace(r.Form.Get("hp"))
	form.Memo = r.Form.Get("memo")
	form.Attack = strings.TrimSpace(r.Form.Get("attack"))
	form.Defense = strings.TrimSpace(r.Form.Get("defense"))
	form.MoveSpeed = strings.TrimSpace(r.Form.Get("moveSpeed"))
	form.DropMode = strings.TrimSpace(r.Form.Get("dropMode"))
	form.EnemySkillIDs = append([]string{}, r.Form["enemySkillIds"]...)
	form.EquipmentText = r.Form.Get("equipmentText")
	form.DropsText = r.Form.Get("dropsText")
	errs := map[string]string{}

	input := enemies.SaveInput{
		ID:            form.ID,
		MobType:       form.MobType,
		Name:          form.Name,
		HP:            parseRequiredFloat(errs, "hp", form.HP),
		Memo:          form.Memo,
		Attack:        parseOptionalFloat(errs, "attack", form.Attack),
		Defense:       parseOptionalFloat(errs, "defense", form.Defense),
		MoveSpeed:     parseOptionalFloat(errs, "moveSpeed", form.MoveSpeed),
		Equipment:     parseEquipment(errs, form.EquipmentText),
		EnemySkillIDs: append([]string{}, form.EnemySkillIDs...),
		DropMode:      form.DropMode,
		Drops:         parseEnemyDrops(errs, form.DropsText),
	}
	return form, input, errs
}

func formatEquipmentText(equipment enemies.Equipment) string {
	lines := make([]string, 0, 6)
	appendSlot := func(name string, slot *enemies.EquipmentSlot) {
		if slot == nil {
			return
		}
		dropChance := ""
		if slot.DropChance != nil {
			dropChance = strconv.FormatFloat(*slot.DropChance, 'f', -1, 64)
		}
		lines = append(lines, strings.Join([]string{name, slot.Kind, slot.RefID, strconv.Itoa(slot.Count), dropChance}, ","))
	}
	appendSlot("mainhand", equipment.Mainhand)
	appendSlot("offhand", equipment.Offhand)
	appendSlot("head", equipment.Head)
	appendSlot("chest", equipment.Chest)
	appendSlot("legs", equipment.Legs)
	appendSlot("feet", equipment.Feet)
	return strings.Join(lines, "\n")
}

func formatEnemyDropsText(drops []enemies.DropRef) string {
	lines := make([]string, 0, len(drops))
	for _, drop := range drops {
		countMin := ""
		countMax := ""
		if drop.CountMin != nil {
			countMin = strconv.FormatFloat(*drop.CountMin, 'f', -1, 64)
		}
		if drop.CountMax != nil {
			countMax = strconv.FormatFloat(*drop.CountMax, 'f', -1, 64)
		}
		lines = append(lines, strings.Join([]string{
			drop.Kind,
			drop.RefID,
			strconv.FormatFloat(drop.Weight, 'f', -1, 64),
			countMin,
			countMax,
		}, ","))
	}
	return strings.Join(lines, "\n")
}
