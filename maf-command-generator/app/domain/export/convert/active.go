package export_convert

import (
	"fmt"

	activeModel "maf_command_editor/app/domain/model/active"
)

// ActiveToBook は独自形式のアクティブのデータをマインクラフトの本に変換する
func ActiveToBook(entry activeModel.Active) string {
	return activeSpellBookModel(entry).ToGiveItem()
}

func activeSpellBookModel(entry activeModel.Active) spellBookModel {
	return spellBookModel{
		itemName: entry.Title,
		lore: []string{
			entry.Description,
			fmt.Sprintf("消費MP:%d 詠唱時間:%d", entry.MPCost, entry.CastTime),
		},
		customData: spellCustomData(entry),
	}
}

func spellCustomData(entry activeModel.Active) string {
	return fmt.Sprintf("{maf:{active_id:%s}}", SNBTString(entry.ID))
}

func ActiveCastingDataSNBT(entry activeModel.Active) string {
	return spellCastingDataSNBT("active", entry.ID, nil, entry.MPCost, entry.CastTime, entry.CoolTime, entry.Title, entry.Description)
}

func spellCastingDataSNBT(kind string, id string, slot *int, mpCost, castTime, coolTime int, title, description string) string {
	if slot != nil {
		return fmt.Sprintf(
			"{kind:%s,id:%s,slot:%d,cost:%d,cast:%d,cooltime:%d,title:%s,description:%s}",
			SNBTString(kind),
			SNBTString(id),
			*slot,
			mpCost,
			castTime,
			coolTime,
			SNBTString(title),
			SNBTString(description),
		)
	}
	return fmt.Sprintf(
		"{kind:%s,id:%s,cost:%d,cast:%d,cooltime:%d,title:%s,description:%s}",
		SNBTString(kind),
		SNBTString(id),
		mpCost,
		castTime,
		coolTime,
		SNBTString(title),
		SNBTString(description),
	)
}
