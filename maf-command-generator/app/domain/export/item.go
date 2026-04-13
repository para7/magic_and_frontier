package export

import (
	"path/filepath"

	ec "maf_command_editor/app/domain/export/convert"
	bowModel "maf_command_editor/app/domain/model/bow"
	grimoireModel "maf_command_editor/app/domain/model/grimoire"
	passiveModel "maf_command_editor/app/domain/model/passive"
)

type ItemGiveFunction struct {
	ID   string
	Body string
}

func BuildItemArtifacts(master DBMaster) ([]ItemGiveFunction, error) {
	if master == nil {
		return []ItemGiveFunction{}, nil
	}

	items := master.ListItems()
	grimoires := master.ListGrimoires()
	grimoiresByID := make(map[string]grimoireModel.Grimoire, len(grimoires))
	for _, entry := range grimoires {
		grimoiresByID[entry.ID] = entry
	}

	passives := master.ListPassives()
	passivesByID := make(map[string]passiveModel.Passive, len(passives))
	for _, entry := range passives {
		passivesByID[entry.ID] = entry
	}
	bows := master.ListBows()
	bowsByID := make(map[string]bowModel.BowPassive, len(bows))
	for _, entry := range bows {
		bowsByID[entry.ID] = entry
	}

	results := make([]ItemGiveFunction, 0, len(items))
	for _, entry := range items {
		body, err := ec.ItemToGiveCommand(entry, grimoiresByID, passivesByID, bowsByID)
		if err != nil {
			return nil, err
		}
		results = append(results, ItemGiveFunction{
			ID:   entry.ID,
			Body: body,
		})
	}
	return results, nil
}

func WriteItemArtifacts(dir string, artifacts []ItemGiveFunction) error {
	for _, entry := range artifacts {
		path := filepath.Join(dir, entry.ID+".mcfunction")
		if err := writeFunctionFile(path, entry.Body); err != nil {
			return err
		}
	}
	return nil
}
