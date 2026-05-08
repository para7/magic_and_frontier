package export

import (
	"fmt"
	"path/filepath"
	"strings"

	ec "maf_command_editor/app/domain/export/convert"
	activeModel "maf_command_editor/app/domain/model/active"
)

type ActiveEffectFunction struct {
	ID        string
	Body      string
	Book      string
	SpellBody string
}

func BuildActiveArtifacts(master DBMaster) []ActiveEffectFunction {
	if master == nil {
		return []ActiveEffectFunction{}
	}

	actives := master.ListActives()
	entries := make([]ActiveEffectFunction, 0, len(actives))

	for _, entry := range actives {
		entries = append(entries, ActiveEffectFunction{
			ID:        entry.ID,
			Body:      strings.Join(entry.Script, "\n"),
			Book:      ec.ActiveToBook(entry),
			SpellBody: activeSpellLoaderBody(entry),
		})
	}

	return entries
}

func WriteActiveSpellArtifacts(spellDir string, effects []ActiveEffectFunction) error {
	for _, entry := range effects {
		path := filepath.Join(spellDir, entry.ID+".mcfunction")
		if err := writeFunctionFile(path, entry.SpellBody); err != nil {
			return err
		}
	}
	return nil
}

func WriteActiveArtifacts(spellEffectDir string, effects []ActiveEffectFunction) error {
	for _, entry := range effects {
		path := filepath.Join(spellEffectDir, entry.ID+".mcfunction")
		if err := writeFunctionFile(path, entry.Body); err != nil {
			return err
		}
	}
	return nil
}

func activeSpellLoaderBody(entry activeModel.Active) string {
	return fmt.Sprintf(
		"data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.magic.casting set value %s",
		ec.ActiveCastingDataSNBT(entry),
	)
}

// デバッグ用の give コマンドを生成
func WriteActiveDebugArtifacts(debugDir string, effects []ActiveEffectFunction) error {
	for _, entry := range effects {
		path := filepath.Join(debugDir, entry.ID+".mcfunction")
		script := fmt.Sprintf("give @p %s 1", entry.Book)
		if err := writeFunctionFile(path, script); err != nil {
			return err
		}
	}
	return nil
}
