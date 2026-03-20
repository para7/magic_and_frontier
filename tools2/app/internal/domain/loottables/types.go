package loottables

import "tools2/app/internal/domain/treasures"

type LootTableEntry struct {
	ID        string              `json:"id"`
	LootPools []treasures.DropRef `json:"lootPools"`
	UpdatedAt string              `json:"updatedAt"`
}

type SaveInput struct {
	ID        string              `json:"id"`
	LootPools []treasures.DropRef `json:"lootPools" validate:"min=1,dive"`
}
