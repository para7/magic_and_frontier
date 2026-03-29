package item

type Item struct {
	ID                  string `json:"id"                   validate:"trimmed_required"`
	ItemID              string `json:"itemId"               validate:"trimmed_required"`
	SkillID             string `json:"skillId,omitempty"`
	CustomName          string `json:"customName"`
	Lore                string `json:"lore"`
	Enchantments        string `json:"enchantments"`
	Unbreakable         bool   `json:"unbreakable"`
	CustomModelData     string `json:"customModelData"`
	RepairCost          string `json:"repairCost"`
	HideFlags           string `json:"hideFlags"`
	PotionID            string `json:"potionId"`
	CustomPotionColor   string `json:"customPotionColor"`
	CustomPotionEffects string `json:"customPotionEffects"`
	AttributeModifiers  string `json:"attributeModifiers"`
	CustomNBT           string `json:"customNbt"`
	NBT                 string `json:"nbt"`
	UpdatedAt           string `json:"updatedAt"`
}
