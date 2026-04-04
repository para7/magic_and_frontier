package export_convert

import (
	"fmt"
	"strings"

	passiveModel "maf_command_editor/app/domain/model/passive"
)

const PassiveCastTime = 200
const PassiveMPCost = 10

func PassiveDerivedCastID(baseCastID, slot int) int {
	return baseCastID*10 + slot
}

func PassiveToBook(entry passiveModel.Passive, slot int) string {
	itemName := fmt.Sprintf(`{text:%s}`, JsonString(passiveBookTitle(entry, slot)))
	lore := fmt.Sprintf(
		`[{text:%s},{text:%s}]`,
		JsonString("右クリックでパッシブを設定"),
		JsonString(fmt.Sprintf("passive=%s slot=%d cast=%d cost=%d", entry.ID, slot, PassiveCastTime, PassiveMPCost)),
	)
	consumable := `{consume_seconds:99999,animation:"bow",has_consume_particles:false}`
	customData := passiveSpellCustomData(entry, slot)

	return fmt.Sprintf(
		`minecraft:book[minecraft:item_name=%s,minecraft:lore=%s,minecraft:consumable=%s,minecraft:custom_data=%s]`,
		itemName, lore, consumable, customData,
	)
}

func passiveBookTitle(entry passiveModel.Passive, slot int) string {
	name := strings.TrimSpace(entry.Name)
	if name == "" {
		name = entry.ID
	}
	return fmt.Sprintf("%s[S%d]%d", name, slot, PassiveCastTime)
}

func passiveBookDescription(entry passiveModel.Passive) string {
	if text := strings.TrimSpace(entry.Description); text != "" {
		return text
	}
	return fmt.Sprintf("condition=%s", strings.TrimSpace(entry.Condition))
}

func passiveSpellCustomData(entry passiveModel.Passive, slot int) string {
	return fmt.Sprintf(
		`{maf:{passive:{id:%s,slot:%d,condition:%s},spell:{castid:%d,cost:%d,cast:%d,cooltime:%d,title:%s,description:%s}}}`,
		JsonString(entry.ID),
		slot,
		JsonString(strings.TrimSpace(entry.Condition)),
		PassiveDerivedCastID(entry.CastID, slot),
		PassiveMPCost,
		PassiveCastTime,
		0,
		JsonString(passiveBookTitle(entry, slot)),
		JsonString(passiveBookDescription(entry)),
	)
}

func passiveLootComponents(entry passiveModel.Passive, slot int) map[string]any {
	return map[string]any{
		"minecraft:item_name": map[string]any{"text": passiveBookTitle(entry, slot)},
		"minecraft:lore": []any{
			map[string]any{"text": "右クリックでパッシブを設定"},
			map[string]any{"text": fmt.Sprintf("passive=%s slot=%d cast=%d cost=%d", entry.ID, slot, PassiveCastTime, PassiveMPCost)},
		},
		"minecraft:consumable": map[string]any{
			"consume_seconds":       99999.0,
			"animation":             "bow",
			"has_consume_particles": false,
		},
	}
}
