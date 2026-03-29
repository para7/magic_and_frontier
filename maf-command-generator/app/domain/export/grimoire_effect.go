package export

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
			Body:         entry.Script,
			SelectScript: selectScript,
			Book:         grimoireToBook(entry),
		})
	}

	return entries
}

func WriteGrimoireArtifacts(spellEffectDir string, selectExecFile string, effects []GrimoireEffectFunction) error {
	os.MkdirAll(spellEffectDir, 0o755)

	selectLines := make([]string, 0, len(effects))
	for _, entry := range effects {
		path := getEffectFunctionName(spellEffectDir, entry.ID)
		err := os.WriteFile(path, []byte(entry.Body), 0o755)
		fmt.Println("write... " + path)
		if err != nil {
			return err
		}
		selectLines = append(selectLines, entry.SelectScript)
	}

	selectExec := strings.Join(selectLines, "\n")
	return os.WriteFile(selectExecFile, []byte(selectExec), 0o755)
}

// デバッグ用の give コマンドを生成
func WriteGrimoireDebugArtifacts(debugDir string, effects []GrimoireEffectFunction) error {
	os.MkdirAll(debugDir, 0o755)
	for _, entry := range effects {
		path := filepath.Join(debugDir, entry.ID+".mcfunction")
		script := fmt.Sprintf("give @p %s 1", entry.Book)
		if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
			return err
		}
	}
	return nil
}

func getEffectFunctionName(spellEffectDir string, grimoireId string) string {
	return filepath.Join(spellEffectDir, grimoireId+".mcfunction")
}

// minecraft で認識される function の名前を取得する
func functionRefName(logicalDir, baseName string) string {
	dir := strings.Trim(filepath.ToSlash(logicalDir), "/")
	if dir == "" {
		return "maf:" + baseName
	}
	return "maf:" + dir + "/" + baseName
}
