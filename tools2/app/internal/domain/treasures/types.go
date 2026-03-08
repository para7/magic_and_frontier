package treasures

type DropRef struct {
	Kind     string   `json:"kind" validate:"trimmed_required,trimmed_oneof=item grimoire"`
	RefID    string   `json:"refId" validate:"required,uuid_any"`
	Weight   float64  `json:"weight" validate:"gte=1,lte=100000"`
	CountMin *float64 `json:"countMin,omitempty" validate:"omitempty,gte=1,lte=64"`
	CountMax *float64 `json:"countMax,omitempty" validate:"omitempty,gte=1,lte=64"`
}

type TreasureEntry struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	LootPools []DropRef `json:"lootPools"`
	UpdatedAt string    `json:"updatedAt"`
}

type SaveInput struct {
	ID        string    `json:"id" validate:"required,uuid_any"`
	Name      string    `json:"name" validate:"trimmed_required,trimmed_min=1,trimmed_max=80"`
	LootPools []DropRef `json:"lootPools" validate:"min=1,dive"`
}
