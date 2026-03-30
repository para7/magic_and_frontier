package export

import (
	"fmt"
	"path/filepath"
	"strings"

	ec "maf_command_editor/app/domain/export/convert"
)

type GrimoireEffectFunction struct {
	ID           string
	Body         string
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
		selectScript := fmt.Sprintf("execute if entity @s[scores={mafEffectID=%d}] run function %s", entry.CastID, functionRefName(effectDir, entry.ID))
		entries = append(entries, GrimoireEffectFunction{
			ID:           entry.ID,
			Body:         strings.Join(entry.Script, "\n"),
			SelectScript: selectScript,
			Book:         ec.GrimoireToBook(entry),
		})
	}

	return entries
}

func WriteGrimoireArtifacts(spellEffectDir string, selectExecFile string, effects []GrimoireEffectFunction) error {
	selectLines := make([]string, 0, len(effects))
	for _, entry := range effects {
		path := filepath.Join(spellEffectDir, entry.ID+".mcfunction")
		if err := writeFunctionFile(path, entry.Body); err != nil {
			return err
		}
		selectLines = append(selectLines, entry.SelectScript)
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
