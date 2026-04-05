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
		itemName: fmt.Sprintf("%s%d", entry.Title, entry.CastTime),
		lore: []string{
			"右クリックで詠唱を開始",
			fmt.Sprintf("effect=%d cast=%d cost=%d", entry.CastID, entry.CastTime, entry.MPCost),
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
