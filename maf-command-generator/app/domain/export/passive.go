package export

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	ec "maf_command_editor/app/domain/export/convert"
)

type PassiveEffectFunction struct {
	ID   string
	Body string
}

type PassiveBowFunction struct {
	ID   string
	Body string
}

type PassiveGrimoireFunction struct {
	PassiveID  string
	Slot       int
	FunctionID string
	GiveBody   string
	ApplyBody  string
	Book       string
}

func BuildPassiveArtifacts(master DBMaster) ([]PassiveEffectFunction, []PassiveBowFunction, []PassiveGrimoireFunction, error) {
	if master == nil {
		return []PassiveEffectFunction{}, []PassiveBowFunction{}, []PassiveGrimoireFunction{}, nil
	}

	passives := master.ListPassives()
	effects := make([]PassiveEffectFunction, 0, len(passives))
	bows := make([]PassiveBowFunction, 0, len(passives))
	grimoires := make([]PassiveGrimoireFunction, 0, len(passives))

	for _, entry := range passives {
		effectBody := strings.Join(entry.Script, "\n")
		if entry.Condition == "bow" {
			lifeSub := 1200 // デフォルト: 弓ヒット後の矢の生存時間 (tick)。BowConfig.LifeSub で上書き可
			if entry.Bow != nil && entry.Bow.LifeSub != nil {
				lifeSub = *entry.Bow.LifeSub
			}
			effectBody = buildBowEffectBody(entry.ID, lifeSub)
			bows = append(bows, PassiveBowFunction{
				ID:   entry.ID,
				Body: strings.Join(entry.Script, "\n"),
			})
		}
		effects = append(effects, PassiveEffectFunction{
			ID:   entry.ID,
			Body: effectBody,
		})

		slots := make([]int, len(entry.Slots))
		copy(slots, entry.Slots)
		sort.Ints(slots)

		for _, slot := range slots {
			functionID := fmt.Sprintf("%s_slot%d", entry.ID, slot)
			book := ec.PassiveToBook(entry, slot)
			displayName := entry.ID // スペースのみの場合は ID にフォールバック
			if trimmed := strings.TrimSpace(entry.Name); trimmed != "" {
				displayName = trimmed
			}
			applyBody := passiveApplyBody(slot, entry.ID, displayName)
			grimoires = append(grimoires, PassiveGrimoireFunction{
				PassiveID:  entry.ID,
				Slot:       slot,
				FunctionID: functionID,
				GiveBody:   fmt.Sprintf("give @p %s 1", book),
				ApplyBody:  applyBody,
				Book:       book,
			})
		}
	}

	return effects, bows, grimoires, nil
}

func WritePassiveArtifacts(effectDir, bowDir, giveDir, applyDir string, effects []PassiveEffectFunction, bows []PassiveBowFunction, grimoires []PassiveGrimoireFunction) error {
	for _, entry := range effects {
		path := filepath.Join(effectDir, entry.ID+".mcfunction")
		if err := writeFunctionFile(path, entry.Body); err != nil {
			return err
		}
	}
	for _, entry := range bows {
		path := filepath.Join(bowDir, entry.ID+".mcfunction")
		if err := writeFunctionFile(path, entry.Body); err != nil {
			return err
		}
	}
	for _, entry := range grimoires {
		givePath := filepath.Join(giveDir, entry.FunctionID+".mcfunction")
		if err := writeFunctionFile(givePath, entry.GiveBody); err != nil {
			return err
		}
		applyPath := filepath.Join(applyDir, entry.FunctionID+".mcfunction")
		if err := writeFunctionFile(applyPath, entry.ApplyBody); err != nil {
			return err
		}
	}
	return nil
}

func passiveApplyBody(slot int, passiveID string, displayName string) string {
	setMessage := fmt.Sprintf("[slot%d]に[%s]を設定しました", slot, displayName)
	return strings.Join([]string{
		"function #oh_my_dat:please",
		fmt.Sprintf("data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot%d.id set value %s", slot, ec.JsonString(passiveID)),
		fmt.Sprintf(`tellraw @s [{"text":%s}]`, ec.JsonString(setMessage)),
	}, "\n")
}

func buildBowEffectBody(id string, lifeSub int) string {
	lifeValue := 1200 - lifeSub
	return strings.Join([]string{
		"execute unless score @s mafBowUsed matches 1.. run return 0",
		"execute store result storage maf:tmp bow_player_id int 1 run scoreboard players get @s mafPlayerID",
		fmt.Sprintf(`execute as @e[type=arrow,distance=..2,nbt=!{inGround:1b},sort=nearest,limit=1] run function maf:passive/tag_passive_arrow {passive_id:%s,life:%d}`, ec.JsonString(id), lifeValue),
	}, "\n")
}
