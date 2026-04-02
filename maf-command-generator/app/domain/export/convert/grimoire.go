package export_convert

import (
	"fmt"

	grimoireModel "maf_command_editor/app/domain/model/grimoire"
)

// GrimoireToBook は独自形式のグリモアのデータをマインクラフトの本に変換する
func GrimoireToBook(grimoire grimoireModel.Grimoire) string {
	itemName := fmt.Sprintf(`{text:"%s%d"}`, grimoire.Title, grimoire.CastTime)
	lore := fmt.Sprintf(`[{text:"右クリックで詠唱を開始"},{text:"effect=%d cast=%d cost=%d"}]`,
		grimoire.CastID, grimoire.CastTime, grimoire.MPCost)
	consumable := `{consume_seconds:99999,animation:"bow",has_consume_particles:false}`
	customData := fmt.Sprintf(`{maf:{grimoire_id:"%s",spell:{castid:%d,cost:%d,cast:%d,cooltime:%d,title:"%s",description:"%s"}}}`,
		grimoire.ID, grimoire.CastID, grimoire.MPCost, grimoire.CastTime, grimoire.CoolTime, grimoire.Title, grimoire.Description)

	return fmt.Sprintf(
		`minecraft:book[minecraft:item_name=%s,minecraft:lore=%s,minecraft:consumable=%s,minecraft:custom_data=%s]`,
		itemName, lore, consumable, customData,
	)
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
