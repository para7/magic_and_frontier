package model

// DropRef はアイテム・グリモア・バニラアイテムへの参照とドロップ設定。
// Enemy, Treasure, LootTable などで共通して使用する。
type DropRef struct {
	Kind     string   `json:"kind"               validate:"trimmed_required,trimmed_oneof=minecraft_item item grimoire passive"`
	RefID    string   `json:"refId"              validate:"trimmed_required,trimmed_min=1,trimmed_max=200"`
	Slot     *int     `json:"slot,omitempty"     validate:"omitempty,gte=1,lte=3"`
	Weight   float64  `json:"weight"             validate:"gte=1,lte=100000"`
	CountMin *float64 `json:"countMin,omitempty" validate:"omitempty,gte=1,lte=64"`
	CountMax *float64 `json:"countMax,omitempty" validate:"omitempty,gte=1,lte=64"`
}

// EquipmentSlot はエネミーの装備スロット。
type EquipmentSlot struct {
	Kind       string   `json:"kind"               validate:"trimmed_required,trimmed_oneof=minecraft_item item"`
	RefID      string   `json:"refId"              validate:"trimmed_required,trimmed_min=1,trimmed_max=200"`
	Count      int      `json:"count"              validate:"gte=1,lte=64"`
	DropChance *float64 `json:"dropChance,omitempty" validate:"omitempty,gte=0,lte=1"`
}

// Equipment はエネミーの全装備スロットをまとめた構造体。
type Equipment struct {
	Mainhand *EquipmentSlot `json:"mainhand,omitempty"`
	Offhand  *EquipmentSlot `json:"offhand,omitempty"`
	Head     *EquipmentSlot `json:"head,omitempty"`
	Chest    *EquipmentSlot `json:"chest,omitempty"`
	Legs     *EquipmentSlot `json:"legs,omitempty"`
	Feet     *EquipmentSlot `json:"feet,omitempty"`
}

// ReplacementEntry はスポーンテーブルのモブ差し替えエントリ。
type ReplacementEntry struct {
	EnemyID string `json:"enemyId"`
	Weight  int    `json:"weight"`
}
