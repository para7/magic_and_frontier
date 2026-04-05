package export

import (
	"fmt"
	"path/filepath"
	"strings"

	ec "maf_command_editor/app/domain/export/convert"
)

const (
	grimoireEffectRefMapStorage      = "p7:maf grimoire.effect_ref_map"
	grimoireDispatchStorage          = "p7:maf grimoire.dispatch"
	grimoireDispatchHelperLogicalDir = "magic/cast/dispatch"
	grimoireSetupMapFunctionID       = "setup_effect_ref_map"
	grimoireReadEffectFunctionID     = "read_effect_ref"
	grimoireRunEffectFunctionID      = "run_effect_ref"
)

type GrimoireEffectFunction struct {
	ID           string
	CastID       int
	Body         string
	FunctionRef  string
	SelectScript string
	Book         string
}

func BuildGrimoireArtifacts(master DBMaster, effectDir string) []GrimoireEffectFunction {
	if master == nil {
		return []GrimoireEffectFunction{}
	}

	grimoires := master.ListGrimoires()
	entries := make([]GrimoireEffectFunction, 0, len(grimoires))

	for _, entry := range grimoires {
		functionRef := functionRefName(effectDir, entry.ID)
		selectScript := fmt.Sprintf("execute if entity @s[scores={mafEffectID=%d}] run function %s", entry.CastID, functionRef)
		entries = append(entries, GrimoireEffectFunction{
			ID:           entry.ID,
			CastID:       entry.CastID,
			Body:         strings.Join(entry.Script, "\n"),
			FunctionRef:  functionRef,
			SelectScript: selectScript,
			Book:         ec.GrimoireToBook(entry),
		})
	}

	return entries
}

func WriteGrimoireArtifacts(spellEffectDir string, selectExecFile string, effects []GrimoireEffectFunction, selectLines []string) error {
	if selectLines == nil {
		selectLines = make([]string, 0, len(effects))
		for _, entry := range effects {
			selectLines = append(selectLines, entry.SelectScript)
		}
	}
	for _, entry := range effects {
		path := filepath.Join(spellEffectDir, entry.ID+".mcfunction")
		if err := writeFunctionFile(path, entry.Body); err != nil {
			return err
		}
	}

	grimoireDir := filepath.Dir(selectExecFile)
	if err := writeFunctionFile(filepath.Join(grimoireDir, grimoireSetupMapFunctionID+".mcfunction"), buildGrimoireSetupMapBody(effects)); err != nil {
		return err
	}

	selectExec := strings.Join(selectLines, "\n")
	return writeFunctionFile(selectExecFile, selectExec)
}

// デバッグ用の give コマンドを生成
func WriteGrimoireDebugArtifacts(debugDir string, effects []GrimoireEffectFunction) error {
	for _, entry := range effects {
		path := filepath.Join(debugDir, entry.ID+".mcfunction")
		script := fmt.Sprintf("give @p %s 1", entry.Book)
		if err := writeFunctionFile(path, script); err != nil {
			return err
		}
	}
	return nil
}

func buildGrimoireSetupMapBody(effects []GrimoireEffectFunction) string {
	lines := make([]string, 0, len(effects)+1)
	lines = append(lines, "data remove storage "+grimoireEffectRefMapStorage)
	for _, entry := range effects {
		lines = append(lines, fmt.Sprintf(
			`data modify storage %s."%d" set value %s`,
			grimoireEffectRefMapStorage,
			entry.CastID,
			ec.JsonString(entry.FunctionRef),
		))
	}
	return strings.Join(lines, "\n")
}
