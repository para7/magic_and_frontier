package enemies

type DropRef struct {
	Kind     string   `json:"kind" validate:"trimmed_required,trimmed_oneof=minecraft_item item grimoire"`
	RefID    string   `json:"refId" validate:"trimmed_required,trimmed_min=1,trimmed_max=200"`
	Weight   float64  `json:"weight" validate:"gte=1,lte=100000"`
	CountMin *float64 `json:"countMin,omitempty" validate:"omitempty,gte=1,lte=64"`
	CountMax *float64 `json:"countMax,omitempty" validate:"omitempty,gte=1,lte=64"`
}

type EquipmentSlot struct {
	Kind       string   `json:"kind" validate:"trimmed_required,trimmed_oneof=minecraft_item item"`
	RefID      string   `json:"refId" validate:"trimmed_required,trimmed_min=1,trimmed_max=200"`
	Count      int      `json:"count" validate:"gte=1,lte=64"`
	DropChance *float64 `json:"dropChance,omitempty" validate:"omitempty,gte=0,lte=1"`
}

type Equipment struct {
	Mainhand *EquipmentSlot `json:"mainhand,omitempty"`
	Offhand  *EquipmentSlot `json:"offhand,omitempty"`
	Head     *EquipmentSlot `json:"head,omitempty"`
	Chest    *EquipmentSlot `json:"chest,omitempty"`
	Legs     *EquipmentSlot `json:"legs,omitempty"`
	Feet     *EquipmentSlot `json:"feet,omitempty"`
}

type EnemyEntry struct {
	ID            string    `json:"id"`
	MobType       string    `json:"mobType"`
	Name          string    `json:"name,omitempty"`
	HP            float64   `json:"hp"`
	Memo          string    `json:"memo,omitempty"`
	Attack        *float64  `json:"attack,omitempty"`
	Defense       *float64  `json:"defense,omitempty"`
	MoveSpeed     *float64  `json:"moveSpeed,omitempty"`
	Equipment     Equipment `json:"equipment"`
	EnemySkillIDs []string  `json:"enemySkillIds"`
	DropMode      string    `json:"dropMode"`
	Drops         []DropRef `json:"drops"`
	UpdatedAt     string    `json:"updatedAt"`
}

type SaveInput struct {
	ID            string    `json:"id"`
	MobType       string    `json:"mobType" validate:"trimmed_required,trimmed_min=3,trimmed_max=120"`
	Name          string    `json:"name" validate:"trimmed_max=80"`
	HP            float64   `json:"hp" validate:"gte=1,lte=100000"`
	Memo          string    `json:"memo" validate:"trimmed_max=400"`
	Attack        *float64  `json:"attack,omitempty" validate:"omitempty,gte=0,lte=100000"`
	Defense       *float64  `json:"defense,omitempty" validate:"omitempty,gte=0,lte=100000"`
	MoveSpeed     *float64  `json:"moveSpeed,omitempty" validate:"omitempty,gte=0,lte=100000"`
	Equipment     Equipment `json:"equipment"`
	EnemySkillIDs []string  `json:"enemySkillIds"`
	DropMode      string    `json:"dropMode" validate:"trimmed_required,trimmed_oneof=append replace"`
	Drops         []DropRef `json:"drops" validate:"dive"`
}
