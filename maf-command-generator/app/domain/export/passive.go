package export

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	ec "maf_command_editor/app/domain/export/convert"
)

type PassiveEffectFunction struct {
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

func BuildPassiveArtifacts(master DBMaster) ([]PassiveEffectFunction, []PassiveGrimoireFunction, error) {
	if master == nil {
		return []PassiveEffectFunction{}, []PassiveGrimoireFunction{}, nil
	}

	passives := master.ListPassives()
	effects := make([]PassiveEffectFunction, 0, len(passives))
	grimoires := make([]PassiveGrimoireFunction, 0, len(passives))

	for _, entry := range passives {
		effectBody := strings.Join(entry.Script, "\n")
		effects = append(effects, PassiveEffectFunction{
			ID:   entry.ID,
			Body: effectBody,
		})

		if entry.GenerateGrimoire != nil && *entry.GenerateGrimoire {
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
				applyBody := passiveApplyBody(slot, entry.ID, entry.Condition, displayName)
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
	}

	return effects, grimoires, nil
}

func WritePassiveArtifacts(effectDir, giveDir, applyDir string, effects []PassiveEffectFunction, grimoires []PassiveGrimoireFunction) error {
	for _, entry := range effects {
		path := filepath.Join(effectDir, entry.ID+".mcfunction")
		if err := writeFunctionFile(path, entry.Body); err != nil {
			return err
		}
	}
	if err := removePassiveSlotFunctionFiles(giveDir); err != nil {
		return err
	}
	if err := removePassiveSlotFunctionFiles(applyDir); err != nil {
		return err
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

func removePassiveSlotFunctionFiles(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*_slot*.mcfunction"))
	if err != nil {
		return err
	}
	for _, path := range files {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func passiveApplyBody(slot int, passiveID string, condition string, displayName string) string {
	setMessage := fmt.Sprintf("[slot%d]に[%s]を設定しました", slot, displayName)
	return strings.Join([]string{
		"function #oh_my_dat:please",
		fmt.Sprintf("data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot%d.id set value %s", slot, ec.JsonString(passiveID)),
		fmt.Sprintf("data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot%d.condition set value %s", slot, ec.JsonString(condition)),
		fmt.Sprintf(`tellraw @s [{"text":%s}]`, ec.JsonString(setMessage)),
	}, "\n")
}
