package loottables

import "tools2/app/internal/domain/treasures"

type LootTableEntry struct {
	ID        string              `json:"id"`
	TablePath string              `json:"tablePath"`
	LootPools []treasures.DropRef `json:"lootPools"`
	UpdatedAt string              `json:"updatedAt"`
}

type SaveInput struct {
	ID        string              `json:"id"`
	TablePath string              `json:"tablePath" validate:"trimmed_required,trimmed_min=1,trimmed_max=200"`
	LootPools []treasures.DropRef `json:"lootPools" validate:"min=1,dive"`
}
