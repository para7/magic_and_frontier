package web

import (
	"net/http"
	"strings"

	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/webui"
	"tools2/app/views"
)

func (a App) itemsPage(w http.ResponseWriter, r *http.Request) {
	notice := consumeFlashNotice(w, r)
	state, err := a.deps.ItemRepo.LoadItemState()
	if err != nil {
		a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error())})
		return
	}
	skillState, err := a.deps.SkillRepo.LoadState()
	if err != nil {
		a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error())})
		return
	}
	a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Entries: state.Items, Notice: notice, Form: defaultItemForm(skillOptions(skillState.Entries))})
}

func (a App) itemsNewPage(w http.ResponseWriter, r *http.Request) {
	returnTo := queryReturnTo(r, itemMeta().CurrentPath)
	skillState, err := a.deps.SkillRepo.LoadState()
	if err != nil {
		form := defaultItemForm(nil)
		form.ReturnTo = returnTo
		a.renderItemForm(w, r, webui.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	form := defaultItemForm(skillOptions(skillState.Entries))
	form.ReturnTo = returnTo
	a.renderItemForm(w, r, webui.ItemsPageData{Meta: itemMeta(), Form: form})
}

func (a App) itemsEditPage(w http.ResponseWriter, r *http.Request) {
	returnTo := queryReturnTo(r, itemMeta().CurrentPath)
	state, err := a.deps.ItemRepo.LoadItemState()
	if err != nil {
		a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error())})
		return
	}
	skillState, err := a.deps.SkillRepo.LoadState()
	if err != nil {
		a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if entry, ok := findEntry(state.Items, id, func(entry items.ItemEntry) string { return entry.ID }); ok {
		form := itemEntryToForm(entry, skillOptions(skillState.Entries))
		form.ReturnTo = returnTo
		a.renderItemForm(w, r, webui.ItemsPageData{Meta: itemMeta(), Form: form})
		return
	}
	a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Entries: state.Items, Notice: errorNotice("Item not found.")})
}

func (a App) itemsSubmit(w http.ResponseWriter, r *http.Request) {
	a.itemsSave(w, r, false)
}

func (a App) itemsEditSubmit(w http.ResponseWriter, r *http.Request) {
	a.itemsSave(w, r, true)
}

func (a App) itemsSave(w http.ResponseWriter, r *http.Request, editing bool) {
	_ = r.ParseForm()
	returnTo := submittedReturnTo(r, itemMeta().CurrentPath)
	state, err := a.deps.ItemRepo.LoadItemState()
	if err != nil {
		form := defaultItemForm(nil)
		form.ReturnTo = returnTo
		a.renderItemForm(w, r, webui.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	skillState, err := a.deps.SkillRepo.LoadState()
	if err != nil {
		form := defaultItemForm(nil)
		form.ReturnTo = returnTo
		a.renderItemForm(w, r, webui.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	form, input, parseErrs := parseItemForm(r, skillState.Entries)
	form.IsEditing = editing
	form.ReturnTo = returnTo
	if editing {
		if _, ok := findEntry(state.Items, form.ID, func(entry items.ItemEntry) string { return entry.ID }); !ok {
			a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Entries: state.Items, Notice: errorNotice("Item not found.")})
			return
		}
	} else if _, ok := findEntry(state.Items, form.ID, func(entry items.ItemEntry) string { return entry.ID }); ok {
		parseErrs["id"] = "この ID は既に使用されています。"
	}
	result := items.ValidateSave(input, toIDSet(skillState.Entries, func(entry skills.SkillEntry) string { return entry.ID }), a.deps.Now())
	errors := mergeFieldErrors(parseErrs, result.FieldErrors)
	if len(errors) > 0 {
		form.FieldErrors = errors
		form.FormError = formErrorText(result.FormError)
		a.renderItemForm(w, r, webui.ItemsPageData{Meta: itemMeta(), Form: form})
		return
	}
	nextState, mode := items.Upsert(state, *result.Entry)
	if err := a.deps.ItemRepo.SaveItemState(nextState); err != nil {
		a.renderItemForm(w, r, webui.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	notice := successNotice(noticeText("Item", mode))
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Entries: nextState.Items, Notice: notice})
}

func (a App) itemsDelete(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	returnTo := submittedReturnTo(r, itemMeta().CurrentPath)
	id := strings.TrimSpace(r.PathValue("id"))
	state, err := a.deps.ItemRepo.LoadItemState()
	if err != nil {
		a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error())})
		return
	}
	nextState, ok := items.Delete(state, id)
	if !ok {
		a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Entries: state.Items, Notice: errorNotice("Item not found.")})
		return
	}
	if err := a.deps.ItemRepo.SaveItemState(nextState); err != nil {
		a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Entries: state.Items, Notice: errorNotice(err.Error())})
		return
	}
	notice := successNotice("Item deleted.")
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	a.renderItems(w, r, webui.ItemsPageData{Meta: itemMeta(), Entries: nextState.Items, Notice: notice})
}

