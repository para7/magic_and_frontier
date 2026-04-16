package export

import (
	"path/filepath"

	ec "maf_command_editor/app/domain/export/convert"
)

type ItemGiveFunction struct {
	ID   string
	Body string
}

func BuildItemArtifacts(master DBMaster) ([]ItemGiveFunction, error) {
	if master == nil {
		return []ItemGiveFunction{}, nil
	}
	lookups := buildMasterEntityLookups(master)
	results := make([]ItemGiveFunction, 0, len(lookups.items))
	for _, entry := range lookups.items {
		body, err := ec.ItemToGiveCommand(entry, lookups.grimoiresByID, lookups.passivesByID, lookups.bowsByID)
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
