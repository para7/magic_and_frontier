package export_convert

import (
	"fmt"

	grimoireModel "maf_command_editor/app/domain/model/grimoire"
)

// GrimoireToBook は独自形式のグリモアのデータをマインクラフトの本に変換する
func GrimoireToBook(entry grimoireModel.Grimoire) string {
	return grimoireSpellBookModel(entry).ToGiveItem()
}

func grimoireSpellBookModel(entry grimoireModel.Grimoire) spellBookModel {
	return spellBookModel{
		itemName: entry.Title,
		lore: []string{
			entry.Description,
			fmt.Sprintf("消費MP:%d 詠唱時間:%d", entry.MPCost, entry.CastTime),
		},
		customData: spellCustomData(entry),
	}
}

func spellCustomData(entry grimoireModel.Grimoire) string {
	return fmt.Sprintf(
		"{maf:{grimoire_id:%s,spell:{castid:%d,cost:%d,cast:%d,cooltime:%d,title:%s,description:%s}}}",
		JsonString(entry.ID),
		entry.CastID,
		entry.MPCost,
		entry.CastTime,
		entry.CoolTime,
		JsonString(entry.Title),
		JsonString(entry.Description),
	)
}
