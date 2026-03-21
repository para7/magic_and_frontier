package treasures

type DropRef struct {
	Kind     string   `json:"kind" validate:"trimmed_required,trimmed_oneof=minecraft_item item grimoire"`
	RefID    string   `json:"refId" validate:"trimmed_required,trimmed_min=1,trimmed_max=200"`
	Weight   float64  `json:"weight" validate:"gte=1,lte=100000"`
	CountMin *float64 `json:"countMin,omitempty" validate:"omitempty,gte=1,lte=64"`
	CountMax *float64 `json:"countMax,omitempty" validate:"omitempty,gte=1,lte=64"`
}

type TreasureEntry struct {
	ID        string    `json:"id"`
	TablePath string    `json:"tablePath"`
	LootPools []DropRef `json:"lootPools"`
	UpdatedAt string    `json:"updatedAt"`
}

type SaveInput struct {
	ID        string    `json:"id"`
	TablePath string    `json:"tablePath" validate:"trimmed_required,trimmed_min=1,trimmed_max=200"`
	LootPools []DropRef `json:"lootPools" validate:"min=1,dive"`
}
