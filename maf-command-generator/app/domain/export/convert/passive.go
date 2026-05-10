package export_convert

import (
	"fmt"
	"strings"

	passiveModel "maf_command_editor/app/domain/model/passive"
)

const PassiveCastTime = 200
const PassiveMPCost = 10

func PassiveToBook(entry passiveModel.Passive, slot int) string {
	return passiveSpellBookModel(entry, slot).ToGiveItem()
}

func passiveItemName(entry passiveModel.Passive) string {
	name := strings.TrimSpace(entry.Name)
	if name == "" {
		name = entry.ID
	}
	return fmt.Sprintf("[パッシブ設定書] %s", name)
}

func passiveSpellTitle(entry passiveModel.Passive, slot int) string {
	name := strings.TrimSpace(entry.Name)
	if name == "" {
		name = entry.ID
	}
	return fmt.Sprintf("%s [スロット%d]", name, slot)
}

func passiveBookDescription(entry passiveModel.Passive) string {
	if text := strings.TrimSpace(entry.Description); text != "" {
		return text
	}
	return fmt.Sprintf("condition=%s", strings.TrimSpace(entry.Condition))
}

func passiveLoreLine(entry passiveModel.Passive) string {
	if lore := strings.TrimSpace(entry.Lore); lore != "" {
		return lore
	}
	return passiveBookDescription(entry)
}

func passiveSpellCustomData(entry passiveModel.Passive, slot int) string {
	return fmt.Sprintf(
		`{maf:{passive:{id:%s,slot:%d,condition:%s}}}`,
		SNBTString(entry.ID),
		slot,
		SNBTString(strings.TrimSpace(entry.Condition)),
	)
}

func PassiveCastingDataSNBT(entry passiveModel.Passive, slot int) string {
	return spellCastingDataSNBT("passive", entry.ID, &slot, PassiveMPCost, PassiveCastTime, 0, passiveSpellTitle(entry, slot), passiveBookDescription(entry))
}

func passiveSpellBookModel(entry passiveModel.Passive, slot int) spellBookModel {
	return spellBookModel{
		itemName: passiveItemName(entry),
		lore: []string{
			passiveLoreLine(entry),
			fmt.Sprintf("パッシブスキル / スロット%d", slot),
		},
		customData: passiveSpellCustomData(entry, slot),
	}
}
