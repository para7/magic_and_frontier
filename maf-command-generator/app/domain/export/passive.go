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

type PassiveGrimoireFunction struct {
	PassiveID  string
	Slot       int
	CastID     int
	FunctionID string
	ApplyRef   string
	GiveBody   string
	ApplyBody  string
	Book       string
}

const passiveApplyByUUIDFunctionID = "set_slot_by_uuid"
const passiveApplyByUUIDLogicalDir = "passive/apply"

func BuildPassiveArtifacts(master DBMaster, applyDir string) ([]PassiveEffectFunction, []PassiveGrimoireFunction, error) {
	if master == nil {
		return []PassiveEffectFunction{}, []PassiveGrimoireFunction{}, nil
	}

	passives := master.ListPassives()
	effects := make([]PassiveEffectFunction, 0, len(passives))
	grimoires := make([]PassiveGrimoireFunction, 0, len(passives))

	for _, entry := range passives {
		effects = append(effects, PassiveEffectFunction{
			ID:   entry.ID,
			Body: strings.Join(entry.Script, "\n"),
		})

		slots := make([]int, len(entry.Slots))
		copy(slots, entry.Slots)
		sort.Ints(slots)

		for _, slot := range slots {
			castID := ec.PassiveDerivedCastID(entry.CastID, slot)
			functionID := passiveGrimoireFunctionID(entry.ID, slot)
			applyRef := functionRefName(applyDir, functionID)
			book := ec.PassiveToBook(entry, slot)
			displayName := passiveDisplayName(entry.Name, entry.ID)
			applyBody := passiveApplyBody(slot, entry.ID, displayName)
			grimoires = append(grimoires, PassiveGrimoireFunction{
				PassiveID:  entry.ID,
				Slot:       slot,
				CastID:     castID,
				FunctionID: functionID,
				ApplyRef:   applyRef,
				GiveBody:   fmt.Sprintf("give @p %s 1", book),
				ApplyBody:  applyBody,
				Book:       book,
			})
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

func BuildSelectExecLines(grimoires []GrimoireEffectFunction, passiveGrimoires []PassiveGrimoireFunction) ([]string, error) {
	seen := map[int]string{}
	lines := make([]string, 0, len(grimoires)+len(passiveGrimoires)+4)
	readEffectRef := functionRefName(grimoireDispatchHelperLogicalDir, grimoireReadEffectFunctionID)
	runEffectRef := functionRefName(grimoireDispatchHelperLogicalDir, grimoireRunEffectFunctionID)

	for _, entry := range grimoires {
		source := "grimoire(" + entry.ID + ")"
		if prev, ok := seen[entry.CastID]; ok {
			return nil, fmt.Errorf("duplicate castid %d between %s and %s", entry.CastID, prev, source)
		}
		seen[entry.CastID] = source
	}
	if len(grimoires) > 0 {
		lines = append(lines,
			"data remove storage "+grimoireDispatchStorage,
			"execute store result storage "+grimoireDispatchStorage+".castid int 1 run scoreboard players get @s mafEffectID",
			"function "+readEffectRef+" with storage "+grimoireDispatchStorage,
			"execute if data storage "+grimoireDispatchStorage+".ref run function "+runEffectRef+" with storage "+grimoireDispatchStorage,
		)
	}

	for _, entry := range passiveGrimoires {
		source := fmt.Sprintf("passive(%s,slot=%d)", entry.PassiveID, entry.Slot)
		if prev, ok := seen[entry.CastID]; ok {
			return nil, fmt.Errorf("duplicate castid %d between %s and %s", entry.CastID, prev, source)
		}
		seen[entry.CastID] = source
		lines = append(lines, castSelectExecLine(entry.CastID, entry.ApplyRef))
	}

	return lines, nil
}

func castSelectExecLine(castID int, functionRef string) string {
	return fmt.Sprintf("execute if entity @s[scores={mafEffectID=%d}] run function %s", castID, functionRef)
}

func passiveGrimoireFunctionID(passiveID string, slot int) string {
	return fmt.Sprintf("%s_slot%d", passiveID, slot)
}

func passiveDisplayName(name, fallbackID string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return fallbackID
	}
	return trimmed
}

func passiveApplyBody(slot int, passiveID string, displayName string) string {
	setMessage := fmt.Sprintf("[slot%d]に[%s]を設定しました", slot, displayName)
	applyByUUIDRef := functionRefName(passiveApplyByUUIDLogicalDir, passiveApplyByUUIDFunctionID)
	return strings.Join([]string{
		"data remove storage p7:maf passive.tmp",
		"data modify storage p7:maf passive.tmp.uuid set from entity @s UUID",
		"execute store result storage p7:maf passive.tmp.u0 int 1 run data get storage p7:maf passive.tmp.uuid[0]",
		"execute store result storage p7:maf passive.tmp.u1 int 1 run data get storage p7:maf passive.tmp.uuid[1]",
		"execute store result storage p7:maf passive.tmp.u2 int 1 run data get storage p7:maf passive.tmp.uuid[2]",
		"execute store result storage p7:maf passive.tmp.u3 int 1 run data get storage p7:maf passive.tmp.uuid[3]",
		fmt.Sprintf("data modify storage p7:maf passive.tmp.slot set value %d", slot),
		fmt.Sprintf("data modify storage p7:maf passive.tmp.id set value %s", ec.JsonString(passiveID)),
		fmt.Sprintf("function %s with storage p7:maf passive.tmp", applyByUUIDRef),
		fmt.Sprintf(`tellraw @s [{"text":%s}]`, ec.JsonString(setMessage)),
	}, "\n")
}
