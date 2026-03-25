package web

import (
	"errors"
	"net/http"
	"strings"

	"maf-command-editor/app/internal/domain/common"
	"maf-command-editor/app/internal/domain/entity/items"
	"maf-command-editor/app/internal/domain/entity/skills"
	dmaster "maf-command-editor/app/internal/domain/master"
	"maf-command-editor/app/internal/web/views"
)

func (a App) itemsPage(w http.ResponseWriter, r *http.Request) {
	notice := consumeFlashNotice(w, r)
	state, err := a.loadItemStateFromMaster()
	if err != nil {
		a.renderItems(w, r, views.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error())})
		return
	}
	skillState, err := a.loadSkillStateFromMaster()
	if err != nil {
		a.renderItems(w, r, views.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error())})
		return
	}
	a.renderItems(w, r, views.ItemsPageData{Meta: itemMeta(), Entries: state.Entries, Notice: notice, Form: defaultItemForm(skillOptions(skillState.Entries))})
}

func (a App) itemsNewPage(w http.ResponseWriter, r *http.Request) {
	returnTo := queryReturnTo(r, itemMeta().CurrentPath)
	skillState, err := a.loadSkillStateFromMaster()
	if err != nil {
		form := defaultItemForm(nil)
		form.ReturnTo = returnTo
		a.renderItemForm(w, r, views.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	form := defaultItemForm(skillOptions(skillState.Entries))
	form.ReturnTo = returnTo
	a.renderItemForm(w, r, views.ItemsPageData{Meta: itemMeta(), Form: form})
}

func (a App) itemsEditPage(w http.ResponseWriter, r *http.Request) {
	returnTo := queryReturnTo(r, itemMeta().CurrentPath)
	state, err := a.loadItemStateFromMaster()
	if err != nil {
		a.renderItems(w, r, views.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error())})
		return
	}
	skillState, err := a.loadSkillStateFromMaster()
	if err != nil {
		a.renderItems(w, r, views.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error())})
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if entry, ok := findEntry(state.Entries, id, func(entry items.ItemEntry) string { return entry.ID }); ok {
		form := itemEntryToForm(entry, skillOptions(skillState.Entries))
		form.ReturnTo = returnTo
		a.renderItemForm(w, r, views.ItemsPageData{Meta: itemMeta(), Form: form})
		return
	}
	a.renderItems(w, r, views.ItemsPageData{Meta: itemMeta(), Entries: state.Entries, Notice: errorNotice("Item not found.")})
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
	master, err := a.masterOrErr()
	if err != nil {
		form := defaultItemForm(nil)
		form.ReturnTo = returnTo
		a.renderItemForm(w, r, views.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	state, err := a.loadItemStateFromMaster()
	if err != nil {
		form := defaultItemForm(nil)
		form.ReturnTo = returnTo
		a.renderItemForm(w, r, views.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	skillState, err := a.loadSkillStateFromMaster()
	if err != nil {
		form := defaultItemForm(nil)
		form.ReturnTo = returnTo
		a.renderItemForm(w, r, views.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	form, input, parseErrs := parseItemForm(r, skillState.Entries)
	form.IsEditing = editing
	form.ReturnTo = returnTo
	if editing {
		if _, ok := findEntry(state.Entries, form.ID, func(entry items.ItemEntry) string { return entry.ID }); !ok {
			a.renderItems(w, r, views.ItemsPageData{Meta: itemMeta(), Entries: state.Entries, Notice: errorNotice("Item not found.")})
			return
		}
	} else if _, ok := findEntry(state.Entries, form.ID, func(entry items.ItemEntry) string { return entry.ID }); ok {
		parseErrs["id"] = "この ID は既に使用されています。"
	}
	result := master.Items().Validate(input, master)
	fieldErrs := mergeFieldErrors(parseErrs, result.FieldErrors)
	if len(fieldErrs) > 0 {
		form.FieldErrors = fieldErrs
		form.ShowEnchantmentsDetail = fieldErrs["enchantments"] != ""
		form.FormError = formErrorText(result.FormError)
		a.renderItemForm(w, r, views.ItemsPageData{Meta: itemMeta(), Form: form})
		return
	}
	mode := common.SaveModeCreated
	if editing {
		mode = common.SaveModeUpdated
		if err := master.Items().Update(*result.Entry, master); err != nil {
			form.FormError = formErrorText(err.Error())
			a.renderItemForm(w, r, views.ItemsPageData{Meta: itemMeta(), Form: form})
			return
		}
	} else {
		if err := master.Items().Create(*result.Entry, master); err != nil {
			if errors.Is(err, dmaster.ErrDuplicateID) {
				form.FieldErrors = mergeFieldErrors(form.FieldErrors, map[string]string{"id": "この ID は既に使用されています。"})
			} else {
				form.FormError = formErrorText(err.Error())
			}
			a.renderItemForm(w, r, views.ItemsPageData{Meta: itemMeta(), Form: form})
			return
		}
	}
	if err := master.Items().Save(); err != nil {
		a.renderItemForm(w, r, views.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	nextState, err := a.loadItemStateFromMaster()
	if err != nil {
		a.renderItemForm(w, r, views.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error()), Form: form})
		return
	}
	notice := successNotice(noticeText("Item", mode))
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	a.renderItems(w, r, views.ItemsPageData{Meta: itemMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) itemsDelete(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	returnTo := submittedReturnTo(r, itemMeta().CurrentPath)
	id := strings.TrimSpace(r.PathValue("id"))
	master, err := a.masterOrErr()
	if err != nil {
		a.renderItems(w, r, views.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error())})
		return
	}
	state, err := a.loadItemStateFromMaster()
	if err != nil {
		a.renderItems(w, r, views.ItemsPageData{Meta: itemMeta(), Notice: errorNotice(err.Error())})
		return
	}
	if err := master.Items().Delete(id, master); err != nil {
		a.renderItems(w, r, views.ItemsPageData{Meta: itemMeta(), Entries: state.Entries, Notice: errorNotice("Item not found.")})
		return
	}
	if err := master.Items().Save(); err != nil {
		a.renderItems(w, r, views.ItemsPageData{Meta: itemMeta(), Entries: state.Entries, Notice: errorNotice(err.Error())})
		return
	}
	nextState, err := a.loadItemStateFromMaster()
	if err != nil {
		a.renderItems(w, r, views.ItemsPageData{Meta: itemMeta(), Entries: state.Entries, Notice: errorNotice(err.Error())})
		return
	}
	notice := successNotice("Item deleted.")
	setToast(w, notice.Text)
	if redirectWithNotice(w, r, returnTo, notice) {
		return
	}
	a.renderItems(w, r, views.ItemsPageData{Meta: itemMeta(), Entries: nextState.Entries, Notice: notice})
}

func (a App) renderItems(w http.ResponseWriter, r *http.Request, data views.ItemsPageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.ItemsShell(data))
		return
	}
	a.renderComponent(w, views.ItemsPage(data))
}

func (a App) renderItemForm(w http.ResponseWriter, r *http.Request, data views.ItemsPageData) {
	data.Meta = applyPageMeta(r, data.Meta)
	if isHX(r) {
		a.renderComponent(w, views.ItemFormShell(data))
		return
	}
	a.renderComponent(w, views.ItemFormPage(data))
}

func itemMeta() views.PageMeta {
	return views.PageMeta{Title: "Items", CurrentPath: "/items"}
}

func defaultItemForm(options []views.ReferenceOption) views.ItemFormData {
	enchantmentOptions, selectedEnchantments := itemFormEnchantmentsFromText("")
	return views.ItemFormData{
		ID:                   "",
		ItemID:               "minecraft:stone",
		SkillOptions:         options,
		FieldErrors:          map[string]string{},
		CustomName:           "",
		Lore:                 "",
		Enchantments:         "",
		CustomModelData:      "",
		RepairCost:           "",
		HideFlags:            "",
		PotionID:             "",
		CustomPotionColor:    "",
		CustomPotionEffects:  "",
		AttributeModifiers:   "",
		CustomNBT:            "",
		EnchantmentOptions:   enchantmentOptions,
		SelectedEnchantments: selectedEnchantments,
	}
}

func itemEntryToForm(entry items.ItemEntry, options []views.ReferenceOption) views.ItemFormData {
	enchantmentOptions, selectedEnchantments := itemFormEnchantmentsFromText(entry.Enchantments)
	return views.ItemFormData{
		ID:                   entry.ID,
		ItemID:               entry.ItemID,
		SkillID:              entry.SkillID,
		SkillOptions:         options,
		CustomName:           entry.CustomName,
		Lore:                 entry.Lore,
		Enchantments:         entry.Enchantments,
		Unbreakable:          entry.Unbreakable,
		CustomModelData:      entry.CustomModelData,
		RepairCost:           entry.RepairCost,
		HideFlags:            entry.HideFlags,
		PotionID:             entry.PotionID,
		CustomPotionColor:    entry.CustomPotionColor,
		CustomPotionEffects:  entry.CustomPotionEffects,
		AttributeModifiers:   entry.AttributeModifiers,
		CustomNBT:            entry.CustomNBT,
		EnchantmentOptions:   enchantmentOptions,
		SelectedEnchantments: selectedEnchantments,
		FieldErrors:          map[string]string{},
		IsEditing:            true,
	}
}

func parseItemForm(r *http.Request, skills []skills.SkillEntry) (views.ItemFormData, items.SaveInput, map[string]string) {
	form := defaultItemForm(skillOptions(skills))
	form.ID = strings.TrimSpace(r.Form.Get("id"))
	form.ItemID = strings.TrimSpace(r.Form.Get("itemId"))
	form.SkillID = strings.TrimSpace(r.Form.Get("skillId"))
	form.CustomName = r.Form.Get("customName")
	form.Lore = r.Form.Get("lore")
	form.Enchantments, form.EnchantmentOptions, form.SelectedEnchantments = itemFormEnchantmentsFromRequest(r)
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

func itemOptions(entries []items.ItemEntry) []views.ReferenceOption {
	options := make([]views.ReferenceOption, 0, len(entries))
	for _, entry := range entries {
		options = append(options, views.ReferenceOption{ID: entry.ID, Label: entry.ItemID})
	}
	return options
}

func skillOptions(entries []skills.SkillEntry) []views.ReferenceOption {
	options := make([]views.ReferenceOption, 0, len(entries))
	for _, entry := range entries {
		label := entry.Name
		if strings.TrimSpace(label) == "" {
			label = entry.ID
		}
		options = append(options, views.ReferenceOption{ID: entry.ID, Label: label})
	}
	return options
}
