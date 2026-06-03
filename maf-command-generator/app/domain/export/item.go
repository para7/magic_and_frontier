package export

import (
	"path/filepath"
	"strconv"

	ec "maf_command_editor/app/domain/export/convert"
)

type ItemGiveFunction struct {
	ID   string
	Body string
}

type ItemUpdateFunction struct {
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
		body, err := ec.ItemToGiveCommand(entry, lookups.activesByID, lookups.passivesByID, lookups.bowsByID)
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

func BuildItemUpdateArtifacts(master DBMaster) ([]ItemUpdateFunction, error) {
	if master == nil {
		return []ItemUpdateFunction{}, nil
	}
	lookups := buildMasterEntityLookups(master)
	results := make([]ItemUpdateFunction, 0, len(lookups.items))
	for _, entry := range lookups.items {
		body, err := ec.ItemToUpdateCommand(entry, lookups.activesByID, lookups.passivesByID, lookups.bowsByID)
		if err != nil {
			return nil, err
		}
		results = append(results, ItemUpdateFunction{
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

func WriteItemUpdateArtifacts(dir string, ver int64, artifacts []ItemUpdateFunction) error {
	if err := writeFunctionFile(filepath.Join(dir, "timestamp.mcfunction"), "scoreboard players set #maf_item_ver tmp "+strconv.FormatInt(ver, 10)); err != nil {
		return err
	}
	for _, entry := range artifacts {
		path := filepath.Join(dir, entry.ID+".mcfunction")
		if err := writeFunctionFile(path, entry.Body); err != nil {
			return err
		}
	}
	return nil
}
