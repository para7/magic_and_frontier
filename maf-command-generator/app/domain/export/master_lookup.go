package export

import (
	activeModel "maf_command_editor/app/domain/model/active"
	bowModel "maf_command_editor/app/domain/model/bow"
	itemModel "maf_command_editor/app/domain/model/item"
	passiveModel "maf_command_editor/app/domain/model/passive"
)

type masterEntityLookups struct {
	items        []itemModel.Item
	itemsByID    map[string]itemModel.Item
	activesByID  map[string]activeModel.Active
	passivesByID map[string]passiveModel.Passive
	bowsByID     map[string]bowModel.BowPassive
}

func buildMasterEntityLookups(master DBMaster) masterEntityLookups {
	items := master.ListItems()
	actives := master.ListActives()
	passives := master.ListPassives()
	bows := master.ListBows()

	return masterEntityLookups{
		items:        items,
		itemsByID:    indexByID(items, func(entry itemModel.Item) string { return entry.ID }),
		activesByID:  indexByID(actives, func(entry activeModel.Active) string { return entry.ID }),
		passivesByID: indexByID(passives, func(entry passiveModel.Passive) string { return entry.ID }),
		bowsByID:     indexByID(bows, func(entry bowModel.BowPassive) string { return entry.ID }),
	}
}

func indexByID[T any](entries []T, idOf func(T) string) map[string]T {
	result := make(map[string]T, len(entries))
	for _, entry := range entries {
		result[idOf(entry)] = entry
	}
	return result
}
