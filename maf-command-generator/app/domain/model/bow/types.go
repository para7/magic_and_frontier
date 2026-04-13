package bow

type BowPassive struct {
	ID           string   `json:"id"                     validate:"trimmed_required,maf_slug_id"`
	Name         string   `json:"name"                   validate:"trimmed_max=80"`
	Role         string   `json:"role"                   validate:"trimmed_max=200"`
	Slots        []int    `json:"slots"`
	LifeSub      *int     `json:"life_sub,omitempty"     validate:"omitempty,gte=0,lte=1200"`
	ScriptHit    []string `json:"script_hit,omitempty"`
	ScriptFired  []string `json:"script_fired,omitempty"`
	ScriptFlying []string `json:"script_flying,omitempty"`
	ScriptGround []string `json:"script_ground,omitempty"`
}
