package export

import (
	"fmt"
	"path/filepath"
	"strings"

	ec "maf_command_editor/app/domain/export/convert"
	grimoireModel "maf_command_editor/app/domain/model/grimoire"
)

type GrimoireEffectFunction struct {
	ID        string
	Body      string
	Book      string
	SpellBody string
}

func BuildGrimoireArtifacts(master DBMaster) []GrimoireEffectFunction {
	if master == nil {
		return []GrimoireEffectFunction{}
	}

	grimoires := master.ListGrimoires()
	entries := make([]GrimoireEffectFunction, 0, len(grimoires))

	for _, entry := range grimoires {
		entries = append(entries, GrimoireEffectFunction{
			ID:        entry.ID,
			Body:      strings.Join(entry.Script, "\n"),
			Book:      ec.GrimoireToBook(entry),
			SpellBody: grimoireSpellLoaderBody(entry),
		})
	}

	return entries
}

func WriteGrimoireSpellArtifacts(spellDir string, effects []GrimoireEffectFunction) error {
	for _, entry := range effects {
		path := filepath.Join(spellDir, entry.ID+".mcfunction")
		if err := writeFunctionFile(path, entry.SpellBody); err != nil {
			return err
		}
	}
	return nil
}

func WriteGrimoireArtifacts(spellEffectDir string, effects []GrimoireEffectFunction) error {
	for _, entry := range effects {
		path := filepath.Join(spellEffectDir, entry.ID+".mcfunction")
		if err := writeFunctionFile(path, entry.Body); err != nil {
			return err
		}
	}
	return nil
}

func grimoireSpellLoaderBody(entry grimoireModel.Grimoire) string {
	return fmt.Sprintf(
		"data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.magic.casting set value %s",
		ec.GrimoireCastingDataSNBT(entry),
	)
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
