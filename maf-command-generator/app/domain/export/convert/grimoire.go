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
	return fmt.Sprintf("{maf:{grimoire_id:%s,%s}}", JsonString(entry.ID), grimoireSpellFragment(entry))
}

func grimoireSpellFragment(entry grimoireModel.Grimoire) string {
	return spellFragment("grimoire", entry.ID, nil, entry.MPCost, entry.CastTime, entry.CoolTime, entry.Title, entry.Description)
}

func spellFragment(kind string, id string, slot *int, mpCost, castTime, coolTime int, title, description string) string {
	if slot != nil {
		return fmt.Sprintf(
			"spell:{kind:%s,id:%s,slot:%d,cost:%d,cast:%d,cooltime:%d,title:%s,description:%s}",
			JsonString(kind),
			JsonString(id),
			*slot,
			mpCost,
			castTime,
			coolTime,
			JsonString(title),
			JsonString(description),
		)
	}
	return fmt.Sprintf(
		"spell:{kind:%s,id:%s,cost:%d,cast:%d,cooltime:%d,title:%s,description:%s}",
		JsonString(kind),
		JsonString(id),
		mpCost,
		castTime,
		coolTime,
		JsonString(title),
		JsonString(description),
	)
}
