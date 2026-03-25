package loottables

import "tools2/app/internal/domain/entity/treasures"

type LootTableEntry struct {
	ID        string              `json:"id"`
	Memo      string              `json:"memo,omitempty"`
	LootPools []treasures.DropRef `json:"lootPools"`
	UpdatedAt string              `json:"updatedAt"`
}

type SaveInput struct {
	ID        string              `json:"id"`
	Memo      string              `json:"memo" validate:"trimmed_max=400"`
	LootPools []treasures.DropRef `json:"lootPools" validate:"min=1,dive"`
}
