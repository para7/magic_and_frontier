package web

import (
	"net/http"
	"strconv"
	"strings"

	"tools2/app/internal/application"
	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/webui"
	"tools2/app/views"
)

func (a App) enemiesPage(w http.ResponseWriter, r *http.Request) {
	notice := consumeFlashNotice(w, r)
	state, err := a.deps.EnemyRepo.LoadState()
	if err != nil {
		a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Entries: state.Entries, Notice: notice})
}

func (a App) enemiesNewPage(w http.ResponseWriter, r *http.Request) {
	enemySkillState, err := a.deps.EnemySkillRepo.LoadState()
	if err != nil {
		a.renderEnemyForm(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: defaultEnemyForm(nil)})
		return
	}
	a.renderEnemyForm(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Form: defaultEnemyForm(enemySkillState.Entries)})
}

func (a App) enemiesEditPage(w http.ResponseWriter, r *http.Request) {
	state, err := a.deps.EnemyRepo.LoadState()
	if err != nil {
		a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	enemySkillState, err := a.deps.EnemySkillRepo.LoadState()
	if err != nil {
		a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if entry, ok := findEntry(state.Entries, id, func(entry enemies.EnemyEntry) string { return entry.ID }); ok {
		form := enemyEntryToForm(entry, enemySkillOptions(enemySkillState.Entries))
		a.renderEnemyForm(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Form: form})
		return
	}
	a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Entries: state.Entries, Notice: errorNotice("Enemy not found.")})
}

func (a App) enemiesSubmit(w http.ResponseWriter, r *http.Request) {
	a.enemiesSave(w, r, false)
}

func (a App) enemiesEditSubmit(w http.ResponseWriter, r *http.Request) {
	a.enemiesSave(w, r, true)
}

func (a App) enemiesSave(w http.ResponseWriter, r *http.Request, editing bool) {
	_ = r.ParseForm()
	enemyState, err := a.deps.EnemyRepo.LoadState()
	if err != nil {
		form := defaultEnemyForm(nil)
		form.IsEditing = editing
		a.renderEnemyForm(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	enemySkillState, err := a.deps.EnemySkillRepo.LoadState()
	if err != nil {
		form := defaultEnemyForm(nil)
		form.IsEditing = editing
		a.renderEnemyForm(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	itemState, err := a.deps.ItemRepo.LoadItemState()
	if err != nil {
		form := defaultEnemyForm(enemySkillState.Entries)
		form.IsEditing = editing
		a.renderEnemyForm(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	grimoireState, err := a.deps.GrimoireRepo.LoadGrimoireState()
	if err != nil {
		form := defaultEnemyForm(enemySkillState.Entries)
		form.IsEditing = editing
		a.renderEnemyForm(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	form, input, parseErrs := parseEnemyForm(r, enemySkillState.Entries)
	form.IsEditing = editing
	if editing {
		if _, ok := findEntry(enemyState.Entries, form.ID, func(entry enemies.EnemyEntry) string { return entry.ID }); !ok {
			a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Entries: enemyState.Entries, Notice: errorNotice("Enemy not found.")})
			return
		}
	} else if strings.TrimSpace(input.ID) == "" {
		id, allocErr := application.NewService(a.cfg, a.deps).AllocateID("enemy")
		if allocErr != nil {
			a.renderEnemyForm(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(allocErr.Error()), Form: form})
			return
		}
		input.ID = id
		form.ID = id
	}
	result := enemies.ValidateSave(input, enemySkillIDSet(enemySkillState), itemIDSet(itemState), grimoireIDSet(grimoireState), a.deps.Now())
	errors := mergeFieldErrors(parseErrs, mapFieldErrors(result.FieldErrors, mapEnemyField))
	if len(errors) > 0 {
		form.FieldErrors = errors
		form.FormError = formErrorText(result.FormError)
		a.renderEnemyForm(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Form: form})
		return
	}
	nextState, mode := common.UpsertEntries(enemyState, *result.Entry, func(entry enemies.EnemyEntry) string { return entry.ID })
	if err := a.deps.EnemyRepo.SaveState(nextState); err != nil {
		a.renderEnemyForm(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	notice := successNotice(noticeText("Enemy", mode))
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, "/enemies", notice) {
		return
	}
	a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) enemiesDelete(w http.ResponseWriter, r *http.Request) {
	state, err := a.deps.EnemyRepo.LoadState()
	if err != nil {
		a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.PathValue("id"))
	nextState, ok := common.DeleteEntries(state, id, func(entry enemies.EnemyEntry) string { return entry.ID })
	if !ok {
		a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Entries: state.Entries, Notice: errorNotice("Enemy not found.")})
		return
	}
	if err := a.deps.EnemyRepo.SaveState(nextState); err != nil {
		a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Entries: state.Entries, Notice: errorNotice(err.Error())})
		return
	}
	notice := successNotice("Enemy deleted.")
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, "/enemies", notice) {
		return
	}
	a.renderEnemies(w, r, webui.EnemiesPageData{Meta: enemiesMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) renderEnemies(w http.ResponseWriter, r *http.Request, data webui.EnemiesPageData) {
	if isHX(r) {
		a.renderComponent(w, views.EnemiesShell(data))
		return
	}
	a.renderComponent(w, views.EnemiesPage(data))
}

func (a App) renderEnemyForm(w http.ResponseWriter, r *http.Request, data webui.EnemiesPageData) {
	if isHX(r) {
		a.renderComponent(w, views.EnemyFormShell(data))
		return
	}
	a.renderComponent(w, views.EnemyFormPage(data))
}

func enemiesMeta() webui.PageMeta {
	return webui.PageMeta{Title: "Enemies", CurrentPath: "/enemies", Description: "enemy の mob type、装備、drop mode、直接ドロップを管理します。"}
}

func defaultEnemyForm(entries []enemyskills.EnemySkillEntry) webui.EnemyFormData {
	return webui.EnemyFormData{
		ID:                "",
		MobType:           "minecraft:zombie",
		HP:                "20",
		DropMode:          "replace",
		FieldErrors:       map[string]string{},
		EnemySkillOptions: enemySkillOptions(entries),
	}
}

func enemyEntryToForm(entry enemies.EnemyEntry, options []webui.ReferenceOption) webui.EnemyFormData {
	form := webui.EnemyFormData{
		ID:                entry.ID,
		MobType:           entry.MobType,
		Name:              entry.Name,
		HP:                strconv.FormatFloat(entry.HP, 'f', -1, 64),
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

func parseEnemyForm(r *http.Request, enemySkillEntries []enemyskills.EnemySkillEntry) (webui.EnemyFormData, enemies.SaveInput, map[string]string) {
	form := defaultEnemyForm(enemySkillEntries)
	form.ID = strings.TrimSpace(r.Form.Get("id"))
	form.MobType = strings.TrimSpace(r.Form.Get("mobType"))
	form.Name = r.Form.Get("name")
	form.HP = strings.TrimSpace(r.Form.Get("hp"))
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
