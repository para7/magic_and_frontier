package passive

type BowConfig struct {
	LifeSub *int `json:"life_sub,omitempty" validate:"omitempty,gte=0,lte=1200"`
}

type Passive struct {
	ID          string     `json:"id"          validate:"trimmed_required,maf_slug_id"`
	Name        string     `json:"name"        validate:"trimmed_max=80"`
	Role        string     `json:"role"        validate:"trimmed_max=200"`
	Condition   string     `json:"condition"   validate:"trimmed_required,trimmed_oneof=always on_sword_hit bow"`
	Slots       []int      `json:"slots"       validate:"min=1,unique,dive,gte=1,lte=3"`
	Description string     `json:"description" validate:"trimmed_max=400"`
	Script      []string   `json:"script"      validate:"min=1"`
	Bow         *BowConfig `json:"bow,omitempty" validate:"omitempty"`
}