func (a App) renderItems(w http.ResponseWriter, r *http.Request, data webui.ItemsPageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.ItemsShell(data))
		return
	}
	a.renderComponent(w, views.ItemsPage(data))
}

func (a App) renderItemForm(w http.ResponseWriter, r *http.Request, data webui.ItemsPageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.ItemFormShell(data))
		return
	}
	a.renderComponent(w, views.ItemFormPage(data))
}

func itemMeta() webui.PageMeta {
	return webui.PageMeta{Title: "Items", CurrentPath: "/items", Description: "アイテム出力を作成・管理します。複雑な NBT 項目はこの移行段階ではテキストのまま扱います。"}
}

func defaultItemForm(options []webui.ReferenceOption) webui.ItemFormData {
	return webui.ItemFormData{
		ID:                  "",
		ItemID:              "minecraft:stone",
		SkillOptions:        options,
		FieldErrors:         map[string]string{},
		CustomName:          "",
		Lore:                "",
		Enchantments:        "",
		CustomModelData:     "",
		RepairCost:          "",
		HideFlags:           "",
		PotionID:            "",
		CustomPotionColor:   "",
		CustomPotionEffects: "",
		AttributeModifiers:  "",
		CustomNBT:           "",
	}
}

func itemEntryToForm(entry items.ItemEntry, options []webui.ReferenceOption) webui.ItemFormData {
	return webui.ItemFormData{
		ID:                  entry.ID,
		ItemID:              entry.ItemID,
		SkillID:             entry.SkillID,
		SkillOptions:        options,
		CustomName:          entry.CustomName,
		Lore:                entry.Lore,
		Enchantments:        entry.Enchantments,
		Unbreakable:         entry.Unbreakable,
		CustomModelData:     entry.CustomModelData,
		RepairCost:          entry.RepairCost,
		HideFlags:           entry.HideFlags,
		PotionID:            entry.PotionID,
		CustomPotionColor:   entry.CustomPotionColor,
		CustomPotionEffects: entry.CustomPotionEffects,
		AttributeModifiers:  entry.AttributeModifiers,
		CustomNBT:           entry.CustomNBT,
		FieldErrors:         map[string]string{},
		IsEditing:           true,
	}
}

func parseItemForm(r *http.Request, skills []skills.SkillEntry) (webui.ItemFormData, items.SaveInput, map[string]string) {
	form := defaultItemForm(skillOptions(skills))
	form.ID = strings.TrimSpace(r.Form.Get("id"))
	form.ItemID = strings.TrimSpace(r.Form.Get("itemId"))
	form.SkillID = strings.TrimSpace(r.Form.Get("skillId"))
	form.CustomName = r.Form.Get("customName")
	form.Lore = r.Form.Get("lore")
	form.Enchantments = r.Form.Get("enchantments")
	form.Unbreakable = r.Form.Get("unbreakable") != ""
	form.CustomModelData = r.Form.Get("customModelData")
	form.RepairCost = r.Form.Get("repairCost")
	form.HideFlags = r.Form.Get("hideFlags")
	form.PotionID = r.Form.Get("potionId")
	form.CustomPotionColor = r.Form.Get("customPotionColor")
	form.CustomPotionEffects = r.Form.Get("customPotionEffects")
	form.AttributeModifiers = r.Form.Get("attributeModifiers")
	form.CustomNBT = r.Form.Get("customNbt")
	errs := map[string]string{}
	input := items.SaveInput{
		ID:                  form.ID,
		ItemID:              form.ItemID,
		SkillID:             form.SkillID,
		CustomName:          form.CustomName,
		Lore:                form.Lore,
		Enchantments:        form.Enchantments,
		Unbreakable:         form.Unbreakable,
		CustomModelData:     form.CustomModelData,
		RepairCost:          form.RepairCost,
		HideFlags:           form.HideFlags,
		PotionID:            form.PotionID,
		CustomPotionColor:   form.CustomPotionColor,
		CustomPotionEffects: form.CustomPotionEffects,
		AttributeModifiers:  form.AttributeModifiers,
		CustomNBT:           form.CustomNBT,
	}
	return form, input, errs
}

func itemOptions(entries []items.ItemEntry) []webui.ReferenceOption {
	options := make([]webui.ReferenceOption, 0, len(entries))
	for _, entry := range entries {
		options = append(options, webui.ReferenceOption{ID: entry.ID, Label: entry.ItemID})
	}
	return options
}

func skillOptions(entries []skills.SkillEntry) []webui.ReferenceOption {
	options := make([]webui.ReferenceOption, 0, len(entries))
	for _, entry := range entries {
		label := entry.Name
		if strings.TrimSpace(label) == "" {
			label = entry.ID
		}
		options = append(options, webui.ReferenceOption{ID: entry.ID, Label: label})
	}
	return options
}

func itemIDSet(state items.ItemState) map[string]struct{} {
	return toIDSet(state.Items, func(entry items.ItemEntry) string { return entry.ID })
}
