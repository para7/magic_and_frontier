package item

type Item struct {
	ID        string        `json:"id"        validate:"trimmed_required"`
	Maf       ItemMaf       `json:"maf,omitempty"`
	Minecraft MinecraftItem `json:"minecraft"`
}

type ItemMaf struct {
	GrimoireID  string `json:"grimoireId,omitempty"`
	PassiveID   string `json:"passiveId,omitempty"`
	PassiveSlot int    `json:"passiveSlot,omitempty" validate:"omitempty,gte=1,lte=3"`
	BowID       string `json:"bowId,omitempty"`
}

type MinecraftItem struct {
	ItemID     string            `json:"itemId"               validate:"trimmed_required"`
	Components map[string]string `json:"components,omitempty" validate:"omitempty,dive,keys,trimmed_required,endkeys,trimmed_required"`
}
