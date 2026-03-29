package export

import (
	"fmt"

	grimoireModel "maf_command_editor/app/domain/model/grimoire"
)

// 独自形式のグリモアのデータをマインクラフトの本に変換する
func grimoireToBook(grimoire grimoireModel.Grimoire) string {
	//example:
	// minecraft:book[minecraft:item_name={text:"動作確認: スウィープ10"},minecraft:lore=[{text:"右クリックで詠唱を開始"},{text:"effect=13 cast=10 cost=0"}],minecraft:consumable={consume_seconds:99999,animation:"bow",has_consume_particles:false},minecraft:custom_data={maf:{grimoire_id:"dev_book2",spell:{castid:13,cost:10,cast:10,title:"スウィープ",description:"近くのアイテムをかき集める"}}}]

	itemName := fmt.Sprintf(`{text:"%s%d"}`, grimoire.Title, grimoire.CastTime)
	lore := fmt.Sprintf(`[{text:"右クリックで詠唱を開始"},{text:"effect=%d cast=%d cost=%d"}]`,
		grimoire.CastID, grimoire.CastTime, grimoire.MPCost)
	consumable := `{consume_seconds:99999,animation:"bow",has_consume_particles:false}`
	customData := fmt.Sprintf(`{maf:{grimoire_id:"%s",spell:{castid:%d,cost:%d,cast:%d,title:"%s",description:"%s"}}}`,
		grimoire.ID, grimoire.CastID, grimoire.MPCost, grimoire.CastTime, grimoire.Title, grimoire.Description)

	return fmt.Sprintf(
		`minecraft:book[minecraft:item_name=%s,minecraft:lore=%s,minecraft:consumable=%s,minecraft:custom_data=%s]`,
		itemName, lore, consumable, customData,
	)
}
