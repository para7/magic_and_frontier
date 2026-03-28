package export

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const namespace = "maf"

type GrimoireEffectFunction struct {
	ID   string
	Body string
}

type GrimoireEffectArtifacts struct {
	Effects    []GrimoireEffectFunction
	SelectExec string
}

func BuildGrimoireArtifacts(master DBMaster, spellEffectDir string) GrimoireEffectArtifacts {
	entries := []GrimoireEffectFunction{}
	dispatchLines := []string{}

	if master != nil {
		grimoires := master.ListGrimoires()
		entries = make([]GrimoireEffectFunction, 0, len(grimoires))
		dispatchLines = make([]string, 0, len(grimoires))

		for _, entry := range grimoires {
			entries = append(entries, GrimoireEffectFunction{
				ID:   entry.ID,
				Body: entry.Script,
			})
			dispatchLines = append(dispatchLines, "execute if entity @s[scores={mafEffectID="+strconv.Itoa(entry.CastID)+"}] run function "+functionResourceID(namespace, spellEffectDir, entry.ID))
		}
	}

	return GrimoireEffectArtifacts{
		Effects:    entries,
		SelectExec: strings.Join(dispatchLines, "\n") + "\n",
	}
}
func WriteGrimoireArtifacts(spellEffectDir string, artifacts GrimoireEffectArtifacts) error {
	os.MkdirAll(spellEffectDir, 0o755)

	for _, entry := range artifacts.Effects {
		path := getEffectFunctionName(spellEffectDir, entry.ID)
		err := os.WriteFile(path, []byte(entry.Body), 0o755)
		fmt.Println("write... " + path)
		if err != nil {
			return err
		}
	}
	return nil
}

func getEffectFunctionName(spellEffectDir string, grimoireId string) string {
	return filepath.Join(spellEffectDir, grimoireId+".mcfunction")
}

func functionResourceID(namespace string, relativeDir, baseName string) string {
	resourcePath := strings.TrimPrefix(filepath.ToSlash(relativeDir), "data/maf/function/")
	resourcePath = strings.Trim(resourcePath, "/")
	if resourcePath == "" {
		return namespace + ":" + baseName
	}
	return namespace + ":" + resourcePath + "/" + baseName
}
